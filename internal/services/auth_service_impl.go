package services

import (
	"context"
	"errors"
	"os"

	"backend-kaffa.ai/configs"
	"backend-kaffa.ai/internal/dto"
	"backend-kaffa.ai/internal/sqlc/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	UsersQueries *users.Queries
}

func NewAuthService(usersQueries *users.Queries) AuthService {
	return &AuthServiceImpl{
		UsersQueries: usersQueries,
	}
}

func (s *AuthServiceImpl) LoginUser(ctx context.Context, loginRequest *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.UsersQueries.GetUserByEmailOrUsername(ctx, loginRequest.Username)

	if err != nil {
		configs.Log.Error("failed to get user", zap.Error(err))
		return nil, errors.New("USER_NOT_FOUND")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		configs.Log.Error("invalid password", zap.Error(err))
		return nil, errors.New("INVALID_PASSWORD")
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"username": user.Username,
		"role_id":  user.RoleID,
	})
	keyData, _ := os.ReadFile("keys/private.pem")
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	ss, err := accessToken.SignedString(privKey)
	if err != nil {
		configs.Log.Error("error signing token", zap.Error(err))
		return nil, errors.New("TOKEN_GENERATION_FAILED")
	}

	loginResponse := &dto.LoginResponse{
		AccessToken: ss,
	}
	configs.Log.Info("user logged in successfully", zap.String("username", user.Username))
	return loginResponse, nil
}
func (s *AuthServiceImpl) RegisterUser(ctx context.Context, registerRequest *dto.RegisterRequest) (*dto.RegisterResponse, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		configs.Log.Error("error hashing password", zap.Error(err))
		return nil, errors.New("PASSWORD_HASHING_FAILED")
	}
	user, err := s.UsersQueries.CreateUser(ctx, users.CreateUserParams{
		ID:       ulid.Make().String(),
		Username: registerRequest.Username,
		Email:    registerRequest.Email,
		Password: string(hashedPassword),
		RoleID:   "admin",
	})
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == "23505" {
			configs.Log.Error("user already exists", zap.String("username", registerRequest.Username), zap.String("email", registerRequest.Email))
			return nil, errors.New("USER_ALREADY_EXISTS")
		}
	}
	configs.Log.Info("user registered successfully", zap.String("username", user.Username))
	return &dto.RegisterResponse{

		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt.Time.String(),
		UpdatedAt: user.UpdatedAt.Time.String(),
	}, nil
}
