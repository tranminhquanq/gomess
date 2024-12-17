package handler

import (
	"net/http"

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
	// filters := utils.ExtractFilterQuery(r)
	// sorts := utils.ExtractSortQuery(r)
	page, limit := utils.ParsePagination(r)

	result, err := h.userUsecase.UsersWithPagination() // TODO: Add pagination
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err)
	}

	return sendJSON(w, http.StatusOK, NewPaginationResponse(result.Items, NewPaginationMeta(result.Count, page, limit)))
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) error {
	user, err := h.userUsecase.UserDetails()
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err)
	}

	return sendJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) error {
	user, err := h.userUsecase.UserDetails()
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err)
	}

	return sendJSON(w, http.StatusOK, user)
}
