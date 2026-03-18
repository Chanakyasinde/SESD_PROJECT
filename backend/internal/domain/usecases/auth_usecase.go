package usecases

import "context"

type RegisterInput struct {
	Name     string
	Email    string
	Password string
	Role     string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthOutput struct {
	Token string `json:"token"`
	User  any    `json:"user"`
}

type AuthUsecase interface {
	Register(ctx context.Context, input RegisterInput) (*AuthOutput, error)
	Login(ctx context.Context, input LoginInput) (*AuthOutput, error)
}
