package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"math"
	"strconv"

	"UASBE/app/model"
	"UASBE/app/repository"
)

type StudentService struct {
	studentRepo     repository.StudentRepository
	lecturerRepo    repository.LecturerRepository
	achievementRepo repository.AchievementRepository
	userRepo        repository.UserRepository
	validate        *validator.Validate
}

func NewStudentService(
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
	achievementRepo repository.AchievementRepository,
	userRepo repository.UserRepository,
) *StudentService {
	return &StudentService{
		studentRepo:     studentRepo,
		lecturerRepo:    lecturerRepo,
		achievementRepo: achievementRepo,
		userRepo:        userRepo,
		validate:        validator.New(),
	}
}

//
// ==================== GET ALL STUDENTS (GET /students) ======================
// Sesuai SRS Section 5.5: GET /api/v1/students
//

func (s *StudentService) GetAllStudents(c *fiber.Ctx) error {
	// Get user dari context untuk authorization
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

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

	// Authorization: Dosen Wali hanya bisa lihat advisees
	if claims.Role == "Dosen Wali" {
		lecturer, err := s.lecturerRepo.FindByUserID(claims.UserID)
		if err != nil {
			return c.Status(404).JSON(model.APIResponse{
				Status: "error",
				Error:  "lecturer profile not found",
			})
		}

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
			if student.AdvisorID != nil && *student.AdvisorID == lecturer.ID {
				advisees = append(advisees, student)
			}
		}

		// Apply pagination
		total := len(advisees)
		start := offset
		end := offset + pageSize
		if start > total {
			start = total
		}
		if end > total {
			end = total
		}

		paginatedStudents := advisees[start:end]

		// Build response
		var responses []map[string]interface{}
		for _, student := range paginatedStudents {
			user, _ := s.userRepo.FindByID(student.ID)
			
			responses = append(responses, map[string]interface{}{
				"id":            student.ID,
				"student_id":    student.StudentID,
				"program_study": student.ProgramStudy,
				"academic_year": student.AcademicYear,
				"advisor_id":    student.AdvisorID,
				"user": map[string]interface{}{
					"username":  user.Username,
					"email":     user.Email,
					"full_name": user.FullName,
					"is_active": user.IsActive,
				},
				"created_at": student.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

		return c.JSON(model.APIResponse{
			Status: "success",
			Data: fiber.Map{
				"students":    responses,
				"total":       total,
				"page":        page,
				"page_size":   pageSize,
				"total_pages": totalPages,
			},
		})
	}

	// Admin & Mahasiswa bisa lihat semua
	students, err := s.studentRepo.GetAll(pageSize, offset)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch students",
		})
	}

	total, err := s.studentRepo.CountAll()
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to count students",
		})
	}

	// Build response with user details
	var responses []map[string]interface{}
	for _, student := range students {
		user, err := s.userRepo.FindByID(student.ID)
		if err != nil {
			continue
		}

		studentData := map[string]interface{}{
			"id":            student.ID,
			"student_id":    student.StudentID,
			"program_study": student.ProgramStudy,
			"academic_year": student.AcademicYear,
			"advisor_id":    student.AdvisorID,
			"user": map[string]interface{}{
				"username":  user.Username,
				"email":     user.Email,
				"full_name": user.FullName,
				"is_active": user.IsActive,
			},
			"created_at": student.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		// Include advisor details if exists
		if student.AdvisorID != nil {
			advisor, err := s.lecturerRepo.FindByID(*student.AdvisorID)
			if err == nil {
				advisorUser, _ := s.userRepo.FindByID(advisor.ID)
				studentData["advisor"] = map[string]interface{}{
					"id":          advisor.ID,
					"lecturer_id": advisor.LecturerID,
					"full_name":   advisorUser.FullName,
					"department":  advisor.Department,
				}
			}
		}

		responses = append(responses, studentData)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.JSON(model.APIResponse{
		Status: "success",
		Data: fiber.Map{
			"students":    responses,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
		},
	})
}

//
// ==================== GET STUDENT BY ID (GET /students/:id) ======================
// Sesuai SRS Section 5.5: GET /api/v1/students/:id
//

func (s *StudentService) GetStudentByID(c *fiber.Ctx) error {
	studentID := c.Params("id")

	// Get user dari context untuk authorization
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "student not found",
		})
	}

	// Authorization check
	if claims.Role == "Mahasiswa" {
		// Mahasiswa hanya bisa lihat profile sendiri
		currentStudent, _ := s.studentRepo.FindByUserID(claims.UserID)
		if currentStudent == nil || currentStudent.ID != studentID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden: you can only view your own profile",
			})
		}
	} else if claims.Role == "Dosen Wali" {
		// Dosen wali hanya bisa lihat advisees
		lecturer, _ := s.lecturerRepo.FindByUserID(claims.UserID)
		if lecturer == nil || student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden: you can only view your advisees",
			})
		}
	}
	// Admin bisa lihat semua

	// Get user details
	user, err := s.userRepo.FindByID(student.ID)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch user details",
		})
	}

	response := map[string]interface{}{
		"id":            student.ID,
		"student_id":    student.StudentID,
		"program_study": student.ProgramStudy,
		"academic_year": student.AcademicYear,
		"advisor_id":    student.AdvisorID,
		"user": map[string]interface{}{
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.FullName,
			"is_active": user.IsActive,
		},
		"created_at": student.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// Include advisor details if exists
	if student.AdvisorID != nil {
		advisor, err := s.lecturerRepo.FindByID(*student.AdvisorID)
		if err == nil {
			advisorUser, _ := s.userRepo.FindByID(advisor.ID)
			response["advisor"] = map[string]interface{}{
				"id":          advisor.ID,
				"lecturer_id": advisor.LecturerID,
				"full_name":   advisorUser.FullName,
				"department":  advisor.Department,
			}
		}
	}

	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   response,
	})
}

