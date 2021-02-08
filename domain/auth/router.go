package auth

import "github.com/gofiber/fiber/v2"

// Router for auth domain
func Router(app *fiber.App) {
	g := app.Group("/auth")

	g.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("pages/post/index", fiber.Map{
			"Title": "Show login page",
		}, "layouts/main")
	})
}