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

func TestCreateUser_Student_Success(t *testing.T) {
	// Setup
	userRepo := new(mocks.MockUserRepository)
	roleRepo := new(mocks.MockRoleRepository)
	stuRepo := new(mocks.MockStudentRepository)
	svc := service.NewUserService(userRepo, roleRepo, nil, stuRepo, nil)

	app := fiber.New()
	app.Post("/users", svc.CreateUser)

	// Mocking logic
	userRepo.On("FindByUsername", "newstudent").Return(nil, nil)
	userRepo.On("FindByEmail", "student@mail.com").Return(nil, nil)
	roleRepo.On("GetRoleByName", "Mahasiswa").Return(&model.Role{ID: "role-mhs", Name: "Mahasiswa"}, nil)
	userRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)
	
	// Mock untuk pengecekan student_id unik
	stuRepo.On("FindByStudentID", "2025001").Return(nil, nil)
	stuRepo.On("Create", mock.AnythingOfType("*model.Student")).Return(nil)

	// ⭐ TAMBAHKAN INI: Mock untuk buildUserResponse di akhir CreateUser
	// Fungsi buildUserResponse memanggil FindByUserID untuk menyusun profil di JSON response
	mockStudentResult := &model.Student{
		ID:           "generated-uuid",
		StudentID:    "2025001",
		ProgramStudy: "Informatika",
		AcademicYear: 2025,
	}
	stuRepo.On("FindByUserID", mock.Anything).Return(mockStudentResult, nil)

	// Payload
	reqBody := model.UserCreateRequest{
		Username: "newstudent",
		Email:    "student@mail.com",
		Password: "Password123!",
		FullName: "Mahasiswa Baru",
		RoleName: "Mahasiswa",
		StudentProfile: &model.StudentProfileRequest{
			StudentID:    "2025001",
			ProgramStudy: "Informatika",
			AcademicYear: 2025,
		},
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, 201, resp.StatusCode)
}