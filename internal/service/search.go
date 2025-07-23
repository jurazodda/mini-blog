package service

type SearchQuery struct {
	Query    string `form:"q" json:"q"`
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
}

func (s *SearchQuery) Validate() error {
	if s.Query == "" {
		return ErrContentRequired
	}
	return nil
}

func (s *SearchQuery) SetDefaults() {
	if s.Page <= 0 {
		s.Page = 1
	}
	if s.PageSize <= 0 {
		s.PageSize = 10
	}
	if s.PageSize > 100 {
		s.PageSize = 100
	}
}

func (s *SearchQuery) GetOffset() int {
	return (s.Page - 1) * s.PageSize
}

func (s *SearchQuery) GetLimit() int {
	return s.PageSize
}
