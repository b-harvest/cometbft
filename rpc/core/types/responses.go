package coretypes

import (
	"encoding/json"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/libs/bytes"
	"github.com/cometbft/cometbft/p2p"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/types"
)

// List of blocks
type ResultBlockchainInfo struct {
	LastHeight int64              `json:"last_height"`
	BlockMetas []*types.BlockMeta `json:"block_metas"`
}

// Genesis file
type ResultGenesis struct {
	Genesis *types.GenesisDoc `json:"genesis"`
}

// ResultGenesisChunk is the output format for the chunked/paginated
// interface. These chunks are produced by converting the genesis
// document to JSON and then splitting the resulting payload into
// 16 megabyte blocks and then base64 encoding each block.
type ResultGenesisChunk struct {
	ChunkNumber int    `json:"chunk"`
	TotalChunks int    `json:"total"`
	Data        string `json:"data"`
}

// Single block (with meta)
type ResultBlock struct {
	BlockID types.BlockID `json:"block_id"`
	Block   *types.Block  `json:"block"`
}

// ResultHeader represents the response for a Header RPC Client query
type ResultHeader struct {
	Header *types.Header `json:"header"`
}

// Commit and Header
type ResultCommit struct {
	types.SignedHeader `json:"signed_header"`
	CanonicalCommit    bool `json:"canonical"`
}

// ABCI results from a block
type ResultBlockResults struct {
	Height                int64                     `json:"height"`
	TxsResults            []*abci.ResponseDeliverTx `json:"txs_results"`
	BeginBlockEvents      []abci.Event              `json:"begin_block_events"`
	EndBlockEvents        []abci.Event              `json:"end_block_events"`
	ValidatorUpdates      []abci.ValidatorUpdate    `json:"validator_updates"`
	ConsensusParamUpdates *cmtproto.ConsensusParams `json:"consensus_param_updates"`
}

// NewResultCommit is a helper to initialize the ResultCommit with
// the embedded struct
func NewResultCommit(header *types.Header, commit *types.Commit,
	canonical bool) *ResultCommit {

	return &ResultCommit{
		SignedHeader: types.SignedHeader{
			Header: header,
			Commit: commit,
		},
		CanonicalCommit: canonical,
	}
}

// Info about the node's syncing state
type SyncInfo struct {
	LatestBlockHash   bytes.HexBytes `json:"latest_block_hash"`
	LatestAppHash     bytes.HexBytes `json:"latest_app_hash"`
	LatestBlockHeight int64          `json:"latest_block_height"`
	LatestBlockTime   time.Time      `json:"latest_block_time"`

	EarliestBlockHash   bytes.HexBytes `json:"earliest_block_hash"`
	EarliestAppHash     bytes.HexBytes `json:"earliest_app_hash"`
	EarliestBlockHeight int64          `json:"earliest_block_height"`
	EarliestBlockTime   time.Time      `json:"earliest_block_time"`

	CatchingUp bool `json:"catching_up"`
}

// Info about the node's validator
type ValidatorInfo struct {
	Address     bytes.HexBytes `json:"address"`
	PubKey      crypto.PubKey  `json:"pub_key"`
	VotingPower int64          `json:"voting_power"`
}

// Node Status
type ResultStatus struct {
	NodeInfo      p2p.DefaultNodeInfo `json:"node_info"`
	SyncInfo      SyncInfo            `json:"sync_info"`
	ValidatorInfo ValidatorInfo       `json:"validator_info"`
}

// Is TxIndexing enabled
func (s *ResultStatus) TxIndexEnabled() bool {
	if s == nil {
		return false
	}
	return s.NodeInfo.Other.TxIndex == "on"
}

// Info about peer connections
type ResultNetInfo struct {
	Listening bool     `json:"listening"`
	Listeners []string `json:"listeners"`
	NPeers    int      `json:"n_peers"`
	Peers     []Peer   `json:"peers"`
}

// Log from dialing seeds
type ResultDialSeeds struct {
	Log string `json:"log"`
}

// Log from dialing peers
type ResultDialPeers struct {
	Log string `json:"log"`
}

// A peer
type Peer struct {
	NodeInfo         p2p.DefaultNodeInfo  `json:"node_info"`
	IsOutbound       bool                 `json:"is_outbound"`
	ConnectionStatus p2p.ConnectionStatus `json:"connection_status"`
	RemoteIP         string               `json:"remote_ip"`
}

// Validators for a height.
type ResultValidators struct {
	BlockHeight int64              `json:"block_height"`
	Validators  []*types.Validator `json:"validators"`
	// Count of actual validators in this result
	Count int `json:"count"`
	// Total number of validators
	Total int `json:"total"`
}

// ConsensusParams for given height
type ResultConsensusParams struct {
	BlockHeight     int64                 `json:"block_height"`
	ConsensusParams types.ConsensusParams `json:"consensus_params"`
}

