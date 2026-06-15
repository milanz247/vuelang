package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"vuelang/app/models"
	"vuelang/app/repositories"
	"vuelang/config"
	"vuelang/internal/framework/hash"
	jwtpkg "vuelang/internal/framework/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already in use")
	ErrTokenInvalid       = errors.New("token is invalid or expired")
	ErrUserNotFound       = errors.New("user not found")
	ErrAccountInactive    = errors.New("account is inactive")
)

// AuthService handles all authentication business logic.
type AuthService struct {
	userRepo *repositories.UserRepository
	roleRepo *repositories.RoleRepository
	hasher   hash.Hasher
	jwt      jwtpkg.Service
	cfg      *config.App
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	roleRepo *repositories.RoleRepository,
	hasher hash.Hasher,
	jwt jwtpkg.Service,
	cfg *config.App,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		hasher:   hasher,
		jwt:      jwt,
		cfg:      cfg,
	}
}

// RegisterResult contains the newly created user and their initial tokens.
type RegisterResult struct {
	User   *models.User
	Tokens *jwtpkg.TokenPair
}

// Register creates a new user, assigns the default role, and returns tokens.
func (s *AuthService) Register(ctx context.Context, name, email, password string) (*RegisterResult, error) {
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailExists
	}

	hashed, err := s.hasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.userRepo.Create(ctx, name, email, hashed)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := s.roleRepo.AssignRole(ctx, user.ID, models.RoleUser); err != nil {
		// Non-fatal in dev; log and continue
		_ = err
	}

	roles, _ := s.roleRepo.GetUserRoles(ctx, user.ID)
	tokens, err := s.jwt.GeneratePair(user.ID, user.Email, roles)
	if err != nil {
		return nil, err
	}

	if err := s.storeRefreshToken(ctx, user.ID, tokens.RefreshToken); err != nil {
		return nil, err
	}

	user.Roles = roles
	return &RegisterResult{User: user, Tokens: tokens}, nil
}

// Login validates credentials and returns tokens.
func (s *AuthService) Login(ctx context.Context, email, password string) (*jwtpkg.TokenPair, *models.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil || !s.hasher.Verify(password, user.Password) {
		return nil, nil, ErrInvalidCredentials
	}
	if !user.IsActive {
		return nil, nil, ErrAccountInactive
	}

	roles, _ := s.roleRepo.GetUserRoles(ctx, user.ID)
	tokens, err := s.jwt.GeneratePair(user.ID, user.Email, roles)
	if err != nil {
		return nil, nil, err
	}

	if err := s.storeRefreshToken(ctx, user.ID, tokens.RefreshToken); err != nil {
		return nil, nil, err
	}

	user.Roles = roles
	return tokens, user, nil
}

// Refresh validates a refresh token and issues a new token pair.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*jwtpkg.TokenPair, error) {
	userID, valid, err := s.userRepo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if !valid || userID == 0 {
		return nil, ErrTokenInvalid
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	if err := s.userRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	roles, _ := s.roleRepo.GetUserRoles(ctx, user.ID)
	tokens, err := s.jwt.GeneratePair(user.ID, user.Email, roles)
	if err != nil {
		return nil, err
	}

	if err := s.storeRefreshToken(ctx, user.ID, tokens.RefreshToken); err != nil {
		return nil, err
	}

	return tokens, nil
}

// Logout invalidates a specific refresh token.
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.userRepo.DeleteRefreshToken(ctx, refreshToken)
}

// ForgotPassword creates a password reset token and (in production) sends an email.
func (s *AuthService) ForgotPassword(ctx context.Context, email string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		// Return success anyway to prevent email enumeration
		return "", nil
	}

	token, err := s.jwt.GenerateRefreshToken() // reuse UUID generator
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().UTC().Add(1 * time.Hour)
	if err := s.userRepo.CreatePasswordReset(ctx, email, token, expiresAt); err != nil {
		return "", err
	}

	// TODO: send email via mail service
	// e.g. mailSvc.Send(email, "Password Reset", resetLink(token))

	return token, nil
}

// ResetPassword validates the token and sets a new password.
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	email, valid, err := s.userRepo.FindPasswordReset(ctx, token)
	if err != nil {
		return err
	}
	if !valid || email == "" {
		return ErrTokenInvalid
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	hashed, err := s.hasher.Hash(newPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(ctx, user.ID, hashed); err != nil {
		return err
	}

	// Invalidate all refresh tokens and the used reset token
	_ = s.userRepo.DeletePasswordReset(ctx, email)
	_ = s.userRepo.DeleteUserRefreshTokens(ctx, user.ID)

	return nil
}

// Me returns the authenticated user with their roles.
func (s *AuthService) Me(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	roles, _ := s.roleRepo.GetUserRoles(ctx, user.ID)
	user.Roles = roles
	return user, nil
}

func (s *AuthService) storeRefreshToken(ctx context.Context, userID uint, token string) error {
	expiresAt := time.Now().UTC().AddDate(0, 0, s.cfg.JWTRefreshTTLDay)
	return s.userRepo.CreateRefreshToken(ctx, userID, token, expiresAt)
}
