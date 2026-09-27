package domain

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/bagusyanuar/pharmacy-be/pkg/response"
)

// SupervisorOverrideLog represents supervisor authorization logs at POS cashiers (BR-PHARM-AUTH-12).
type SupervisorOverrideLog struct {
	ID               string     `json:"id"`
	BranchID         string     `json:"branch_id"`
	CashierUserID    string     `json:"cashier_user_id"`
	SupervisorUserID string     `json:"supervisor_user_id"`
	ActionType       string     `json:"action_type"`
	ReasonCategory   string     `json:"reason_category"`
	ReasonNotes      *string    `json:"reason_notes,omitempty"`
	TargetEntityType *string    `json:"target_entity_type,omitempty"`
	TargetEntityID   *string    `json:"target_entity_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	CreatedBy        *string    `json:"created_by,omitempty"`
	DeletedAt        *time.Time `json:"-"`
}

// NewSupervisorOverrideLog creates a new SupervisorOverrideLog entity.
func NewSupervisorOverrideLog(
	branchID, cashierUserID, supervisorUserID, actionType, reasonCategory string,
	reasonNotes, targetEntityType, targetEntityID *string,
) *SupervisorOverrideLog {
	now := time.Now()
	return &SupervisorOverrideLog{
		ID:               uuid.NewString(),
		BranchID:         branchID,
		CashierUserID:    cashierUserID,
		SupervisorUserID: supervisorUserID,
		ActionType:       actionType,
		ReasonCategory:   reasonCategory,
		ReasonNotes:      reasonNotes,
		TargetEntityType: targetEntityType,
		TargetEntityID:   targetEntityID,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// SupervisorOverrideLogRepository defines the persistence contract for SupervisorOverrideLog.
type SupervisorOverrideLogRepository interface {
	GetByID(ctx context.Context, id string) (*SupervisorOverrideLog, error)
	ListByBranch(ctx context.Context, branchID string, params response.PaginationParams) ([]*SupervisorOverrideLog, int64, error)
	Create(ctx context.Context, log *SupervisorOverrideLog) error
}
