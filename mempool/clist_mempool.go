package mempool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/libs/clist"
	"github.com/cometbft/cometbft/libs/log"
	cmtmath "github.com/cometbft/cometbft/libs/math"
	cmtsync "github.com/cometbft/cometbft/libs/sync"
	"github.com/cometbft/cometbft/proxy"
	"github.com/cometbft/cometbft/types"
)

// CListMempool is an ordered in-memory pool for transactions before they are
// proposed in a consensus round. Transaction validity is checked using the
// CheckTx abci message before the transaction is added to the pool. The
// mempool uses a concurrent list structure for storing transactions that can
// be efficiently accessed by multiple concurrent readers.
type CListMempool struct {
	height   atomic.Int64 // the last block Update()'d to
	txsBytes atomic.Int64 // total size of mempool, in bytes

	// notify listeners (ie. consensus) when txs are available
	notifiedTxsAvailable atomic.Bool
	txsAvailable         chan struct{} // fires once for each height, when the mempool is not empty

	config *config.MempoolConfig

	// Exclusive mutex for Update method to prevent concurrent execution of
	// CheckTx or ReapMaxBytesMaxGas(ReapMaxTxs) methods.
	updateMtx cmtsync.RWMutex
	preCheck  PreCheckFunc
	postCheck PostCheckFunc

	chReqCheckTx chan *requestCheckTxAsync

	txs          *clist.CList // concurrent linked-list of good txs
	proxyAppConn proxy.AppConnMempool

	// Map for quick access to txs to record sender in CheckTx.
	// txsMap: txKey -> CElement
	txsMap sync.Map

	// Keep a cache of already-seen txs.
	// This reduces the pressure on the proxyApp.
	cache TxCache

	logger  log.Logger
	metrics *Metrics
}

type requestCheckTxAsync struct {
	tx        types.Tx
	txInfo    TxInfo
	prepareCb func(error)
	checkTxCb func(*abci.Response)
}

var _ Mempool = &CListMempool{}

// CListMempoolOption sets an optional parameter on the mempool.
type CListMempoolOption func(*CListMempool)

// NewCListMempool returns a new mempool with the given configuration and
// connection to an application.
func NewCListMempool(
	cfg *config.MempoolConfig,
	proxyAppConn proxy.AppConnMempool,
	height int64,
	options ...CListMempoolOption,
) *CListMempool {
	mp := &CListMempool{
		config:       cfg,
		proxyAppConn: proxyAppConn,
		txs:          clist.New(),
		chReqCheckTx: make(chan *requestCheckTxAsync, cfg.Size),
		logger:       log.NewNopLogger(),
		metrics:      NopMetrics(),
	}
	mp.height.Store(height)

	if cfg.CacheSize > 0 {
		mp.cache = NewLRUTxCache(cfg.CacheSize)
	} else {
		mp.cache = NopTxCache{}
	}
	proxyAppConn.SetResponseCallback(mp.globalCb)

	for _, option := range options {
		option(mp)
	}
	go mp.checkTxAsyncReactor()
	return mp
}

func (mem *CListMempool) getCElement(txKey types.TxKey) (*clist.CElement, bool) {
	if e, ok := mem.txsMap.Load(txKey); ok {
		return e.(*clist.CElement), true
	}
	return nil, false
}

func (mem *CListMempool) getMemTx(txKey types.TxKey) *mempoolTx {
	if e, ok := mem.getCElement(txKey); ok {
		return e.Value.(*mempoolTx)
	}
	return nil
}

func (mem *CListMempool) removeAllTxs() {
	for e := mem.txs.Front(); e != nil; e = e.Next() {
		mem.txs.Remove(e)
		e.DetachPrev()
	}

	mem.txsMap.Range(func(key, _ interface{}) bool {
		mem.txsMap.Delete(key)
		return true
	})
}

// NOTE: not thread safe - should only be called once, on startup
func (mem *CListMempool) EnableTxsAvailable() {
	mem.txsAvailable = make(chan struct{}, 1)
}

// SetLogger sets the Logger.
func (mem *CListMempool) SetLogger(l log.Logger) {
	mem.logger = l
}

