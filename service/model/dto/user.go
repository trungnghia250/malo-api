package dto

type LoginRequest struct {
	Email    string `json:"email" query:"email"`
	Password string `json:"password" query:"password"`
}

type LoginResponse struct {
}

type UserInfo struct {
	UserID       string `json:"user_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
