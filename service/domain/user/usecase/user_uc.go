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
	"time"
)

type IUserUseCase interface {
	UserLogin(ctx *fiber.Ctx, req dto.LoginRequest) (dto.UserInfo, error)
	UserProfile(ctx *fiber.Ctx, email string) (dto.UserInfo, error)
	UserLogout(ctx *fiber.Ctx, email string) error
	ListUser(ctx *fiber.Ctx, req dto.ListUserRequest) ([]model.User, error)
	CreateUser(ctx *fiber.Ctx, data *model.User) (*model.User, error)
	UpdateUser(ctx *fiber.Ctx, req model.User) error
	DeleteUsers(ctx *fiber.Ctx, userIDs []string) error
	UserDetail(ctx *fiber.Ctx, id string) (dto.UserInfo, error)
}

type userUseCase struct {
	repo repo.IRepo
}

func NewUserUseCase(repo repo.IRepo) IUserUseCase {
	return &userUseCase{
		repo: repo,
	}
}

func (u *userUseCase) UserLogin(ctx *fiber.Ctx, req dto.LoginRequest) (dto.UserInfo, error) {
	user, err := u.repo.NewUserRepo().GetUserByEmail(ctx, req.Email)
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

	err = u.repo.NewUserRepo().UpdateUser(ctx, user)
	if err != nil {
		return dto.UserInfo{}, errors.New("update token, refresh token failed")
	}
	userInfo := transform.UserToUserInfo(user)

	return userInfo, nil
}

func (u *userUseCase) UserProfile(ctx *fiber.Ctx, email string) (dto.UserInfo, error) {
	user, err := u.repo.NewUserRepo().GetUserByEmail(ctx, email)
	if err != nil {
		return dto.UserInfo{}, errors.New(fmt.Sprintf("User with email %s does not exists", email))
	}

	userInfo := transform.UserToUserInfo(user)
	return userInfo, nil
}

func (u *userUseCase) UserLogout(ctx *fiber.Ctx, email string) error {
	user, err := u.repo.NewUserRepo().GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	err = u.repo.NewUserRepo().RemoveUserData(ctx, user.UserID, bson.M{"token": "", "refresh_token": ""})
	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) ListUser(ctx *fiber.Ctx, req dto.ListUserRequest) ([]model.User, error) {
	users, err := u.repo.NewUserRepo().ListUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userUseCase) CreateUser(ctx *fiber.Ctx, data *model.User) (*model.User, error) {
	data.CreatedAt = time.Now()
	userID, err := u.repo.NewCounterRepo().GetSequenceNextValue(ctx, "user_id")
	if err != nil {
		return nil, err
	}
	data.UserID = fmt.Sprintf("U%d", userID)
	err = u.repo.NewUserRepo().CreateUser(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *userUseCase) UpdateUser(ctx *fiber.Ctx, req model.User) error {
	err := u.repo.NewUserRepo().UpdateUser(ctx, &req)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) DeleteUsers(ctx *fiber.Ctx, userIDs []string) error {
	err := u.repo.NewUserRepo().DeleteUsers(ctx, userIDs)

	return err
}

func (u *userUseCase) UserDetail(ctx *fiber.Ctx, id string) (dto.UserInfo, error) {
	user, err := u.repo.NewUserRepo().GetUserByID(ctx, id)
	if err != nil {
		return dto.UserInfo{}, errors.New(fmt.Sprintf("User with id %s does not exists", id))
	}

	userInfo := transform.UserToUserInfo(user)
	return userInfo, nil
}
