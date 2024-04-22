package types

import (
	fmt "fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// GmpDecoder is the payload sent from Axelar to IBC middleware.
// It needs to be decoded using the ABI.
type GmpDecoder struct {
	AssetNames      [][32]byte
	ContractAddress common.Address
	CommandSelector [4]byte
	CommandParams   []byte
	Timestamp       *big.Int
}

// abiSpec is the ABI specification for the GMP data.
var decoderSpec = abi.Arguments{
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

// NewGmpDecoder decodes a payload from GMP given a byte array
func NewGmpDecoder(payload []byte) (GmpDecoder, error) {
	args, err := decoderSpec.Unpack(payload)
	if err != nil {
		return GmpDecoder{}, err
	}

	// check to make sure each argument is the correct type
	//nolint: all
	if assetNames, ok := args[0].([][32]byte); !ok {
		return GmpDecoder{}, fmt.Errorf("invalid asset names type: %T", args[0])
	} else if contractAddress, ok := args[1].(common.Address); !ok {
		return GmpDecoder{}, fmt.Errorf("invalid contract address type: %T", args[1])
	} else if commandSelector, ok := args[2].([4]byte); !ok {
		return GmpDecoder{}, fmt.Errorf("invalid command selector type: %T", args[2])
	} else if commandParams, ok := args[3].([]byte); !ok {
		return GmpDecoder{}, fmt.Errorf("invalid command params type: %T", args[3])
	} else if timestamp, ok := args[4].(*big.Int); !ok {
		return GmpDecoder{}, fmt.Errorf("invalid timestamp type: %T", args[4])
	} else {
		return GmpDecoder{
			AssetNames:      assetNames,
			ContractAddress: contractAddress,
			CommandSelector: commandSelector,
			CommandParams:   commandParams,
			Timestamp:       timestamp,
		}, nil
	}
}

func (g GmpDecoder) GetDenoms() []string {
	denoms := make([]string, len(g.AssetNames))
	for i, name := range g.AssetNames {
		denoms[i] = strings.TrimRight(string(name[:]), "\x00")
	}
	return denoms
}
