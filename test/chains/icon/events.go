package icon

import (
	"context"
	"fmt"
	"sync"

	"github.com/icon-project/centralized-relay/test/chains"
	icontypes "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon/types"
)

type WebSocketListener struct {
	height       uint64
	contract     string
	chain        *IconRemotenet
	Shutdown     chan struct{}
	Events       []map[string][]string
	eventMapLock sync.RWMutex
	wg           sync.WaitGroup
}

func NewIconEventListener(c *IconRemotenet, contract string) *WebSocketListener {
	height, _ := c.Height(context.Background())
	return &WebSocketListener{
		chain:    c,
		contract: c.GetContractAddress(contract),
		height:   height,
		Shutdown: make(chan struct{}),
		Events:   []map[string][]string{},
	}
}

func (w *WebSocketListener) Start() {

}

func (w *WebSocketListener) Stop() {
	fmt.Println("Stopped.")
}

func (w *WebSocketListener) FindEvent(filters chains.Filter) (chains.Event, error) {
	signature := filters["signature"].(string)
	index := filters["index"].([]*string)
	ctx := context.Background()
	event, err := w.chain.FindEvent(ctx, w.height, w.contract, signature, index)
	if err != nil {
		return nil, err
	}
	intHeight, _ := event.Height.Int()
	block, _ := w.chain.IconClient.GetBlockByHeight(&icontypes.BlockHeightParam{Height: icontypes.NewHexInt(int64(intHeight - 1))})
	i, _ := event.Index.Int()
	tx := block.NormalTransactions[i]
	trResult, _ := w.chain.TransactionResult(ctx, string(tx.TxHash))
	eventIndex, _ := event.Events[0].Int()

	var result = make(chains.Event)
	result["data"] = trResult.EventLogs[eventIndex].Data
	result["indexed"] = trResult.EventLogs[eventIndex].Indexed
	return result, nil
}
