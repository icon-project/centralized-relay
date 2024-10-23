package socket

import (
	"math/big"
	"net"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/icon-project/centralized-relay/relayer/types"
)

type Event string

type Request struct {
	ID    string `json:"id"`
	Event Event  `json:"event"`
	Data  any    `json:"data"`
}

type Response struct {
	ID      string `json:"id"`
	Event   Event  `json:"event"`
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func (r *Response) SetError(err error) *Response {
	r.Message = err.Error()
	return r
}

func (r *Response) SetData(data any) *Response {
	r.Success = true
	r.Data = data
	return r
}

type Server struct {
	listener  net.Listener
	startedAt int64
	rly       *relayer.Relayer
}

type Pagination struct {
	Page  uint
	Limit uint
}

type ReqMessageList struct {
	Chain string `json:"chain"`
	Limit uint   `json:"pagination"`
}

type ReqGetBlock struct {
	Chain string `json:"chain"`
}

type ReqRelayMessage struct {
	Chain  string `json:"chain"`
	Height uint64 `json:"height"`
	TxHash string `json:"txHash"`
}

type ReqMessageRemove struct {
	Chain string   `json:"chain"`
	Sn    *big.Int `json:"sn"`
}

type ResMessageRemove struct {
	Sn     *big.Int `json:"sn"`
	Chain  string   `json:"chain"`
	Dst    string   `json:"dst"`
	Height uint64   `json:"height"`
	Event  string   `json:"event"`
}

type ResMessageList struct {
	Message []*types.RouteMessage `json:"message"`
	Total   int                   `json:"total"`
}

type ResGetBlock struct {
	Chain            string `json:"chain"`
	CheckPointHeight uint64 `json:"checkPointHeight"`
	LatestHeight     uint64 `json:"latestHeight"`
}

type ReqPruneDB struct {
	ID    string `json:"id"`
	Chain string `json:"chain"`
}

type ResPruneDB struct {
	Status string `json:"status"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

type ReqRevertMessage struct {
	Chain string `json:"chain"`
	Sn    uint64 `json:"sn"`
}

type ResRevertMessage struct {
	Sn uint64 `json:"sn"`
}

type ReqGetFee struct {
	Chain    string `json:"chain"`
	Network  string `json:"network"`
	Response bool   `json:"response"`
}

type ResGetFee struct {
	Chain    string `json:"chain"`
	Fee      uint64 `json:"fee"`
	Response bool   `json:"response"`
}

// ReqSetFee sends SetFee event to socket
type ReqSetFee struct {
	Chain   string   `json:"chain"`
	Network string   `json:"network"`
	MsgFee  *big.Int `json:"msg_fee"`
	ResFee  *big.Int `json:"res_fee"`
}

// ResSetFee sends SetFee event to socket
type ResSetFee struct {
	Status string `json:"status"`
}

// ReqClaimFee sends ClaimFee event to socket
type ReqClaimFee struct {
	Chain string `json:"chain"`
}

// ResClaimFee sends ClaimFee event to socket
type ResClaimFee struct {
	Status string `json:"status"`
}

type ReqChainHeight struct {
	Chain string `json:"chain"`
}

type ResChainHeight struct {
	Chain  string `json:"chain"`
	Height uint64 `json:"height"`
}

type ReqProcessedBlock struct {
	Chain string `json:"chain"`
}

type ReqRangeBlockQuery struct {
	Chain      string `json:"chain"`
	FromHeight uint64 `json:"from_height"`
	ToHeight   uint64 `json:"to_height"`
}

type ResRangeBlockQuery struct {
	Chain string           `json:"chain"`
	Msgs  []*types.Message `json:"messages"`
}

type ReqListChain struct {
	Chains []string `json:"chains,omitempty"`
}

type ResChainInfo struct {
	Name           string            `json:"name"`
	NID            string            `json:"nid"`
	Address        string            `json:"address"`
	Type           string            `json:"type"`
	Contracts      map[string]string `json:"contracts"`
	LatestHeight   uint64            `json:"latestHeight"`
	LastCheckPoint uint64            `json:"lastCheckPoint"`
}

type ReqGetBalance struct {
	Chain   string `json:"chain"`
	Address string `json:"address"`
}

type ResGetBalance struct {
	Chain   string      `json:"chain"`
	Address string      `json:"address"`
	Balance *types.Coin `json:"balance"`
	Value   string      `json:"value"`
}

type ReqRelayInfo struct{}

type ResRelayInfo struct {
	Version string `json:"version"`
	Uptime  int64  `json:"uptime"`
}

type ReqMessageReceived struct {
	Chain string `json:"chain"`
	Sn    uint64 `json:"sn"`
}

type ResMessageReceived struct {
	Chain    string `json:"chain"`
	Sn       uint64 `json:"sn"`
	Received bool   `json:"received"`
}

type ReqGetBlockEvents struct {
	Height uint64 `json:"height,omitempty"`
	TxHash string `json:"txHash,omitempty"`
}

type ResGetBlockEvents struct {
	Event     string `json:"event"`
	Height    uint64 `json:"height"`
	Executed  bool   `json:"executed"`
	TxHash    string `json:"txHash"`
	ChainInfo struct {
		NID       string                  `json:"nid"`
		Name      string                  `json:"name"`
		Type      string                  `json:"type"`
		Contracts types.ContractConfigMap `json:"contracts"`
	} `json:"chainInfo"`
}

type ChainProviderError struct {
	Message string
}
