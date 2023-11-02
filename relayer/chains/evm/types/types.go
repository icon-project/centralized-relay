package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"nhooyr.io/websocket"
)

type EventLog struct {
	Addr    Address
	Indexed [][]byte
	Data    [][]byte
}

type EventLogStr struct {
	Addr    Address  `json:"scoreAddress"`
	Indexed []string `json:"indexed"`
	Data    []string `json:"data"`
}

type TransactionResult struct {
	To                 Address       `json:"to"`
	CumulativeStepUsed HexInt        `json:"cumulativeStepUsed"`
	StepUsed           HexInt        `json:"stepUsed"`
	StepPrice          HexInt        `json:"stepPrice"`
	EventLogs          []EventLogStr `json:"eventLogs"`
	LogsBloom          HexBytes      `json:"logsBloom"`
	Status             HexInt        `json:"status"`
	Failure            *struct {
		CodeValue    HexInt `json:"code"`
		MessageValue string `json:"message"`
	} `json:"failure,omitempty"`
	SCOREAddress Address  `json:"scoreAddress,omitempty"`
	BlockHash    HexBytes `json:"blockHash" validate:"required,t_hash"`
	BlockHeight  HexInt   `json:"blockHeight" validate:"required,t_int"`
	TxIndex      HexInt   `json:"txIndex" validate:"required,t_int"`
	TxHash       HexBytes `json:"txHash" validate:"required,t_int"`
}

type TransactionParamForEstimate struct {
	Version     HexInt   `json:"version" validate:"required,t_int"`
	FromAddress Address  `json:"from" validate:"required,t_addr_eoa"`
	ToAddress   Address  `json:"to" validate:"required,t_addr"`
	Value       HexInt   `json:"value,omitempty" validate:"optional,t_int"`
	Timestamp   HexInt   `json:"timestamp" validate:"required,t_int"`
	NetworkID   HexInt   `json:"nid" validate:"required,t_int"`
	Nonce       HexInt   `json:"nonce,omitempty" validate:"optional,t_int"`
	DataType    string   `json:"dataType,omitempty" validate:"optional,call|deploy|message|deposit"`
	Data        CallData `json:"data,omitempty"`
}
type TransactionParam struct {
	Version     HexInt   `json:"version" validate:"required,t_int"`
	FromAddress Address  `json:"from" validate:"required,t_addr_eoa"`
	ToAddress   Address  `json:"to" validate:"required,t_addr"`
	Value       HexInt   `json:"value,omitempty" validate:"optional,t_int"`
	StepLimit   HexInt   `json:"stepLimit" validate:"required,t_int"`
	Timestamp   HexInt   `json:"timestamp" validate:"required,t_int"`
	NetworkID   HexInt   `json:"nid" validate:"required,t_int"`
	Nonce       HexInt   `json:"nonce,omitempty" validate:"optional,t_int"`
	Signature   string   `json:"signature" validate:"required,t_sig"`
	DataType    string   `json:"dataType,omitempty" validate:"optional,call|deploy|message"`
	Data        CallData `json:"data,omitempty"`
	TxHash      HexBytes `json:"-"`
}

type BlockHeaderResult struct {
	StateHash        []byte
	PatchReceiptHash []byte
	ReceiptHash      common.Hash
	ExtensionData    []byte
}

type TxResult struct {
	Status             int64
	To                 []byte
	CumulativeStepUsed []byte
	StepUsed           []byte
	StepPrice          []byte
	LogsBloom          []byte
	EventLogs          []EventLog
	ScoreAddress       []byte
	EventLogsHash      common.Hash
	TxIndex            HexInt
	BlockHeight        HexInt
}

