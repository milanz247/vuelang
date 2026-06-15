package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims are embedded in the access token.
type Claims struct {
	UserID uint     `json:"uid"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// TokenPair holds both an access token and a refresh token.
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// Service wraps JWT operations.
type Service interface {
	GeneratePair(userID uint, email string, roles []string) (*TokenPair, error)
	ValidateAccess(tokenString string) (*Claims, error)
	GenerateRefreshToken() (string, error)
}

type service struct {
	secret        []byte
	accessTTL     time.Duration
	refreshTTLDay int
}

func NewService(secret string, accessTTL time.Duration, refreshTTLDay int) Service {
	return &service{
		secret:        []byte(secret),
		accessTTL:     accessTTL,
		refreshTTLDay: refreshTTLDay,
	}
}

func (s *service) GeneratePair(userID uint, email string, roles []string) (*TokenPair, error) {
	expiresAt := time.Now().UTC().Add(s.accessTTL)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    "vuelang",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refreshToken, err := s.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  signed,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *service) ValidateAccess(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}

// GenerateRefreshToken returns a cryptographically random UUID string.
func (s *service) GenerateRefreshToken() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	return id.String(), nil
}

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")
)