// WithPreCheck sets a filter for the mempool to reject a tx if f(tx) returns
// false. This is ran before CheckTx. Only applies to the first created block.
// After that, Update overwrites the existing value.
func WithPreCheck(f PreCheckFunc) CListMempoolOption {
	return func(mem *CListMempool) { mem.preCheck = f }
}

// WithPostCheck sets a filter for the mempool to reject a tx if f(tx) returns
// false. This is ran after CheckTx. Only applies to the first created block.
// After that, Update overwrites the existing value.
func WithPostCheck(f PostCheckFunc) CListMempoolOption {
	return func(mem *CListMempool) { mem.postCheck = f }
}

// WithMetrics sets the metrics.
func WithMetrics(metrics *Metrics) CListMempoolOption {
	return func(mem *CListMempool) { mem.metrics = metrics }
}

// Safe for concurrent use by multiple goroutines.
func (mem *CListMempool) Lock() {
	mem.updateMtx.Lock()
}

// Safe for concurrent use by multiple goroutines.
func (mem *CListMempool) Unlock() {
	mem.updateMtx.Unlock()
}

// Safe for concurrent use by multiple goroutines.
func (mem *CListMempool) Size() int {
	return mem.txs.Len()
}

// Safe for concurrent use by multiple goroutines.
func (mem *CListMempool) SizeBytes() int64 {
	return mem.txsBytes.Load()
}

// Lock() must be help by the caller during execution.
func (mem *CListMempool) FlushAppConn() error {
	err := mem.proxyAppConn.Flush(context.TODO())
	if err != nil {
		return ErrFlushAppConn{Err: err}
	}

	return nil
}

// XXX: Unsafe! Calling Flush may leave mempool in inconsistent state.
func (mem *CListMempool) Flush() {
	mem.updateMtx.RLock()
	defer mem.updateMtx.RUnlock()

	mem.txsBytes.Store(0)
	mem.cache.Reset()

	mem.removeAllTxs()
}

// TxsFront returns the first transaction in the ordered list for peer
// goroutines to call .NextWait() on.
// FIXME: leaking implementation details!
//
// Safe for concurrent use by multiple goroutines.
func (mem *CListMempool) TxsFront() *clist.CElement {
	return mem.txs.Front()
}

// TxsWaitChan returns a channel to wait on transactions. It will be closed
// once the mempool is not empty (ie. the internal `mem.txs` has at least one
// element)
//
// Safe for concurrent use by multiple goroutines.
func (mem *CListMempool) TxsWaitChan() <-chan struct{} {
	return mem.txs.WaitChan()
}

// It blocks if we're waiting on Update() or Reap().
// Safe for concurrent use by multiple goroutines.
func (mem *CListMempool) CheckTxSync(tx types.Tx, txInfo TxInfo) (res *abci.Response, err error) {
	mem.updateMtx.RLock()
	// use defer to unlock mutex because application (*local client*) might panic
	defer mem.updateMtx.RUnlock()

	if err = mem.prepareCheckTx(tx, txInfo); err != nil {
		return res, err
	}

	// CONTRACT: `app.CheckTxSync()` should check whether `GasWanted` is valid (0 <= GasWanted <= block.masGas)
	var r *abci.ResponseCheckTx
	r, err = mem.proxyAppConn.CheckTxSync(context.TODO(), &abci.RequestCheckTx{Tx: tx})
	if err != nil {
		return res, err
	}

	// TODO refactor to pass a `pointer` directly
	res = abci.ToResponseCheckTx(r)
	mem.reqResCb(tx, txInfo, res, nil)
	return res, err
}

// cb: A callback from the CheckTx command.
//
//	It gets called from another goroutine.
//
// Safe for concurrent use by multiple goroutines.
func (mem *CListMempool) CheckTxAsync(tx types.Tx, txInfo TxInfo, prepareCb func(error), checkTxCb func(*abci.Response)) {
	mem.chReqCheckTx <- &requestCheckTxAsync{tx: tx, txInfo: txInfo, prepareCb: prepareCb, checkTxCb: checkTxCb}
}

func (mem *CListMempool) checkTxAsyncReactor() {
	for req := range mem.chReqCheckTx {
		mem.checkTxAsync(req.tx, req.txInfo, req.prepareCb, req.checkTxCb)
	}
}

