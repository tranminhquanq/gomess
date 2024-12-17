package utils

import (
	"net/http"
	"strconv"
	"strings"
)

type Pagination struct {
	Page   int `json:"page"`   // The current page (starting from 1)
	Limit  int `json:"limit"`  // The number of items per page
	Offset int `json:"offset"` // Optional: The offset value for pagination (used in some approaches)
}

// ParsePagination extracts pagination parameters from the query string
func ParsePagination(r *http.Request) (page, limit int) {
	// Default values
	page = 1
	limit = 10

	query := r.URL.Query()

	if p, ok := query["page"]; ok {
		if pageNum, err := strconv.Atoi(p[0]); err == nil {
			page = pageNum
		}
	}

	if l, ok := query["limit"]; ok {
		if limitNum, err := strconv.Atoi(l[0]); err == nil {
			limit = limitNum
		}
	}

	return page, limit
}

// Filter represents a single filter condition
type Filter struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Sort struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

// ParseFilter extracts filter parameters from the query string
func ParseFilter(r *http.Request) (filters []Filter) {
	query := r.URL.Query()

	for key, values := range query {
		if key == "filter" {
			for _, value := range values {
				parts := strings.Split(value, ":")
				if len(parts) >= 3 { // field:operator:value
					filters = append(filters, Filter{
						Field:    parts[0],
						Operator: parts[1],
						Value:    strings.Join(parts[2:], ":"), // in case the value contains ":"
					})
				}
			}
		}
	}

	return filters
}

// ParseSort extracts sort parameters from the query string
func ParseSort(r *http.Request) (sorts []Sort) {
	query := r.URL.Query()

	for key, values := range query {
		if key == "sort" {
			for _, value := range values {
				parts := strings.Split(value, ":")
				if len(parts) == 2 { // field:order
					sorts = append(sorts, Sort{
						Field: parts[0],
						Order: parts[1],
					})
				}
			}
		}
	}

	return sorts
}
