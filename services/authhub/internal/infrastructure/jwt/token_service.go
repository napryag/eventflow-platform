package jwt

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/napryag/eventflow-platform/pkg/auth"
	"github.com/napryag/eventflow-platform/pkg/errs"
	"github.com/napryag/eventflow-platform/services/authhub/config"
	"github.com/napryag/eventflow-platform/services/authhub/internal/application"
)

type TokenService struct {
	cfg config.Config
}

type accessClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type refreshClaims struct {
	jwt.RegisteredClaims
}

func NewService(cfg config.Config) TokenService {
	return TokenService{cfg: cfg}
}

func (ts *TokenService) GeneratePair(ctx context.Context, userID uuid.UUID, email string) (*application.TokenPair, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	accessClaims := accessClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ts.cfg.Token.AccessTTLMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    ts.cfg.Token.Issuer,
			Audience:  jwt.ClaimStrings{ts.cfg.Token.Audience},
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	refreshClaims := refreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ts.cfg.Token.RefreshTTLHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedAccessToken, err := accessToken.SignedString([]byte(ts.cfg.Token.AccessSecret))
	if err != nil {
		return nil, errs.New("failed to sign access token").Wrap(err)
	}

	signedRefreshToken, err := refreshToken.SignedString([]byte(ts.cfg.Token.RefreshSecret))
	if err != nil {
		return nil, errs.New("failed to sign refresh token").Wrap(err)
	}

	return &application.TokenPair{AccessToken: signedAccessToken, RefreshToken: signedRefreshToken}, nil
}

func (ts *TokenService) ValidateAccessToken(ctx context.Context, token string) (*auth.Claims, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	parser := jwt.NewParser(
		jwt.WithIssuer(ts.cfg.Token.Issuer),
		jwt.WithAudience(ts.cfg.Token.Audience),
		jwt.WithValidMethods([]string{"HS256"}),
	)

	parsed, err := parser.ParseWithClaims(token, &accessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.New("unexpected signing method").Arg("alg", token.Header["alg"])
		}
		return []byte(ts.cfg.Token.AccessSecret), nil
	})
	if err != nil {
		return nil, errs.New("invalid token").Wrap(err)
	}

	claims, ok := parsed.Claims.(*accessClaims)
	if !ok || !parsed.Valid {
		return nil, errs.New("invalid token claims")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, errs.New("invalid subject UUID").Wrap(err)
	}

	return &auth.Claims{
		UserID: userID,
		Email:  claims.Email,
	}, nil
}