type CallData struct {
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

type CallParam struct {
	FromAddress Address   `json:"from" validate:"optional,t_addr_eoa"`
	ToAddress   Address   `json:"to" validate:"required,t_addr_score"`
	DataType    string    `json:"dataType" validate:"required,call"`
	Data        *CallData `json:"data"`
	Height      *big.Int  `json:"height,omitempty"`
}

// Added to implement RelayerMessage interface
func (c *CallParam) Type() string {
	return c.DataType
}

func (c *CallParam) MsgBytes() ([]byte, error) {
	return nil, nil
}

type AddressParam struct {
	Address Address `json:"address" validate:"required,t_addr"`
	Height  HexInt  `json:"height,omitempty" validate:"optional,t_int"`
}

type TransactionHashParam struct {
	Hash HexBytes `json:"txHash" validate:"required,t_hash"`
}

type BlockHeightParam struct {
	Height *big.Int `json:"height" validate:"required,t_int"`
}
type DataHashParam struct {
	Hash HexBytes `json:"hash" validate:"required,t_hash"`
}
type ProofResultParam struct {
	BlockHash HexBytes `json:"hash" validate:"required,t_hash"`
	Index     HexInt   `json:"index" validate:"required,t_int"`
}
type ProofEventsParam struct {
	BlockHash HexBytes `json:"hash" validate:"required,t_hash"`
	Index     HexInt   `json:"index" validate:"required,t_int"`
	Events    []HexInt `json:"events"`
}

type BlockRequest struct {
	Height       HexInt         `json:"height"`
	EventFilters []*EventFilter `json:"eventFilters,omitempty"`
}

type EventFilter struct {
	Addr      Address   `json:"addr,omitempty"`
	Signature string    `json:"event"`
	Indexed   []*string `json:"indexed,omitempty"`
	Data      []*string `json:"data,omitempty"`
}

type BlockNotification struct {
	Hash    HexBytes     `json:"hash"`
	Height  HexInt       `json:"height"`
	Indexes [][]HexInt   `json:"indexes,omitempty"`
	Events  [][][]HexInt `json:"events,omitempty"`
}

type EventRequest struct {
	EventFilter
	Height HexInt `json:"height"`
}

type EventNotification struct {
	Hash   HexBytes `json:"hash"`
	Height uint64   `json:"height"`
	Index  HexInt   `json:"index"`
	Events []HexInt `json:"events,omitempty"`
}

type WSEvent string

const (
	WSEventInit WSEvent = "WSEventInit"
)

type WSResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// T_BIN_DATA, T_HASH
type HexBytes string

func (hs HexBytes) Value() ([]byte, error) {
	if hs == "" {
		return nil, nil
	}
	return hex.DecodeString(string(hs[2:]))
}

func NewHexBytes(b []byte) HexBytes {
	return HexBytes("0x" + hex.EncodeToString(b))
}

// T_INT
type HexInt string

func (i HexInt) Value() (int64, error) {
	s := string(i)
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	return strconv.ParseInt(s, 16, 64)
}

func (i HexInt) Int() (int, error) {
	s := string(i)
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	v, err := strconv.ParseInt(s, 16, 32)
	return int(v), err
}

func (i HexInt) BigInt() (*big.Int, error) {
	bi := new(big.Int)

	if err := ParseBigInt(bi, string(i)); err != nil {
		return nil, err
	} else {
		return bi, nil
	}
}

func decodeHexNumber(s string) (bool, []byte, error) {
	negative := false
	if len(s) > 0 && s[0] == '-' {
		negative = true
		s = s[1:]
	}
	if len(s) > 2 && s[0:2] == "0x" {
		s = s[2:]
	}
	if (len(s) % 2) == 1 {
		s = "0" + s
	}
	bs, err := hex.DecodeString(s)
	return negative, bs, err
}

func ParseBigInt(i *big.Int, s string) error {
	neg, bs, err := decodeHexNumber(s)
	if err != nil {
		return err
	}
	i.SetBytes(bs)
	if neg {
		i.Neg(i)
	}
	return nil
}

func NewHexInt(v int64) HexInt {
	return HexInt("0x" + strconv.FormatInt(v, 16))
}

// T_ADDR_EOA, T_ADDR_SCORE
type Address string

func (a Address) Value() ([]byte, error) {
	var b [21]byte
	switch a[:2] {
	case "cx":
		b[0] = 1
	case "hx":
	default:
		return nil, fmt.Errorf("invalid prefix %s", a[:2])
	}
	n, err := hex.Decode(b[1:], []byte(a[2:]))
	if err != nil {
		return nil, err
	}
	if n != 20 {
		return nil, fmt.Errorf("invalid length %d", n)
	}
	return b[:], nil
}

func NewAddress(b []byte) Address {
	if len(b) != 21 {
		return ""
	}
	switch b[0] {
	case 1:
		return Address("cx" + hex.EncodeToString(b[1:]))
	case 0:
		return Address("hx" + hex.EncodeToString(b[1:]))
	default:
		return ""
	}
}

type NormalTransactions struct {
	TxHash   HexBytes        `json:"txHash"`
	From     Address         `json:"from"`
	To       Address         `json:"to"`
	DataType string          `json:"dataType,omitempty"`
	Data     json.RawMessage `json:"data,omitempty"`
}

type Block struct {
	Height             uint64 `json:"height" validate:"required,t_int"`
	Timestamp          uint64 `json:"time_stamp" validate:"required,t_int"`
	NormalTransactions []NormalTransactions
}

type RpcReadCallback func(interface{}) error
