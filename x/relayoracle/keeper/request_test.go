package keeper_test

import "github.com/ojo-network/ojo/x/relayoracle/types"

func (s *IntegrationTestSuite) TestAddRequest() {
	req := types.Request{
		RequestCallData: nil,
		ClientID:        "",
		RequestHeight:   0,
		RequestTime:     0,
		IBCChannel:      nil,
	}

	id := s.app.RelayOracle.AddRequest(s.ctx, req)

	s.Require().EqualValues(1, id)

	storedReq, err := s.app.RelayOracle.GetRequest(s.ctx, id)
	s.Require().NoError(err)
	s.Require().Equal(req, storedReq)

	// check pending list
	reqIds := s.app.RelayOracle.GetPendingRequestList(s.ctx)
	s.Require().Len(reqIds, 1)

	s.Require().EqualValues(reqIds[0], 1)
}

func (s *IntegrationTestSuite) TestGetSetRequestCount() {
	s.app.RelayOracle.SetRequestCount(s.ctx, 5)

	count := s.app.RelayOracle.GetRequestCount(s.ctx)
	s.Require().Equal(uint64(5), count)
}

func (s *IntegrationTestSuite) TestSetGetRequest() {
	req := types.Request{
		RequestCallData: nil,
		ClientID:        "",
		RequestHeight:   0,
		RequestTime:     0,
		IBCChannel:      nil,
	}

	s.app.RelayOracle.SetRequest(s.ctx, 1, req)

	storedReq, err := s.app.RelayOracle.GetRequest(s.ctx, 1)
	s.Require().NoError(err)
	s.Require().Equal(req, storedReq)
}

func (s *IntegrationTestSuite) TestDeleteRequest() {
	req := types.Request{
		RequestCallData: nil,
		ClientID:        "",
		RequestHeight:   0,
		RequestTime:     0,
		IBCChannel:      nil,
	}

	s.app.RelayOracle.SetRequest(s.ctx, 1, req)
	s.app.RelayOracle.DeleteRequest(s.ctx, 1)

	_, err := s.app.RelayOracle.GetRequest(s.ctx, 1)
	s.Require().ErrorIs(err, types.ErrRequestNotFound)
}

func (s *IntegrationTestSuite) TestSetGetResult() {
	res := types.Result{
		RequestID:       100,
		RequestCallData: nil,
		ClientID:        "",
		RequestHeight:   0,
		RequestTime:     0,
		Status:          0,
		Result:          nil,
	}

	s.app.RelayOracle.SetResult(s.ctx, res)
	storedResBz := s.ctx.KVStore(s.app.GetKey(types.StoreKey)).Get(types.ResultStoreKey(res.RequestID))

	var storedRes types.Result
	s.app.AppCodec().MustUnmarshal(storedResBz, &storedRes)

	s.Require().Equal(res, storedRes)
}

func (s *IntegrationTestSuite) TestAddRequestIDToPendingList() {
	reqID := uint64(1)
	s.app.RelayOracle.AddRequestIDToPendingList(s.ctx, reqID)

	pendingRequestIds := s.app.RelayOracle.GetPendingRequestList(s.ctx)

	s.Require().Contains(pendingRequestIds, reqID)
}

func (s *IntegrationTestSuite) TestFlushPendingRequestList() {
	reqID := uint64(1)
	s.app.RelayOracle.AddRequestIDToPendingList(s.ctx, reqID)

	s.app.RelayOracle.FlushPendingRequestList(s.ctx)

	// Try to get the pending request list from the store
	storedPendingBz := s.ctx.KVStore(s.app.GetKey(types.StoreKey)).Get(types.PendingRequestListKey)
	s.Require().Nil(storedPendingBz)

	pendingRequestIds := s.app.RelayOracle.GetPendingRequestList(s.ctx)
	s.Require().Empty(pendingRequestIds)
}

func (s *IntegrationTestSuite) TestProcessResult() {
	req := types.Request{
		RequestCallData: nil,
		ClientID:        "",
		RequestHeight:   0,
		RequestTime:     0,
		IBCChannel:      nil,
	}

	id := s.app.RelayOracle.AddRequest(s.ctx, req)

	status := types.RESOLVE_STATUS_SUCCESS
	resultData := []byte("dummyResult")

	s.app.RelayOracle.ProcessResult(s.ctx, id, status, resultData)
	storedResBz := s.ctx.KVStore(s.app.GetKey(types.StoreKey)).Get(types.ResultStoreKey(id))

	var storedRes types.Result
	s.app.AppCodec().MustUnmarshal(storedResBz, &storedRes)

	s.Require().Equal(id, storedRes.RequestID)
	s.Require().Equal(status, storedRes.Status)
	s.Require().Equal(resultData, storedRes.Result)
}
