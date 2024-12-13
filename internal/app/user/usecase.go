package user

import "context"

type UserUsecase interface {
	Details(ctx context.Context) error
	ListWithPagination(ctx context.Context) error
	PartialUpdate(ctx context.Context) error
	UpdateStatus(ctx context.Context) error
}
