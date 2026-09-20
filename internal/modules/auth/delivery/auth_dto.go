package delivery

import (
	userdomain "github.com/bagusyanuar/pharmacy-be/internal/modules/user/domain"
)

// RegisterRequest is the body payload for POST /auth/register.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,strongpassword"`
}

// LoginRequest is the body payload for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// TokenResponse is returned upon successful authentication.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// ProfileResponse represents the user's authenticated profile.
type ProfileResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func toProfileResponse(user *userdomain.User) ProfileResponse {
	return ProfileResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		IsActive: user.IsActive,
	}
}
