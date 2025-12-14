package routes

import (
	"UASBE/middleware"
	"UASBE/app/service"

	"github.com/gofiber/fiber/v2"
)

//
// ==================== AUTH ROUTES ======================
//

func AuthRoutes(app *fiber.App, authService *service.AuthService) {
	auth := app.Group("/api/v1/auth")

	auth.Post("/login", authService.Login)
	auth.Post("/refresh", authService.Refresh)

	// Protected routes
	protected := auth.Group("/", middleware.AuthRequired)
	protected.Get("/profile", authService.Profile)
	protected.Post("/logout", authService.Logout)
}

//
// ==================== USER ROUTES (ADMIN ONLY) ======================
//

func UserRoutes(app *fiber.App, userService *service.UserService) {
	users := app.Group("/api/v1/users")

	// Semua endpoint user butuh auth + permission "user:manage"
	users.Use(middleware.AuthRequired)
	users.Use(middleware.RequirePermission("user:manage"))

	users.Get("/", userService.GetUsers)          // GET /api/v1/users
	users.Get("/:id", userService.GetUserByID)    // GET /api/v1/users/:id
	users.Post("/", userService.CreateUser)       // POST /api/v1/users
	users.Put("/:id", userService.UpdateUser)     // PUT /api/v1/users/:id
	users.Delete("/:id", userService.DeleteUser)  // DELETE /api/v1/users/:id
	users.Put("/:id/role", userService.AssignRole) // PUT /api/v1/users/:id/role
}

//
// ==================== STUDENT ROUTES ======================
//

func StudentRoutes(app *fiber.App, studentService *service.StudentService) {
	students := app.Group("/api/v1/students")

	// Auth required untuk semua endpoint
	students.Use(middleware.AuthRequired)

	students.Get("/", studentService.GetAllStudents)           // GET /api/v1/students
	students.Get("/:id", studentService.GetStudentByID)        // GET /api/v1/students/:id
	students.Put("/:id/advisor", studentService.SetAdvisor)    // PUT /api/v1/students/:id/advisor
	
	// TODO: GET /api/v1/students/:id/achievements (nanti saat implement achievements)
}

//
// ==================== ACHIEVEMENT ROUTES ======================
//

func AchievementRoutes(app *fiber.App, achievementService *service.AchievementService) {
	achievements := app.Group("/api/v1/achievements")
	
	// All achievement routes require authentication
	achievements.Use(middleware.AuthRequired)
	
	// GET /achievements - List achievements (filtered by role)
	// Mahasiswa: own achievements
	// Dosen Wali: advisees' achievements (FR-006)
	// Admin: all achievements (FR-010)
	achievements.Get("/", achievementService.GetAchievements)
	
	// GET /achievements/:id - Get achievement detail
	achievements.Get("/:id", achievementService.GetAchievementByID)
	
	// POST /achievements - Create achievement (Mahasiswa only) (FR-003)
	achievements.Post("/",
		middleware.RequirePermission("achievement:create"),
		achievementService.CreateAchievement,
	)
	
	// PUT /achievements/:id - Update achievement (Mahasiswa only, draft only)
	achievements.Put("/:id",
		middleware.RequirePermission("achievement:update"),
		achievementService.UpdateAchievement,
	)
	
	// DELETE /achievements/:id - Delete achievement (Mahasiswa only, draft only) (FR-005)
	achievements.Delete("/:id",
		middleware.RequirePermission("achievement:delete"),
		achievementService.DeleteAchievement,
	)
	
	// POST /achievements/:id/submit - Submit for verification (Mahasiswa) (FR-004)
	achievements.Post("/:id/submit",
		middleware.RequirePermission("achievement:create"),
		achievementService.SubmitForVerification,
	)
	
	// POST /achievements/:id/verify - Verify achievement (Dosen Wali) (FR-007)
	achievements.Post("/:id/verify",
		middleware.RequirePermission("achievement:verify"),
		achievementService.VerifyAchievement,
	)
	
	// POST /achievements/:id/reject - Reject achievement (Dosen Wali) (FR-008)
	achievements.Post("/:id/reject",
		middleware.RequirePermission("achievement:verify"),
		achievementService.RejectAchievement,
	)
	
	// POST /achievements/:id/attachments - Upload attachment
	achievements.Post("/:id/attachments",
		middleware.RequirePermission("achievement:create"),
		achievementService.UploadAttachment,
	)
}
