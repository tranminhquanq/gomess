package repository

import (
	"github.com/gofrs/uuid"
	"github.com/tranminhquanq/gomess/internal/app/domain"
)

type UserRepository interface {
	SaveUser(domain.User) (domain.User, error)
	FindUserById(id uuid.UUID) (domain.User, error)
	FindUsersWithPagination() (domain.ListResult[domain.User], error)
	UpdateUser(domain.User) (domain.User, error)
}
