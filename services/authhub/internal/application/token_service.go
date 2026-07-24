package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/napryag/eventflow-platform/pkg/auth"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenService interface {
	GeneratePair(ctx context.Context, userID uuid.UUID, email string) (*TokenPair, error)
	ValidateAccessToken(ctx context.Context, token string) (*auth.Claims, error)
}
