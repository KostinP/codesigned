// internal/shared/pagination/http.go
package pagination

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// PaginationQueryParams — DTO для парсинга query-параметров из запроса
type PaginationQueryParams struct {
	Limit  int
	Offset int
	SortBy string
	Order  string
}

// ParsePaginationParams — парсит query-параметры из echo.Context
func ParsePaginationParams(c echo.Context) PaginationQueryParams {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	sortBy := strings.TrimSpace(c.QueryParam("sort_by"))
	order := strings.ToLower(c.QueryParam("order"))

	if limit <= 0 || limit > MaxLimit {
		limit = DefaultLimit
	}
	if offset < 0 {
		offset = DefaultOffset
	}
	if order != "asc" && order != "desc" {
		order = DefaultOrder
	}

	return PaginationQueryParams{
		Limit:  limit,
		Offset: offset,
		SortBy: sortBy,
		Order:  order,
	}
}

// ToDomainParams — конвертирует HTTP DTO в доменный тип Params
func (p PaginationQueryParams) ToDomainParams() Params {
	return Params{
		Limit:  p.Limit,
		Offset: p.Offset,
		SortBy: p.SortBy,
		Order:  p.Order,
	}
}
