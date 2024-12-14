package handler

import (
	"net/http"

	"github.com/tranminhquanq/gomess/internal/app/usecase"
)

func (hdl *APIHandler) GetUser(w http.ResponseWriter, r *http.Request) error {
	user, err := usecase.UserDetails(userRepository)
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err)
	}

	return sendJSON(w, http.StatusOK, user)
}

func (hdl *APIHandler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := usecase.UsersWithPagination(userRepository)
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err)
	}

	return sendJSON(w, http.StatusOK, users)
}
