package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserBranch represents user assignment to a pharmacy branch (BR-PHARM-AUTH-02).
type UserBranch struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	BranchID      string     `json:"branch_id"`
	IsDefault     bool       `json:"is_default"`
	CanOperatePOS bool       `json:"can_operate_pos"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CreatedBy     *string    `json:"created_by,omitempty"`
	DeletedAt     *time.Time `json:"-"`
}

// NewUserBranch creates a new UserBranch assignment entity.
func NewUserBranch(userID, branchID string, isDefault, canOperatePOS bool) *UserBranch {
	now := time.Now()
	return &UserBranch{
		ID:            uuid.NewString(),
		UserID:        userID,
		BranchID:      branchID,
		IsDefault:     isDefault,
		CanOperatePOS: canOperatePOS,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// UserBranchRepository defines the persistence contract for UserBranch.
type UserBranchRepository interface {
	GetByUserIDAndBranchID(ctx context.Context, userID, branchID string) (*UserBranch, error)
	ListByUserID(ctx context.Context, userID string) ([]*UserBranch, error)
	Assign(ctx context.Context, userBranch *UserBranch) error
	Revoke(ctx context.Context, userID, branchID string) error
}
