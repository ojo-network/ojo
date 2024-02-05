package types

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// AxelarMessage is a copy of the Message type on Axelar
type AxelarMessage struct {
	DestinationChain   string     `json:"destination_chain"`
	DestinationAddress string     `json:"destination_address"`
	Payload            []byte     `json:"payload"`
	Type               int64      `json:"type"`
	Fee                *AxelarFee `json:"fee"` // Optional
}

// AxelarFee is a copy of the Fee type on Axelar
type AxelarFee struct {
	Amount          string  `json:"amount"`
	Recipient       string  `json:"recipient"`
	RefundRecipient *string `json:"refund_recipient"`
}

// TestUnmarshalIntoAxelarMessage tests that the message payload marshaled into an
// int array will unmarshal into a byte array and still be identical.
func TestUnmarshalIntoAxelarMessage(t *testing.T) {
	var commandSelector [4]byte
	copy(commandSelector[:], []byte{})

	encoder := NewGMPEncoder(
		[]PriceData{},
		[]string{"ATOM"},
		common.HexToAddress(""),
		commandSelector,
		[]byte{},
	)
	payload, err := encoder.GMPEncode()
	require.NoError(t, err)

	message := GmpMessage{
		DestinationChain:   "base",
		DestinationAddress: "0xa97Abf29D45DEDD18612f181383fC127da1cAa8d",
		Payload:            payload,
		Type:               1,
		Fee: &GmpFee{
			Amount:    "100000000000",
			Recipient: "axelar1zl3rxpp70lmte2xr6c4lgske2fyuj3hupcsvcd",
		},
	}
	bz, err := json.Marshal(&message)
	require.NoError(t, err)

	var axlMsg AxelarMessage
	err = json.Unmarshal(bz, &axlMsg)
	require.NoError(t, err)
	require.Equal(t, payload, axlMsg.Payload)
}
