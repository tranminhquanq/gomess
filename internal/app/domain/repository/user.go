package repository

import "github.com/tranminhquanq/gomess/internal/app/domain"

type UserRepository interface {
	SaveUser(domain.User) (domain.User, error)
	FindUser() (domain.User, error)
	FindUsersWithPagination() (domain.ListResult[domain.User], error)
	UpdateUser(domain.User) (domain.User, error)
}
