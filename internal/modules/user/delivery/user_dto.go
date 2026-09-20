package delivery

import (
	"time"

	"github.com/bagusyanuar/pharmacy-be/internal/modules/user/domain"
)

// UserResponse represents a safe JSON representation of a user.
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func toUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}
}

func toUserResponses(users []*domain.User) []UserResponse {
	out := make([]UserResponse, len(users))
	for i, u := range users {
		out[i] = toUserResponse(u)
	}
	return out
}
