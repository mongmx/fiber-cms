package post

import "github.com/gofiber/fiber/v2"

// Router for post domain
func Router(app *fiber.App)  {
	g := app.Group("/post")

	g.Get("/list", func(c *fiber.Ctx) error {
		return c.Render("pages/post/index", fiber.Map{
			"Title": "Show post list page",
		}, "layouts/main")
	})
}