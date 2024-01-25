package icon

import (
	"context"
	"sort"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icon-project/centralized-relay/relayer/chains/icon/types"
	providerTypes "github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/codec"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	maxRetires = 5
)

type btpBlockResponse struct {
	Height    int64
	Hash      common.HexHash
	Header    *types.BlockHeader
	EventLogs []*types.EventLog
}

type btpBlockRequest struct {
	height   int64
	hash     types.HexBytes
	indexes  [][]types.HexInt
	events   [][][]types.HexInt
	err      error
	retry    uint8
	response *btpBlockResponse
}

// TODO: check for balance and if the balance is low show info balance is low
// starting listener
func (icp *IconProvider) Listener(ctx context.Context, lastSavedHeight uint64, incoming chan providerTypes.BlockInfo) error {
	errCh := make(chan error)                                            // error channel
	reconnectCh := make(chan struct{}, 1)                                // reconnect channel
	btpBlockNotifCh := make(chan *types.BlockNotification, 100)          // block notification channel
	btpBlockRespCh := make(chan *btpBlockResponse, cap(btpBlockNotifCh)) // block result channel

	reconnect := func() {
		select {
		case reconnectCh <- struct{}{}:
		default:
		}
		for len(btpBlockRespCh) > 0 || len(btpBlockNotifCh) > 0 {
			select {
			case <-btpBlockRespCh: // clear block result channel
			case <-btpBlockNotifCh: // clear block notification channel
			}
		}
	}

	processedheight, err := icp.StartFromHeight(ctx, lastSavedHeight)
	if err != nil {
		return errors.Wrapf(err, "failed to calculate start height")
	}

	icp.log.Info("Start querying from height", zap.Int64("height", processedheight))
	// subscribe to monitor block
	ctxMonitorBlock, cancelMonitorBlock := context.WithCancel(ctx)
	reconnect()

	blockReq := &types.BlockRequest{
		Height:       types.NewHexInt(int64(processedheight)),
		EventFilters: GetMonitorEventFilters(icp.PCfg.Contracts[providerTypes.ConnectionContract], MonitorEventsList),
	}

loop:
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			return err

		case <-reconnectCh:
			cancelMonitorBlock()
			ctxMonitorBlock, cancelMonitorBlock = context.WithCancel(ctx)

			go func(ctx context.Context, cancel context.CancelFunc) {
				blockReq.Height = types.NewHexInt(int64(processedheight))
				icp.log.Debug("try to reconnect from", zap.Int64("height", processedheight))
				err := icp.client.MonitorBlock(ctx, blockReq, func(conn *websocket.Conn, v *types.BlockNotification) error {
					if !errors.Is(ctx.Err(), context.Canceled) {
						btpBlockNotifCh <- v
					}
					return nil
				}, func(conn *websocket.Conn) {
				}, func(conn *websocket.Conn, err error) {})
				if err != nil {
					if errors.Is(err, context.Canceled) {
						return
					}
					time.Sleep(time.Second * 5)
					reconnect()
					icp.log.Warn("error occured during monitor block", zap.Error(err))
				}
			}(ctxMonitorBlock, cancelMonitorBlock)
		case br := <-btpBlockRespCh:
			for ; br != nil; processedheight++ {
				icp.log.Debug("block notification received", zap.Int64("height", int64(processedheight)))

				// note: because of monitorLoop height should be subtract by 1
				height := br.Height - 1

				messages := icp.parseMessagesFromEventlogs(icp.log, br.EventLogs, uint64(height))

				// TODO: check for the concurrency
				incoming <- providerTypes.BlockInfo{
					Messages: messages,
					Height:   uint64(height),
				}

				if br = nil; len(btpBlockRespCh) > 0 {
					br = <-btpBlockRespCh
				}
			}
			// remove unprocessed blockResponses
			for len(btpBlockRespCh) > 0 {
				<-btpBlockRespCh
			}

		default:
			select {
			default:
			case bn := <-btpBlockNotifCh:
				requestCh := make(chan *btpBlockRequest, cap(btpBlockNotifCh))
				for i := int64(0); bn != nil; i++ {
					height, err := bn.Height.Value()

					if err != nil {
						return err
					} else if height != processedheight+i {
						icp.log.Warn("Reconnect: missing block notification",
							zap.Int64("got", height),
							zap.Int64("expected", processedheight+i),
						)
						reconnect()
						continue loop
					}

					requestCh <- &btpBlockRequest{
						height:  height,
						hash:    bn.Hash,
						indexes: bn.Indexes,
						events:  bn.Events,
						retry:   providerTypes.MaxTxRetry,
					}
					if bn = nil; len(btpBlockNotifCh) > 0 && len(requestCh) < cap(requestCh) {
						bn = <-btpBlockNotifCh
					}
				}

				brs := make([]*btpBlockResponse, 0, len(requestCh))
				for request := range requestCh {
					switch {
					case request.err != nil:
						if request.retry > 0 {
							request.retry--
							request.response, request.err = nil, nil
							requestCh <- request
							continue
						}
						icp.log.Info("Request error ",
							zap.Any("height", request.height),
							zap.Error(request.err))
						brs = append(brs, nil)
						if len(brs) == cap(brs) {
							close(requestCh)
						}
					case request.response != nil:
						brs = append(brs, request.response)
						if len(brs) == cap(brs) {
							close(requestCh)
						}
					default:
						go icp.handleBTPBlockRequest(request, requestCh)

					}
				}
				// filter nil
				_brs, brs := brs, brs[:0]
				for _, v := range _brs {
					if v != nil {
						brs = append(brs, v)
					}
				}

				// sort and forward notifications
				if len(brs) > 0 {
					sort.SliceStable(brs, func(i, j int) bool {
						return brs[i].Height < brs[j].Height
					})
					for i, d := range brs {
						if d.Height == processedheight+int64(i) {
							btpBlockRespCh <- d
						}
					}
				}

			}
		}
	}
}

