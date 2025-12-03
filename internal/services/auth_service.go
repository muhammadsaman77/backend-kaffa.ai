package services

import (
	"context"

	"backend-kaffa.ai/internal/dto"
)

type AuthService interface {
	LoginUser(ctx context.Context, loginRequest *dto.LoginRequest) (*dto.LoginResponse, error)
	RegisterUser(ctx context.Context, registerRequest *dto.RegisterRequest) (*dto.RegisterResponse, error)
}
