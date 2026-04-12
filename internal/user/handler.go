package user

import (
	"encoding/json"
	"errors"
	"golang/pkg/constants"
	"golang/pkg/response"
	"net/http"
)

type UserHandler struct {
	service UserService
	appURL  string
}

func NewUserHandler(service UserService, appURL string) *UserHandler {
	return &UserHandler{service: service, appURL: appURL}
}

// RegisterHandler godoc
// @Summary      Register a new user
// @Description  Register a new user with email, username, and password
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "User registration data"
// @Success      201 {object} response.SuccessResponse{data=UserResponse}
// @Failure      400 {object} response.ErrorResponse
// @Failure      409 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /api/v1/users/register [post]
func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if validationErr := validate.Struct(req); validationErr != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrBadRequest, validationErr)
		return
	}

	res, err := h.service.Register(r.Context(), req, h.appURL)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailAlreadyExists), errors.Is(err, ErrUsernameAlreadyExists):
			response.Error(w, http.StatusConflict, constants.ErrConflict, nil)
		case errors.Is(err, ErrRoleNotFound):
			response.Error(w, http.StatusInternalServerError, constants.ErrInternal, nil)
		default:
			response.Error(w, http.StatusInternalServerError, constants.ErrInternal, nil)
		}
		return
	}

	response.Success(w, http.StatusCreated, "User registered successfully", res)
}

// ConfirmEmailHandler godoc
// @Summary      Confirm user email
// @Description  Confirm a user's email address using a verification token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        token query string true "Email verification token"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /api/v1/users/confirm-email [get]
func (h *UserHandler) ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	token := r.URL.Query().Get("token")

	err := h.service.ConfirmEmail(r.Context(), token)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid or expired token", nil)
		return
	}

	response.Success(w, http.StatusOK, "Email confirmed successfully", nil)
}
