// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	types "github.com/cometbft/cometbft/abci/types"
	mock "github.com/stretchr/testify/mock"
)

// Application is an autogenerated mock type for the Application type
type Application struct {
	mock.Mock
}

// ApplySnapshotChunk provides a mock function with given fields: _a0, _a1
func (_m *Application) ApplySnapshotChunk(_a0 context.Context, _a1 *types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ApplySnapshotChunk")
	}

	var r0 *types.ResponseApplySnapshotChunk
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestApplySnapshotChunk) *types.ResponseApplySnapshotChunk); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseApplySnapshotChunk)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestApplySnapshotChunk) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BeginRecheckTx provides a mock function with given fields: _a0, _a1
func (_m *Application) BeginRecheckTx(_a0 context.Context, _a1 *types.RequestBeginRecheckTx) (*types.ResponseBeginRecheckTx, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for BeginRecheckTx")
	}

	var r0 *types.ResponseBeginRecheckTx
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestBeginRecheckTx) (*types.ResponseBeginRecheckTx, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestBeginRecheckTx) *types.ResponseBeginRecheckTx); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseBeginRecheckTx)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestBeginRecheckTx) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckTxAsyncForApp provides a mock function with given fields: _a0, _a1, _a2
func (_m *Application) CheckTxAsyncForApp(_a0 context.Context, _a1 *types.RequestCheckTx, _a2 types.CheckTxCallback) {
	_m.Called(_a0, _a1, _a2)
}

// CheckTxSyncForApp provides a mock function with given fields: _a0, _a1
func (_m *Application) CheckTxSyncForApp(_a0 context.Context, _a1 *types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CheckTxSyncForApp")
	}

	var r0 *types.ResponseCheckTx
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestCheckTx) (*types.ResponseCheckTx, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestCheckTx) *types.ResponseCheckTx); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseCheckTx)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestCheckTx) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Commit provides a mock function with given fields: _a0, _a1
func (_m *Application) Commit(_a0 context.Context, _a1 *types.RequestCommit) (*types.ResponseCommit, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Commit")
	}

	var r0 *types.ResponseCommit
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestCommit) (*types.ResponseCommit, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestCommit) *types.ResponseCommit); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseCommit)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestCommit) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EndRecheckTx provides a mock function with given fields: _a0, _a1
func (_m *Application) EndRecheckTx(_a0 context.Context, _a1 *types.RequestEndRecheckTx) (*types.ResponseEndRecheckTx, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for EndRecheckTx")
	}

	var r0 *types.ResponseEndRecheckTx
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestEndRecheckTx) (*types.ResponseEndRecheckTx, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestEndRecheckTx) *types.ResponseEndRecheckTx); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseEndRecheckTx)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestEndRecheckTx) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExtendVote provides a mock function with given fields: _a0, _a1
func (_m *Application) ExtendVote(_a0 context.Context, _a1 *types.RequestExtendVote) (*types.ResponseExtendVote, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ExtendVote")
	}

	var r0 *types.ResponseExtendVote
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestExtendVote) (*types.ResponseExtendVote, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestExtendVote) *types.ResponseExtendVote); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseExtendVote)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestExtendVote) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FinalizeBlock provides a mock function with given fields: _a0, _a1
func (_m *Application) FinalizeBlock(_a0 context.Context, _a1 *types.RequestFinalizeBlock) (*types.ResponseFinalizeBlock, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for FinalizeBlock")
	}

	var r0 *types.ResponseFinalizeBlock
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestFinalizeBlock) (*types.ResponseFinalizeBlock, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestFinalizeBlock) *types.ResponseFinalizeBlock); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseFinalizeBlock)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestFinalizeBlock) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Info provides a mock function with given fields: _a0, _a1
func (_m *Application) Info(_a0 context.Context, _a1 *types.RequestInfo) (*types.ResponseInfo, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Info")
	}

	var r0 *types.ResponseInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestInfo) (*types.ResponseInfo, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestInfo) *types.ResponseInfo); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestInfo) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InitChain provides a mock function with given fields: _a0, _a1
func (_m *Application) InitChain(_a0 context.Context, _a1 *types.RequestInitChain) (*types.ResponseInitChain, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for InitChain")
	}

	var r0 *types.ResponseInitChain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestInitChain) (*types.ResponseInitChain, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestInitChain) *types.ResponseInitChain); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseInitChain)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestInitChain) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListSnapshots provides a mock function with given fields: _a0, _a1
func (_m *Application) ListSnapshots(_a0 context.Context, _a1 *types.RequestListSnapshots) (*types.ResponseListSnapshots, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ListSnapshots")
	}

	var r0 *types.ResponseListSnapshots
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestListSnapshots) (*types.ResponseListSnapshots, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestListSnapshots) *types.ResponseListSnapshots); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseListSnapshots)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestListSnapshots) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadSnapshotChunk provides a mock function with given fields: _a0, _a1
func (_m *Application) LoadSnapshotChunk(_a0 context.Context, _a1 *types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for LoadSnapshotChunk")
	}

	var r0 *types.ResponseLoadSnapshotChunk
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestLoadSnapshotChunk) *types.ResponseLoadSnapshotChunk); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseLoadSnapshotChunk)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestLoadSnapshotChunk) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OfferSnapshot provides a mock function with given fields: _a0, _a1
func (_m *Application) OfferSnapshot(_a0 context.Context, _a1 *types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for OfferSnapshot")
	}

	var r0 *types.ResponseOfferSnapshot
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestOfferSnapshot) *types.ResponseOfferSnapshot); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseOfferSnapshot)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestOfferSnapshot) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrepareProposal provides a mock function with given fields: _a0, _a1
func (_m *Application) PrepareProposal(_a0 context.Context, _a1 *types.RequestPrepareProposal) (*types.ResponsePrepareProposal, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PrepareProposal")
	}

	var r0 *types.ResponsePrepareProposal
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestPrepareProposal) (*types.ResponsePrepareProposal, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestPrepareProposal) *types.ResponsePrepareProposal); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponsePrepareProposal)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestPrepareProposal) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProcessProposal provides a mock function with given fields: _a0, _a1
func (_m *Application) ProcessProposal(_a0 context.Context, _a1 *types.RequestProcessProposal) (*types.ResponseProcessProposal, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ProcessProposal")
	}

	var r0 *types.ResponseProcessProposal
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestProcessProposal) (*types.ResponseProcessProposal, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestProcessProposal) *types.ResponseProcessProposal); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseProcessProposal)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestProcessProposal) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: _a0, _a1
func (_m *Application) Query(_a0 context.Context, _a1 *types.RequestQuery) (*types.ResponseQuery, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Query")
	}

	var r0 *types.ResponseQuery
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestQuery) (*types.ResponseQuery, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestQuery) *types.ResponseQuery); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseQuery)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestQuery) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyVoteExtension provides a mock function with given fields: _a0, _a1
func (_m *Application) VerifyVoteExtension(_a0 context.Context, _a1 *types.RequestVerifyVoteExtension) (*types.ResponseVerifyVoteExtension, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for VerifyVoteExtension")
	}

	var r0 *types.ResponseVerifyVoteExtension
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestVerifyVoteExtension) (*types.ResponseVerifyVoteExtension, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestVerifyVoteExtension) *types.ResponseVerifyVoteExtension); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseVerifyVoteExtension)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestVerifyVoteExtension) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewApplication creates a new instance of Application. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApplication(t interface {
	mock.TestingT
	Cleanup(func())
}) *Application {
	mock := &Application{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