// Info about the consensus state.
// UNSTABLE
type ResultDumpConsensusState struct {
	RoundState json.RawMessage `json:"round_state"`
	Peers      []PeerStateInfo `json:"peers"`
}

// UNSTABLE
type PeerStateInfo struct {
	NodeAddress string          `json:"node_address"`
	PeerState   json.RawMessage `json:"peer_state"`
}

// UNSTABLE
type ResultConsensusState struct {
	RoundState json.RawMessage `json:"round_state"`
}

// CheckTx result
type ResultBroadcastTx struct {
	Code      uint32         `json:"code"`
	Data      bytes.HexBytes `json:"data"`
	Log       string         `json:"log"`
	Codespace string         `json:"codespace"`

	Hash bytes.HexBytes `json:"hash"`
}

// CheckTx and DeliverTx results
type ResultBroadcastTxCommit struct {
	CheckTx   abci.ResponseCheckTx   `json:"check_tx"`
	DeliverTx abci.ResponseDeliverTx `json:"deliver_tx"`
	Hash      bytes.HexBytes         `json:"hash"`
	Height    int64                  `json:"height"`
}

// ResultCheckTx wraps abci.ResponseCheckTx.
type ResultCheckTx struct {
	abci.ResponseCheckTx
}

// Result of querying for a tx
type ResultTx struct {
	Hash     bytes.HexBytes         `json:"hash"`
	Height   int64                  `json:"height"`
	Index    uint32                 `json:"index"`
	TxResult abci.ResponseDeliverTx `json:"tx_result"`
	Tx       types.Tx               `json:"tx"`
	Proof    types.TxProof          `json:"proof,omitempty"`
}

// Result of searching for txs
type ResultTxSearch struct {
	Txs        []*ResultTx `json:"txs"`
	TotalCount int         `json:"total_count"`
}

// ResultBlockSearch defines the RPC response type for a block search by events.
type ResultBlockSearch struct {
	Blocks     []*ResultBlock `json:"blocks"`
	TotalCount int            `json:"total_count"`
}

// List of mempool txs
type ResultUnconfirmedTxs struct {
	Count      int        `json:"n_txs"`
	Total      int        `json:"total"`
	TotalBytes int64      `json:"total_bytes"`
	Txs        []types.Tx `json:"txs"`
}

// Info abci msg
type ResultABCIInfo struct {
	Response abci.ResponseInfo `json:"response"`
}

// Query abci msg
type ResultABCIQuery struct {
	Response abci.ResponseQuery `json:"response"`
}

// Result of broadcasting evidence
type ResultBroadcastEvidence struct {
	Hash []byte `json:"hash"`
}

// empty results
type (
	ResultUnsafeFlushMempool struct{}
	ResultUnsafeProfile      struct{}
	ResultSubscribe          struct{}
	ResultUnsubscribe        struct{}
	ResultHealth             struct{}
)

// Event data from a subscription
type ResultEvent struct {
	Query  string              `json:"query"`
	Data   types.TMEventData   `json:"data"`
	Events map[string][]string `json:"events"`
}

// RequestEvents is the argument for the "/events" RPC endpoint.
type RequestEvents struct {
	// Optional filter spec. If nil or empty, all items are eligible.
	Filter *EventFilter `json:"filter"`

	// The maximum number of eligible items to return.
	// If zero or negative, the server will report a default number.
	MaxItems int `json:"maxItems"`

	// Return only items after this cursor. If empty, the limit is just
	// before the the beginning of the event log.
	After string `json:"after"`

	// Return only items before this cursor.  If empty, the limit is just
	// after the head of the event log.
	Before string `json:"before"`

	// Wait for up to this long for events to be available.
	WaitTime time.Duration `json:"waitTime"`
}

// An EventFilter specifies which events are selected by an /events request.
type EventFilter struct {
	Query string `json:"query"`
}

// ResultEvents is the response from the "/events" RPC endpoint.
type ResultEvents struct {
	// The items matching the request parameters, from newest
	// to oldest, if any were available within the timeout.
	Items []*EventItem `json:"items"`

	// This is true if there is at least one older matching item
	// available in the log that was not returned.
	More bool `json:"more"`

	// The cursor of the oldest item in the log at the time of this reply,
	// or "" if the log is empty.
	Oldest string `json:"oldest"`

	// The cursor of the newest item in the log at the time of this reply,
	// or "" if the log is empty.
	Newest string `json:"newest"`
}

type EventItem struct {
	// The cursor of this item.
	Cursor string `json:"cursor"`

	// The event label of this item (for example, "Vote").
	Event string `json:"event,omitempty"`

	// The encoded event data for this item. The content is a JSON object with
	// the following structure:
	//
	//   <json-encoded-value>
	//
	// The known type tags are defined by the tendermint/types package.
	Data json.RawMessage `json:"data"`
}
