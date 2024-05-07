package types

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/big"

	"github.com/fardream/go-bcs/bcs"
)

// Uint64 is like `u64` in move.
type Uint64 struct {
	lo uint32
	hi uint32
}

var (
	_ bcs.Marshaler    = (*Uint64)(nil)
	_ bcs.Unmarshaler  = (*Uint64)(nil)
	_ json.Marshaler   = (*Uint64)(nil)
	_ json.Unmarshaler = (*Uint64)(nil)
)

var maxU64 = (&big.Int{}).Lsh(big.NewInt(1), 64)

func checkUint64(bigI *big.Int) error {
	if bigI.Sign() < 0 {
		return fmt.Errorf("%s is negative", bigI.String())
	}

	if bigI.Cmp(maxU64) >= 0 {
		return fmt.Errorf("%s is greater than Max Uint 64", bigI.String())
	}

	return nil
}

// 32 ones
const ones32 uint32 = (1 << 32) - 1

// 1 << 32
var oneLsh32 = big.NewInt(0).Lsh(big.NewInt(1), 32)

func NewBigIntFromUint32(i uint32) *big.Int {
	r := big.NewInt(int64(i & ones32))
	if i > ones32 {
		r = r.Add(r, oneLsh32)
	}
	return r
}

func (i Uint64) Big() *big.Int {
	loBig := NewBigIntFromUint32(i.lo)
	hiBig := NewBigIntFromUint32(i.hi)
	hiBig = hiBig.Lsh(hiBig, 32)

	return hiBig.Add(hiBig, loBig)
}

func (i *Uint64) SetBigInt(bigI *big.Int) error {
	if err := checkUint64(bigI); err != nil {
		return err
	}

	r := make([]byte, 0, 8)
	bs := bigI.Bytes()
	for i := 0; i+len(bs) < 8; i++ {
		r = append(r, 0)
	}
	r = append(r, bs...)

	hi := binary.BigEndian.Uint32(r[0:4])
	lo := binary.BigEndian.Uint32(r[4:])

	i.hi = hi
	i.lo = lo

	return nil
}

func NewUint64FromBigInt(bigI *big.Int) (*Uint64, error) {
	i := &Uint64{}

	if err := i.SetBigInt(bigI); err != nil {
		return nil, err
	}

	return i, nil
}

func (i Uint64) MarshalBCS() ([]byte, error) {
	r := make([]byte, 8)

	binary.LittleEndian.PutUint32(r, i.lo)
	binary.LittleEndian.PutUint32(r[4:], i.hi)

	return r, nil
}

func (i *Uint64) UnmarshalBCS(r io.Reader) (int, error) {
	buf := make([]byte, 8)
	n, err := r.Read(buf)
	if err != nil {
		return n, err
	}
	if n != 8 {
		return n, fmt.Errorf("failed to read 16 bytes for Uint64 (read %d bytes)", n)
	}

	i.lo = binary.LittleEndian.Uint32(buf[0:4])
	i.hi = binary.LittleEndian.Uint32(buf[4:8])

	return n, nil
}

func (i Uint64) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Big().String())
}

func (i *Uint64) UnmarshalJSON(data []byte) error {
	var dataStr string
	if err := json.Unmarshal(data, &dataStr); err != nil {
		return err
	}

	bigI := &big.Int{}
	_, ok := bigI.SetString(dataStr, 10)
	if !ok {
		return fmt.Errorf("failed to parse %s as an integer", dataStr)
	}

	return i.SetBigInt(bigI)
}

func (u Uint64) String() string {
	return u.Big().String()
}
