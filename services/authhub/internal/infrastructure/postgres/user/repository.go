package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/napryag/eventflow-platform/services/authhub/internal/application"
	"github.com/napryag/eventflow-platform/services/authhub/internal/domain/user"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user user.User) error {
	userModel := toModel(user)
	if err := r.db.WithContext(ctx).Create(userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return application.ErrUserAlreadyExists
		}

		return err
	}
	return nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var userModel model

	if err := r.db.WithContext(ctx).Where("email = ?", email).
		First(&userModel).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, application.ErrUserNotFound
		}

		return nil, err
	}

	return toDomain(userModel)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var userModel model

	if err := r.db.WithContext(ctx).Where("id = ?", id).
		First(&userModel).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, application.ErrUserNotFound
		}

		return nil, err
	}

	return toDomain(userModel)
}
