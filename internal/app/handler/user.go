package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tranminhquanq/gomess/internal/app/usecase"
	"github.com/tranminhquanq/gomess/internal/config"
	"github.com/tranminhquanq/gomess/internal/utils"
)

type UserHandler struct {
	globalConfig *config.GlobalConfiguration
	userUsecase  *usecase.UserUsecase
	// TODO: Add other usecases here
}

func NewUserHandler(
	globalConfig *config.GlobalConfiguration,
	userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		globalConfig: globalConfig,
		userUsecase:  userUsecase,
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
	ctx := r.Context()
	claims := getClaims(ctx)

	if claims == nil {
		return internalServerError("No claims found in context")
	}

	aud := h.requestAud(ctx, r)
	audienceFromClaims, _ := claims.GetAudience()
	if len(audienceFromClaims) == 0 || aud != audienceFromClaims[0] {
		return badRequestError(ErrorCodeValidationFailed, "Token audience doesn't match request audience")
	}

	user, err := h.userUsecase.UserDetails(claims.ID)
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
