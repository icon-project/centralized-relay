package types

type TxSearchParam struct {
	Query   string
	Prove   bool
	Page    *int
	PerPage *int
	OrderBy string
}
