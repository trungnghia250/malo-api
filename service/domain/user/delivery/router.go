package delivery

import "github.com/gofiber/fiber/v2"

func (u *UserHandler) InternalUserAPIRoute(router fiber.Router) {
	router.Post("/login", u.UserLogin)
	router.Get("/me", u.UserProfile)
	router.Post("/logout", u.UserLogout)

	router.Get("/user/list", u.ListUser)
	router.Post("/user", u.CreateUser)
	router.Put("/user", u.UpdateUser)
	router.Delete("/user", u.DeleteUsers)
	router.Get("/user/detail", u.GetUserDetail)
}
