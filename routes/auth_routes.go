package routes

import (
	"UASBE/app/service"
	"UASBE/middleware"
	"UASBE/app/model"
	"UASBE/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) SetupRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	// Public routes
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.RefreshToken)

	// Protected routes
	auth.Get("/profile", middleware.AuthRequired, h.GetProfile)
	auth.Post("/logout", middleware.AuthRequired, h.Logout)
}

// Login handler
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validate request
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Username and password are required",
		})
	}

	response, err := h.authService.Login(req)
	if err != nil {
		statusCode := fiber.StatusUnauthorized
		if err.Error() == "user account is inactive" {
			statusCode = fiber.StatusForbidden
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

// RefreshToken handler
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req model.RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Refresh token is required",
		})
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

// GetProfile handler
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	// Get claims from context (set by middleware)
	claims := c.Locals("user").(*utils.Claims)

	profile, err := h.authService.GetProfile(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get profile",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   profile,
	})
}

// Logout handler
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// In stateless JWT, logout is handled client-side
	// by deleting the tokens
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Logged out successfully",
	})
}