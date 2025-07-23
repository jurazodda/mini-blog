package service

type PaginationQuery struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

func (p *PaginationQuery) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

func (p *PaginationQuery) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationQuery) GetLimit() int {
	return p.PageSize
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

func NewPaginatedResponse(data interface{}, pagination PaginationQuery, total int64) *PaginatedResponse {
	totalPages := int((total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize))
	return &PaginatedResponse{
		Data:       data,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
