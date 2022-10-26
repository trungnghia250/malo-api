package delivery

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/domain/user/usecase"
	"github.com/trungnghia250/malo-api/service/model/dto"
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
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	if reqToken == "" {
		return errors.New("token must required")
	}
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	// claims are actually a map[string]interface{}
	userInfo, err := u.userUseCase.UserProfile(ctx, claims["eoc"].(string))
	if err != nil {
		return err
	}

	return ctx.JSON(userInfo)
}

func (u *UserHandler) UserLogout(ctx *fiber.Ctx) error {
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	token, _ := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	if err := u.userUseCase.UserLogout(ctx, claims["eoc"].(string)); err != nil {
		return ctx.JSON(errors.New("logout failed"))
	}
	return ctx.JSON(nil)
}
