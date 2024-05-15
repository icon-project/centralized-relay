package solana

import (
	"bytes"

	bin "github.com/gagliardetto/binary"
)

func BorshEncode(data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := bin.NewBorshEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func BorshDecode(data []byte, dest interface{}) error {
	dec := bin.NewBorshDecoder(data)
	return dec.Decode(dest)
}
