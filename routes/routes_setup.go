package routes

import (
	"database/sql"
	"UASBE/app/repository"
	"UASBE/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	// API v1 group
	api := app.Group("/api/v1")

	// Initialize repositories
	authRepo := repository.NewAuthRepository(db)

	// Initialize services
	authService := service.NewAuthService(authRepo)

	// Initialize handlers
	authHandler := NewAuthHandler(authService)

	// Setup routes
	authHandler.SetupRoutes(api)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Server is running",
		})
	})

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to Sistem Pelaporan Prestasi Mahasiswa API",
			"version": "1.0",
		})
	})
}