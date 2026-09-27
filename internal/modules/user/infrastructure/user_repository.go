package infrastructure

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/bagusyanuar/pharmacy-be/internal/modules/user/domain"
	"github.com/bagusyanuar/pharmacy-be/pkg/response"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository constructs a GORM-backed UserRepository implementation.
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Select("id", "email", "password_hash", "pin_hash", "barcode_card", "role_id", "is_active", "created_at", "updated_at", "created_by", "deleted_at").
		Where("id = ?", id).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Select("id", "email", "password_hash", "pin_hash", "barcode_card", "role_id", "is_active", "created_at", "updated_at", "created_by", "deleted_at").
		Where("email = ?", email).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByBarcode(ctx context.Context, barcode string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Select("id", "email", "password_hash", "pin_hash", "barcode_card", "role_id", "is_active", "created_at", "updated_at", "created_by", "deleted_at").
		Where("barcode_card = ?", barcode).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByIDs(ctx context.Context, ids []string) ([]*domain.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var users []*domain.User
	err := r.db.WithContext(ctx).
		Select("id", "email", "password_hash", "pin_hash", "barcode_card", "role_id", "is_active", "created_at", "updated_at", "created_by", "deleted_at").
		Where("id IN ?", ids).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) List(ctx context.Context, params response.PaginationParams) ([]*domain.User, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []*domain.User
	err := r.db.WithContext(ctx).
		Select("id", "email", "barcode_card", "role_id", "is_active", "created_at", "updated_at").
		Offset(params.Offset()).
		Limit(params.PerPage).
		Order("created_at DESC").
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.User{}).Error
}
