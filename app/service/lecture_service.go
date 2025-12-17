package service

import (
	"math"
	"strconv"

	"UASBE/app/model"
	"UASBE/app/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type LecturerService struct {
	lecturerRepo    repository.LecturerRepository
	studentRepo     repository.StudentRepository
	achievementRepo repository.AchievementRepository
	userRepo        repository.UserRepository
	validate        *validator.Validate
}

func NewLecturerService(
	lecturerRepo repository.LecturerRepository,
	studentRepo repository.StudentRepository,
	achievementRepo repository.AchievementRepository,
	userRepo repository.UserRepository,
) *LecturerService {
	return &LecturerService{
		lecturerRepo:    lecturerRepo,
		studentRepo:     studentRepo,
		achievementRepo: achievementRepo,
		userRepo:        userRepo,
		validate:        validator.New(),
	}
}

//
// ==================== GET ALL LECTURERS (GET /lecturers) ======================
// Sesuai SRS Section 5.5: GET /api/v1/lecturers
//

func (s *LecturerService) GetAllLecturers(c *fiber.Ctx) error {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get lecturers
	lecturers, err := s.lecturerRepo.GetAll(pageSize, offset)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch lecturers",
		})
	}

	// Count total
	total, err := s.lecturerRepo.CountAll()
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to count lecturers",
		})
	}

	// Build response with user details
	var responses []map[string]interface{}
	for _, lecturer := range lecturers {
		// Get user details
		user, err := s.userRepo.FindByID(lecturer.ID)
		if err != nil {
			continue
		}

		responses = append(responses, map[string]interface{}{
			"id":           lecturer.ID,
			"lecturer_id":  lecturer.LecturerID,
			"department":   lecturer.Department,
			"user": map[string]interface{}{
				"username":  user.Username,
				"email":     user.Email,
				"full_name": user.FullName,
				"is_active": user.IsActive,
			},
			"created_at": lecturer.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.JSON(model.APIResponse{
		Status: "success",
		Data: fiber.Map{
			"lecturers":   responses,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
		},
	})
}

//
// ==================== GET LECTURER'S ADVISEES (GET /lecturers/:id/advisees) ======================
// Sesuai SRS Section 5.5: GET /api/v1/lecturers/:id/advisees
// FR-006: Dosen wali melihat daftar prestasi mahasiswa bimbingannya
//

func (s *LecturerService) GetLecturerAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	// Get user dari context untuk authorization
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Authorization check:
	// Dosen Wali hanya bisa lihat advisees sendiri
	// Admin bisa lihat semua
	if claims.Role == "Dosen Wali" {
		lecturer, err := s.lecturerRepo.FindByUserID(claims.UserID)
		if err != nil || lecturer.ID != lecturerID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden: you can only view your own advisees",
			})
		}
	}

	// Verify lecturer exists
	lecturer, err := s.lecturerRepo.FindByID(lecturerID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "lecturer not found",
		})
	}

	// Parse query params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	includeAchievements := c.Query("include_achievements", "false") == "true"

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get all students
	allStudents, err := s.studentRepo.GetAll(1000, 0)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch students",
		})
	}

	// Filter students yang advisornya adalah lecturer ini
	var advisees []model.Student
	for _, student := range allStudents {
		if student.AdvisorID != nil && *student.AdvisorID == lecturerID {
			advisees = append(advisees, student)
		}
	}

	total := len(advisees)

	// Apply pagination
	start := offset
	end := offset + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedAdvisees := advisees[start:end]

	// Build response
	var responses []map[string]interface{}
	for _, student := range paginatedAdvisees {
		// Get user details
		user, err := s.userRepo.FindByID(student.ID)
		if err != nil {
			continue
		}

		studentData := map[string]interface{}{
			"id":            student.ID,
			"student_id":    student.StudentID,
			"program_study": student.ProgramStudy,
			"academic_year": student.AcademicYear,
			"user": map[string]interface{}{
				"username":  user.Username,
				"email":     user.Email,
				"full_name": user.FullName,
				"is_active": user.IsActive,
			},
			"created_at": student.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		// Include achievements if requested
		if includeAchievements {
			// Get achievement count by status
			totalAchievements, _ := s.achievementRepo.CountReferencesByStudentID(student.ID, "")
			submittedCount, _ := s.achievementRepo.CountReferencesByStudentID(student.ID, "submitted")
			verifiedCount, _ := s.achievementRepo.CountReferencesByStudentID(student.ID, "verified")
			rejectedCount, _ := s.achievementRepo.CountReferencesByStudentID(student.ID, "rejected")

			studentData["achievements_summary"] = map[string]interface{}{
				"total":     totalAchievements,
				"submitted": submittedCount,
				"verified":  verifiedCount,
				"rejected":  rejectedCount,
			}
		}

		responses = append(responses, studentData)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	// Get lecturer user details for response
	lecturerUser, _ := s.userRepo.FindByID(lecturer.ID)

	return c.JSON(model.APIResponse{
		Status: "success",
		Data: fiber.Map{
			"lecturer": map[string]interface{}{
				"id":          lecturer.ID,
				"lecturer_id": lecturer.LecturerID,
				"department":  lecturer.Department,
				"user": map[string]interface{}{
					"username":  lecturerUser.Username,
					"email":     lecturerUser.Email,
					"full_name": lecturerUser.FullName,
				},
			},
			"advisees":    responses,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
		},
	})
}