// It blocks if we're waiting on Update() or Reap().
func (mem *CListMempool) checkTxAsync(tx types.Tx, txInfo TxInfo, prepareCb func(error), checkTxCb func(*abci.Response)) {
	mem.updateMtx.RLock()
	defer func() {
		if r := recover(); r != nil {
			mem.updateMtx.RUnlock()
			panic(r)
		}
	}()

	err := mem.prepareCheckTx(tx, txInfo)
	if prepareCb != nil {
		prepareCb(err)
	}
	if err != nil {
		mem.updateMtx.RUnlock()
		return
	}

	// CONTRACT: `app.CheckTxAsync()` should check whether `GasWanted` is valid (0 <= GasWanted <= block.masGas)
	reqRes, _ := mem.proxyAppConn.CheckTxAsync(context.TODO(), &abci.RequestCheckTx{Tx: tx})
	reqRes.SetCallback(func(res *abci.Response) {
		mem.reqResCb(tx, txInfo, res, func(response *abci.Response) {
			if checkTxCb != nil {
				checkTxCb(response)
			}
			mem.updateMtx.RUnlock()
		})
	})
}

// It blocks if we're waiting on Update() or Reap().
// cb: A callback from the CheckTx command.
//
//	It gets called from another goroutine.
//
// CONTRACT: Either cb will get called, or err returned.
//
// Safe for concurrent use by multiple goroutines.
func (mem *CListMempool) prepareCheckTx(tx types.Tx, txInfo TxInfo) error {
	txSize := len(tx)

	if err := mem.isFull(txSize); err != nil {
		return err
	}

	if txSize > mem.config.MaxTxBytes {
		return ErrTxTooLarge{
			Max:    mem.config.MaxTxBytes,
			Actual: txSize,
		}
	}

	if mem.preCheck != nil {
		if err := mem.preCheck(tx); err != nil {
			return ErrPreCheck{Err: err}
		}
	}

	// NOTE: proxyAppConn may error if tx buffer is full
	if err := mem.proxyAppConn.Error(); err != nil {
		return ErrAppConnMempool{Err: err}
	}

	if !mem.cache.Push(tx) { // if the transaction already exists in the cache
		// Record a new sender for a tx we've already seen.
		// Note it's possible a tx is still in the cache but no longer in the mempool
		// (eg. after committing a block, txs are removed from mempool but not cache),
		// so we only record the sender for txs still in the mempool.
		if memTx := mem.getMemTx(tx.Key()); memTx != nil {
			memTx.addSender(txInfo.SenderID)
			// TODO: consider punishing peer for dups,
			// its non-trivial since invalid txs can become valid,
			// but they can spam the same tx with little cost to them atm.
		}
		return ErrTxInCache
	}

	return nil
}

// Global callback that will be called after every ABCI response.
// Having a single global callback avoids needing to set a callback for each request.
// However, processing the checkTx response requires the peerID (so we can track which txs we heard from who),
// and peerID is not included in the ABCI request, so we have to set request-specific callbacks that
// include this information. If we're not in the midst of a recheck, this function will just return,
// so the request specific callback can do the work.
//
// When rechecking, we don't need the peerID, so the recheck callback happens
// here.
// Global callback that will be called after every ABCI response.
func (mem *CListMempool) globalCb(req *abci.Request, res *abci.Response) {
	checkTxReq := req.GetCheckTx()
	if checkTxReq == nil {
		return
	}

	if checkTxReq.Type == abci.CheckTxType_Recheck {
		mem.metrics.RecheckTimes.Add(1)
		mem.resCbRecheck(req, res)

		// update metrics
		mem.metrics.Size.Set(float64(mem.Size()))
	}
}

// Request specific callback that should be set on individual reqRes objects
// to incorporate local information when processing the response.
// This allows us to track the peer that sent us this tx, so we can avoid sending it back to them.
// NOTE: alternatively, we could include this information in the ABCI request itself.
//
// External callers of CheckTx, like the RPC, can also pass an externalCb through here that is called
// when all other response processing is complete.
//
// Used in CheckTx to record PeerID who sent us the tx.
func (mem *CListMempool) reqResCb(
	tx []byte,
	txInfo TxInfo,
	res *abci.Response,
	externalCb func(*abci.Response),
) {
	mem.resCbFirstTime(tx, txInfo, res)

	// update metrics
	mem.metrics.Size.Set(float64(mem.Size()))
	mem.metrics.SizeBytes.Set(float64(mem.SizeBytes()))

	// passed in by the caller of CheckTx, eg. the RPC
	if externalCb != nil {
		externalCb(res)
	}
}

