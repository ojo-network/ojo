package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	appparams "github.com/ojo-network/ojo/app/params"
	ojosim "github.com/ojo-network/ojo/util/sim"
	"github.com/ojo-network/ojo/x/airdrop/keeper"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

const (
	OpWeightMsgCreateAirdropAccount = "op_weight_msg_create_airdrop_account" //nolint: gosec
	OpWeightMsgClaimAirdrop         = "op_weight_msg_claim_airdrop"          //nolint: gosec

	DefaultWeightMsgSend = 100 // from simappparams.DefaultWeightMsgSend
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONCodec,
	ak types.AccountKeeper,
	bk bankkeeper.Keeper,
	k keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgCreateAirdropAccount int
		weightMsgClaimAirdrop         int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgCreateAirdropAccount, &weightMsgCreateAirdropAccount, nil,
		func(_ *rand.Rand) {
			weightMsgCreateAirdropAccount = DefaultWeightMsgSend * 2
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgClaimAirdrop, &weightMsgClaimAirdrop, nil,
		func(_ *rand.Rand) {
			weightMsgClaimAirdrop = DefaultWeightMsgSend * 2
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgClaimAirdrop,
			SimulateMsgClaimAirdrop(ak, bk, k),
		),
	}
}

func SimulateMsgClaimAirdrop(ak types.AccountKeeper, bk bankkeeper.Keeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		airdropAccounts := k.GetAllAirdropAccounts(ctx)
		if len(airdropAccounts) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgClaimAirdrop, "no Airdrop Accounts"), nil, nil
		}
		airdropAccount := randomAirdropAccount(r, airdropAccounts)

		originAddr, err := sdk.AccAddressFromBech32(airdropAccount.OriginAddress)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgClaimAirdrop, "bad Airdrop Account Address"), nil, err
		}

		originSimAcct, found := simtypes.FindAccount(accs, originAddr)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgClaimAirdrop, "simulation acct does not exist"), nil, nil
		}

		claimSimAddr := simtypes.RandomAccounts(r, 1)[0]
		claimAddr := ak.GetAccount(ctx, claimSimAddr.Address)

		msg := types.NewMsgClaimAirdrop(
			airdropAccount.OriginAddress,
			claimAddr.String(),
		)

		return deliver(r, app, ctx, ak, bk, originSimAcct, msg, nil)
	}
}

func randomAirdropAccount(r *rand.Rand, airdropAccounts []*types.AirdropAccount) *types.AirdropAccount {
	return airdropAccounts[r.Intn(len(airdropAccounts))]
}

func deliver(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, ak simulation.AccountKeeper,
	bk bankkeeper.Keeper, from simtypes.Account, msg sdk.Msg, coins sdk.Coins,
) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
	cfg := testutil.MakeTestEncodingConfig()
	o := simulation.OperationInput{
		R:               r,
		App:             app,
		TxGen:           cfg.TxConfig,
		Cdc:             cfg.Codec.(*codec.ProtoCodec),
		Msg:             msg,
		MsgType:         sdk.MsgTypeURL(msg),
		Context:         ctx,
		SimAccount:      from,
		AccountKeeper:   ak,
		Bankkeeper:      bk,
		ModuleName:      types.ModuleName,
		CoinsSpentInMsg: coins,
	}

	return ojosim.GenAndDeliver(bk, o, appparams.DefaultGasLimit*50)
}
