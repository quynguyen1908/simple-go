package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name           string    `gorm:"type:varchar(255);not null" json:"name"`
	NormalizedName string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"-"`
}

type User struct {
	ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RoleID             uuid.UUID      `gorm:"type:uuid;not null" json:"role_id"`
	Username           string         `gorm:"type:varchar(255);not null" json:"username"`
	NormalizedUsername string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"-"`
	Email              string         `gorm:"type:varchar(255);not null" json:"email"`
	NormalizedEmail    string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"-"`
	EmailConfirmed     bool           `gorm:"default:false" json:"email_confirmed"`
	PasswordHash       string         `gorm:"type:text;not null" json:"-"`
	Status             string         `gorm:"type:varchar(50);not null" json:"status"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	Role    Role        `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"role"`
	Profile UserProfile `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"profile"`
	Tokens  []UserToken `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tokens"`
}

type UserToken struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	LoginProvider string    `gorm:"type:varchar(255);not null" json:"login_provider"`
	TokenType     string    `gorm:"type:varchar(255);not null" json:"token_type"`
	TokenValue    string    `gorm:"type:text;not null" json:"token_value"`
	ExpiresAt     time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type UserProfile struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID               uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	FirstName            string    `gorm:"type:varchar(255)" json:"first_name"`
	LastName             string    `gorm:"type:varchar(255)" json:"last_name"`
	AvatarURL            string    `gorm:"type:text" json:"avatar_url"`
	PhoneNumber          string    `gorm:"type:varchar(20)" json:"phone_number"`
	PhoneNumberConfirmed bool      `gorm:"default:false" json:"phone_number_confirmed"`
	Street               string    `gorm:"type:varchar(255)" json:"street"`
	City                 string    `gorm:"type:varchar(255)" json:"city"`
	State                string    `gorm:"type:varchar(255)" json:"state"`
	Country              string    `gorm:"type:varchar(255)" json:"country"`
	ZipCode              string    `gorm:"type:varchar(20)" json:"zip_code"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