// Called from:
//   - resCbFirstTime (lock not held) if tx is valid
func (mem *CListMempool) addTx(memTx *mempoolTx) {
	e := mem.txs.PushBack(memTx)
	mem.txsMap.Store(memTx.tx.Key(), e)
	mem.txsBytes.Add(int64(len(memTx.tx)))
	mem.metrics.TxSizeBytes.Observe(float64(len(memTx.tx)))
}

// RemoveTxByKey removes a transaction from the mempool by its TxKey index.
// Called from:
//   - Update (lock held) if tx was committed
//   - resCbRecheck (lock not held) if tx was invalidated
func (mem *CListMempool) RemoveTxByKey(txKey types.TxKey) error {
	if elem, ok := mem.getCElement(txKey); ok {
		mem.txs.Remove(elem)
		elem.DetachPrev()
		mem.txsMap.Delete(txKey)
		tx := elem.Value.(*mempoolTx).tx
		mem.txsBytes.Add(int64(-len(tx)))
		return nil
	}
	return ErrTxNotFound
}

func (mem *CListMempool) isFull(txSize int) error {
	memSize := mem.Size()
	txsBytes := mem.SizeBytes()
	if memSize >= mem.config.Size || int64(txSize)+txsBytes > mem.config.MaxTxsBytes {
		return ErrMempoolIsFull{
			NumTxs:      memSize,
			MaxTxs:      mem.config.Size,
			TxsBytes:    txsBytes,
			MaxTxsBytes: mem.config.MaxTxsBytes,
		}
	}

	return nil
}

// callback, which is called after the app checked the tx for the first time.
//
// The case where the app checks the tx for the second and subsequent times is
// handled by the resCbRecheck callback.
func (mem *CListMempool) resCbFirstTime(
	tx []byte,
	txInfo TxInfo,
	res *abci.Response,
) {
	switch r := res.Value.(type) {
	case *abci.Response_CheckTx:
		var postCheckErr error
		if mem.postCheck != nil {
			postCheckErr = mem.postCheck(tx, r.CheckTx)
		}
		if (r.CheckTx.Code == abci.CodeTypeOK) && postCheckErr == nil {
			// Check mempool isn't full again to reduce the chance of exceeding the
			// limits.
			if err := mem.isFull(len(tx)); err != nil {
				// remove from cache (mempool might have a space later)
				mem.cache.Remove(tx)
				mem.logger.Error(err.Error())
				return
			}

			// Check transaction not already in the mempool
			if e, ok := mem.txsMap.Load(types.Tx(tx).Key()); ok {
				memTx := e.(*clist.CElement).Value.(*mempoolTx)
				memTx.addSender(txInfo.SenderID)
				mem.logger.Debug(
					"transaction already there, not adding it again",
					"tx", types.Tx(tx).Hash(),
					"res", r,
					"height", mem.height.Load(),
					"total", mem.Size(),
				)
				return
			}

			memTx := &mempoolTx{
				height:    mem.height.Load(),
				gasWanted: r.CheckTx.GasWanted,
				tx:        tx,
			}
			memTx.addSender(txInfo.SenderID)
			mem.addTx(memTx)
			mem.logger.Debug(
				"added good transaction",
				"tx", types.Tx(tx).Hash(),
				"res", r,
				"height", mem.height.Load(),
				"total", mem.Size(),
			)
			mem.notifyTxsAvailable()
		} else {
			// ignore bad transaction
			mem.logger.Debug(
				"rejected bad transaction",
				"tx", types.Tx(tx).Hash(),
				"peerID", txInfo.SenderP2PID,
				"res", r,
				"err", postCheckErr,
			)
			mem.metrics.FailedTxs.Add(1)

			if !mem.config.KeepInvalidTxsInCache {
				// remove from cache (it might be good later)
				mem.cache.Remove(tx)
			}
		}

	default:
		// ignore other messages
	}
}

