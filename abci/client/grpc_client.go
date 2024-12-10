package abcicli

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cometbft/cometbft/abci/types"
	cmtnet "github.com/cometbft/cometbft/libs/net"
	"github.com/cometbft/cometbft/libs/service"
)

var _ Client = (*grpcClient)(nil)

// A stripped copy of the remoteClient that makes
// synchronous calls using grpc
type grpcClient struct {
	service.BaseService
	mustConnect bool

	client types.ABCIClient
	conn   *grpc.ClientConn

	mtx  sync.Mutex
	addr string
	err  error

	globalCbMtx sync.Mutex
	globalCb    func(*types.Request, *types.Response) // listens to all callbacks
}

func NewGRPCClient(addr string, mustConnect bool) Client {
	cli := &grpcClient{
		addr:        addr,
		mustConnect: mustConnect,
	}
	cli.BaseService = *service.NewBaseService(nil, "grpcClient", cli)
	return cli
}

func dialerFunc(_ context.Context, addr string) (net.Conn, error) {
	return cmtnet.Connect(addr)
}

func (cli *grpcClient) OnStart() error {
	if err := cli.BaseService.OnStart(); err != nil {
		return err
	}

RETRY_LOOP:
	for {
		conn, err := grpc.NewClient(cli.addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(dialerFunc),
		)
		if err != nil {
			if cli.mustConnect {
				return err
			}
			cli.Logger.Error(fmt.Sprintf("abci.grpcClient failed to connect to %v.  Retrying...\n", cli.addr), "err", err)
			time.Sleep(time.Second * dialRetryIntervalSeconds)
			continue RETRY_LOOP
		}

		cli.Logger.Info("Dialed server. Waiting for echo.", "addr", cli.addr)
		client := types.NewABCIClient(conn)
		cli.conn = conn

	ENSURE_CONNECTED:
		for {
			_, err := client.Echo(context.Background(), &types.RequestEcho{Message: "hello"}, grpc.WaitForReady(true))
			if err == nil {
				break ENSURE_CONNECTED
			}
			cli.Logger.Error("Echo failed", "err", err)
			time.Sleep(time.Second * echoRetryIntervalSeconds)
		}

		cli.client = client
		return nil
	}
}

func (cli *grpcClient) OnStop() {
	cli.BaseService.OnStop()

	if cli.conn != nil {
		cli.conn.Close()
	}
}

func (cli *grpcClient) StopForError(err error) {
	if !cli.IsRunning() {
		return
	}

	cli.mtx.Lock()
	if cli.err == nil {
		cli.err = err
	}
	cli.mtx.Unlock()

	cli.Logger.Error(fmt.Sprintf("Stopping abci.grpcClient for error: %v", err.Error()))
	if err := cli.Stop(); err != nil {
		cli.Logger.Error("Error stopping abci.grpcClient", "err", err)
	}
}

func (cli *grpcClient) Error() error {
	cli.mtx.Lock()
	defer cli.mtx.Unlock()
	return cli.err
}

func (cli *grpcClient) SetGlobalCallback(globalCb GlobalCallback) {
	cli.globalCbMtx.Lock()
	cli.globalCb = globalCb
	cli.globalCbMtx.Unlock()
}

func (cli *grpcClient) GetGlobalCallback() (cb GlobalCallback) {
	cli.globalCbMtx.Lock()
	cb = cli.globalCb
	cli.globalCbMtx.Unlock()
	return cb
}

//----------------------------------------

func (cli *grpcClient) finishAsyncCall(req *types.Request, res *types.Response, cb ResponseCallback) *ReqRes {
	reqRes := NewReqRes(req, cb)

	// goroutine for callbacks
	go func() {
		set := reqRes.SetDone(res)
		if set {
			// Notify client listener if set
			if globalCb := cli.GetGlobalCallback(); globalCb != nil {
				globalCb(req, res)
			}
		}
	}()

	return reqRes
}

//----------------------------------------

func (cli *grpcClient) Flush(ctx context.Context) error {
	_, err := cli.client.Flush(ctx, types.ToRequestFlush().GetFlush(), grpc.WaitForReady(true))
	return err
}

func (cli *grpcClient) Echo(ctx context.Context, msg string) (*types.ResponseEcho, error) {
	return cli.client.Echo(ctx, types.ToRequestEcho(msg).GetEcho(), grpc.WaitForReady(true))
}

func (cli *grpcClient) Info(ctx context.Context, req *types.RequestInfo) (*types.ResponseInfo, error) {
	return cli.client.Info(ctx, req, grpc.WaitForReady(true))
}

func (cli *grpcClient) Query(ctx context.Context, req *types.RequestQuery) (*types.ResponseQuery, error) {
	return cli.client.Query(ctx, types.ToRequestQuery(req).GetQuery(), grpc.WaitForReady(true))
}

func (cli *grpcClient) Commit(ctx context.Context, _ *types.RequestCommit) (*types.ResponseCommit, error) {
	return cli.client.Commit(ctx, types.ToRequestCommit().GetCommit(), grpc.WaitForReady(true))
}

func (cli *grpcClient) InitChain(ctx context.Context, req *types.RequestInitChain) (*types.ResponseInitChain, error) {
	return cli.client.InitChain(ctx, types.ToRequestInitChain(req).GetInitChain(), grpc.WaitForReady(true))
}

func (cli *grpcClient) ListSnapshots(ctx context.Context, req *types.RequestListSnapshots) (*types.ResponseListSnapshots, error) {
	return cli.client.ListSnapshots(ctx, types.ToRequestListSnapshots(req).GetListSnapshots(), grpc.WaitForReady(true))
}

