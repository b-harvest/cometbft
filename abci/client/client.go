package abcicli

import (
	"context"
	"fmt"
	"sync"

	"github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/service"
	cmtsync "github.com/cometbft/cometbft/libs/sync"
)

const (
	dialRetryIntervalSeconds = 3
	echoRetryIntervalSeconds = 1
)

//go:generate ../../scripts/mockery_generate.sh Client

// Client defines the interface for an ABCI client.
//
// NOTE these are client errors, eg. ABCI socket connectivity issues.
// Application-related errors are reflected in response via ABCI error codes
// and (potentially) error response.
type Client interface {
	service.Service
	types.Application

	// TODO: remove as each method now returns an error
	Error() error
	// TODO: remove as this is not implemented
	Flush(context.Context) error
	Echo(context.Context, string) (*types.ResponseEcho, error)

	// FIXME: All other operations are run synchronously and rely
	// on the caller to dictate concurrency (i.e. run a go routine),
	// with the exception of `CheckTxAsync` which we maintain
	// for the v0 mempool. We should explore refactoring the
	// mempool to remove this vestige behavior.
	SetGlobalCallback(GlobalCallback)
	GetGlobalCallback() GlobalCallback

	CheckTxSync(context.Context, *types.RequestCheckTx) (*types.ResponseCheckTx, error)
	BeginRecheckTxSync(context.Context, *types.RequestBeginRecheckTx) (*types.ResponseBeginRecheckTx, error) // Signals the beginning of rechecking
	EndRecheckTxSync(context.Context, *types.RequestEndRecheckTx) (*types.ResponseEndRecheckTx, error)       // Signals the end of rechecking

	CheckTxAsync(context.Context, *types.RequestCheckTx, ResponseCallback) (*ReqRes, error)
	BeginRecheckTxAsync(context.Context, *types.RequestBeginRecheckTx, ResponseCallback) (*ReqRes, error)
	EndRecheckTxAsync(context.Context, *types.RequestEndRecheckTx, ResponseCallback) (*ReqRes, error)
}

//----------------------------------------

// NewClient returns a new ABCI client of the specified transport type.
// It returns an error if the transport is not "socket" or "grpc"
func NewClient(addr, transport string, mustConnect bool) (client Client, err error) {
	switch transport {
	case "socket":
		client = NewSocketClient(addr, mustConnect)
	case "grpc":
		client = NewGRPCClient(addr, mustConnect)
	default:
		err = fmt.Errorf("unknown abci transport %s", transport)
	}
	return
}

type GlobalCallback func(*types.Request, *types.Response)
type ResponseCallback func(*types.Response)

type ReqRes struct {
	*types.Request
	*types.Response // Not set atomically, so be sure to use WaitGroup.

	mtx cmtsync.Mutex

	wg   *sync.WaitGroup
	done bool
	cb   func(*types.Response) // A single callback that may be set.
}

func NewReqRes(req *types.Request, cb ResponseCallback) *ReqRes {
	return &ReqRes{
		Request:  req,
		Response: nil,

		wg:   waitGroup1(),
		done: false,
		cb:   cb,
	}
}

// InvokeCallback invokes a thread-safe execution of the configured callback
// if non-nil.
func (reqRes *ReqRes) InvokeCallback() {
	reqRes.mtx.Lock()
	defer reqRes.mtx.Unlock()

	if reqRes.cb != nil {
		reqRes.cb(reqRes.Response)
	}
}

func (r *ReqRes) SetDone(res *types.Response) (set bool) {
	r.mtx.Lock()
	// TODO should we panic if it's already done?
	set = !r.done
	if set {
		r.Response = res
		r.done = true
		r.wg.Done()
	}
	r.mtx.Unlock()

	// NOTE `r.cb` is immutable so we're safe to access it at here without `mtx`
	if set && r.cb != nil {
		r.cb(res)
	}

	return set
}

func (r *ReqRes) Wait() {
	r.wg.Wait()
}

func waitGroup1() (wg *sync.WaitGroup) {
	wg = &sync.WaitGroup{}
	wg.Add(1)
	return
}