// callback, which is called after the app rechecked the tx.
//
// The case where the app checks the tx for the first time is handled by the
// resCbFirstTime callback.
func (mem *CListMempool) resCbRecheck(req *abci.Request, res *abci.Response) {
	switch r := res.Value.(type) {
	case *abci.Response_CheckTx:
		tx := req.GetCheckTx().Tx
		txKey := types.Tx(tx).Key()
		_, ok := mem.txsMap.Load(txKey)
		if !ok {
			mem.logger.Debug("re-CheckTx transaction does not exist", "expected", types.Tx(tx).Hash())
			return
		}

		var postCheckErr error
		if r.CheckTx.Code == abci.CodeTypeOK {
			if mem.postCheck == nil {
				return
			}
			postCheckErr = mem.postCheck(tx, r.CheckTx)
			if postCheckErr == nil {
				return
			}
		}
		// Tx became invalidated due to newly committed block.
		mem.logger.Debug("tx is no longer valid", "tx", types.Tx(tx).Hash(), "res", r, "err", postCheckErr)

		// We remove the invalid tx from the cache because it might be good later
		if err := mem.RemoveTxByKey(txKey); err != nil {
			mem.logger.Debug("Transaction could not be removed from mempool", "err", err)
		}
		if !mem.config.KeepInvalidTxsInCache {
			mem.cache.Remove(tx)
		}
	default:
		// ignore other messages
	}
}

// Safe for concurrent use by multiple goroutines.
func (mem *CListMempool) TxsAvailable() <-chan struct{} {
	return mem.txsAvailable
}

func (mem *CListMempool) notifyTxsAvailable() {
	if mem.Size() == 0 {
		panic("notified txs available but mempool is empty!")
	}
	if mem.txsAvailable != nil && mem.notifiedTxsAvailable.CompareAndSwap(false, true) {
		// channel cap is 1, so this will send once
		select {
		case mem.txsAvailable <- struct{}{}:
		default:
		}
	}
}

// Safe for concurrent use by multiple goroutines.
func (mem *CListMempool) ReapMaxBytesMaxGas(maxBytes, maxGas int64) types.Txs {
	mem.updateMtx.RLock()
	defer mem.updateMtx.RUnlock()

	var (
		totalGas    int64
		runningSize int64
	)

	// TODO: we will get a performance boost if we have a good estimate of avg
	// size per tx, and set the initial capacity based off of that.
	// txs := make([]types.Tx, 0, cmtmath.MinInt(mem.txs.Len(), max/mem.avgTxSize))
	txs := make([]types.Tx, 0, mem.txs.Len())
	for e := mem.txs.Front(); e != nil; e = e.Next() {
		memTx := e.Value.(*mempoolTx)

		txs = append(txs, memTx.tx)

		dataSize := types.ComputeProtoSizeForTxs([]types.Tx{memTx.tx})

		// Check total size requirement
		if maxBytes > -1 && runningSize+dataSize > maxBytes {
			return txs[:len(txs)-1]
		}

		runningSize += dataSize

		// Check total gas requirement.
		// If maxGas is negative, skip this check.
		// Since newTotalGas < masGas, which
		// must be non-negative, it follows that this won't overflow.
		newTotalGas := totalGas + memTx.gasWanted
		if maxGas > -1 && newTotalGas > maxGas {
			return txs[:len(txs)-1]
		}
		totalGas = newTotalGas
	}
	return txs
}

// Safe for concurrent use by multiple goroutines.
func (mem *CListMempool) ReapMaxTxs(max int) types.Txs {
	mem.updateMtx.RLock()
	defer mem.updateMtx.RUnlock()

	if max < 0 {
		max = mem.txs.Len()
	}

	txs := make([]types.Tx, 0, cmtmath.MinInt(mem.txs.Len(), max))
	for e := mem.txs.Front(); e != nil && len(txs) <= max; e = e.Next() {
		memTx := e.Value.(*mempoolTx)
		txs = append(txs, memTx.tx)
	}
	return txs
}

