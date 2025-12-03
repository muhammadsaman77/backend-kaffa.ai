package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"backend-kaffa.ai/internal/dto"
	"backend-kaffa.ai/internal/sqlc/users"
	"github.com/golang-jwt/jwt/v5"
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
	fmt.Printf("Retrieved user: %+v\n", err)
	if err != nil {
		return nil, errors.New("USER_NOT_FOUND")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
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
		log.Println("Error signing token:", err)
		return nil, errors.New("TOKEN_GENERATION_FAILED")
	}

	loginResponse := &dto.LoginResponse{
		AccessToken: ss,
	}

	return loginResponse, nil
}
func (s *AuthServiceImpl) RegisterUser(ctx context.Context, registerRequest *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	return nil, nil // Implement registration logic here
}
