package keeper

//import (
//	"testing"
//
//	"github.com/ojo-network/ojo/x/qoracle/keeper"
//	"github.com/ojo-network/ojo/x/qoracle/types"
//
//	"github.com/cosmos/cosmos-sdk/codec"
//	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
//	"github.com/cosmos/cosmos-sdk/store"
//	storetypes "github.com/cosmos/cosmos-sdk/store/types"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
//	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
//	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
//	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
//	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
//	"github.com/stretchr/testify/require"
//	"github.com/cometbft/cometbft/libs/log"
//	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
//	tmdb "github.com/cometbft/cometbft-db"
//)
//
//// qoracleChannelKeeper is a stub of cosmosibckeeper.ChannelKeeper.
//type qoracleChannelKeeper struct{}
//
//func (qoracleChannelKeeper) GetChannel(ctx sdk.Context, portID, channelID string) (channeltypes.Channel, bool) {
//	return channeltypes.Channel{}, false
//}
//
//func (qoracleChannelKeeper) GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool) {
//	return 0, false
//}
//
//func (qoracleChannelKeeper) SendPacket(
//    ctx sdk.Context,
//    channelCap *capabilitytypes.Capability,
//    sourcePort string,
//    sourceChannel string,
//    timeoutHeight clienttypes.Height,
//    timeoutTimestamp uint64,
//    data []byte,
//) (uint64, error) {
//    return 0, nil
//}
//
//func (qoracleChannelKeeper) ChanCloseInit(ctx sdk.Context, portID, channelID string, chanCap *capabilitytypes.Capability) error {
//	return nil
//}
//
//// qoracleportKeeper is a stub of cosmosibckeeper.PortKeeper
//type qoraclePortKeeper struct{}
//
//func (qoraclePortKeeper) BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability {
//	return &capabilitytypes.Capability{}
//}
//
//
//
//func QoracleKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
//	logger := log.NewNopLogger()
//
//	storeKey := sdk.NewKVStoreKey(types.StoreKey)
//	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
//
//	db := tmdb.NewMemDB()
//	stateStore := store.NewCommitMultiStore(db)
//	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
//	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
//	require.NoError(t, stateStore.LoadLatestVersion())
//
//	registry := codectypes.NewInterfaceRegistry()
//	appCodec := codec.NewProtoCodec(registry)
//	capabilityKeeper := capabilitykeeper.NewKeeper(appCodec, storeKey, memStoreKey)
//
//	paramsSubspace := typesparams.NewSubspace(appCodec,
//		types.Amino,
//		storeKey,
//		memStoreKey,
//		"QoracleParams",
//	)
//	k := keeper.NewKeeper(
//        appCodec,
//        storeKey,
//        memStoreKey,
//        paramsSubspace,
//        qoracleChannelKeeper{},
//        qoraclePortKeeper{},
//        capabilityKeeper.ScopeToModule("QoracleScopedKeeper"),
//    )
//
//	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, logger)
//
//	// Initialize params
//	k.SetParams(ctx, types.DefaultParams())
//
//	return k, ctx
//}
