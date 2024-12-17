package usecase

import "github.com/tranminhquanq/gomess/internal/app/domain/repository"

type AuthUsecase struct {
	repository repository.UserRepository
}

func NewAuthUsecase(repository repository.UserRepository) *AuthUsecase {
	return &AuthUsecase{
		repository: repository,
	}
}
