package delivery

import "github.com/gofiber/fiber/v2"

func (u *UserHandler) InternalUserAPIRoute(router fiber.Router) {
	router.Post("/login", u.UserLogin)
	router.Get("/me", u.UserProfile)
	router.Post("/logout", u.UserLogout)
}
