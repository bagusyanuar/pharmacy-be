package domain

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrAccountInactive     = errors.New("account is inactive")
)
