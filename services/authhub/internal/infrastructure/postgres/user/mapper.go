package user

import "github.com/napryag/eventflow-platform/services/authhub/internal/domain/user"

func toModel(u user.User) model {
	return model{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func toDomain(m model) (*user.User, error) {
	return user.NewUser(
		m.ID,
		m.Email,
		m.PasswordHash,
		m.CreatedAt,
		m.UpdatedAt,
	)
}
