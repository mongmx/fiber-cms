package post

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mongmx/fiber-cms/middleware"
)

// Router for post domain
func Router(app *fiber.App) {
	g := app.Group("/post")

	g.Get("/list", mustLogin(), func(c *fiber.Ctx) error {
		return c.Render("pages/post/index", fiber.Map{
			"Title": "Show post list page",
		}, "layouts/main")
	})
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
