package post

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mongmx/fiber-cms/middleware"
)

// Router for post domain
func Router(app *fiber.App) {
	g := app.Group("/post")

	g.Get("/list", middleware.MustLogin(), func(c *fiber.Ctx) error {
		return c.Render("pages/post/index", fiber.Map{
			"Title": "Show post list page",
		}, "layouts/main")
	})
}
