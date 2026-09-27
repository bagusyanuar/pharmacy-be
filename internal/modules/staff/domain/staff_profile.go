package domain

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/bagusyanuar/pharmacy-be/pkg/response"
)

// StaffProfile represents the physical HR employee profile entity (BR-PHARM-AUTH-13).
type StaffProfile struct {
	ID           string     `json:"id"`
	UserID       *string    `json:"user_id,omitempty"`
	EmployeeCode string     `json:"employee_code"`
	NIK          *string    `json:"nik,omitempty"`
	FullName     string     `json:"full_name"`
	Gender       *string    `json:"gender,omitempty"`
	Phone        string     `json:"phone"`
	Email        *string    `json:"email,omitempty"`
	Address      *string    `json:"address,omitempty"`
	HireDate     time.Time  `json:"hire_date"`
	JobPosition  string     `json:"job_position"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CreatedBy    *string    `json:"created_by,omitempty"`
	DeletedAt    *time.Time `json:"-"`
}

// NewStaffProfile creates a new StaffProfile instance.
func NewStaffProfile(employeeCode, fullName, phone, jobPosition string, hireDate time.Time) *StaffProfile {
	now := time.Now()
	return &StaffProfile{
		ID:           uuid.NewString(),
		EmployeeCode: employeeCode,
		FullName:     fullName,
		Phone:        phone,
		JobPosition:  jobPosition,
		HireDate:     hireDate,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// StaffProfileRepository defines the persistence contract for StaffProfile.
type StaffProfileRepository interface {
	GetByID(ctx context.Context, id string) (*StaffProfile, error)
	GetByUserID(ctx context.Context, userID string) (*StaffProfile, error)
	GetByEmployeeCode(ctx context.Context, code string) (*StaffProfile, error)
	GetByNIK(ctx context.Context, nik string) (*StaffProfile, error)
	List(ctx context.Context, params response.PaginationParams) ([]*StaffProfile, int64, error)
	Create(ctx context.Context, profile *StaffProfile) error
	Update(ctx context.Context, profile *StaffProfile) error
	Delete(ctx context.Context, id string) error
}
