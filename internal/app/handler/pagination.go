package handler

type PaginationMeta struct {
	ItemCount  int64 `json:"item_count"`
	TotalPages int64 `json:"total_pages"`
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

func NewPaginationMeta(totalCount, page, limit int64) PaginationMeta {
	totalPages := totalCount / limit
	if totalCount%limit > 0 {
		totalPages++
	}

	return PaginationMeta{
		ItemCount:  totalCount,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

type PaginationResponse[T any] struct {
	Data T              `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

func NewPaginationResponse[T any](data T, meta PaginationMeta) PaginationResponse[T] {
	return PaginationResponse[T]{
		Data: data,
		Meta: meta,
	}
}
