package service_test

import (
	"UASBE/app/model"
	"UASBE/app/service"
	"UASBE/test/mocks"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetLecturerAdvisees_Success(t *testing.T) {
	// Setup Mocks
	lecRepo := new(mocks.MockLecturerRepository)
	stuRepo := new(mocks.MockStudentRepository)
	userRepo := new(mocks.MockUserRepository)
	svc := service.NewLecturerService(lecRepo, stuRepo, nil, userRepo)

	app := fiber.New()
	app.Get("/lecturers/:id/advisees", func(c *fiber.Ctx) error {
		// Mock login sebagai Dosen Wali pemilik data
		c.Locals("user", &model.JWTClaims{UserID: "user-lec-1", Role: "Dosen Wali"})
		return svc.GetLecturerAdvisees(c)
	})

	// Mock Data
	lecturerID := "lec-1"
	mockLecturer := &model.Lecturer{ID: lecturerID}
	mockLecUser := &model.User{ID: lecturerID, FullName: "Dosen Wali"}
	
	// Mahasiswa yang memiliki AdvisorID sama dengan lecturerID
	mockStudents := []model.Student{
		{ID: "std-1", StudentID: "2021001", AdvisorID: &lecturerID},
	}
	mockStdUser := &model.User{FullName: "Advisee Name"}

	// Expectations
	lecRepo.On("FindByUserID", "user-lec-1").Return(mockLecturer, nil)
	lecRepo.On("FindByID", lecturerID).Return(mockLecturer, nil)
	stuRepo.On("GetAll", 1000, 0).Return(mockStudents, nil)
	userRepo.On("FindByID", "std-1").Return(mockStdUser, nil)
	userRepo.On("FindByID", lecturerID).Return(mockLecUser, nil)

	req := httptest.NewRequest("GET", "/lecturers/lec-1/advisees", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}