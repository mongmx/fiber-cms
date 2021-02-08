package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	_ "github.com/mongmx/fiber-cms/docs"
	"github.com/mongmx/fiber-cms/middleware"
	"golang.org/x/sync/errgroup"

	"github.com/markbates/pkger"
)

// @title Fiber CMS
// @version 1.0
// @description CMS system made from Fiber
// @termsOfService http://fiber-cms.io/terms/
// @contact.name Support
// @contact.email support@fiber-cms.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

func main() {
	err := prepareEnvFile()
	if err != nil {
		log.Fatal("error prepare .env file")
	}
	err = godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	cfgDoc, err := strconv.ParseBool(os.Getenv("APP_DOC"))
	if err != nil {
		log.Fatal(err)
	}
	cfgMonitor, err := strconv.ParseBool(os.Getenv("APP_MONITOR"))
	if err != nil {
		log.Fatal(err)
	}
	fiberMain := mainApp()
	fiberDoc := docApp()
	fiberMetrics := metricsApp(fiberMain)
	var eg errgroup.Group
	eg.Go(func() error {
		return fiberMain.Listen(":8080")
	})
	if cfgDoc {
		eg.Go(func() error {
			return fiberDoc.Listen(":8081")
		})
	}
	if cfgMonitor {
		eg.Go(func() error {
			return fiberMetrics.Listen(":8082")
		})
	}
	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := fiberMain.Shutdown(); err != nil {
		log.Fatal(err)
	}
	if cfgDoc {
		if err := fiberDoc.Shutdown(); err != nil {
			log.Fatal(err)
		}
	}
	if cfgMonitor {
		if err := fiberMetrics.Shutdown(); err != nil {
			log.Fatal(err)
		}
	}
}

func mainApp() *fiber.App {
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
		return c.Render("pages/index", fiber.Map{
			"Title": "Hello, World!",
		}, "layouts/main")
	})
	app.Get("/post/list", func(c *fiber.Ctx) error {
		return c.Render("pages/post/index", fiber.Map{
			"Title": "Hello, World!",
		}, "layouts/main")
	})
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).Render("errors/404", fiber.Map{})
	})
	return app
}

func docApp() *fiber.App {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html", http.StatusMovedPermanently)
	})
	app.Get("/swagger/*", swagger.Handler)
	return app
}

func metricsApp(mainApp *fiber.App) *fiber.App {
	app := fiber.New()
	promMiddleware := middleware.NewPromMiddleware("fiber", "http")
	promMiddleware.Register(mainApp)
	promMiddleware.SetupPath(app)
	return app
}

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
