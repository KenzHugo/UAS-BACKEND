package service_test

import (
	"UASBE/app/model"
	"UASBE/app/service"
	"UASBE/test/mocks"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetStatistics_Admin_Success(t *testing.T) {
	// 1. Setup Mocks
	reportRepo := new(mocks.MockReportRepository)
	achRepo := new(mocks.MockAchievementRepository)
	stuRepo := new(mocks.MockStudentRepository)
	lecRepo := new(mocks.MockLecturerRepository)
	userRepo := new(mocks.MockUserRepository)

	svc := service.NewReportService(reportRepo, achRepo, stuRepo, lecRepo, userRepo)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		// Mock Admin Login
		c.Locals("user", &model.JWTClaims{UserID: "admin-1", Role: "Admin"})
		return c.Next()
	})
	app.Get("/reports/statistics", svc.GetStatistics)

	// 2. Mock Data & Expectations (Scope Admin: studentID & advisorID = nil)
	reportRepo.On("GetTotalByType", (*string)(nil), (*string)(nil)).Return(map[string]int{"competition": 5}, nil)
	reportRepo.On("GetTotalByPeriod", (*string)(nil), (*string)(nil)).Return([]model.PeriodStats{{Period: "2025-01", Count: 5}}, nil)
	reportRepo.On("GetTopStudents", 10, (*string)(nil)).Return([]model.TopStudent{{StudentName: "John Doe"}}, nil)
	reportRepo.On("GetCompetitionLevelDistribution", (*string)(nil), (*string)(nil)).Return(map[string]int{"national": 3}, nil)
	reportRepo.On("GetStatusBreakdown", (*string)(nil), (*string)(nil)).Return(map[string]int{"verified": 5}, nil)

	// 3. Execution
	req := httptest.NewRequest("GET", "/reports/statistics", nil)
	resp, _ := app.Test(req)

	// 4. Assertions
	assert.Equal(t, 200, resp.StatusCode)
	
	var apiResp model.APIResponse
	json.NewDecoder(resp.Body).Decode(&apiResp)
	assert.Equal(t, "success", apiResp.Status)
}

func TestGetStudentReport_Forbidden_Mahasiswa(t *testing.T) {
	// Skenario: Mahasiswa A mencoba melihat report Mahasiswa B
	stuRepo := new(mocks.MockStudentRepository)
	svc := service.NewReportService(nil, nil, stuRepo, nil, nil)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{UserID: "user-A", Role: "Mahasiswa"})
		return c.Next()
	})
	app.Get("/reports/student/:id", svc.GetStudentReport)

	// Mahasiswa A (ID: student-A) mencoba akses student-B
	stuRepo.On("FindByID", "student-B").Return(&model.Student{ID: "student-B"}, nil)
	stuRepo.On("FindByUserID", "user-A").Return(&model.Student{ID: "student-A"}, nil)

	req := httptest.NewRequest("GET", "/reports/student/student-B", nil)
	resp, _ := app.Test(req)

	// Harus 403 Forbidden sesuai logic di report_service.go
	assert.Equal(t, 403, resp.StatusCode)
}