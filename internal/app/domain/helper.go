package domain

type ListResult[T any] struct {
	Items []T
	Count int64
}
