package orchestrator

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	appparams "github.com/ojo-network/ojo/app/params"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

var (
	oracleAcceptList = []oracletypes.Denom{
		{BaseDenom: "uumee", SymbolDenom: "UMEE", Exponent: 6},
		{BaseDenom: "ibc/1", SymbolDenom: "ATOM", Exponent: 6},
		{BaseDenom: "ibc/2", SymbolDenom: "USDC", Exponent: 6},
		{BaseDenom: "ibc/3", SymbolDenom: "DAI", Exponent: 18},
		{BaseDenom: "ibc/4", SymbolDenom: "ETH", Exponent: 18},
		{BaseDenom: "ibc/5", SymbolDenom: "BTC", Exponent: 8},
		{BaseDenom: "ibc/6", SymbolDenom: "BNB", Exponent: 18},
		{BaseDenom: "ibc/7", SymbolDenom: "stATOM", Exponent: 6},
		{BaseDenom: "ibc/8", SymbolDenom: "stOSMO", Exponent: 6},
		{BaseDenom: "ibc/9", SymbolDenom: "USDT", Exponent: 6},
		{BaseDenom: "ibc/10", SymbolDenom: "IST", Exponent: 6},
	}

	oracleMandatoryList = []oracletypes.Denom{
		{BaseDenom: "ibc/1", SymbolDenom: "ATOM", Exponent: 6},
		{BaseDenom: "ibc/4", SymbolDenom: "ETH", Exponent: 18},
		{BaseDenom: "ibc/5", SymbolDenom: "BTC", Exponent: 8},
		{BaseDenom: "ibc/9", SymbolDenom: "USDT", Exponent: 6},
	}

	oracleRewardBands = []oracletypes.RewardBand{
		{SymbolDenom: "UMEE", RewardBand: math.LegacyMustNewDecFromStr("1.0")},
		{SymbolDenom: "ATOM", RewardBand: math.LegacyMustNewDecFromStr("1.0")},
		{SymbolDenom: "USDC", RewardBand: math.LegacyMustNewDecFromStr("1.0")},
		{SymbolDenom: "DAI", RewardBand: math.LegacyMustNewDecFromStr("1.0")},
		{SymbolDenom: "ETH", RewardBand: math.LegacyMustNewDecFromStr("1.0")},
		{SymbolDenom: "BTC", RewardBand: math.LegacyMustNewDecFromStr("1.0")},
		{SymbolDenom: "BNB", RewardBand: math.LegacyMustNewDecFromStr("1.0")},
		{SymbolDenom: "stATOM", RewardBand: math.LegacyMustNewDecFromStr("1.0")},
		{SymbolDenom: "stOSMO", RewardBand: math.LegacyMustNewDecFromStr("1.0")},
		{SymbolDenom: "USDT", RewardBand: math.LegacyMustNewDecFromStr("1.0")},
		{SymbolDenom: "IST", RewardBand: math.LegacyMustNewDecFromStr("1.0")},
	}
)

var (
	minGasPrice            = appparams.ProtocolMinGasPrice.String()
	minorityValidatorStake = initStakeAmount("100000000000")
	majorityValidatorStake = initStakeAmount("500000000000")
)

func initStakeAmount(amount string) sdk.Coin {
	stakeAmount, _ := math.NewIntFromString(amount)
	return sdk.NewCoin(appparams.BondDenom, stakeAmount)
}
