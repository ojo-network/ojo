package gmpmiddleware

import (
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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

// GmpData is the payload sent from Axelar to IBC middleware.
// It needs to be decoded using the ABI.
type GmpData struct {
	AssetNames      [][32]byte     // bytes32 in Solidity is represented as [32]byte in Go
	ContractAddress common.Address // address in Solidity is represented as common.Address in Go
	CommandSelector [4]byte        // bytes4 in Solidity is represented as [4]byte in Go
	CommandParams   []byte         // bytes in Solidity is represented as []byte in Go
	Timestamp       *big.Int       // uint256 in Solidity is represented as *big.Int in Go
	AbiEncodedData  []byte         // the ABI encoded data
}

var (
	assetNamesType, _      = abi.NewType("bytes32[]", "bytes32[]", nil)
	contractAddressType, _ = abi.NewType("address", "address", nil)
	commandSelectorType, _ = abi.NewType("bytes4", "bytes4", nil)
	commandParamsType, _   = abi.NewType("bytes", "bytes", nil)
	timestampType, _       = abi.NewType("uint256", "uint256", nil)
)

var abiSpec = abi.Arguments{
	{
		Type: assetNamesType,
	},
	{
		Type: contractAddressType,
	},
	{
		Type: commandSelectorType,
	},
	{
		Type: commandParamsType,
	},
	{
		Type: timestampType,
	},
}

// NewGmpData decodes a payload from GMP given a byte array
func NewGmpData(payload []byte) (GmpData, error) {
	args, err := abiSpec.Unpack(payload)
	if err != nil {
		return GmpData{}, err
	}

	// check to make sure each argument is the correct type
	if assetNames, ok := args[0].([][32]byte); !ok {
		return GmpData{}, fmt.Errorf("invalid asset names type: %T", args[0])
	} else if contractAddress, ok := args[1].(common.Address); !ok {
		return GmpData{}, fmt.Errorf("invalid contract address type: %T", args[1])
	} else if commandSelector, ok := args[2].([4]byte); !ok {
		return GmpData{}, fmt.Errorf("invalid command selector type: %T", args[2])
	} else if commandParams, ok := args[3].([]byte); !ok {
		return GmpData{}, fmt.Errorf("invalid command params type: %T", args[3])
	} else if timestamp, ok := args[4].(*big.Int); !ok {
		return GmpData{}, fmt.Errorf("invalid timestamp type: %T", args[4])
	} else {
		return GmpData{
			AssetNames:      assetNames,
			ContractAddress: contractAddress,
			CommandSelector: commandSelector,
			CommandParams:   commandParams,
			Timestamp:       timestamp,
		}, nil
	}
}

// assets takes a GmpData and returns the asset names as a slice of strings
func (g GmpData) assets() []string {
	var assetNames []string
	for _, assetName := range g.AssetNames {
		assetNames = append(assetNames, string(assetName[:]))
	}
	return assetNames
}

func (g GmpData) Encode() ([]byte, error) {
	return abiSpec.Pack(
		g.AssetNames,
		g.ContractAddress,
		g.CommandSelector,
		g.CommandParams,
		g.Timestamp,
	)
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