func (cli *grpcClient) OfferSnapshot(ctx context.Context, req *types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error) {
	return cli.client.OfferSnapshot(ctx, types.ToRequestOfferSnapshot(req).GetOfferSnapshot(), grpc.WaitForReady(true))
}

func (cli *grpcClient) LoadSnapshotChunk(ctx context.Context, req *types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error) {
	return cli.client.LoadSnapshotChunk(ctx, types.ToRequestLoadSnapshotChunk(req).GetLoadSnapshotChunk(), grpc.WaitForReady(true))
}

func (cli *grpcClient) ApplySnapshotChunk(ctx context.Context, req *types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error) {
	return cli.client.ApplySnapshotChunk(ctx, types.ToRequestApplySnapshotChunk(req).GetApplySnapshotChunk(), grpc.WaitForReady(true))
}

func (cli *grpcClient) PrepareProposal(ctx context.Context, req *types.RequestPrepareProposal) (*types.ResponsePrepareProposal, error) {
	return cli.client.PrepareProposal(ctx, types.ToRequestPrepareProposal(req).GetPrepareProposal(), grpc.WaitForReady(true))
}

func (cli *grpcClient) ProcessProposal(ctx context.Context, req *types.RequestProcessProposal) (*types.ResponseProcessProposal, error) {
	return cli.client.ProcessProposal(ctx, types.ToRequestProcessProposal(req).GetProcessProposal(), grpc.WaitForReady(true))
}

func (cli *grpcClient) ExtendVote(ctx context.Context, req *types.RequestExtendVote) (*types.ResponseExtendVote, error) {
	return cli.client.ExtendVote(ctx, types.ToRequestExtendVote(req).GetExtendVote(), grpc.WaitForReady(true))
}

func (cli *grpcClient) VerifyVoteExtension(ctx context.Context, req *types.RequestVerifyVoteExtension) (*types.ResponseVerifyVoteExtension, error) {
	return cli.client.VerifyVoteExtension(ctx, types.ToRequestVerifyVoteExtension(req).GetVerifyVoteExtension(), grpc.WaitForReady(true))
}

func (cli *grpcClient) FinalizeBlock(ctx context.Context, req *types.RequestFinalizeBlock) (*types.ResponseFinalizeBlock, error) {
	return cli.client.FinalizeBlock(ctx, types.ToRequestFinalizeBlock(req).GetFinalizeBlock(), grpc.WaitForReady(true))
}

func (cli *grpcClient) CheckTxSync(ctx context.Context, req *types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	return cli.client.CheckTx(ctx, req, grpc.WaitForReady(true))
}

func (cli *grpcClient) BeginRecheckTxSync(ctx context.Context, params *types.RequestBeginRecheckTx) (*types.ResponseBeginRecheckTx, error) {
	reqres, _ := cli.BeginRecheckTxAsync(ctx, params, nil)
	reqres.Wait()
	return reqres.Response.GetBeginRecheckTx(), cli.Error()
}

func (cli *grpcClient) EndRecheckTxSync(ctx context.Context, params *types.RequestEndRecheckTx) (*types.ResponseEndRecheckTx, error) {
	reqres, _ := cli.EndRecheckTxAsync(ctx, params, nil)
	reqres.Wait()
	return reqres.Response.GetEndRecheckTx(), cli.Error()
}

func (cli *grpcClient) CheckTxAsync(ctx context.Context, req *types.RequestCheckTx, cb ResponseCallback) (*ReqRes, error) {
	res, err := cli.client.CheckTx(ctx, req, grpc.WaitForReady(true))
	if err != nil {
		cli.StopForError(err)
		return nil, err
	}
	return cli.finishAsyncCall(types.ToRequestCheckTx(req), &types.Response{Value: &types.Response_CheckTx{CheckTx: res}}, cb), nil
}

func (cli *grpcClient) BeginRecheckTxAsync(ctx context.Context, params *types.RequestBeginRecheckTx, cb ResponseCallback) (*ReqRes, error) {
	req := types.ToRequestBeginRecheckTx(params)
	res, err := cli.client.BeginRecheckTx(ctx, req.GetBeginRecheckTx(), grpc.WaitForReady(true))
	if err != nil {
		cli.StopForError(err)
	}
	return cli.finishAsyncCall(req, &types.Response{Value: &types.Response_BeginRecheckTx{BeginRecheckTx: res}}, cb), nil
}

func (cli *grpcClient) EndRecheckTxAsync(ctx context.Context, params *types.RequestEndRecheckTx, cb ResponseCallback) (*ReqRes, error) {
	req := types.ToRequestEndRecheckTx(params)
	res, err := cli.client.EndRecheckTx(ctx, req.GetEndRecheckTx(), grpc.WaitForReady(true))
	if err != nil {
		cli.StopForError(err)
	}
	return cli.finishAsyncCall(req, &types.Response{Value: &types.Response_EndRecheckTx{EndRecheckTx: res}}, cb), nil
}

func (cli *grpcClient) CheckTxSyncForApp(context.Context, *types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	panic("not implemented")
}

func (cli *grpcClient) CheckTxAsyncForApp(context.Context, *types.RequestCheckTx, types.CheckTxCallback) {
	panic("not implemented")
}

func (cli *grpcClient) BeginRecheckTx(ctx context.Context, params *types.RequestBeginRecheckTx) (*types.ResponseBeginRecheckTx, error) {
	panic("not implemented")
}

func (cli *grpcClient) EndRecheckTx(ctx context.Context, params *types.RequestEndRecheckTx) (*types.ResponseEndRecheckTx, error) {
	panic("not implemented")
}
