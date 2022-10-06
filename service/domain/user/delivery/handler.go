package delivery

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/domain/user/usecase"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"strings"
)

type UserHandler struct {
	userUseCase usecase.IUserUseCase
}

func NewUserHandler(userUseCase usecase.IUserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (u *UserHandler) UserLogin(ctx *fiber.Ctx) error {
	req := new(dto.LoginRequest)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	status, err := u.userUseCase.UserLogin(ctx, *req)
	if err != nil {
		return err
	}

	return ctx.JSON(status)
}

func (u *UserHandler) UserProfile(ctx *fiber.Ctx) error {
	user := ctx.Locals("x-access-token").(*jwt.Token)
	claims := user.Claims.(*model.JWTProfileClaims)
	userInfo, err := u.userUseCase.UserProfile(ctx, claims.Email)
	if err != nil {
		return err
	}

	return ctx.JSON(userInfo)
}

func (u *UserHandler) UserLogout(ctx *fiber.Ctx) error {
	user := ctx.Locals("x-access-token").(*jwt.Token)
	claims := user.Claims.(*model.JWTProfileClaims)
	rawToken := ctx.GetRespHeader("Authorization")
	rawToken = strings.TrimSpace(strings.Replace(rawToken, "Bearer", "", -1))

	if err := u.userUseCase.UserLogout(ctx, claims.Email); err != nil {
		return ctx.JSON(errors.New("logout failed"))
	}
	return ctx.JSON(nil)
}
