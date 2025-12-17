package service

import (
	"UASBE/app/model"
	"UASBE/app/repository"

	"github.com/gofiber/fiber/v2"
)

type ReportService struct {
	reportRepo      repository.ReportRepository
	achievementRepo repository.AchievementRepository
	studentRepo     repository.StudentRepository
	lecturerRepo    repository.LecturerRepository
	userRepo        repository.UserRepository
}

func NewReportService(
	reportRepo repository.ReportRepository,
	achievementRepo repository.AchievementRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
	userRepo repository.UserRepository,
) *ReportService {
	return &ReportService{
		reportRepo:      reportRepo,
		achievementRepo: achievementRepo,
		studentRepo:     studentRepo,
		lecturerRepo:    lecturerRepo,
		userRepo:        userRepo,
	}
}

//
// ==================== GET STATISTICS (GET /reports/statistics) ======================
// FR-011: Achievement Statistics
// Actor: Mahasiswa (own), Dosen Wali (advisee), Admin (all)
//

func (s *ReportService) GetStatistics(c *fiber.Ctx) error {
	// Get user dari context
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	var studentID *string
	var advisorID *string

	// Determine scope based on role
	if claims.Role == "Mahasiswa" {
		// Mahasiswa: hanya statistik sendiri
		student, err := s.studentRepo.FindByUserID(claims.UserID)
		if err != nil {
			return c.Status(404).JSON(model.APIResponse{
				Status: "error",
				Error:  "student profile not found",
			})
		}
		studentID = &student.ID

	} else if claims.Role == "Dosen Wali" {
		// Dosen Wali: statistik mahasiswa bimbingannya
		lecturer, err := s.lecturerRepo.FindByUserID(claims.UserID)
		if err != nil {
			return c.Status(404).JSON(model.APIResponse{
				Status: "error",
				Error:  "lecturer profile not found",
			})
		}
		advisorID = &lecturer.ID
	}
	// Admin: studentID & advisorID tetap nil (semua data)

	// Get statistics
	stats := &model.AchievementStatistics{}

	// 1. Total by type
	totalByType, err := s.reportRepo.GetTotalByType(studentID, advisorID)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to get statistics by type",
		})
	}
	stats.TotalByType = totalByType

	// 2. Total by period
	totalByPeriod, err := s.reportRepo.GetTotalByPeriod(studentID, advisorID)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to get statistics by period",
		})
	}
	stats.TotalByPeriod = totalByPeriod

	// 3. Top students (kecuali untuk Mahasiswa)
	if claims.Role != "Mahasiswa" {
		topStudents, err := s.reportRepo.GetTopStudents(10, advisorID)
		if err == nil {
			stats.TopStudents = topStudents
		}
	}

	// 4. Competition level distribution
	competitionDist, err := s.reportRepo.GetCompetitionLevelDistribution(studentID, advisorID)
	if err == nil {
		stats.CompetitionLevelDistribution = competitionDist
	}

	// 5. Status breakdown
	statusBreakdown, err := s.reportRepo.GetStatusBreakdown(studentID, advisorID)
	if err == nil {
		stats.StatusBreakdown = statusBreakdown
	}

	// 6. Calculate total achievements
	totalAchievements := 0
	for _, count := range stats.StatusBreakdown {
		totalAchievements += count
	}
	stats.TotalAchievements = totalAchievements

	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   stats,
	})
}

//
// ==================== GET STUDENT REPORT (GET /reports/student/:id) ======================
// FR-011: Detail report untuk satu mahasiswa
// Actor: Mahasiswa (own), Dosen Wali (advisee), Admin (all)
//

func (s *ReportService) GetStudentReport(c *fiber.Ctx) error {
	studentID := c.Params("id")

	// Get user dari context
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
		// Mahasiswa hanya bisa lihat report sendiri
		currentStudent, _ := s.studentRepo.FindByUserID(claims.UserID)
		if currentStudent == nil || currentStudent.ID != studentID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden: you can only view your own report",
			})
		}
	} else if claims.Role == "Dosen Wali" {
		// Dosen wali hanya bisa lihat report advisees
		lecturer, _ := s.lecturerRepo.FindByUserID(claims.UserID)
		if lecturer == nil || student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden: you can only view reports of your advisees",
			})
		}
	}
	// Admin bisa lihat semua

	// Build student info
	user, err := s.userRepo.FindByID(student.ID)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch user details",
		})
	}

	studentInfo := model.StudentInfo{
		ID:           student.ID,
		StudentID:    student.StudentID,
		FullName:     user.FullName,
		Email:        user.Email,
		ProgramStudy: student.ProgramStudy,
		AcademicYear: student.AcademicYear,
	}

	// Get advisor name if exists
	if student.AdvisorID != nil {
		advisor, err := s.lecturerRepo.FindByID(*student.AdvisorID)
		if err == nil {
			advisorUser, err := s.userRepo.FindByID(advisor.ID)
			if err == nil {
				studentInfo.AdvisorName = &advisorUser.FullName
			}
		}
	}

	// Get summary statistics
	summary, err := s.reportRepo.GetStudentSummary(studentID)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to get student summary",
		})
	}

	// Get achievements by type
	achievementsByType, err := s.reportRepo.GetStudentAchievementsByType(studentID)
	if err != nil {
		achievementsByType = make(map[string]int)
	}

	// Get achievements by status
	achievementsByStatus, err := s.reportRepo.GetStudentAchievementsByStatus(studentID)
	if err != nil {
		achievementsByStatus = make(map[string]int)
	}

	// Get recent achievements (last 10)
	references, err := s.achievementRepo.GetReferencesByStudentID(studentID, "", 10, 0)
	if err != nil {
		references = []model.AchievementReference{}
	}

	var recentAchievements []model.AchievementResponse
	for _, ref := range references {
		achievement, err := s.achievementRepo.GetAchievementByID(ref.MongoAchievementID)
		if err != nil {
			continue
		}

		response := s.buildAchievementResponse(achievement, &ref)
		recentAchievements = append(recentAchievements, *response)
	}

	// Get timeline
	timeline, err := s.reportRepo.GetStudentTimeline(studentID)
	if err != nil {
		timeline = []model.PeriodStats{}
	}

	// Build report
	report := model.StudentReport{
		Student:              studentInfo,
		Summary:              *summary,
		AchievementsByType:   achievementsByType,
		AchievementsByStatus: achievementsByStatus,
		RecentAchievements:   recentAchievements,
		Timeline:             timeline,
	}

	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   report,
	})
}

//
// ==================== HELPER: BUILD ACHIEVEMENT RESPONSE ======================
//

func (s *ReportService) buildAchievementResponse(
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
