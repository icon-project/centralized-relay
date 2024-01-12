package types

import "io"

type TxSearchParam struct {
	Query   string
	Prove   bool
	Page    *int
	PerPage *int
	OrderBy string
}

type KeyringPassword string

func (kp KeyringPassword) Read(p []byte) (n int, err error) {
	copy(p, kp)
	return len(kp), io.EOF
}
