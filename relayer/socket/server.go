package socket

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"os"
	"path"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/icon-project/centralized-relay/relayer"
	"github.com/icon-project/centralized-relay/relayer/store"
	"github.com/icon-project/centralized-relay/relayer/types"
)

var (
	SocketPath = getEnvOrFallback("SOCKET_PATH", path.Join(os.TempDir(), "relayer.sock"))
	network    = "unix"
)

func NewSocket(rly *relayer.Relayer) (*Server, error) {
	l, err := net.Listen(network, SocketPath)
	if err != nil {
		return nil, err
	}
	return &Server{listener: l, startedAt: time.Now().Unix(), rly: rly}, nil
}

// Listen to socket
func (s *Server) Listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.server(conn)
	}
}

// Send sends message to socket
func (s *Server) server(c net.Conn) {
	for {
		buf := make([]byte, 1024*100)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}
		message, err := s.parse(buf[:nr])
		if err != nil {
			message = makeError(err)
		}
		if err := s.send(c, message); err != nil {
			return
		}
	}
}

// Parse message from socket
func (s *Server) parse(data []byte) ([]byte, error) {
	msg := new(Request)
	if err := jsoniter.Unmarshal(data, msg); err != nil {
		return makeError(err), nil
	}
	payload := s.parseEvent(msg)
	return jsoniter.Marshal(payload)
}

// makeError for the client to write to socket
func makeError(err error) []byte {
	message := &Response{Event: EventError, Message: err.Error()}
	data, err := jsoniter.Marshal(message)
	if err != nil {
		return []byte(fmt.Sprintf(`{"event":"%s","message":"%s","success":false}`, EventError, err.Error()))
	}
	return data
}

