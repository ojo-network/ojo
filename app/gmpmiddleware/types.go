package gmpmiddleware

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

type GeneralMessageHandler interface {
	HandleGeneralMessage(ctx sdk.Context, srcChain, srcAddress string, destAddress string, payload []byte, sender string, channel string) error
	HandleGeneralMessageWithToken(ctx sdk.Context, srcChain, srcAddress string, destAddress string, payload []byte, sender string, channel string, coin sdk.Coin) error
}

// Message is attached in ICS20 packet memo field
type Message struct {
	SourceChain   string `json:"source_chain"`
	SourceAddress string `json:"source_address"`
	Payload       []byte `json:"payload"`
	Type          int64  `json:"type"`
}

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
