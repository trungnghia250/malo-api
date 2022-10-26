package usecase

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/model/transform"
	"github.com/trungnghia250/malo-api/service/repo"
	"github.com/trungnghia250/malo-api/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type IUserUseCase interface {
	UserLogin(ctx *fiber.Ctx, req dto.LoginRequest) (dto.UserInfo, error)
	UserProfile(ctx *fiber.Ctx, email string) (dto.UserInfo, error)
	UserLogout(ctx *fiber.Ctx, email string) error
}

type userUseCase struct {
	repo repo.IRepo
}

func NewUserUseCase(repo repo.IRepo) IUserUseCase {
	return &userUseCase{
		repo: repo,
	}
}

func (c *userUseCase) UserLogin(ctx *fiber.Ctx, req dto.LoginRequest) (dto.UserInfo, error) {
	user, err := c.repo.NewUserRepo().GetUserByEmail(ctx, req.Email)
	if err != nil {
		return dto.UserInfo{}, errors.New(fmt.Sprintf("User with email %s does not exists", req.Email))
	}

	if req.Password != user.Password {
		return dto.UserInfo{}, errors.New(fmt.Sprintf("Password incorrect"))
	}

	tokenContent := model.TokenContent{
		AccountType: "user",
		Email:       user.Email,
		Name:        user.Name,
		Role:        user.Role,
	}

	token, err := utils.GenToken(tokenContent)
	refreshToken := utils.GenRefreshToken()
	user.Token = token
	user.RefreshToken = refreshToken

	err = c.repo.NewUserRepo().UpdateUser(ctx, user)
	if err != nil {
		return dto.UserInfo{}, errors.New("update token, refresh token failed")
	}
	userInfo := transform.UserToUserInfo(user)

	return userInfo, nil
}

func (c *userUseCase) UserProfile(ctx *fiber.Ctx, email string) (dto.UserInfo, error) {
	user, err := c.repo.NewUserRepo().GetUserByEmail(ctx, email)
	if err != nil {
		return dto.UserInfo{}, errors.New(fmt.Sprintf("User with email %s does not exists", email))
	}

	userInfo := transform.UserToUserInfo(user)
	return userInfo, nil
}

func (c *userUseCase) UserLogout(ctx *fiber.Ctx, email string) error {
	user, err := c.repo.NewUserRepo().GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	err = c.repo.NewUserRepo().RemoveUserData(ctx, user.UserID, bson.M{"token": "", "refresh_token": ""})
	if err != nil {
		return err
	}

	return nil
}
