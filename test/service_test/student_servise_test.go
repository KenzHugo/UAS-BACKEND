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
)

func TestSetAdvisor_Success(t *testing.T) {
	// Setup Mocks
	stuRepo := new(mocks.MockStudentRepository)
	lecRepo := new(mocks.MockLecturerRepository)
	userRepo := new(mocks.MockUserRepository)
	svc := service.NewStudentService(stuRepo, lecRepo, nil, userRepo)

	app := fiber.New()
	app.Put("/students/:id/advisor", svc.SetAdvisor)

	// Mock Data
	studentID := "std-1"
	lecturerID := "lec-1"
	mockStudent := &model.Student{ID: studentID, StudentID: "2021001"}
	mockLecturer := &model.Lecturer{ID: lecturerID, LecturerID: "1980001"}
	mockUser := &model.User{FullName: "Student Name"}
	mockLecUser := &model.User{FullName: "Lecturer Name"}

	// Expectations sesuai alur FR-009
	stuRepo.On("FindByID", studentID).Return(mockStudent, nil)
	lecRepo.On("FindByID", lecturerID).Return(mockLecturer, nil)
	stuRepo.On("SetAdvisor", studentID, lecturerID).Return(nil)
	userRepo.On("FindByID", studentID).Return(mockUser, nil)
	userRepo.On("FindByID", lecturerID).Return(mockLecUser, nil)

	// Payload request
	reqBody := model.SetAdvisorRequest{AdvisorID: lecturerID}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/students/std-1/advisor", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}