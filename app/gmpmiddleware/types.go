package gmpmiddleware

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gmptypes "github.com/ojo-network/ojo/x/gmp/types"
)

type GeneralMessageHandler interface {
	HandleGeneralMessage(
		ctx sdk.Context,
		srcChain,
		srcAddress string,
		destAddress string,
		payload []byte,
		sender string,
		channel string,
	) error
	HandleGeneralMessageWithToken(
		ctx sdk.Context,
		srcChain,
		srcAddress string,
		destAddress string,
		payload []byte,
		sender string,
		channel string,
		coin sdk.Coin,
	) error
}

// Message is attached in ICS20 packet memo field
type Message struct {
	SourceChain   string `json:"source_chain"`
	SourceAddress string `json:"source_address"`
	Payload       []byte `json:"payload"`
	Type          int64  `json:"type"`
}

// AxelarPayload is the payload sent from Axelar to IBC middleware.
// It needs to be decoded using the ABI encoded data.
type GmpPayload struct {
	AssetNames      []string
	ContractAddress string
	CommandSelector [4]byte
	CommandParams   []byte
	Timestamp       uint64
}

// axelarPayloadSpec is the ABI spec for the AxelarPayload struct.
//
//nolint:lll
const gmpPayloadSpec = `[{"constant":false,"inputs":[{"name":"assetNames","type":"bytes32[]"},{"name":"contractAddress","type":"address"},{"name":"commandSelector","type":"bytes4"},{"name":"commandParams","type":"bytes"},{"name":"timestamp","type":"uint256"}],"name":"EncodedData","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`

// parseDenom convert denom to receiver chain representation
func parseDenom(packet channeltypes.Packet, denom string) string {
	if types.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), denom) {
		// remove prefix added by sender chain
		voucherPrefix := types.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())
		unprefixedDenom := denom[len(voucherPrefix):]

		// coin denomination used in sending from the escrow address
		denom = unprefixedDenom

		// The denomination used to send the coins is either the native denom or the hash of the path
		// if the denomination is not native.
		denomTrace := types.ParseDenomTrace(unprefixedDenom)
		if denomTrace.Path != "" {
			denom = denomTrace.IBCDenom()
		}

		return denom
	}

	prefixedDenom := types.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel()) + denom
	denom = types.ParseDenomTrace(prefixedDenom).IBCDenom()

	return denom
}

// parsePayload takes an encoded payload and decodes it into an EncodedData object.
func parsePayload(payload []byte) (GmpPayload, error) {
	parsedABI, err := abi.JSON(strings.NewReader(gmpPayloadSpec))
	if err != nil {
		return GmpPayload{}, err
	}

	var decodedData GmpPayload
	err = parsedABI.UnpackIntoInterface(&decodedData, "EncodedData", payload)
	if err != nil {
		return GmpPayload{}, err
	}

	return decodedData, nil
}

func verifyParams(params gmptypes.Params, sender string, channel string) error {
	if !strings.EqualFold(params.GmpAddress, sender) {
		return fmt.Errorf("invalid sender address: %s", sender)
	}
	if !strings.EqualFold(params.GmpChannel, channel) {
		return fmt.Errorf("invalid channel: %s", channel)
	}
	return nil
}
