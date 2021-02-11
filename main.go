package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	_ "github.com/mongmx/fiber-cms/docs"
	"github.com/mongmx/fiber-cms/domain/auth"
	"github.com/mongmx/fiber-cms/domain/post"
	"github.com/pagongamedev/godd"

	"github.com/markbates/pkger"
)

// @title Fiber CMS
// @version 1.0
// @description CMS system made from Fiber
// @termsOfService https://fiber-cms.io/terms/
// @contact.name Support
// @contact.email support@fiber-cms.io
// @license.name MIT License
// @license.url https://github.com/mongmx/fiber-cms/blob/main/LICENSE
// @host localhost:8080
// @BasePath /

func main() {
	// Set Environment
	err := prepareEnvFile()
	godd.MustError(err)

	err = godotenv.Load()
	godd.MustError(err, "error loading .env file")

	cfgDoc, err := strconv.ParseBool(os.Getenv("APP_DOC"))
	godd.MustError(err)

	cfgMonitor, err := strconv.ParseBool(os.Getenv("APP_MONITOR"))
	godd.MustError(err)

	// Set Portal
	portal := godd.NewPortal()
	appMain := appMain()

	portal.AppendApp(appMain, ":8180")
	if cfgDoc {
		portal.AppendApp(godd.AppAPIDocument(), ":8181")
	}
	if cfgMonitor {
		portal.AppendApp(godd.AppMetricsPrometheus(appMain), ":8182")
	}

	portal.StartServer()
}

func appMain() *fiber.App {
	engine := html.NewFileSystem(pkger.Dir("/views"), ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root: pkger.Dir("/views/assets"),
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("pages/auth/login", fiber.Map{})
	})
	auth.Router(app)
	post.Router(app)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).Render("errors/404", fiber.Map{})
	})
	return app
}

// func docApp() *fiber.App {
// 	app := fiber.New()
// 	app.Get("/", func(c *fiber.Ctx) error {
// 		return c.Redirect("/swagger/index.html", http.StatusMovedPermanently)
// 	})
// 	app.Get("/swagger/*", swagger.Handler)
// 	return app
// }

// func metricsApp(mainApp *fiber.App) *fiber.App {
// 	app := fiber.New()
// 	promMiddleware := middleware.NewPromMiddleware("fiber", "http")
// 	promMiddleware.Register(mainApp)
// 	promMiddleware.SetupPath(app)
// 	return app
// }

func prepareEnvFile() error {
	envFile, err := os.OpenFile(".env", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = envFile.WriteString("APP_MODE=localhost\nAPP_DOC=false\nAPP_MONITOR=false")
	if err != nil {
		return err
	}
	defer func() {
		err = envFile.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	return nil
}
