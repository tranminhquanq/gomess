package usecase

import (
	"github.com/tranminhquanq/gomess/internal/app/domain"
	"github.com/tranminhquanq/gomess/internal/app/domain/repository"
)

type UserUsecase struct {
	repository repository.UserRepository
}

func NewUserUsecase(repository repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		repository: repository,
	}
}

func (u *UserUsecase) UserDetails() (domain.User, error) {
	return u.repository.FindUser()
}

func (u *UserUsecase) UsersWithPagination() (domain.ListResult[domain.User], error) {
	return u.repository.FindUsersWithPagination()
}

func (u *UserUsecase) PartialUpdateUser(user domain.User) (domain.User, error) {
	return u.repository.UpdateUser(user)
}

func (u *UserUsecase) UpdateUserStatus() (interface{}, error) {
	return nil, nil
}
