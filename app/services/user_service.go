package services

import (
	"context"
	"errors"

	"vuelang/app/models"
	"vuelang/app/repositories"
	"vuelang/internal/framework/hash"
)

var ErrEmailTaken = errors.New("email is already in use")

// UserService handles user management business logic.
type UserService struct {
	userRepo *repositories.UserRepository
	hasher   hash.Hasher
}

func NewUserService(userRepo *repositories.UserRepository, hasher hash.Hasher) *UserService {
	return &UserService{userRepo: userRepo, hasher: hasher}
}

func (s *UserService) All(ctx context.Context) ([]models.User, error) {
	return s.userRepo.All(ctx)
}

func (s *UserService) FindByID(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, name, email, password string) (*models.User, error) {
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}
	hashed, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}
	return s.userRepo.Create(ctx, name, email, hashed)
}

func (s *UserService) Update(ctx context.Context, id uint, name, email string, isActive bool) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil || user == nil {
		return nil, err
	}
	// Check email uniqueness if changed
	if user.Email != email {
		existing, _ := s.userRepo.FindByEmail(ctx, email)
		if existing != nil && existing.ID != id {
			return nil, ErrEmailTaken
		}
	}
	return s.userRepo.Update(ctx, id, name, email, isActive)
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}
