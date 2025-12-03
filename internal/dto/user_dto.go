package dto

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=2,max=30"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Email    string `json:"email" validate:"required,min=2,max=50"`
	RoleID   int    `json:"role_id" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterResponse struct {
	AccessToken string `json:"access_token"`
	UserID      int    `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	RoleID      int    `json:"role_id"`
}
