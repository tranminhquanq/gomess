package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func sendJSON(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error encoding json response: %v", obj))
	}
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}

func sendText(w http.ResponseWriter, status int, text string) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	_, err := w.Write([]byte(text))
	return err
}

type PaginationMeta struct {
	ItemCount  int64 `json:"item_count"`
	TotalPages int64 `json:"total_pages"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

func NewPaginationMeta(totalCount int64, page, limit int) PaginationMeta {
	totalPages := totalCount / int64(limit)
	if totalCount%int64(limit) > 0 {
		totalPages++
	}

	return PaginationMeta{
		ItemCount:  totalCount,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
		HasNext:    int64(page) < totalPages,
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
