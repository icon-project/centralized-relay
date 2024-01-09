package types

import "fmt"

var QueryOperator = struct {
	Eq  string
	Gt  string
	Gte string
	Lt  string
	Lte string
}{
	"=",
	">",
	">=",
	"<",
	"<=",
}

type QueryExpression interface {
	GetQuery() string
}

type Query struct {
	Field    string
	Operator string
	Value    interface{}
}

func (q Query) GetQuery() string {
	return fmt.Sprintf("%s%s%v", q.Field, q.Operator, q.Value)
}

type CompositeQuery struct {
	Or      bool
	Queries []QueryExpression
}

func (cq CompositeQuery) GetQuery() string {
	merger := "AND"
	if cq.Or {
		merger = "OR"
	}

	resultQuery := ""
	for i, q := range cq.Queries {
		if i == 0 {
			resultQuery = q.GetQuery()
		} else {
			resultQuery = fmt.Sprintf("%s %s %s", resultQuery, merger, q.GetQuery())
		}
	}

	return resultQuery
}
