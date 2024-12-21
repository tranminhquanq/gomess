package repository

import (
	"github.com/gofrs/uuid"
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

func (repo *UserRepositoryImpl) FindUserById(id uuid.UUID) (domain.User, error) {
	// userModel, err := findUser(repo.db, "id = ?", id)

	// if err != nil {
	// 	return domain.User{}, err
	// }

	// return userFactory.CreateUser(
	// 	userModel.ID.String(),
	// 	userModel.GetName(),
	// 	userModel.Email.String(),
	// ), nil

	return userFactory.CreateUser(
		uuid.Must(uuid.NewV4()).String(),
		"Quang Tran",
		"tranminhquang4421@gmail.com",
	), nil
}

func (repo *UserRepositoryImpl) UpdateUser(user domain.User) (domain.User, error) {
	return user, nil
}

func (repo *UserRepositoryImpl) FindWithPagination() (domain.ListResult[domain.User], error) {
	users := []domain.User{
		userFactory.CreateUser(uuid.Must(uuid.NewV4()).String(), "John Doe", "john.doe@exampl.com"),
		userFactory.CreateUser(uuid.Must(uuid.NewV4()).String(), "Quang Tran", "quang.tran@example.com"),
	}

	return domain.ListResult[domain.User]{Items: users, Count: int64(len(users))}, nil
}

// func findUser(tx *storage.Connection, query string, args ...interface{}) (*models.User, error) {
// 	user := &models.User{}

// 	if err := tx.Eager().Q().Where(query, args...).First(user); err != nil {
// 		if errors.Cause(err) == sql.ErrNoRows {
// 			return nil, models.UserNotFoundError{}
// 		}
// 		return nil, errors.Wrap(err, "failed to find user")
// 	}

// 	return user, nil
// }
