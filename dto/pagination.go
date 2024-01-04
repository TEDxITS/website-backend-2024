package dto

type (
	PaginationQuery struct {
		Search  string `form:"search"`
		Page    int    `form:"page"`
		PerPage int    `form:"per_page"`
	}

	PaginationMetadata struct {
		Page    int   `json:"page"`
		PerPage int   `json:"per_page"`
		MaxPage int64 `json:"max_page"`
		Count   int64 `json:"count"`
	}
)

func (p *PaginationQuery) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}

func (p *PaginationMetadata) GetLimit() int {
	return p.PerPage
}

func (p *PaginationMetadata) GetPage() int {
	return p.Page
}
