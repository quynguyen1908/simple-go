package user

import (
	"context"
	"strings"

	"golang/pkg/constants"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req RegisterRequest) (*UserResponse, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, req RegisterRequest) (*UserResponse, error) {

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
		Status:             constants.StatusActive,
	}

	if err := s.repo.CreateUser(ctx, &newUser); err != nil {
		return nil, err
	}

	res := &UserResponse{
		ID:       newUser.ID,
		Username: newUser.Username,
		Email:    newUser.Email,
		Role:     role.Name,
	}

	return res, nil
}
