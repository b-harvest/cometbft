// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	abcicli "github.com/cometbft/cometbft/abci/client"

	mock "github.com/stretchr/testify/mock"

	types "github.com/cometbft/cometbft/abci/types"
)

// AppConnMempool is an autogenerated mock type for the AppConnMempool type
type AppConnMempool struct {
	mock.Mock
}

// BeginRecheckTx provides a mock function with given fields: _a0, _a1
func (_m *AppConnMempool) BeginRecheckTx(_a0 context.Context, _a1 *types.RequestBeginRecheckTx) (*types.ResponseBeginRecheckTx, error) {
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

// BeginRecheckTxAsync provides a mock function with given fields: _a0, _a1
func (_m *AppConnMempool) BeginRecheckTxAsync(_a0 context.Context, _a1 *types.RequestBeginRecheckTx) (*abcicli.ReqRes, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for BeginRecheckTxAsync")
	}

	var r0 *abcicli.ReqRes
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestBeginRecheckTx) (*abcicli.ReqRes, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestBeginRecheckTx) *abcicli.ReqRes); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestBeginRecheckTx) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckTx provides a mock function with given fields: _a0, _a1
func (_m *AppConnMempool) CheckTx(_a0 context.Context, _a1 *types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CheckTx")
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

// CheckTxAsync provides a mock function with given fields: _a0, _a1
func (_m *AppConnMempool) CheckTxAsync(_a0 context.Context, _a1 *types.RequestCheckTx) (*abcicli.ReqRes, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CheckTxAsync")
	}

	var r0 *abcicli.ReqRes
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestCheckTx) (*abcicli.ReqRes, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestCheckTx) *abcicli.ReqRes); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestCheckTx) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EndRecheckTx provides a mock function with given fields: _a0, _a1
func (_m *AppConnMempool) EndRecheckTx(_a0 context.Context, _a1 *types.RequestEndRecheckTx) (*types.ResponseEndRecheckTx, error) {
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

// EndRecheckTxAsync provides a mock function with given fields: _a0, _a1
func (_m *AppConnMempool) EndRecheckTxAsync(_a0 context.Context, _a1 *types.RequestEndRecheckTx) (*abcicli.ReqRes, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for EndRecheckTxAsync")
	}

	var r0 *abcicli.ReqRes
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestEndRecheckTx) (*abcicli.ReqRes, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.RequestEndRecheckTx) *abcicli.ReqRes); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*abcicli.ReqRes)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.RequestEndRecheckTx) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Error provides a mock function with no fields
func (_m *AppConnMempool) Error() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Error")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Flush provides a mock function with given fields: _a0
func (_m *AppConnMempool) Flush(_a0 context.Context) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Flush")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetResponseCallback provides a mock function with given fields: _a0
func (_m *AppConnMempool) SetResponseCallback(_a0 abcicli.Callback) {
	_m.Called(_a0)
}

// NewAppConnMempool creates a new instance of AppConnMempool. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAppConnMempool(t interface {
	mock.TestingT
	Cleanup(func())
}) *AppConnMempool {
	mock := &AppConnMempool{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
