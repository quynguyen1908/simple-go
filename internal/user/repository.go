package user

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	SeedRoles(ctx context.Context) error
	CheckExists(ctx context.Context, email, username string) (bool, bool, error)
	ConfirmUserEmail(ctx context.Context, userID uuid.UUID) error

	GetRoleByName(ctx context.Context, roleName string) (*Role, error)
	GetToken(ctx context.Context, tokenValue, tokenType string) (*UserToken, error)

	CreateUser(ctx context.Context, user *User) error
	CreateUserToken(ctx context.Context, token *UserToken) error

	DeleteToken(ctx context.Context, tokenID uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) SeedRoles(ctx context.Context) error {
	roles := []Role{
		{Name: "Admin", NormalizedName: "admin"},
		{Name: "User", NormalizedName: "user"},
	}

	for _, role := range roles {
		err := r.db.WithContext(ctx).
			Where("normalized_name = ?", role.NormalizedName).
			FirstOrCreate(&role).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepository) CheckExists(ctx context.Context, email, username string) (emailExists bool, usernameExists bool, err error) {
	var count int64

	// Check if email exists
	err = r.db.WithContext(ctx).Model(&User{}).
		Where("normalized_email = ?", strings.ToLower(email)).
		Count(&count).Error
	if err != nil {
		return false, false, err
	}
	emailExists = count > 0

	// Check if username exists
	err = r.db.WithContext(ctx).Model(&User{}).
		Where("normalized_username = ?", strings.ToLower(username)).
		Count(&count).Error
	if err != nil {
		return false, false, err
	}
	usernameExists = count > 0

	return emailExists, usernameExists, nil
}

func (r *userRepository) ConfirmUserEmail(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&User{}).
		Where("id = ?", userID).
		Update("email_confirmed", true).Error
}

func (r *userRepository) GetRoleByName(ctx context.Context, roleName string) (*Role, error) {
	var role Role

	err := r.db.WithContext(ctx).
		Where("normalized_name = ?", strings.ToLower(roleName)).
		First(&role).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	return &role, nil
}

func (r *userRepository) GetToken(ctx context.Context, tokenValue, tokenType string) (*UserToken, error) {
	var token UserToken

	err := r.db.WithContext(ctx).
		Where("token_value = ? AND token_type = ?", tokenValue, tokenType).
		First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	return &token, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) CreateUserToken(ctx context.Context, token *UserToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *userRepository) DeleteToken(ctx context.Context, tokenID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&UserToken{}, "id = ?", tokenID).Error
}
