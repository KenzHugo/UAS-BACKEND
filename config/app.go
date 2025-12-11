package config

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Application struct {
	App    *fiber.App
	DB     *sql.DB
	Config *Config
}

func NewApplication() *Application {
	// Load environment variables
	LoadEnv()

	// Load config
	cfg := LoadConfig()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
		AppName:      "Sistem Pelaporan Prestasi Mahasiswa",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))

	return &Application{
		App:    app,
		Config: cfg,
	}
}

func (a *Application) SetDB(db *sql.DB) {
	a.DB = db
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": err.Error(),
	})
}