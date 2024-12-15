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

func NewUserRepository(db *storage.Connection) UserRepositoryImpl {
	return UserRepositoryImpl{db: db}
}

func (repo UserRepositoryImpl) SaveUser(user domain.User) domain.User {
	return user
}

func (repo UserRepositoryImpl) FindUser() domain.User {
	return userFactory.CreateUser(1, "John Doe", "john.doe@exampl.com")
}

func (repo UserRepositoryImpl) UpdateUser(user domain.User) domain.User {
	return user
}

func (repo UserRepositoryImpl) FindUsersWithPagination() ([]domain.User, int64) {
	users := []domain.User{
		userFactory.CreateUser(1, "John Doe", "john.doe@exampl.com"),
		userFactory.CreateUser(2, "Quang Tran", "quang.tran@example.com"),
	}

	return users, int64(len(users))
}
