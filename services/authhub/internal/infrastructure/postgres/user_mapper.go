package postgres

import "github.com/napryag/eventflow-platform/services/authhub/internal/domain/user"

func toModel(u *user.User) userModel {
	return userModel{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func toDomain(model userModel) (*user.User, error) {
	return user.NewUser(
		model.ID,
		model.Email,
		model.PasswordHash,
		model.CreatedAt,
		model.UpdatedAt,
	)
}
