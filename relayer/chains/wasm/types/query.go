package types

import "fmt"

var QueryOperator = struct{ Eq, Gt, Gte, Lt, Lte string }{
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

func (q *Query) GetQuery() string {
	operator := QueryOperator.Eq
	if q.Operator != "" {
		operator = q.Operator
	}
	return fmt.Sprintf("%s%s%v", q.Field, operator, q.Value)
}

type CompositeQuery struct {
	Or      bool
	Queries []QueryExpression
}

func (cq *CompositeQuery) GetQuery() string {
	merger := "AND"
	if cq.Or {
		merger = "OR"
	}

	var resultQuery string
	for _, q := range cq.Queries {
		if q.GetQuery() != "" && resultQuery != "" {
			resultQuery = fmt.Sprintf("%s %s %s", resultQuery, merger, q.GetQuery())
		} else if q.GetQuery() != "" {
			resultQuery = q.GetQuery()
		}
	}

	return resultQuery
}
