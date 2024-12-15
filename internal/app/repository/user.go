package repository

import (
	"github.com/tranminhquanq/gomess/internal/app/domain"
	"github.com/tranminhquanq/gomess/internal/app/domain/factory"
	"github.com/tranminhquanq/gomess/internal/storage"
)

var (
	userFactory = factory.UserFactory{}
)

type UserRepositoryImpl struct {
	db *storage.Connection
}

func NewUserRepository(db *storage.Connection) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (repo *UserRepositoryImpl) SaveUser(user domain.User) (domain.User, error) {
	return user, nil
}

func (repo *UserRepositoryImpl) FindUser() (domain.User, error) {
	return userFactory.CreateUser(1, "John Doe", "john.doe@exampl.com"), nil
}

func (repo *UserRepositoryImpl) UpdateUser(user domain.User) (domain.User, error) {
	return user, nil
}

func (repo *UserRepositoryImpl) FindUsersWithPagination() (domain.ListResult[domain.User], error) {
	users := []domain.User{
		userFactory.CreateUser(1, "John Doe", "john.doe@exampl.com"),
		userFactory.CreateUser(2, "Quang Tran", "quang.tran@example.com"),
	}

	return domain.ListResult[domain.User]{Items: users, Count: int64(len(users))}, nil
}