// Send message to socket
func (s *Server) send(conn net.Conn, data []byte) error {
	_, err := conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// parseEvent for the client to write to socket
func (s *Server) parseEvent(msg *Request) *Response {
	data, err := jsoniter.Marshal(msg.Data)
	if err != nil {
		return &Response{ID: msg.ID, Event: EventError, Message: err.Error()}
	}
	response := &Response{ID: msg.ID, Event: msg.Event}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	switch msg.Event {
	case EventGetBlock:
		req := new(ReqGetBlock)
		if err := jsoniter.Unmarshal(data, req); err != nil {
			return response.SetError(err)
		}
		store := s.rly.GetBlockStore()
		var blocks []*ResGetBlock
		if req.Chain == "" {
			for _, chain := range s.rly.GetAllChainsRuntime() {
				latestHeight, err := chain.Provider.QueryLatestHeight(ctx)
				if err != nil {
					return response.SetError(err)
				}
				blocks = append(blocks, &ResGetBlock{
					Chain:            chain.Provider.NID(),
					CheckPointHeight: chain.LastSavedHeight,
					LatestHeight:     latestHeight,
				})
			}
			return response.SetData(blocks)
		}
		checkPointHeight, err := store.GetLastStoredBlock(req.Chain)
		if err != nil {
			return response.SetError(err)
		}
		chain, err := s.rly.FindChainRuntime(req.Chain)
		if err != nil {
			return response.SetError(err)
		}
		latestHeight, err := chain.Provider.QueryLatestHeight(ctx)
		if err != nil {
			return response.SetError(err)
		}
		blocks = append(blocks, &ResGetBlock{Chain: req.Chain, CheckPointHeight: checkPointHeight, LatestHeight: latestHeight})
		return response.SetData(blocks)
	case EventGetMessageList:
		req := new(ReqMessageList)
		if err := jsoniter.Unmarshal(data, req); err != nil {
			return response.SetError(err)
		}
		msgStore := s.rly.GetMessageStore()
		messages, err := msgStore.GetMessages(req.Chain, &store.Pagination{Limit: req.Limit})
		if err != nil {
			return response.SetError(err)
		}
		total, err := msgStore.TotalCountByChain(req.Chain)
		if err != nil {
			return response.SetError(err)
		}
		return response.SetData(&ResMessageList{messages, int(total)})
	case EventMessageRemove:
		req := new(ReqMessageRemove)
		if err := jsoniter.Unmarshal(data, req); err != nil {
			return response.SetError(err)
		}
		store := s.rly.GetMessageStore()
		key := &types.MessageKey{Src: req.Chain, Sn: req.Sn}
		message, err := store.GetMessage(key)
		if err != nil {
			return response.SetError(err)
		}
		if err := store.DeleteMessage(key); err != nil {
			return response.SetError(err)
		}
		return response.SetData(&ResMessageRemove{req.Sn, req.Chain, message.Dst, message.MessageHeight, message.EventType})
	case EventRelayMessage:
		req := new(ReqRelayMessage)
		if err := jsoniter.Unmarshal(data, req); err != nil {
			return response.SetError(err)
		}
		src, err := s.rly.FindChainRuntime(req.Chain)
		if err != nil {
			return response.SetError(err)
		}

		messages := []*types.Message{}
		if req.TxHash != "" {
			msgs, err := src.Provider.FetchTxMessages(ctx, req.TxHash)
			if err != nil {
				return response.SetError(err)
			}
			messages = append(messages, msgs...)
		} else if req.Height != 0 {
			msgs, err := src.Provider.GenerateMessages(ctx, req.Height, req.Height)
			if err != nil {
				return response.SetError(err)
			}
			messages = append(messages, msgs...)
		}
		for _, msg := range messages {
			src.MessageCache.Add(types.NewRouteMessage(msg))
		}
		return response.SetData(messages)
	case EventPruneDB:
		req := new(ReqPruneDB)
		if err := jsoniter.Unmarshal(data, req); err != nil {
			return response.SetError(err)
		}
		if err := s.rly.PruneDB(); err != nil {
			return response.SetError(err)
		}
		return response.SetData(&ResPruneDB{"Success"})
	case EventRevertMessage:
		req := new(ReqRevertMessage)
		if err := jsoniter.Unmarshal(data, req); err != nil {
			return response.SetError(err)
		}
		chain, err := s.rly.FindChainRuntime(req.Chain)
		if err != nil {
			return response.SetError(err)
		}
		if err := chain.Provider.RevertMessage(ctx, new(big.Int).SetUint64(req.Sn)); err != nil {
			return response.SetError(err)
		}
		return response.SetData(&ResRevertMessage{req.Sn})
	case EventGetFee:
		req := new(ReqGetFee)
		if err := jsoniter.Unmarshal(data, req); err != nil {
			return response.SetError(err)
		}
		chain, err := s.rly.FindChainRuntime(req.Chain)
		if err != nil {
			return response.SetError(err)
		}
		fee, err := chain.Provider.GetFee(ctx, req.Network, req.Response)
		if err != nil {
			return response.SetError(err)
		}
		return response.SetData(&ResGetFee{Chain: req.Chain, Fee: fee})
	case EventSetFee:
		req := new(ReqSetFee)
		if err := jsoniter.Unmarshal(data, req); err != nil {
			return response.SetError(err)
		}
		chain, err := s.rly.FindChainRuntime(req.Chain)
		if err != nil {
			return response.SetError(err)
		}
		if err := chain.Provider.SetFee(ctx, req.Network, req.MsgFee, req.ResFee); err != nil {
			return response.SetError(err)
		}
		return response.SetData(&ResSetFee{"Success"})
	case EventClaimFee:
		req := new(ReqClaimFee)
		if err := jsoniter.Unmarshal(data, req); err != nil {
			return response.SetError(err)
		}
		chain, err := s.rly.FindChainRuntime(req.Chain)
		if err != nil {
			return response.SetError(err)
		}
		if err := chain.Provider.ClaimFee(ctx); err != nil {
			return response.SetError(err)
		}
		return response.SetData(&ResClaimFee{"Success"})
	case EventGetConfig:
		req := new(ReqChainHeight)
		if err := jsoniter.Unmarshal(data, req); err != nil {
			return response.SetError(err)
		}
		chain, err := s.rly.FindChainRuntime(req.Chain)
		if err != nil {
			return response.SetError(err)
		}
		return response.SetData(chain.Provider.Config())
	case EventListChainInfo:
		req := new(ReqListChain)
		if err := jsoniter.Unmarshal(data, req); err != nil {
			return response.SetError(err)
		}
		var (
			chainNames []*ResChainInfo
			chains     []*relayer.ChainRuntime
		)
		if len(req.Chains) > 0 {
			for _, chainName := range req.Chains {
				chain, err := s.rly.FindChainRuntime(chainName)
				if err != nil {
					return response.SetError(err)
				}
				chains = append(chains, chain)
			}
		} else {
			chains = s.rly.GetAllChainsRuntime()
		}
		for _, chain := range chains {
			latestHeight, _ := chain.Provider.QueryLatestHeight(ctx)
			chainNames = append(chainNames, &ResChainInfo{
				Name:           chain.Provider.Name(),
				NID:            chain.Provider.NID(),
				Address:        chain.Provider.Config().GetWallet(),
				Type:           chain.Provider.Type(),
				LatestHeight:   latestHeight,
				LastCheckPoint: chain.LastSavedHeight,
				Contracts:      chain.Provider.Config().ContractsAddress(),
			})
		}
		return response.SetData(chainNames)
	case EventGetBalance:
		var reqs []ReqGetBalance
		if err := jsoniter.Unmarshal(data, &reqs); err != nil {
			return response.SetError(err)
		}
		res := make([]*ResGetBalance, 0, len(reqs))
		for _, req := range reqs {
			chain, err := s.rly.FindChainRuntime(req.Chain)
			if err != nil {
				return &Response{ID: msg.ID, Event: EventGetBalance, Data: res, Message: err.Error()}
			}
			balance, err := chain.Provider.QueryBalance(ctx, req.Address)
			if err != nil {
				balance = types.NewCoin("N/A", 0, 0)
			}
			res = append(res, &ResGetBalance{Chain: req.Chain, Address: req.Address, Balance: balance, Value: balance.Calculate()})
		}
		return response.SetData(res)
	case EventMessageReceived:
		req := new(ReqMessageReceived)
		if err := jsoniter.Unmarshal(data, req); err != nil {
			return response.SetError(err)
		}
		chain, err := s.rly.FindChainRuntime(req.Chain)
		if err != nil {
			return response.SetError(err)
		}
		key := &types.MessageKey{Src: req.Chain, Sn: new(big.Int).SetUint64(req.Sn)}
		received, err := chain.Provider.MessageReceived(context.Background(), key)
		if err != nil {
			return response.SetError(err)
		}
		return response.SetData(&ResMessageReceived{Chain: req.Chain, Sn: req.Sn, Received: received})
	case EventRelayerInfo:
		req := new(ReqRelayInfo)
		if err := jsoniter.Unmarshal(data, req); err != nil {
			return response.SetError(err)
		}
		return response.SetData(&ResRelayInfo{Version: relayer.Version, Uptime: s.startedAt})
	case EventGetBlockEvents:
		req := new(ReqGetBlockEvents)
		if err := jsoniter.Unmarshal(data, req); err != nil {
			return response.SetError(err)
		}
		var events []*ResGetBlockEvents
		for _, chain := range s.rly.GetAllChainsRuntime() {
			msgs, _ := chain.Provider.FetchTxMessages(ctx, req.TxHash)
			for _, msg := range msgs {
				msgKey := types.NewMessageKey(msg.Sn, msg.Src, msg.Dst, msg.EventType)
				received, err := chain.Provider.MessageReceived(ctx, msgKey)
				if err != nil {
					return response.SetError(err)
				}
				data := &ResGetBlockEvents{
					Event:    msg.EventType,
					Height:   msg.MessageHeight,
					Executed: received,
					TxHash:   req.TxHash,
					ChainInfo: struct {
						NID       string                  `json:"nid"`
						Name      string                  `json:"name"`
						Type      string                  `json:"type"`
						Contracts types.ContractConfigMap `json:"contracts"`
					}{
						NID:       chain.Provider.NID(),
						Name:      chain.Provider.Name(),
						Type:      chain.Provider.Type(),
						Contracts: chain.Provider.Config().ContractsAddress(),
					},
				}
				events = append(events, data)
			}
		}
		return response.SetData(events)
	default:
		return response.SetError(fmt.Errorf("unknown event %s", msg.Event))
	}
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) IsClosed() bool {
	return s.listener == nil
}

func getEnvOrFallback(key string, fallback string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return fallback
}
