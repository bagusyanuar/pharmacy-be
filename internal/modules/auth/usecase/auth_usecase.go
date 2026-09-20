package usecase

import (
	"context"
	"errors"

	"go.uber.org/zap"

	authdomain "github.com/bagusyanuar/pharmacy-be/internal/modules/auth/domain"
	userdomain "github.com/bagusyanuar/pharmacy-be/internal/modules/user/domain"
	"github.com/bagusyanuar/pharmacy-be/pkg/jwt"
	"github.com/bagusyanuar/pharmacy-be/pkg/password"
)

// TokenPair holds access and refresh tokens.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// AuthUsecase handles user authentication, registration, and token management.
type AuthUsecase struct {
	userRepo       userdomain.UserRepository
	accessManager  *jwt.Manager
	refreshManager *jwt.Manager
	log            *zap.Logger
}

func NewAuthUsecase(
	userRepo userdomain.UserRepository,
	accessManager *jwt.Manager,
	refreshManager *jwt.Manager,
	log *zap.Logger,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:       userRepo,
		accessManager:  accessManager,
		refreshManager: refreshManager,
		log:            log,
	}
}

// Register creates a new user account and returns the created user entity.
func (uc *AuthUsecase) Register(ctx context.Context, email, username, plainPassword string) (*userdomain.User, error) {
	if existing, _ := uc.userRepo.GetByEmail(ctx, email); existing != nil {
		return nil, userdomain.ErrEmailAlreadyExists
	}
	if existing, _ := uc.userRepo.GetByUsername(ctx, username); existing != nil {
		return nil, userdomain.ErrUsernameAlreadyExists
	}

	hash, err := password.Hash(plainPassword)
	if err != nil {
		uc.log.Error("failed to hash password", zap.Error(err))
		return nil, err
	}

	newUser := userdomain.NewUser(email, username, hash)
	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		uc.log.Error("failed to create user in repository", zap.Error(err))
		return nil, err
	}

	return newUser, nil
}

// Login authenticates credentials and generates a token pair.
func (uc *AuthUsecase) Login(ctx context.Context, email, plainPassword string) (*TokenPair, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, userdomain.ErrUserNotFound) {
			return nil, authdomain.ErrInvalidCredentials
		}
		uc.log.Error("failed to query user by email", zap.Error(err))
		return nil, err
	}

	if err := password.Compare(user.PasswordHash, plainPassword); err != nil {
		return nil, authdomain.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, authdomain.ErrAccountInactive
	}

	return uc.issueTokenPair(user.ID)
}

// RefreshToken validates a refresh token and generates a new token pair.
func (uc *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := uc.refreshManager.ParseToken(refreshToken)
	if err != nil {
		return nil, authdomain.ErrInvalidRefreshToken
	}

	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, userdomain.ErrUserNotFound) {
			return nil, authdomain.ErrInvalidRefreshToken
		}
		uc.log.Error("failed to find user during token refresh", zap.Error(err))
		return nil, err
	}

	if !user.IsActive {
		return nil, authdomain.ErrAccountInactive
	}

	return uc.issueTokenPair(user.ID)
}

// Me retrieves the profile of the current authenticated user.
func (uc *AuthUsecase) Me(ctx context.Context, userID string) (*userdomain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		if !errors.Is(err, userdomain.ErrUserNotFound) {
			uc.log.Error("failed to fetch user profile", zap.Error(err))
		}
		return nil, err
	}
	return user, nil
}

func (uc *AuthUsecase) issueTokenPair(userID string) (*TokenPair, error) {
	accessToken, err := uc.accessManager.GenerateToken(userID)
	if err != nil {
		uc.log.Error("failed to generate access token", zap.Error(err))
		return nil, err
	}

	refreshToken, err := uc.refreshManager.GenerateToken(userID)
	if err != nil {
		uc.log.Error("failed to generate refresh token", zap.Error(err))
		return nil, err
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
