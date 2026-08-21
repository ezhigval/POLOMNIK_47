package ports

type Pagination struct {
	Page  int
	Limit int
}

type PageMeta struct {
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	Total   int  `json:"total"`
	HasNext bool `json:"has_next"`
}

func NormalizePagination(page int, limit int) Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return Pagination{Page: page, Limit: limit}
}