//
// ==================== GET STUDENT ACHIEVEMENTS (GET /students/:id/achievements) ======================
// Sesuai SRS Section 5.5: GET /api/v1/students/:id/achievements
// FR-006: Dosen wali melihat daftar prestasi mahasiswa bimbingannya
//

func (s *StudentService) GetStudentAchievements(c *fiber.Ctx) error {
	studentID := c.Params("id")

	// Get user dari context untuk authorization
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Verify student exists
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "student not found",
		})
	}

	// Authorization check
	if claims.Role == "Mahasiswa" {
		// Mahasiswa hanya bisa lihat achievements sendiri
		currentStudent, _ := s.studentRepo.FindByUserID(claims.UserID)
		if currentStudent == nil || currentStudent.ID != studentID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden: you can only view your own achievements",
			})
		}
	} else if claims.Role == "Dosen Wali" {
		// Dosen wali hanya bisa lihat achievements advisees
		lecturer, _ := s.lecturerRepo.FindByUserID(claims.UserID)
		if lecturer == nil || student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden: you can only view achievements of your advisees",
			})
		}
	}
	// Admin bisa lihat semua

	// Parse query params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	status := c.Query("status", "") // filter by status

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get achievement references
	references, err := s.achievementRepo.GetReferencesByStudentID(studentID, status, pageSize, offset)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch achievements",
		})
	}

	// Count total
	total, err := s.achievementRepo.CountReferencesByStudentID(studentID, status)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to count achievements",
		})
	}

	// Fetch details dari MongoDB
	var achievements []model.AchievementResponse
	for _, ref := range references {
		achievement, err := s.achievementRepo.GetAchievementByID(ref.MongoAchievementID)
		if err != nil {
			continue
		}

		response := s.buildAchievementResponse(achievement, &ref)
		achievements = append(achievements, *response)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	// Get student user details
	user, _ := s.userRepo.FindByID(student.ID)

	return c.JSON(model.APIResponse{
		Status: "success",
		Data: fiber.Map{
			"student": map[string]interface{}{
				"id":            student.ID,
				"student_id":    student.StudentID,
				"program_study": student.ProgramStudy,
				"academic_year": student.AcademicYear,
				"user": map[string]interface{}{
					"username":  user.Username,
					"email":     user.Email,
					"full_name": user.FullName,
				},
			},
			"achievements": achievements,
			"total":        total,
			"page":         page,
			"page_size":    pageSize,
			"total_pages":  totalPages,
		},
	})
}

//
// ==================== SET ADVISOR (PUT /students/:id/advisor) ======================
// Sesuai SRS Section 5.5: PUT /api/v1/students/:id/advisor
// FR-009: Admin dapat set advisor untuk mahasiswa
//

func (s *StudentService) SetAdvisor(c *fiber.Ctx) error {
	studentID := c.Params("id")

	// Cari student
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "student not found",
		})
	}

	// Parse request
	req := new(model.SetAdvisorRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "invalid request body",
		})
	}

	// Validate
	if err := s.validate.Struct(req); err != nil {
		return c.Status(422).JSON(model.APIResponse{
			Status: "error",
			Error:  err.Error(),
		})
	}

	// Validasi advisor exists
	advisor, err := s.lecturerRepo.FindByID(req.AdvisorID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "advisor (lecturer) not found",
		})
	}

	// Update advisor
	if err := s.studentRepo.SetAdvisor(studentID, req.AdvisorID); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to set advisor",
		})
	}

	// Refresh student data
	student, _ = s.studentRepo.FindByID(studentID)
	user, _ := s.userRepo.FindByID(student.ID)
	advisorUser, _ := s.userRepo.FindByID(advisor.ID)

	response := map[string]interface{}{
		"id":            student.ID,
		"student_id":    student.StudentID,
		"program_study": student.ProgramStudy,
		"academic_year": student.AcademicYear,
		"user": map[string]interface{}{
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.FullName,
		},
		"advisor": map[string]interface{}{
			"id":          advisor.ID,
			"lecturer_id": advisor.LecturerID,
			"full_name":   advisorUser.FullName,
			"department":  advisor.Department,
		},
		"created_at": student.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "advisor set successfully",
		Data:    response,
	})
}

//
// ==================== HELPER: BUILD ACHIEVEMENT RESPONSE ======================
//

func (s *StudentService) buildAchievementResponse(
	achievement *model.Achievement,
	reference *model.AchievementReference,
) *model.AchievementResponse {
	response := &model.AchievementResponse{
		ID:              reference.ID,
		StudentID:       achievement.StudentID,
		AchievementType: achievement.AchievementType,
		Title:           achievement.Title,
		Description:     achievement.Description,
		Details:         achievement.Details,
		Attachments:     achievement.Attachments,
		Tags:            achievement.Tags,
		Points:          achievement.Points,
		Status:          reference.Status,
		CreatedAt:       achievement.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       achievement.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if reference.SubmittedAt != nil {
		submittedAt := reference.SubmittedAt.Format("2006-01-02 15:04:05")
		response.SubmittedAt = &submittedAt
	}

	if reference.VerifiedAt != nil {
		verifiedAt := reference.VerifiedAt.Format("2006-01-02 15:04:05")
		response.VerifiedAt = &verifiedAt
	}

	response.VerifiedBy = reference.VerifiedBy
	response.RejectionNote = reference.RejectionNote

	return response
}