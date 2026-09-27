package domain

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/bagusyanuar/pharmacy-be/pkg/response"
)

// User represents the IAM digital credential domain entity.
type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	PINHash      *string    `json:"-"`
	BarcodeCard  *string    `json:"barcode_card,omitempty"`
	RoleID       string     `json:"role_id"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CreatedBy    *string    `json:"created_by,omitempty"`
	DeletedAt    *time.Time `json:"-"`
}

// UserSummary is a minimal projection for cross-module consumption (see .agents/rules/architecture.md §7).
type UserSummary struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	RoleID   string `json:"role_id"`
	IsActive bool   `json:"is_active"`
}

// Summary returns a UserSummary projection.
func (u *User) Summary() *UserSummary {
	if u == nil {
		return nil
	}
	return &UserSummary{
		ID:       u.ID,
		Email:    u.Email,
		RoleID:   u.RoleID,
		IsActive: u.IsActive,
	}
}

// NewUser creates a new User entity instance with UUID and timestamps.
func NewUser(email, passwordHash, roleID string) *User {
	now := time.Now()
	return &User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: passwordHash,
		RoleID:       roleID,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// UserRepository is the Domain contract for persisting and querying User records.
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByBarcode(ctx context.Context, barcode string) (*User, error)
	GetByIDs(ctx context.Context, ids []string) ([]*User, error)
	List(ctx context.Context, params response.PaginationParams) ([]*User, int64, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}
