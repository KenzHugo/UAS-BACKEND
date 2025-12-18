package service_test

import (
	"UASBE/app/model"
	"UASBE/app/service"
	"UASBE/test/mocks"
	"UASBE/utils"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestLogin_Success(t *testing.T) {
	// 1. Setup
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	permRepo := new(mocks.MockPermissionRepository)
	authSvc := service.NewAuthService(userRepo, roleRepo, permRepo)
	utils.JwtKey = []byte("test_secret")

	app := fiber.New()
	app.Post("/login", authSvc.Login)

	// 2. Mock Data
	password := "SecurePass123!"
	hashed, _ := utils.HashPassword(password)
	mockUser := &model.User{
		ID:           "uuid-1",
		Username:     "mahasiswa123",
		PasswordHash: hashed,
		IsActive:     true,
		RoleID:       "role-1",
	}
	mockRole := &model.Role{ID: "role-1", Name: "Mahasiswa"}
	mockPerms := []string{"achievement:create", "achievement:read"}

	// 3. Expectations
	userRepo.On("FindByUsername", "mahasiswa123").Return(mockUser, nil)
	roleRepo.On("GetRoleByID", "role-1").Return(mockRole, nil)
	permRepo.On("GetPermissionsByRoleID", "role-1").Return(mockPerms, nil)

	// 4. Request
	loginReq := model.LoginRequest{Username: "mahasiswa123", Password: password}
	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// 5. Execution
	resp, _ := app.Test(req)

	// 6. Assertions
	assert.Equal(t, 200, resp.StatusCode)
	
	var apiResp model.APIResponse
	json.NewDecoder(resp.Body).Decode(&apiResp)
	assert.Equal(t, "success", apiResp.Status)
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	authSvc := service.NewAuthService(userRepo, nil, nil)
	app := fiber.New()
	app.Post("/login", authSvc.Login)

	hashed, _ := utils.HashPassword("correct_pass")
	mockUser := &model.User{Username: "user1", PasswordHash: hashed}
	userRepo.On("FindByUsername", "user1").Return(mockUser, nil)

	loginReq := model.LoginRequest{Username: "user1", Password: "wrong_password"}
	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode) // Unauthorized
}