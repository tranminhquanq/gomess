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

func (u *UserUsecase) UserDetails() (domain.User, error) {
	user, err := u.repository.FindUserById(uuid.FromStringOrNil("73f8606b-78ba-4e43-9fcd-e2a50db4cd71")) // TODO: get user id from context or request
	return user, err
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
