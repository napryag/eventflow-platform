package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/napryag/eventflow-platform/services/authhub/internal/application"
	"github.com/napryag/eventflow-platform/services/authhub/internal/domain/user"
	"gorm.io/gorm"
)

const uniqueViolationCode = "23505"

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Create(ctx context.Context, user *user.User) error {
	var pgErr *pgconn.PgError
	model := toModel(user)
	if err := u.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return application.ErrUserAlreadyExists
		}

		return err
	}
	return nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var model userModel

	if err := u.db.WithContext(ctx).Where("email = ?", email).
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

	if err := u.db.WithContext(ctx).Where("id = ?", id).
		First(&model).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, application.ErrUserNotFound
		}

		return nil, err
	}

	return toDomain(model)
}
