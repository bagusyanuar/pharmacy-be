package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/bagusyanuar/pharmacy-be/pkg/response"
)

// User represents the system user domain entity.
type User struct {
	ID           string         `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	Email        string         `gorm:"column:email;type:varchar(255);uniqueIndex;not null" json:"email"`
	Username     string         `gorm:"column:username;type:varchar(255);uniqueIndex;not null" json:"username"`
	PasswordHash string         `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	IsActive     bool           `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt    time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// UserSummary is a minimal projection for cross-module consumption (see .agents/rules/architecture.md §7).
type UserSummary struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Summary returns a UserSummary projection.
func (u *User) Summary() *UserSummary {
	if u == nil {
		return nil
	}
	return &UserSummary{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}

// NewUser creates a new User entity instance with UUID and timestamps.
func NewUser(email, username, passwordHash string) *User {
	now := time.Now()
	return &User{
		ID:           uuid.NewString(),
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// UserRepository is the Domain contract for persisting and querying User records.
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByIDs(ctx context.Context, ids []string) ([]*User, error)
	List(ctx context.Context, params response.PaginationParams) ([]*User, int64, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}
