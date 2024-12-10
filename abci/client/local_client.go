package abcicli

import (
	"context"

	types "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/service"
	cmtsync "github.com/cometbft/cometbft/libs/sync"
)

// NOTE: use defer to unlock mutex because Application might panic (e.g., in
// case of malicious tx or query). It only makes sense for publicly exposed
// methods like CheckTx (/broadcast_tx_* RPC endpoint) or Query (/abci_query
// RPC endpoint), but defers are used everywhere for the sake of consistency.
type localClient struct {
	service.BaseService

	mtx *cmtsync.Mutex
	types.Application

	globalCbMtx cmtsync.Mutex
	globalCb    GlobalCallback
}

var _ Client = (*localClient)(nil)

// NewLocalClient creates a local client, which wraps the application interface that
// Tendermint as the client will call to the application as the server. The only
// difference, is that the local client has a global mutex which enforces serialization
// of all the ABCI calls from Tendermint to the Application.
func NewLocalClient(mtx *cmtsync.Mutex, app types.Application) Client {
	if mtx == nil {
		mtx = new(cmtsync.Mutex)
	}
	cli := &localClient{
		mtx:         mtx,
		Application: app,
	}
	cli.BaseService = *service.NewBaseService(nil, "localClient", cli)
	return cli
}

func (app *localClient) SetGlobalCallback(globalCb GlobalCallback) {
	app.globalCbMtx.Lock()
	app.globalCb = globalCb
	app.globalCbMtx.Unlock()
}

func (app *localClient) GetGlobalCallback() (cb GlobalCallback) {
	app.globalCbMtx.Lock()
	cb = app.globalCb
	app.globalCbMtx.Unlock()
	return cb
}

func (app *localClient) done(reqRes *ReqRes, res *types.Response) *ReqRes {
	set := reqRes.SetDone(res)
	if set {
		if globalCb := app.GetGlobalCallback(); globalCb != nil {
			globalCb(reqRes.Request, res)
		}
	}
	return reqRes
}

//-------------------------------------------------------

func (app *localClient) Error() error {
	return nil
}

func (app *localClient) Flush(context.Context) error {
	return nil
}

func (app *localClient) Echo(_ context.Context, msg string) (*types.ResponseEcho, error) {
	return &types.ResponseEcho{Message: msg}, nil
}

func (app *localClient) Info(ctx context.Context, req *types.RequestInfo) (*types.ResponseInfo, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.Info(ctx, req)
}

func (app *localClient) Query(ctx context.Context, req *types.RequestQuery) (*types.ResponseQuery, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.Query(ctx, req)
}

func (app *localClient) Commit(ctx context.Context, req *types.RequestCommit) (*types.ResponseCommit, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.Commit(ctx, req)
}

func (app *localClient) InitChain(ctx context.Context, req *types.RequestInitChain) (*types.ResponseInitChain, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.InitChain(ctx, req)
}

func (app *localClient) ListSnapshots(ctx context.Context, req *types.RequestListSnapshots) (*types.ResponseListSnapshots, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.ListSnapshots(ctx, req)
}

func (app *localClient) OfferSnapshot(ctx context.Context, req *types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.OfferSnapshot(ctx, req)
}

func (app *localClient) LoadSnapshotChunk(ctx context.Context,
	req *types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.LoadSnapshotChunk(ctx, req)
}

func (app *localClient) ApplySnapshotChunk(ctx context.Context,
	req *types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.ApplySnapshotChunk(ctx, req)
}

func (app *localClient) PrepareProposal(ctx context.Context, req *types.RequestPrepareProposal) (*types.ResponsePrepareProposal, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.PrepareProposal(ctx, req)
}

func (app *localClient) ProcessProposal(ctx context.Context, req *types.RequestProcessProposal) (*types.ResponseProcessProposal, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.ProcessProposal(ctx, req)
}

func (app *localClient) ExtendVote(ctx context.Context, req *types.RequestExtendVote) (*types.ResponseExtendVote, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.ExtendVote(ctx, req)
}

func (app *localClient) VerifyVoteExtension(ctx context.Context, req *types.RequestVerifyVoteExtension) (*types.ResponseVerifyVoteExtension, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.VerifyVoteExtension(ctx, req)
}

func (app *localClient) FinalizeBlock(ctx context.Context, req *types.RequestFinalizeBlock) (*types.ResponseFinalizeBlock, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.FinalizeBlock(ctx, req)
}

func (app *localClient) CheckTxSync(ctx context.Context, req *types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	// CONTRACT: Application should handle concurrent `CheckTx`
	// In this abci client layer, we don't protect `CheckTx` with a mutex for concurrency
	// app.mtx.Lock()
	// defer app.mtx.Unlock()
	return app.Application.CheckTxSyncForApp(ctx, req)
}

func (app *localClient) BeginRecheckTxSync(ctx context.Context, req *types.RequestBeginRecheckTx) (*types.ResponseBeginRecheckTx, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.BeginRecheckTx(ctx, req)
}

func (app *localClient) EndRecheckTxSync(ctx context.Context, req *types.RequestEndRecheckTx) (*types.ResponseEndRecheckTx, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.EndRecheckTx(ctx, req)
}

func (app *localClient) CheckTxAsync(ctx context.Context, params *types.RequestCheckTx, cb ResponseCallback) (*ReqRes, error) {
	req := types.ToRequestCheckTx(params)
	reqRes := NewReqRes(req, cb)
	app.Application.CheckTxAsyncForApp(ctx, params, func(r *types.ResponseCheckTx) {
		res := types.ToResponseCheckTx(r)
		app.done(reqRes, res)
	})
	return reqRes, nil
}

func (app *localClient) BeginRecheckTxAsync(ctx context.Context, req *types.RequestBeginRecheckTx, cb ResponseCallback) (*ReqRes, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	reqRes := NewReqRes(types.ToRequestBeginRecheckTx(req), cb)
	res, _ := app.Application.BeginRecheckTx(ctx, req)
	return app.done(reqRes, types.ToResponseBeginRecheckTx(res)), nil
}

func (app *localClient) EndRecheckTxAsync(ctx context.Context, req *types.RequestEndRecheckTx, cb ResponseCallback) (*ReqRes, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	reqRes := NewReqRes(types.ToRequestEndRecheckTx(req), cb)
	res, _ := app.Application.EndRecheckTx(ctx, req)
	return app.done(reqRes, types.ToResponseEndRecheckTx(res)), nil
}

func (app *localClient) CheckTxSyncForApp(context.Context, *types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	panic("not implemented")
}

func (app *localClient) CheckTxAsyncForApp(context.Context, *types.RequestCheckTx, types.CheckTxCallback) {
	panic("not implemented")
}

func (app *localClient) BeginRecheckTx(ctx context.Context, params *types.RequestBeginRecheckTx) (*types.ResponseBeginRecheckTx, error) {
	panic("not implemented")
}

func (app *localClient) EndRecheckTx(ctx context.Context, params *types.RequestEndRecheckTx) (*types.ResponseEndRecheckTx, error) {
	panic("not implemented")
}
