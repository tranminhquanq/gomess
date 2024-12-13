package user

import "context"

type UserRepository interface {
	FindUsers(ctx context.Context) error
	FilterUser(ctx context.Context) error
	UpdateUser(ctx context.Context) error
}
