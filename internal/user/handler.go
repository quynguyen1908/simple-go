package user

import (
	"encoding/json"
	"errors"
	"golang/pkg/constants"
	"golang/pkg/response"
	"net/http"
)

type UserHandler struct {
	service   UserService
	appURL    string
	jwtSecret string
}

func NewUserHandler(service UserService, appURL string, jwtSecret string) *UserHandler {
	return &UserHandler{service: service, appURL: appURL, jwtSecret: jwtSecret}
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
		switch {
		case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrTokenNotFound):
			response.Error(w, http.StatusBadRequest, "Invalid or expired token", nil)
		default:
			response.Error(w, http.StatusInternalServerError, constants.ErrInternal, nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "Email confirmed successfully", nil)
}

// ResendConfirmationEmailHandler godoc
// @Summary      Resend confirmation email
// @Description  Resend an email confirmation to the user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body ResendConfirmationRequest true "Resend confirmation email request"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorResponse
// @Router       /api/v1/users/resend-confirmation [post]
func (h *UserHandler) ResendConfirmationEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req ResendConfirmationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if validationErr := validate.Struct(req); validationErr != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrBadRequest, validationErr)
		return
	}

	err := h.service.ResendConfirmationEmail(r.Context(), req, h.appURL)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(w, http.StatusOK, "Verification email resent successfully", nil)
}

// LoginHandler godoc
// @Summary      Login user
// @Description  Login with username and password to get JWT token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "User login credentials"
// @Success      200 {object} response.SuccessResponse{data=LoginResponse}
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /api/v1/users/login [post]
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if validationErr := req.Validate(); len(validationErr) > 0 {
		response.Error(w, http.StatusBadRequest, constants.ErrBadRequest, validationErr)
		return
	}

	res, err := h.service.Login(r.Context(), req, h.jwtSecret)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			response.Error(w, http.StatusUnauthorized, "Invalid username/email or password", nil)
		default:
			response.Error(w, http.StatusInternalServerError, constants.ErrInternal, nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "Login successful", res)
}

// ForgotPasswordHandler godoc
// @Summary      Forgot password
// @Description  Request a password reset email
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body ForgotPasswordRequest true "Forgot password request"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /api/v1/users/forgot-password [post]
func (h *UserHandler) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if validationErr := validate.Struct(req); validationErr != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrBadRequest, validationErr)
		return
	}

	err := h.service.ForgotPassword(r.Context(), req, h.appURL)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(w, http.StatusOK, "Password reset email sent if the email exists", nil)
}

// ResetPasswordHandler godoc
// @Summary      Reset password
// @Description  Reset user password with a valid reset token
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body ResetPasswordRequest true "Reset password request"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /api/v1/users/reset-password [post]
func (h *UserHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if validationErr := validate.Struct(req); validationErr != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrBadRequest, validationErr)
		return
	}

	err := h.service.ResetPassword(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrTokenNotFound):
			response.Error(w, http.StatusBadRequest, "Invalid or expired token", nil)
		case errors.Is(err, ErrTokenExpired):
			response.Error(w, http.StatusBadRequest, "Token expired", nil)
		default:
			response.Error(w, http.StatusInternalServerError, constants.ErrInternal, nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "Password reset successfully", nil)
}
