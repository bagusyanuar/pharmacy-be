package domain

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/bagusyanuar/pharmacy-be/pkg/response"
)

// BranchType defines operational classification of a pharmacy facility.
type BranchType string

const (
	BranchTypeRetailPharmacy   BranchType = "RETAIL_PHARMACY"
	BranchTypeClinicPharmacy   BranchType = "CLINIC_PHARMACY"
	BranchTypeCentralWarehouse BranchType = "CENTRAL_WAREHOUSE"
)

// Branch represents the physical pharmacy store or warehouse entity.
type Branch struct {
	ID             string     `json:"id"`
	BranchCode     string     `json:"branch_code"`
	BranchName     string     `json:"branch_name"`
	BranchType     BranchType `json:"branch_type"`
	SIANumber      *string    `json:"sia_number,omitempty"`
	SIAExpiredDate *time.Time `json:"sia_expired_date,omitempty"`
	Phone          *string    `json:"phone,omitempty"`
	Email          *string    `json:"email,omitempty"`
	Address        string     `json:"address"`
	City           string     `json:"city"`
	Province       string     `json:"province"`
	PostalCode     *string    `json:"postal_code,omitempty"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedBy      *string    `json:"created_by,omitempty"`
	DeletedAt      *time.Time `json:"-"`
}

// NewBranch creates a new Branch instance.
func NewBranch(branchCode, branchName string, branchType BranchType, address, city, province string) *Branch {
	now := time.Now()
	return &Branch{
		ID:         uuid.NewString(),
		BranchCode: branchCode,
		BranchName: branchName,
		BranchType: branchType,
		Address:    address,
		City:       city,
		Province:   province,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// BranchRepository defines the persistence contract for Branch.
type BranchRepository interface {
	GetByID(ctx context.Context, id string) (*Branch, error)
	GetByCode(ctx context.Context, code string) (*Branch, error)
	List(ctx context.Context, params response.PaginationParams) ([]*Branch, int64, error)
	Create(ctx context.Context, branch *Branch) error
	Update(ctx context.Context, branch *Branch) error
}
