package transform

import (
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
)

func UserToUserInfo(user *model.User) dto.UserInfo {
	return dto.UserInfo{
		UserID:       user.UserID,
		Name:         user.Name,
		Email:        user.Email,
		Role:         user.Role,
		RefreshToken: user.RefreshToken,
		Token:        user.Token,
	}
}
