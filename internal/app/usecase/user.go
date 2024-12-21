package usecase

import (
	"github.com/gofrs/uuid"
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

func (u *UserUsecase) UserDetails(userId string) (domain.User, error) {
	user, err := u.repository.FindUserById(uuid.FromStringOrNil(userId))
	return user, err
}

func (u *UserUsecase) UsersWithPagination() (domain.ListResult[domain.User], error) {
	return u.repository.FindWithPagination()
}

func (u *UserUsecase) PartialUpdateUser(user domain.User) (domain.User, error) {
	return u.repository.UpdateUser(user)
}

func (u *UserUsecase) UpdateUserStatus() (interface{}, error) {
	return nil, nil
}
