package service_test

import (
	"UASBE/app/model"
	"UASBE/app/service"
	"UASBE/test/mocks"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAchievement_Success(t *testing.T) {
	// Setup
	achRepo := new(mocks.MockAchievementRepository)
	stuRepo := new(mocks.MockStudentRepository)
	svc := service.NewAchievementService(achRepo, stuRepo, nil, nil)

	app := fiber.New()
	app.Post("/achievements", func(c *fiber.Ctx) error {
		// Mock login as Student
		c.Locals("user", &model.JWTClaims{UserID: "user-123", Role: "Mahasiswa"})
		return svc.CreateAchievement(c)
	})

	// Mock Expectations
	mockStudent := &model.Student{ID: "student-123"}
	stuRepo.On("FindByUserID", "user-123").Return(mockStudent, nil)
	achRepo.On("CreateAchievement", mock.AnythingOfType("*model.Achievement")).Return("mongo-id-1", nil)
	achRepo.On("CreateReference", mock.Anything).Return(nil)

	// Execute
	reqBody := model.AchievementCreateRequest{
		AchievementType: "competition",
		Title:           "Juara 1 Nasional",
		Description:     "Lomba Coding",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/achievements", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, 201, resp.StatusCode)
}

func TestSubmitForVerification_Success(t *testing.T) {
	achRepo := new(mocks.MockAchievementRepository)
	stuRepo := new(mocks.MockStudentRepository)
	svc := service.NewAchievementService(achRepo, stuRepo, nil, nil)

	app := fiber.New()
	app.Post("/achievements/:id/submit", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{UserID: "user-123", Role: "Mahasiswa"})
		return svc.SubmitForVerification(c)
	})

	// Mock Data: Status must be 'draft' to submit
	mockRef := &model.AchievementReference{ID: "ref-1", StudentID: "student-123", Status: "draft"}
	mockStudent := &model.Student{ID: "student-123"}

	achRepo.On("GetReferenceByID", "ref-1").Return(mockRef, nil)
	stuRepo.On("FindByUserID", "user-123").Return(mockStudent, nil)
	achRepo.On("UpdateReference", mock.MatchedBy(func(r *model.AchievementReference) bool {
		return r.Status == "submitted" // Verify status change to 'submitted'
	})).Return(nil)

	req := httptest.NewRequest("POST", "/achievements/ref-1/submit", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}