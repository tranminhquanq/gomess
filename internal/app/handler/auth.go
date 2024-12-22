package handler

import (
	"net/http"

	"github.com/tranminhquanq/gomess/internal/app/usecase"
	"github.com/tranminhquanq/gomess/internal/config"
)

type AuthHandler struct {
	globalConfig *config.GlobalConfiguration
	userUsecase  *usecase.UserUsecase
}

func NewAuthHandler(
	globalConfig *config.GlobalConfiguration,
	userUsecase *usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{
		globalConfig: globalConfig,
		userUsecase:  userUsecase,
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
