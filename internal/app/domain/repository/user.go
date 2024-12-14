package repository

import "github.com/tranminhquanq/gomess/internal/app/domain"

type UserRepository interface {
	FindUser() domain.User
	FindUsersWithPagination() ([]domain.User, int64)
	UpdateUser(domain.User) domain.User
}
