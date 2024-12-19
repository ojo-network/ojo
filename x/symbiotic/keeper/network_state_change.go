package keeper

import (
	"context"
	"cosmossdk.io/math"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ojo-network/ojo/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	symbiotictypes "github.com/ojo-network/ojo/x/symbiotic/types"
)

// Struct to unmarshal the response from the Beacon Chain API
type Block struct {
	Finalized bool `json:"finalized"`
	Data      struct {
		Message struct {
			Body struct {
				ExecutionPayload struct {
					BlockHash string `json:"block_hash"`
				} `json:"execution_payload"`
			} `json:"body"`
		} `json:"message"`
	} `json:"data"`
}

type RPCRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type RPCResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Validator struct {
	Stake    *big.Int
	ConsAddr [32]byte
}

const (
	SLEEP_ON_RETRY                  = 200
	RETRIES                         = 5
	BEACON_GENESIS_TIMESTAMP        = 1695902400
	SLOTS_IN_EPOCH                  = 32
	SLOT_DURATION                   = 12
	INVALID_BLOCKHASH               = "invalid"
	BLOCK_PATH                      = "/eth/v2/beacon/blocks/"
	GET_VALIDATOR_SET_FUNCTION_NAME = "getValidatorSet"
	GET_CURRENT_EPOCH_FUNCTION_NAME = "getCurrentEpoch"
	CONTRACT_ABI                    = `[
		{
			"type": "function",
			"name": "getCurrentEpoch",
			"outputs": [
				{
					"name": "epoch",
					"type": "uint48",
					"internalType": "uint48"
				}
			],
			"stateMutability": "view"
		},
		{
			"type": "function",
			"name": "getValidatorSet",
			"inputs": [
				{
					"name": "epoch",
					"type": "uint48",
					"internalType": "uint48"
				}
			],
			"outputs": [
				{
					"name": "validatorsData",
					"type": "tuple[]",
					"internalType": "struct SimpleMiddleware.ValidatorData[]",
					"components": [
						{
							"name": "stake",
							"type": "uint256",
							"internalType": "uint256"
						},
						{
							"name": "consAddr",
							"type": "bytes32",
							"internalType": "bytes32"
						}
					]
				}
			],
			"stateMutability": "view"
		}
	]`
)

func (k Keeper) SymbioticUpdateValidatorsPower(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(sdkCtx)

	if params.MiddlewareAddress == "" {
		panic("middleware address is not set")
	}

	height := sdkCtx.BlockHeight()

	if height%params.SymbioticSyncPeriod != 0 {
		return nil
	}

	cachedBlockHash, err := k.GetCachedBlockHash(sdkCtx, util.SafeInt64ToUint64(height))
	if err != nil {
		return err
	}

	if cachedBlockHash.BlockHash == INVALID_BLOCKHASH {
		return nil
	}

	var validators []Validator

	for i := 0; i < RETRIES; i++ {
		validators, err = k.getSymbioticValidatorSet(ctx, cachedBlockHash.BlockHash)
		if err == nil {
			break
		}

		if strings.HasSuffix(err.Error(), "is not currently canonical") {
			k.Logger(sdkCtx).Warn("not canonical block hash", "hash", cachedBlockHash.BlockHash)
			err = nil
			break
		}

		k.apiUrls.RotateEthUrl()
		time.Sleep(time.Millisecond * SLEEP_ON_RETRY)
	}

	if err != nil {
		return err
	}

	for _, v := range validators {
		val, err := k.StakingKeeper.GetValidatorByConsAddr(ctx, v.ConsAddr[:20])
		if err != nil {
			if errors.Is(err, stakingtypes.ErrNoValidatorFound) {
				continue
			}
			return err
		}
		tokens := math.NewIntFromBigInt(v.Stake)

		if tokens.GT(val.GetTokens()) {
			k.StakingKeeper.AddValidatorTokensAndShares(ctx, val, tokens.Sub(val.GetTokens()))
		} else {
			k.StakingKeeper.RemoveValidatorTokens(ctx, val, val.GetTokens().Sub(tokens))
		}
	}

	return nil
}

func (k Keeper) GetFinalizedBlockHash(ctx context.Context) (string, error) {
	var err error
	var block Block

	for i := 0; i < RETRIES; i++ {
		slot := k.getSlot(ctx)
		block, err = k.parseBlock(ctx, slot)

		for errors.Is(err, symbiotictypes.ErrSymbioticNotFound) { // some slots on api may be omitted
			for i := 1; i < SLOTS_IN_EPOCH; i++ {
				block, err = k.parseBlock(ctx, slot-i)
				if err == nil {
					break
				}
				if !errors.Is(err, symbiotictypes.ErrSymbioticNotFound) {
					return "", err
				}
			}
		}

		if err == nil {
			break
		}

		k.apiUrls.RotateBeaconUrl()
		time.Sleep(time.Millisecond * SLEEP_ON_RETRY)
	}

	if err != nil {
		return "", err
	}

	if !block.Finalized {
		return INVALID_BLOCKHASH, nil
	}

	return block.Data.Message.Body.ExecutionPayload.BlockHash, nil
}

