package pagination

import (
	"fmt"
	"strings"
)

type Params struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	SortBy string `json:"sort_by"`
	Order  string `json:"order"`
}

const (
	DefaultLimit  = 20
	MaxLimit      = 100
	DefaultOffset = 0
	DefaultOrder  = "asc"
)

// Normalize — корректирует значения Limit, Offset, Order
func (p *Params) Normalize() {
	if p.Limit <= 0 || p.Limit > MaxLimit {
		p.Limit = DefaultLimit
	}
	if p.Offset < 0 {
		p.Offset = DefaultOffset
	}
	if p.Order != "asc" && p.Order != "desc" {
		p.Order = DefaultOrder
	}
}

// SQLWithPagination — генерирует SQL с LIMIT, OFFSET и ORDER BY
func SQLWithPagination(baseQuery string, p Params, allowedSortFields map[string]string) (string, []interface{}) {
	var clauses []string

	if p.SortBy != "" && allowedSortFields != nil {
		if dbField, ok := allowedSortFields[p.SortBy]; ok {
			clauses = append(clauses, fmt.Sprintf("ORDER BY %s %s", dbField, strings.ToUpper(p.Order)))
		}
	}

	clauses = append(clauses, "LIMIT $1 OFFSET $2")

	sql := baseQuery + " " + strings.Join(clauses, " ")
	args := []interface{}{p.Limit, p.Offset}

	return sql, args
}
