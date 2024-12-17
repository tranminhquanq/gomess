package handler

import (
	"net/http"

	"github.com/tranminhquanq/gomess/internal/app/usecase"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
	userUsecase *usecase.UserUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase, userUsecase *usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		userUsecase: userUsecase,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	return sendText(w, http.StatusNotImplemented, "Not implemented")
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	return sendText(w, http.StatusNotImplemented, "Not implemented")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	return sendText(w, http.StatusNotImplemented, "Not implemented")
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) error {
	return sendText(w, http.StatusNotImplemented, "Not implemented")
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) error {
	return sendText(w, http.StatusNotImplemented, "Not implemented")
}
