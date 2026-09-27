package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// PharmacistProfile represents professional pharmaceutical license information (SIPA/STRTTK).
type PharmacistProfile struct {
	ID                 string     `json:"id"`
	StaffID            string     `json:"staff_id"`
	LicenseType        string     `json:"license_type"` // SIPA or STRTTK
	LicenseNumber      string     `json:"license_number"`
	LicenseExpiredDate time.Time  `json:"license_expired_date"`
	IsAPA              bool       `json:"is_apa"`
	AssignedBranchID   *string    `json:"assigned_branch_id,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	CreatedBy          *string    `json:"created_by,omitempty"`
	DeletedAt          *time.Time `json:"-"`
}

// NewPharmacistProfile creates a new PharmacistProfile entity.
func NewPharmacistProfile(staffID, licenseType, licenseNumber string, expiry time.Time, isAPA bool, branchID *string) *PharmacistProfile {
	now := time.Now()
	return &PharmacistProfile{
		ID:                 uuid.NewString(),
		StaffID:            staffID,
		LicenseType:        licenseType,
		LicenseNumber:      licenseNumber,
		LicenseExpiredDate: expiry,
		IsAPA:              isAPA,
		AssignedBranchID:   branchID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// PharmacistProfileRepository defines the persistence contract for PharmacistProfile.
type PharmacistProfileRepository interface {
	GetByID(ctx context.Context, id string) (*PharmacistProfile, error)
	GetByStaffID(ctx context.Context, staffID string) (*PharmacistProfile, error)
	GetByLicenseNumber(ctx context.Context, licenseNumber string) (*PharmacistProfile, error)
	Create(ctx context.Context, profile *PharmacistProfile) error
	Update(ctx context.Context, profile *PharmacistProfile) error
}