func (icp *IconProvider) handleBTPBlockRequest(
	request *btpBlockRequest, requestCh chan *btpBlockRequest,
) {
	defer func() {
		time.Sleep(500 * time.Millisecond)
		requestCh <- request
	}()

	if request.response == nil {
		request.response = &btpBlockResponse{}
	}
	request.response.Height = request.height
	request.response.Hash, request.err = request.hash.Value()
	if request.err != nil {
		request.err = errors.Wrapf(request.err,
			"invalid hash: height=%v, hash=%v, %v", request.height, request.hash, request.err)
		return
	}

	containsEventlogs := len(request.indexes) > 0 && len(request.events) > 0
	if containsEventlogs {
		blockHeader, err := icp.client.GetBlockHeaderByHeight(request.height)
		if err != nil {
			request.err = errors.Wrapf(request.err, "getBlockHeader: %v", err)
			return
		}

		var receiptHash types.BlockHeaderResult
		_, err = codec.RLP.UnmarshalFromBytes(blockHeader.Result, &receiptHash)
		if err != nil {
			request.err = errors.Wrapf(err, "BlockHeaderResult.UnmarshalFromBytes: %v", err)
			return

		}

		var eventlogs []*types.EventLog
		for id := 0; id < len(request.indexes); id++ {
			for i, index := range request.indexes[id] {
				p := &types.ProofEventsParam{
					Index:     index,
					BlockHash: request.hash,
					Events:    request.events[id][i],
				}

				proofs, err := icp.client.GetProofForEvents(p)
				if err != nil {
					request.err = errors.Wrapf(err, "GetProofForEvents: %v", err)
					return

				}

				// Processing receipt index
				serializedReceipt, err := MptProve(index, proofs[0], receiptHash.ReceiptHash)
				if err != nil {
					request.err = errors.Wrapf(err, "MPTProve Receipt: %v", err)
					return

				}
				var result types.TxResult
				_, err = codec.RLP.UnmarshalFromBytes(serializedReceipt, &result)
				if err != nil {
					request.err = errors.Wrapf(err, "Unmarshal Receipt: %v", err)
					return
				}

				for j := 0; j < len(p.Events); j++ {
					serializedEventLog, err := MptProve(p.Events[j], proofs[j+1], common.HexBytes(result.EventLogsHash))
					if err != nil {
						request.err = errors.Wrapf(err, "event.MPTProve: %v", err)
						return
					}
					el := new(types.EventLog)
					_, err = codec.RLP.UnmarshalFromBytes(serializedEventLog, el)
					if err != nil {
						request.err = errors.Wrapf(err, "event.UnmarshalFromBytes: %v", err)
						return
					}
					icp.log.Info("Detected eventlog ", zap.Int64("height", request.height),
						zap.String("eventlog", EventNameToType[string(el.Indexed[0])]))
					eventlogs = append(eventlogs, el)
				}

			}
		}
		request.response.EventLogs = eventlogs
	}
}

func (icp *IconProvider) StartFromHeight(ctx context.Context, lastSavedHeight uint64) (int64, error) {
	latestHeight, err := icp.QueryLatestHeight(ctx)
	if err != nil {
		return 0, err
	}

	if icp.PCfg.StartHeight > latestHeight {
		icp.log.Error("start height provided on config cannot be greater than latest height",
			zap.Uint64("start-height", icp.PCfg.StartHeight),
			zap.Int64("latest-height", int64(latestHeight)),
		)
	}

	// priority2: lastsaveheight from db
	if lastSavedHeight != 0 && lastSavedHeight < latestHeight {
		return int64(lastSavedHeight), nil
	}

	// priority1: startHeight from config
	if icp.PCfg.StartHeight != 0 && icp.PCfg.StartHeight < latestHeight {
		return int64(icp.PCfg.StartHeight), nil
	}

	// priority3: latest height
	return int64(latestHeight), nil
}
