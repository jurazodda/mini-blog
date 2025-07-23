package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaginatedResponse(t *testing.T) {
	data := []string{"a", "b"}
	pagination := PaginationQuery{Page: 2, PageSize: 10}
	resp := NewPaginatedResponse(data, pagination, 25)
	assert.Equal(t, data, resp.Data)
	assert.Equal(t, 2, resp.Page)
	assert.Equal(t, 10, resp.PageSize)
	assert.Equal(t, int64(25), resp.Total)
	assert.Equal(t, 3, resp.TotalPages)
}

func TestPaginationQuery_SetDefaults(t *testing.T) {
	p := PaginationQuery{}
	p.SetDefaults()
	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 10, p.PageSize)
	p.Page = -2
	p.PageSize = 200
	p.SetDefaults()
	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 100, p.PageSize)
}

func TestPaginationQuery_GetOffsetAndLimit(t *testing.T) {
	p := PaginationQuery{Page: 3, PageSize: 20}
	assert.Equal(t, 40, p.GetOffset())
	assert.Equal(t, 20, p.GetLimit())
}
