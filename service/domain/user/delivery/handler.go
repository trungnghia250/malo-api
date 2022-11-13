package delivery

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/trungnghia250/malo-api/service/domain/user/usecase"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"time"
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

func (u *UserHandler) ListUser(ctx *fiber.Ctx) error {
	req := new(dto.ListUserRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	users, err := u.userUseCase.ListUser(ctx, *req)
	if err != nil {
		return err
	}

	count := int32(0)
	if len(users) > 0 {
		count = users[0].TotalCount
	}

	response := dto.ListUserResponse{
		Count: count,
		Data:  users,
	}

	return ctx.JSON(response)
}

func (u *UserHandler) CreateUser(ctx *fiber.Ctx) error {
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	if reqToken == "" {
		return errors.New("token is required")
	}
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)

	req := new(model.User)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	req.ModifiedBy = claims["noc"].(string)
	req.ModifiedAt = time.Now()
	req.CreatedAt = time.Now()
	customer, err := u.userUseCase.CreateUser(ctx, req)
	if err != nil {
		return err
	}

	return ctx.JSON(customer)
}

func (u *UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	reqToken := ctx.GetReqHeaders()["X-Access-Token"]
	if reqToken == "" {
		return errors.New("token is required")
	}
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return errors.New("token not valid")
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	req := new(model.User)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	req.ModifiedBy = claims["noc"].(string)
	err = u.userUseCase.UpdateUser(ctx, *req)
	if err != nil {
		return err
	}

	return ctx.JSON(nil)
}

func (u *UserHandler) DeleteUsers(ctx *fiber.Ctx) error {
	req := new(dto.DeleteUsersRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	err := u.userUseCase.DeleteUsers(ctx, req.UserIDs)
	if err != nil {
		return err
	}

	return ctx.JSON(dto.DefaultResponse{
		StatusCode: fiber.StatusOK,
	})
}

func (u *UserHandler) GetUserDetail(ctx *fiber.Ctx) error {
	req := new(dto.GetUserDetailRequest)
	if err := ctx.QueryParser(req); err != nil {
		return err
	}

	user, err := u.userUseCase.UserDetail(ctx, req.UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(user)
}