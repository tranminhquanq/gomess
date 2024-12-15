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
	return u.repository.FindUser(), nil
}

func (u *UserUsecase) UsersWithPagination() (interface{}, error) {
	users, totalCount := u.repository.FindUsersWithPagination()

	return NewPaginationResponse(users, NewPaginationMeta(totalCount, 1, 10)), nil
}

func (u *UserUsecase) PartialUpdateUser(user domain.User) (domain.User, error) {
	return u.repository.UpdateUser(user), nil
}

func (u *UserUsecase) UpdateUserStatus() (interface{}, error) {
	return nil, nil
}
