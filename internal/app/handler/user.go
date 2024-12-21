package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tranminhquanq/gomess/internal/app/usecase"
	"github.com/tranminhquanq/gomess/internal/utils"
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

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	// filters := utils.ParseFilter(r)
	// sorts := utils.ParseSort(r)
	page, limit := utils.ParsePagination(r)

	result, err := h.userUsecase.UsersWithPagination() // TODO: Add pagination
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err)
	}

	return sendJSON(w, http.StatusOK, NewPaginationResponse(result.Items, NewPaginationMeta(result.Count, page, limit)))
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) error {
	userId := chi.URLParam(r, "userId") // TODO: Get user from token instead of URL
	user, err := h.userUsecase.UserDetails(userId)
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err)
	}

	return sendJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetUserDetails(w http.ResponseWriter, r *http.Request) error {
	userId := chi.URLParam(r, "userId")
	user, err := h.userUsecase.UserDetails(userId)
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err)
	}

	return sendJSON(w, http.StatusOK, user)
}
