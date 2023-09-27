package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (packet *OracleRequestPacketData) ValidateBasic() error {
	return nil
}

func NewIbcChannel(port, channel string) *IBCChannel {
	return &IBCChannel{
		ChannelId: channel,
		PortId:    port,
	}
}

func NewRequest(calldata []byte, clientID string, channel *IBCChannel) Request {
	return Request{
		RequestCallData: calldata,
		ClientID:        clientID,
		IBCChannel:      channel,
	}
}

func NewRequestOracleAcknowledgement(id uint64) *OracleRequestPacketAcknowledgement {
	return &OracleRequestPacketAcknowledgement{
		RequestID: id,
	}
}

func NewOracleResponsePacketData(
	clientID string,
	requestID uint64,
	requestTime,
	resolveTime int64,
	resolveStatus ResolveStatus,
	result []byte,
) *OracleResponsePacketData {
	return &OracleResponsePacketData{
		ClientID:      clientID,
		RequestID:     requestID,
		RequestTime:   requestTime,
		ResolveTime:   resolveTime,
		ResolveStatus: resolveStatus,
		Result:        result,
	}
}

func (o *OracleResponsePacketData) ToBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(o))
}
