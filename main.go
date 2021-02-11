package main

import (
	"crypto/sha256"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/gofiber/template/html"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/markbates/pkger"
	_ "github.com/mongmx/fiber-cms/docs"
	"github.com/mongmx/fiber-cms/domain/auth"
	"github.com/mongmx/fiber-cms/domain/post"
	"github.com/mongmx/fiber-cms/middleware"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
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
		Prefork: false,
		Views:   engine,
	})
	app.Get("/dashboard", monitor.New())
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root: pkger.Dir("/views/assets"),
	}))

	//app.Use(logger.New())
	//app.Use(recover.New())

	postgresDB := initPostgres()
	redisStorage := initRedis()
	store := session.New(session.Config{
		Expiration: 10 * time.Minute,
		Storage:    redisStorage,
		CookieName: "ssid",
		KeyGenerator: func() string {
			id := uuid.NewV4()
			hash := sha256.Sum256(id.Bytes())
			return fmt.Sprintf("%x", hash)
		},
	})
	store.RegisterType("")
	app.Use(middleware.Session(store))
	app.Use(middleware.Auth(redisStorage))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("pages/auth/login", fiber.Map{})
	})

	authRepo := auth.NewRepository(postgresDB, redisStorage)
	authUseCase := auth.NewUseCase(authRepo)
	auth.Router(app, authUseCase)
	post.Router(app)
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

func initPostgres() *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s sslmode=%s",
		"127.0.0.1", "5432", "fiber_cms", "mongmx", "disable",
	)
	postgresDB, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return postgresDB
}

func initRedis() *redis.Storage {
	return redis.New(redis.Config{
		Host:     "127.0.0.1",
		Port:     6379,
		Username: "",
		Password: "",
		Database: 0,
		Reset:    false,
	})
}
