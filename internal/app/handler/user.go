package handler

import (
	"net/http"

	"github.com/tranminhquanq/gomess/internal/app/usecase"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
	// TODO: Add other usecases here
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) error {
	user, err := h.userUsecase.UserDetails()
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err)
	}

	return sendJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.userUsecase.UsersWithPagination()
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err)
	}

	return sendJSON(w, http.StatusOK, users)
}
