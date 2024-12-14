package factory

import "github.com/tranminhquanq/gomess/internal/app/domain"

// User is the factory of domain.User
type UserFactory struct{}

func (u UserFactory) CreateUser(
	id int,
	name string,
	email string,
) domain.User {
	return domain.User{
		ID:    id,
		Name:  name,
		Email: email,
	}
}