// Lock() must be help by the caller during execution.
func (mem *CListMempool) Update(
	block *types.Block,
	txResults []*abci.ExecTxResult,
	preCheck PreCheckFunc,
	postCheck PostCheckFunc,
) error {
	// Set height
	mem.height.Store(block.Height)
	mem.notifiedTxsAvailable.Store(false)

	if preCheck != nil {
		mem.preCheck = preCheck
	}
	if postCheck != nil {
		mem.postCheck = postCheck
	}

	for i, tx := range block.Txs {
		if txResults[i].Code == abci.CodeTypeOK {
			// Add valid committed tx to the cache (if missing).
			_ = mem.cache.Push(tx)
		} else if !mem.config.KeepInvalidTxsInCache {
			// Allow invalid transactions to be resubmitted.
			mem.cache.Remove(tx)
		}

		// Remove committed tx from the mempool.
		//
		// Note an evil proposer can drop valid txs!
		// Mempool before:
		//   100 -> 101 -> 102
		// Block, proposed by an evil proposer:
		//   101 -> 102
		// Mempool after:
		//   100
		// https://github.com/tendermint/tendermint/issues/3322.
		if err := mem.RemoveTxByKey(tx.Key()); err != nil {
			mem.logger.Debug("Committed transaction not in local mempool (not an error)",
				"key", tx.Key(),
				"error", err.Error())
		}
	}

	// Either recheck non-committed txs to see if they became invalid
	// or just notify there're some txs left.
	if mem.Size() > 0 {
		if mem.config.Recheck {
			fName, tFormat := "Mempool.Update", "15:04:05.000"
			mem.logger.Info(fmt.Sprintf("[%s]%s::start recheck txs", time.Now().Format(tFormat), fName), "numtxs", mem.Size(), "height", block.Height)
			mem.logger.Debug("recheck txs", "numtxs", mem.Size(), "height", block.Height)
			res, err := mem.proxyAppConn.BeginRecheckTxSync(context.TODO(), &abci.RequestBeginRecheckTx{
				Header: types.TM2PB.Header(&block.Header),
			})
			if res.Code == abci.CodeTypeOK && err == nil {
				mem.recheckTxs()
				res2, err2 := mem.proxyAppConn.EndRecheckTxSync(context.TODO(), &abci.RequestEndRecheckTx{Height: block.Height})
				if res2.Code != abci.CodeTypeOK {
					return errors.New("the function EndRecheckTxSync does not respond CodeTypeOK")
				}
				if err2 != nil {
					return errors.Join(err2, errors.New("the function EndRecheckTxSync returns an error"))
				}
			} else {
				if res.Code != abci.CodeTypeOK {
					return errors.New("the function BeginRecheckTxSync does not respond CodeTypeOK")
				}
				if err != nil {
					return errors.Join(err, errors.New("the function BeginRecheckTxSync returns an error"))
				}
			}
			mem.logger.Info(fmt.Sprintf("[%s]%s::done recheck txs", time.Now().Format(tFormat), fName))
		}
	}

	if mem.Size() > 0 {
		mem.notifyTxsAvailable()
	}

	// Update metrics
	mem.metrics.Size.Set(float64(mem.Size()))
	mem.metrics.SizeBytes.Set(float64(mem.SizeBytes()))

	return nil
}

func (mem *CListMempool) recheckTxs() {
	if mem.Size() == 0 {
		panic("recheckTxs is called, but the mempool is empty")
	}

	wg := sync.WaitGroup{}

	// Push txs to proxyAppConn
	// NOTE: globalCb may be called concurrently.
	for e := mem.txs.Front(); e != nil; e = e.Next() {
		wg.Add(1)

		memTx := e.Value.(*mempoolTx)
		reqRes, _ := mem.proxyAppConn.CheckTxAsync(context.TODO(), &abci.RequestCheckTx{
			Tx:   memTx.tx,
			Type: abci.CheckTxType_Recheck,
		})
		reqRes.SetCallback(func(res *abci.Response) {
			wg.Done()
		})
	}

	wg.Wait()

	// In <v0.37 we would call FlushAsync at the end of recheckTx forcing the buffer to flush
	// all pending messages to the app. There doesn't seem to be any need here as the buffer
	// will get flushed regularly or when filled.
}
