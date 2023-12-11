package types

import "github.com/ethereum/go-ethereum/accounts/abi"

// These are the types we use to encode and decode data to and from the GMP.
var (
	assetNamesType, _      = abi.NewType("bytes32[]", "bytes32[]", nil)
	contractAddressType, _ = abi.NewType("address", "address", nil)
	commandSelectorType, _ = abi.NewType("bytes4", "bytes4", nil)
	commandParamsType, _   = abi.NewType("bytes", "bytes", nil)
	timestampType, _       = abi.NewType("uint256", "uint256", nil)

	// priceDataType is the ABI specification for the PriceData tuple in Solidity.
	// It is a tuple of (bytes32, uint256, uint256, tuple).
	// It includes MedianData, another tuple of (uint256[], uint256[], uint256[]).
	priceDataType, _ = abi.NewType("tuple[]", "",
		[]abi.ArgumentMarshaling{
			{Name: "AssetName", Type: "bytes32"},
			{Name: "Price", Type: "uint256"},
			{Name: "ResolveTime", Type: "uint256"},
			{
				Name: "MedianData", Type: "tuple", Components: []abi.ArgumentMarshaling{
					{Name: "BlockNums", Type: "uint256[]"},
					{Name: "Medians", Type: "uint256[]"},
					{Name: "Deviations", Type: "uint256[]"},
				},
			},
		},
	)
)
