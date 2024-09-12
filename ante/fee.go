package ante

import (
	"cosmossdk.io/math"
	evidencetypes "cosmossdk.io/x/evidence/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/util"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

// MaxMsgGasUsage defines the maximum gas allowed for an oracle transaction.
const MaxMsgGasUsage = uint64(140_000)

// FeeAndPriority ensures tx has enough fee coins to pay for the gas at the CheckTx time
// to early remove transactions from the mempool without enough attached fee.
// The validator min fee check is ignored if the tx contains only oracle messages and
// tx gas limit is <= MaxMsgGasUsage. Essentially, validators can provide price
// transactison for free as long as the gas per message is in the MaxMsgGasUsage limit.
func FeeAndPriority(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return nil, 0, sdkerrors.ErrTxDecode.Wrap("Tx must be a FeeTx")
	}

	providedFees := feeTx.GetFee()
	gasLimit := feeTx.GetGas()
	msgs := feeTx.GetMsgs()
	isOracle := IsOracleTx(msgs)
	priority := getTxPriority(isOracle, msgs)
	chargeFees := !isOracle || gasLimit > uint64(len(msgs))*MaxMsgGasUsage
	// We also don't charge fees for the genesis transactions.
	if !chargeFees || ctx.BlockHeight() == 0 {
		return sdk.Coins{}, priority, nil
	}

	if ctx.IsCheckTx() {
		return providedFees, priority, checkFees(ctx.MinGasPrices(), providedFees, gasLimit)
	}
	return providedFees, priority, checkFees(nil, providedFees, gasLimit)
}

func checkFees(minGasPrices sdk.DecCoins, fees sdk.Coins, gasLimit uint64) error {
	if minGasPrices != nil {
		// check minGasPrices set by validator
		if err := AssertMinProtocolGasPrice(minGasPrices); err != nil {
			return err
		}
	} else {
		// in deliverTx = use protocol min gas price
		minGasPrices = sdk.DecCoins{appparams.ProtocolMinGasPrice}
	}

	requiredFees := sdk.NewCoins()

	// Determine the required fees by multiplying each required minimum gas
	// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
	// Zero fees are removed.
	glDec := math.LegacyNewDec(util.SafeUint64ToInt64(gasLimit))
	for _, gp := range minGasPrices {
		if gasLimit == 0 || gp.IsZero() {
			continue
		}
		fee := gp.Amount.Mul(glDec)
		requiredFees = append(requiredFees, sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt()))
	}

	if !requiredFees.Empty() && !fees.IsAnyGTE(requiredFees) {
		return sdkerrors.ErrInsufficientFee.Wrapf(
			"insufficient fees; got: %s required: %s", fees, requiredFees)
	}
	return nil
}

// IsOracleTx checks if all messages are oracle messages
func IsOracleTx(msgs []sdk.Msg) bool {
	if len(msgs) == 0 {
		return false
	}
	for _, msg := range msgs {
		switch msg.(type) {
		case *oracletypes.MsgAggregateExchangeRatePrevote,
			*oracletypes.MsgAggregateExchangeRateVote:
			continue
		default:
			return false
		}
	}

	return true
}

// AssertMinProtocolGasPrice returns an error if the provided gasPrices are lower then
// the required by protocol.
func AssertMinProtocolGasPrice(gasPrices sdk.DecCoins) error {
	if gasPrices.AmountOf(appparams.ProtocolMinGasPrice.Denom).LT(appparams.ProtocolMinGasPrice.Amount) {
		return sdkerrors.ErrInsufficientFee.Wrapf(
			"gas price too small; got: %v required min: %v", gasPrices, appparams.ProtocolMinGasPrice)
	}

	return nil
}

// getTxPriority returns naive tx priority based on the lowest fee amount (regardless of the
// denom) and oracle tx check.
// Dirty optimization: since we already check if msgs are oracle message, then we
// don't recomupte it again: isOracle flag takes a precedence over msgs check.
func getTxPriority(isOracle bool, msgs []sdk.Msg) int64 {
	var priority int64
	if isOracle {
		return 100
	}
	for _, msg := range msgs {
		var p int64
		switch msg.(type) {
		case *evidencetypes.MsgSubmitEvidence:
			p = 90
		default:
			// in case there is a non-prioritized message, return 0
			return 0
		}
		if priority == 0 || p < priority {
			priority = p
		}
	}

	return priority
}
