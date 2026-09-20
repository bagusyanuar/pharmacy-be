package usecase

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/bagusyanuar/pharmacy-be/internal/modules/user/domain"
	"github.com/bagusyanuar/pharmacy-be/pkg/response"
)

// UserUsecase orchestrates user business logic.
type UserUsecase struct {
	userRepo domain.UserRepository
	log      *zap.Logger
}

// NewUserUsecase creates a new UserUsecase.
func NewUserUsecase(userRepo domain.UserRepository, log *zap.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
		log:      log,
	}
}

// GetByID retrieves a user by their ID.
func (uc *UserUsecase) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			uc.log.Error("failed to get user by id", zap.String("user_id", id), zap.Error(err))
		}
		return nil, err
	}
	return user, nil
}

// ListUsers retrieves a paginated list of users.
func (uc *UserUsecase) ListUsers(ctx context.Context, params response.PaginationParams) ([]*domain.User, int64, error) {
	users, total, err := uc.userRepo.List(ctx, params)
	if err != nil {
		uc.log.Error("failed to list users", zap.Error(err))
		return nil, 0, err
	}
	return users, total, nil
}
