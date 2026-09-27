package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Standard system role codes
const (
	RoleCodeSuperAdmin      = "SUPER_ADMIN"
	RoleCodeApoteker        = "APOTEKER"
	RoleCodeAsistenApoteker = "ASISTEN_APOTEKER"
	RoleCodeKasir           = "KASIR"
	RoleCodeGudang          = "GUDANG"
	RoleCodeFinanceOwner    = "FINANCE_OWNER"
)

// Role represents the RBAC authorization role entity.
type Role struct {
	ID                    string     `json:"id"`
	RoleCode              string     `json:"role_code"`
	RoleName              string     `json:"role_name"`
	Description           string     `json:"description,omitempty"`
	CanSupervisorOverride bool       `json:"can_supervisor_override"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	CreatedBy             *string    `json:"created_by,omitempty"`
	DeletedAt             *time.Time `json:"-"`
}

// NewRole constructs a new Role entity.
func NewRole(roleCode, roleName, description string, canOverride bool) *Role {
	now := time.Now()
	return &Role{
		ID:                    uuid.NewString(),
		RoleCode:              roleCode,
		RoleName:              roleName,
		Description:           description,
		CanSupervisorOverride: canOverride,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
}

// RoleRepository defines the persistence contract for Role.
type RoleRepository interface {
	GetByID(ctx context.Context, id string) (*Role, error)
	GetByCode(ctx context.Context, code string) (*Role, error)
	List(ctx context.Context) ([]*Role, error)
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
}
