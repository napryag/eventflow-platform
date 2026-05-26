package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/napryag/eventflow-platform/services/authhub/internal/application"
	"github.com/napryag/eventflow-platform/services/authhub/internal/domain/user"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Create(ctx context.Context, user *user.User) error {
	model := toModel(user)
	if err := u.db.Create(model).Error; err != nil {
		if isDuplicateEmailError(err) {
			return application.ErrUserAlreadyExists
		}

		return err
	}
	return nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var model userModel

	if err := u.db.Where("email = ?", email).
		First(&model).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, application.ErrUserNotFound
		}

		return nil, err
	}

	return toDomain(model)
}

func (u *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var model userModel

	if err := u.db.Where("id = ?", id).
		First(&model).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, application.ErrUserNotFound
		}

		return nil, err
	}

	return toDomain(model)
}
