package jwt

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/napryag/eventflow-platform/pkg/errs"
	"github.com/napryag/eventflow-platform/services/authhub/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testConfig() config.Config {
	return config.Config{
		Token: config.TokenConfig{
			AccessSecret:     "test_access_secret",
			RefreshSecret:    "test_refresh_secret",
			AccessTTLMinutes: 15,
			RefreshTTLHours:  720,
			Issuer:           "test_issuer",
			Audience:         "test_audience",
		},
	}
}

func TestGeneratePair(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	email := "test@test.test"
	svc := NewService(testConfig())

	pair, err := svc.GeneratePair(ctx, userID, email)
	require.NoError(t, err)
	require.NotEmpty(t, pair)

	parsedAccess, err := jwt.ParseWithClaims(pair.AccessToken, &accessClaims{}, func(t *jwt.Token) (any, error) {
		// Ensure only HMAC is accepted
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.New("unexpected signing method")
		}
		return []byte(svc.cfg.Token.AccessSecret), nil
	})
	require.NoError(t, err)

	claims, ok := parsedAccess.Claims.(*accessClaims)
	require.True(t, ok)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, userID.String(), claims.Subject)
	assert.Equal(t, jwt.ClaimStrings{"test_audience"}, claims.Audience)
	assert.Equal(t, "test_issuer", claims.Issuer)

	_, err = svc.ValidateAccessToken(ctx, pair.RefreshToken)
	assert.Error(t, err, "refresh token should be rejected as access token")
}

func TestValidateAccessToken_Success(t *testing.T) {
	ctx := context.Background()
	svc := NewService(testConfig())
	userID := uuid.New()
	email := "test@test.test"

	pair, err := svc.GeneratePair(ctx, userID, email)
	require.NoError(t, err)

	claims, err := svc.ValidateAccessToken(ctx, pair.AccessToken)
	require.NoError(t, err)
	require.Equal(t, email, claims.Email)
	require.Equal(t, userID, claims.UserID)
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	ctx := context.Background()
	svc := NewService(testConfig())

	// 1. Not a token
	_, err := svc.ValidateAccessToken(ctx, "not a token")
	assert.Error(t, err, "should be rejected")

	// 2. Wrong secret
	claims := accessClaims{
		Email: "test@test.test",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    svc.cfg.Token.Issuer,
			Audience:  jwt.ClaimStrings{svc.cfg.Token.Audience},
		},
	}
	wrongToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	wrongSignedToken, err := wrongToken.SignedString([]byte("wrong secret"))
	assert.NoError(t, err)

	_, err = svc.ValidateAccessToken(ctx, wrongSignedToken)
	require.Error(t, err, "token signed with wrond secret should be rejected")

	// 3. Token with alg=none
	noneToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	noneSignedToken, err := noneToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	_, err = svc.ValidateAccessToken(ctx, noneSignedToken)
	require.Error(t, err, "token signed with wrong method should be rejected")
}

func TestValidateAccessToken_ExpiredToken(t *testing.T) {
	ctx := context.Background()
	svc := NewService(testConfig())

	expiredClaims := accessClaims{
		Email: "test@test.test",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Minute)),
			Issuer:    svc.cfg.Token.Issuer,
			Audience:  jwt.ClaimStrings{svc.cfg.Token.Audience},
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	signedExpiredToken, err := expiredToken.SignedString([]byte(svc.cfg.Token.AccessSecret))
	require.NoError(t, err)

	_, err = svc.ValidateAccessToken(ctx, signedExpiredToken)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, jwt.ErrTokenExpired), "expected ErrTokenExpired, got %v", err)
}
