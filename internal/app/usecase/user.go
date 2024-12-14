package usecase

import (
	"github.com/tranminhquanq/gomess/internal/app/domain"
	"github.com/tranminhquanq/gomess/internal/app/domain/repository"
)

func UserDetails(repo repository.UserRepository) (domain.User, error) {
	return repo.FindUser(), nil
}

func UsersWithPagination(repo repository.UserRepository) (interface{}, error) {
	users, totalCount := repo.FindUsersWithPagination()

	return NewPaginationResponse(users, NewPaginationMeta(totalCount, 1, 10)), nil
}

func PartialUpdateUser(repo repository.UserRepository, user domain.User) (domain.User, error) {
	return repo.UpdateUser(user), nil
}

func UpdateUserStatus(repo repository.UserRepository) (interface{}, error) {
	return nil, nil
}
