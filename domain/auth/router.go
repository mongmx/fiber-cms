package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mongmx/fiber-cms/middleware"
)

// Router for auth domain
func Router(app *fiber.App, u UseCase) {
	h := NewHandler(u)
	g := app.Group("/auth")
	{
		g.Get("/register", h.getRegister)
		g.Post("/register", h.postRegister)
		g.Get("/login", h.getLogin)
		g.Post("/login", h.postLogin)
		g.Get("/logout", h.getLogout)
		g.Get("/profile", mustLogin(), h.getProfile)
	}
}

func mustLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		u, ok := c.Locals("user").(*middleware.User)
		if !ok || u.ID <= 0 {
			return c.Redirect("/auth/login")
		}
		return c.Next()
	}
}
