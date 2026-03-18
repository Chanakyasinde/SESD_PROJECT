package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"inventory-management-system/backend/internal/domain/entities"
	"inventory-management-system/backend/internal/domain/repositories"
	"inventory-management-system/backend/internal/domain/usecases"
	"inventory-management-system/backend/internal/infrastructure/security"
)

type AuthUsecase struct {
	userRepo repositories.UserRepository
	jwtSvc   *security.JWTService
}

func NewAuthUsecase(userRepo repositories.UserRepository, jwtSvc *security.JWTService) usecases.AuthUsecase {
	return &AuthUsecase{userRepo: userRepo, jwtSvc: jwtSvc}
}

func (u *AuthUsecase) Register(ctx context.Context, input usecases.RegisterInput) (*usecases.AuthOutput, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" || input.Email == "" || len(input.Password) < 6 {
		return nil, entities.ErrInvalidInput
	}

	role := strings.ToLower(strings.TrimSpace(input.Role))
	if role == "" {
		role = entities.RoleCustomer
	}
	if role != entities.RoleAdmin && role != entities.RoleStaff && role != entities.RoleCustomer {
		return nil, entities.ErrInvalidInput
	}

	_, err := u.userRepo.GetByEmail(ctx, input.Email)
	if err == nil {
		return nil, fmt.Errorf("%w: email already registered", entities.ErrConflict)
	}
	if !errors.Is(err, entities.ErrNotFound) {
		return nil, err
	}

	hash, err := security.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &entities.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hash,
		Role:         role,
		TimeStamp: entities.TimeStamp{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	if err = u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := u.jwtSvc.GenerateToken(user.ID.Hex(), user.Role)
	if err != nil {
		return nil, err
	}

	return &usecases.AuthOutput{Token: token, User: user}, nil
}

func (u *AuthUsecase) Login(ctx context.Context, input usecases.LoginInput) (*usecases.AuthOutput, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if input.Email == "" || input.Password == "" {
		return nil, entities.ErrInvalidInput
	}

	user, err := u.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, entities.ErrUnauthorized
	}
	if security.ComparePassword(user.PasswordHash, input.Password) != nil {
		return nil, entities.ErrUnauthorized
	}

	token, err := u.jwtSvc.GenerateToken(user.ID.Hex(), user.Role)
	if err != nil {
		return nil, err
	}

	return &usecases.AuthOutput{Token: token, User: user}, nil
}
