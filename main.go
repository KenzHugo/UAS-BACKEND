package main

import (
	"log"
	"UASBE/app/repository"
	"UASBE/routes"
	"UASBE/app/service"
	"UASBE/config"
	"UASBE/database"
	"UASBE/utils"

	_ "UASBE/docs" // ← TAMBAHKAN INI (Import swagger docs)

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger" // ← TAMBAHKAN INI
)

// @title Sistem Pelaporan Prestasi Mahasiswa API
// @version 1.0
// @description API untuk sistem pelaporan prestasi mahasiswa dengan RBAC
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load config
	config.LoadEnv()

	// Initialize JWT
	utils.InitJWT()

	// Connect PostgreSQL database
	database.ConnectDatabase()
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatal("Failed to get database connection:", err)
	}

	// Connect MongoDB database
	database.ConnectMongoDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(sqlDB)
	roleRepo := repository.NewRoleRepository(sqlDB)
	permRepo := repository.NewPermissionRepository(sqlDB)
	studentRepo := repository.NewStudentRepository(sqlDB)
	lecturerRepo := repository.NewLecturerRepository(sqlDB)
	achievementRepo := repository.NewAchievementRepository(sqlDB, database.MongoDB)
	reportRepo := repository.NewReportRepository(sqlDB, database.MongoDB)

	// Initialize services
	authService := service.NewAuthService(userRepo, roleRepo, permRepo)
	userService := service.NewUserService(userRepo, roleRepo, permRepo, studentRepo, lecturerRepo)
	studentService := service.NewStudentService(studentRepo, lecturerRepo, achievementRepo, userRepo)
	lecturerService := service.NewLecturerService(lecturerRepo, studentRepo, achievementRepo, userRepo)
	achievementService := service.NewAchievementService(achievementRepo, studentRepo, lecturerRepo, userRepo)
	reportService := service.NewReportService(reportRepo, achievementRepo, studentRepo, lecturerRepo, userRepo) 

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"status": "error",
				"error":  err.Error(),
			})
		},
	})

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Serve static files untuk uploads
	app.Static("/uploads", "./uploads")

	// ← TAMBAHKAN INI: Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "API is running",
			"swagger": "/swagger/index.html", // ← Info Swagger
		})
	})

	// Register routes
	routes.AuthRoutes(app, authService)
	routes.UserRoutes(app, userService)
	routes.StudentRoutes(app, studentService)
	routes.LecturerRoutes(app, lecturerService)
	routes.AchievementRoutes(app, achievementService)
	routes.ReportRoutes(app, reportService) 

	// Start server
	port := config.AppConfig.Port
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on port %s", port)
	log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", port) // ← Log Swagger URL
	log.Fatal(app.Listen(":" + port))
}