func (k Keeper) GetBlockByHash(ctx context.Context, blockHash string) (*types.Block, error) {
	var block *types.Block
	client, err := ethclient.Dial(k.apiUrls.GetEthApiUrl())
	if err != nil {
		return nil, err
	}
	defer client.Close()

	for i := 0; i < RETRIES; i++ {
		block, err = client.BlockByHash(ctx, common.HexToHash(blockHash))
		if err == nil {
			break
		}

		k.apiUrls.RotateEthUrl()
		time.Sleep(time.Millisecond * SLEEP_ON_RETRY)
	}

	if err != nil {
		return nil, err
	}

	return block, nil
}

func (k Keeper) GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	var block *types.Block
	client, err := ethclient.Dial(k.apiUrls.GetEthApiUrl())
	if err != nil {
		return nil, err
	}
	defer client.Close()

	for i := 0; i < RETRIES; i++ {
		block, err = client.BlockByNumber(ctx, number)
		if err == nil {
			break
		}

		k.apiUrls.RotateEthUrl()
		time.Sleep(time.Millisecond * SLEEP_ON_RETRY)
	}

	if err != nil {
		return nil, err
	}

	return block, nil
}

func (k Keeper) GetMinBlockTimestamp(ctx context.Context) uint64 {
	return uint64(k.getSlot(ctx)-SLOTS_IN_EPOCH)*12 + BEACON_GENESIS_TIMESTAMP
}

func (k Keeper) getSymbioticValidatorSet(ctx context.Context, blockHash string) ([]Validator, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(sdkCtx)

	client, err := ethclient.Dial(k.apiUrls.GetEthApiUrl())
	if err != nil {
		k.Logger(sdkCtx).With(err).Error("rpc error: ethclient dial error", "url", k.apiUrls.GetEthApiUrl())
		return nil, err
	}
	defer client.Close()

	contractABI, err := abi.JSON(strings.NewReader(CONTRACT_ABI))
	if err != nil {
		return nil, err
	}

	contractAddress := common.HexToAddress(params.MiddlewareAddress)

	data, err := contractABI.Pack(GET_CURRENT_EPOCH_FUNCTION_NAME)
	if err != nil {
		return nil, err
	}

	query := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}
	result, err := client.CallContractAtHash(ctx, query, common.HexToHash(blockHash))
	if err != nil {
		k.Logger(sdkCtx).With(err).Error("rpc error: eth_call error", "url", k.apiUrls.GetEthApiUrl())
		return nil, err
	}

	currentEpoch := new(big.Int).SetBytes(result)

	data, err = contractABI.Pack(GET_VALIDATOR_SET_FUNCTION_NAME, currentEpoch)
	if err != nil {
		return nil, err
	}

	query = ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}
	result, err = client.CallContractAtHash(ctx, query, common.HexToHash(blockHash))
	if err != nil {
		k.Logger(sdkCtx).With(err).Error("rpc error: eth_call error", "url", k.apiUrls.GetEthApiUrl())
		return nil, err
	}

	var validators []Validator
	err = contractABI.UnpackIntoInterface(&validators, GET_VALIDATOR_SET_FUNCTION_NAME, result)
	if err != nil {
		return nil, err
	}

	return validators, nil
}

func (k Keeper) getSlot(ctx context.Context) int {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	slot := (sdkCtx.BlockHeader().Time.Unix() - BEACON_GENESIS_TIMESTAMP) / SLOT_DURATION // get beacon slot
	slot = slot / SLOTS_IN_EPOCH * SLOTS_IN_EPOCH                                         // first slot of epoch
	slot -= 3 * SLOTS_IN_EPOCH                                                            // get finalized slot
	return int(slot)
}

func (k Keeper) parseBlock(ctx context.Context, slot int) (Block, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	url := k.apiUrls.GetBeaconApiUrl() + BLOCK_PATH + strconv.Itoa(slot)

	var block Block
	resp, err := http.Get(url)
	if err != nil {
		k.Logger(sdkCtx).With(err).Error("rpc error: beacon rpc call error", "url", url, "err", err)
		return block, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		k.Logger(sdkCtx).Error("rpc error: beacon rpc call error", "url", k.apiUrls.GetEthApiUrl(), "err", "no err", "status", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusNotFound {
		return block, symbiotictypes.ErrSymbioticNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return block, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return block, fmt.Errorf("error reading response body: %v", err)
	}

	err = json.Unmarshal(body, &block)
	if err != nil {
		return block, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return block, nil
}
