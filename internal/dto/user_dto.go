package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=2,max=30"`
	Email           string `json:"email" binding:"required,min=2,max=50"`
	Password        string `json:"password" binding:"required,min=8,max=100"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterResponse struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	RoleID    string `json:"role_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
