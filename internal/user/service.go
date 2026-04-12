package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang/pkg/constants"
	"golang/pkg/mailer"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req RegisterRequest, appURL string) (*UserResponse, error)
	ConfirmEmail(ctx context.Context, tokenValue string) error
}

type userService struct {
	repo   UserRepository
	mailer mailer.Mailer
}

func NewUserService(repo UserRepository, mail mailer.Mailer) UserService {
	return &userService{repo: repo, mailer: mail}
}

func (s *userService) Register(ctx context.Context, req RegisterRequest, appURL string) (*UserResponse, error) {

	emailExists, usernameExists, err := s.repo.CheckExists(ctx, req.Email, req.Username)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, ErrEmailAlreadyExists
	}
	if usernameExists {
		return nil, ErrUsernameAlreadyExists
	}

	role, err := s.repo.GetRoleByName(ctx, constants.RoleUser)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := User{
		RoleID:             role.ID,
		Username:           req.Username,
		NormalizedUsername: strings.ToLower(req.Username),
		Email:              req.Email,
		NormalizedEmail:    strings.ToLower(req.Email),
		PasswordHash:       string(hashedPassword),
		Status:             constants.UserStatusActive,
	}

	if err := s.repo.CreateUser(ctx, &newUser); err != nil {
		return nil, err
	}

	verifyTokenString := uuid.New().String()

	verifyToken := UserToken{
		UserID:        newUser.ID,
		LoginProvider: constants.ProviderSystem,
		TokenType:     constants.TokenTypeEmailConfirmation,
		TokenValue:    verifyTokenString,
		ExpiresAt:     time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.CreateUserToken(ctx, &verifyToken); err == nil {
		go func() {
			errMail := s.mailer.SendVerificationEmail(
				newUser.Email,
				verifyTokenString,
				appURL,
			)
			if errMail != nil {
				// Log the error (replace with your logging mechanism)
				fmt.Printf("Failed to send verification email to %s: %v\n", newUser.Email, errMail)
			}
		}()
	}

	res := &UserResponse{
		ID:       newUser.ID,
		Username: newUser.Username,
		Email:    newUser.Email,
		Role:     role.Name,
	}

	return res, nil
}

func (s *userService) ConfirmEmail(ctx context.Context, tokenValue string) error {
	if tokenValue == "" {
		return ErrInvalidToken
	}

	token, err := s.repo.GetToken(ctx, tokenValue, constants.TokenTypeEmailConfirmation)
	if err != nil {
		return ErrTokenNotFound
	}

	if time.Now().After(token.ExpiresAt) {
		return ErrTokenExpired
	}

	if err := s.repo.ConfirmUserEmail(ctx, token.UserID); err != nil {
		return errors.New("failed to confirm email")
	}

	_ = s.repo.DeleteToken(ctx, token.ID)

	return nil
}
