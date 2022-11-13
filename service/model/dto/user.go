package dto

import "github.com/trungnghia250/malo-api/service/model"

type LoginRequest struct {
	Email    string `json:"email" query:"email"`
	Password string `json:"password" query:"password"`
}

type LoginResponse struct {
}

type UserInfo struct {
	UserID       string `json:"user_id,omitempty"`
	Name         string `json:"name,omitempty"`
	Email        string `json:"email,omitempty"`
	Role         string `json:"role,omitempty"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type ListUserRequest struct {
	Limit  int32 `json:"limit,omitempty"`
	Offset int32 `json:"offset,omitempty"`
}

type ListUserResponse struct {
	Count int32        `json:"count"`
	Data  []model.User `json:"data"`
}

type DeleteUsersRequest struct {
	UserIDs []string `json:"user_ids" query:"user_ids"`
}
