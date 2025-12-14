package service

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"UASBE/app/model"
	"UASBE/app/repository"
)

type AchievementService struct {
	achievementRepo repository.AchievementRepository
	studentRepo     repository.StudentRepository
	lecturerRepo    repository.LecturerRepository
	userRepo        repository.UserRepository
	validate        *validator.Validate
}

func NewAchievementService(
	achievementRepo repository.AchievementRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
	userRepo repository.UserRepository,
) *AchievementService {
	return &AchievementService{
		achievementRepo: achievementRepo,
		studentRepo:     studentRepo,
		lecturerRepo:    lecturerRepo,
		userRepo:        userRepo,
		validate:        validator.New(),
	}
}

//
// ==================== FR-003: CREATE ACHIEVEMENT (POST /achievements) ======================
// Actor: Mahasiswa
// Flow:
// 1. Mahasiswa mengisi data prestasi
// 2. Mahasiswa upload dokumen pendukung
// 3. Sistem simpan ke MongoDB (achievement) dan PostgreSQL (reference)
// 4. Status awal: 'draft'
// 5. Return achievement data
//

func (s *AchievementService) CreateAchievement(c *fiber.Ctx) error {
	// Get user from JWT
	claims := c.Locals("user").(*model.JWTClaims)
	
	// Verify user is Mahasiswa
	if claims.Role != "Mahasiswa" {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "only students can create achievements",
		})
	}
	
	// Parse request
	req := new(model.CreateAchievementRequest)
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
	
	// Get student profile
	student, err := s.studentRepo.FindByUserID(claims.UserID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "student profile not found",
		})
	}
	
	// Step 1 & 2: Create achievement in MongoDB
	achievement := &model.Achievement{
		StudentID:   student.ID,
		Type:        req.AchievementType,
		Title:       req.Title,
		Description: req.Description,
		Details:     req.Details,
		Tags:        req.Tags,
		Points:      req.Points,
		Attachments: []model.Attachment{}, // Will be added via separate endpoint
	}
	
	mongoID, err := s.achievementRepo.CreateAchievement(achievement)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to create achievement in MongoDB",
		})
	}
	
	// Step 3: Create reference in PostgreSQL
	ref := &model.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          student.ID,
		MongoAchievementID: mongoID,
		Status:             "draft", // Step 4: Status awal draft
	}
	
	if err := s.achievementRepo.CreateReference(ref); err != nil {
		// Rollback: delete from MongoDB
		s.achievementRepo.DeleteAchievement(mongoID)
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to create achievement reference",
		})
	}
	
	// Step 5: Build response
	response := s.buildAchievementResponse(ref, achievement, student, nil, nil)
	
	return c.Status(201).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement created successfully",
		Data:    response,
	})
}

//
// ==================== GET ACHIEVEMENTS LIST (GET /achievements) ======================
// Actor: Mahasiswa (own), Dosen Wali (advisee), Admin (all)
// FR-006 (partial): View Prestasi Mahasiswa Bimbingan
//

func (s *AchievementService) GetAchievements(c *fiber.Ctx) error {
	claims := c.Locals("user").(*model.JWTClaims)
	
	// Parse filters
	filters := &model.AchievementFilters{
		Status:          c.Query("status", ""),
		AchievementType: c.Query("achievement_type", ""),
		Page:            1,
		PageSize:        10,
	}
	
	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		filters.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("page_size", "10")); err == nil && pageSize > 0 {
		filters.PageSize = pageSize
	}
	
	// Apply role-based filters
	switch claims.Role {
	case "Mahasiswa":
		// Mahasiswa hanya lihat prestasi sendiri
		student, err := s.studentRepo.FindByUserID(claims.UserID)
		if err != nil {
			return c.Status(404).JSON(model.APIResponse{
				Status: "error",
				Error:  "student profile not found",
			})
		}
		filters.StudentID = student.ID
		
	case "Dosen Wali":
		// FR-006: Dosen Wali hanya lihat prestasi mahasiswa bimbingannya
		lecturer, err := s.lecturerRepo.FindByUserID(claims.UserID)
		if err != nil {
			return c.Status(404).JSON(model.APIResponse{
				Status: "error",
				Error:  "lecturer profile not found",
			})
		}
		filters.AdvisorID = lecturer.ID
		
	case "Admin":
		// Admin bisa lihat semua
		// No additional filters
		
	default:
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized access",
		})
	}
	
	// Get references from PostgreSQL
	references, err := s.achievementRepo.GetReferences(filters)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch achievements",
		})
	}
	
	// Get total count
	total, err := s.achievementRepo.CountReferences(filters)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to count achievements",
		})
	}
	
	// Get achievement details from MongoDB
	mongoIDs := make([]string, len(references))
	for i, ref := range references {
		mongoIDs[i] = ref.MongoAchievementID
	}
	
	achievements, err := s.achievementRepo.GetAchievements(mongoIDs)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch achievement details",
		})
	}
	
	// Build response
	achievementMap := make(map[string]*model.Achievement)
	for i := range achievements {
		achievementMap[achievements[i].ID] = &achievements[i]
	}
	
	var responses []model.AchievementResponse
	for _, ref := range references {
		achievement := achievementMap[ref.MongoAchievementID]
		if achievement == nil {
			continue
		}
		
		// Get student info
		student, _ := s.studentRepo.FindByID(ref.StudentID)
		
		// Get verified by info
		var verifiedByUser *model.User
		if ref.VerifiedBy != nil {
			verifiedByUser, _ = s.userRepo.FindByID(*ref.VerifiedBy)
		}
		
		response := s.buildAchievementResponse(&ref, achievement, student, verifiedByUser, nil)
		responses = append(responses, *response)
	}
	
	totalPages := int(math.Ceil(float64(total) / float64(filters.PageSize)))
	
	listResponse := model.AchievementListResponse{
		Achievements: responses,
		Total:        total,
		Page:         filters.Page,
		PageSize:     filters.PageSize,
		TotalPages:   totalPages,
	}
	
	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   listResponse,
	})
}

//
// ==================== GET ACHIEVEMENT BY ID (GET /achievements/:id) ======================
//

func (s *AchievementService) GetAchievementByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*model.JWTClaims)
	achievementID := c.Params("id")
	
	// Get reference from PostgreSQL
	ref, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}
	
	// Authorization check
	if err := s.checkAchievementAccess(claims, ref); err != nil {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  err.Error(),
		})
	}
	
	// Get achievement from MongoDB
	achievement, err := s.achievementRepo.GetAchievementByID(ref.MongoAchievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement details not found",
		})
	}
	
	// Get student info
	student, _ := s.studentRepo.FindByID(ref.StudentID)
	
	// Get verified by info
	var verifiedByUser *model.User
	if ref.VerifiedBy != nil {
		verifiedByUser, _ = s.userRepo.FindByID(*ref.VerifiedBy)
	}
	
	response := s.buildAchievementResponse(ref, achievement, student, verifiedByUser, nil)
	
	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   response,
	})
}

//
// ==================== FR-003 (UPDATE): UPDATE ACHIEVEMENT (PUT /achievements/:id) ======================
// Actor: Mahasiswa
// Precondition: Status 'draft'
//

func (s *AchievementService) UpdateAchievement(c *fiber.Ctx) error {
	claims := c.Locals("user").(*model.JWTClaims)
	achievementID := c.Params("id")
	
	// Verify user is Mahasiswa
	if claims.Role != "Mahasiswa" {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "only students can update achievements",
		})
	}
	
	// Get reference
	ref, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}
	
	// Check ownership
	student, err := s.studentRepo.FindByUserID(claims.UserID)
	if err != nil || student.ID != ref.StudentID {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "you can only update your own achievements",
		})
	}
	
	// Check status (can only update draft)
	if ref.Status != "draft" {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "can only update draft achievements",
		})
	}
	
	// Parse request
	req := new(model.UpdateAchievementRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "invalid request body",
		})
	}
	
	// Get existing achievement from MongoDB
	achievement, err := s.achievementRepo.GetAchievementByID(ref.MongoAchievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement details not found",
		})
	}
	
	// Update fields
	if req.Title != "" {
		achievement.Title = req.Title
	}
	if req.Description != "" {
		achievement.Description = req.Description
	}
	if req.Details != nil {
		achievement.Details = req.Details
	}
	if req.Tags != nil {
		achievement.Tags = req.Tags
	}
	if req.Points > 0 {
		achievement.Points = req.Points
	}
	
	// Update in MongoDB
	if err := s.achievementRepo.UpdateAchievement(ref.MongoAchievementID, achievement); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to update achievement",
		})
	}
	
	// Update timestamp in PostgreSQL
	ref.UpdatedAt = time.Now()
	s.achievementRepo.UpdateReference(ref)
	
	response := s.buildAchievementResponse(ref, achievement, student, nil, nil)
	
	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement updated successfully",
		Data:    response,
	})
}

//
// ==================== FR-005: DELETE ACHIEVEMENT (DELETE /achievements/:id) ======================
// Actor: Mahasiswa
// Precondition: Status 'draft'
// Flow:
// 1. Soft delete data di MongoDB
// 2. Update reference di PostgreSQL
// 3. Return success message
//

func (s *AchievementService) DeleteAchievement(c *fiber.Ctx) error {
	claims := c.Locals("user").(*model.JWTClaims)
	achievementID := c.Params("id")
	
	// Verify user is Mahasiswa
	if claims.Role != "Mahasiswa" {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "only students can delete achievements",
		})
	}
	
	// Get reference
	ref, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}
	
	// Check ownership
	student, err := s.studentRepo.FindByUserID(claims.UserID)
	if err != nil || student.ID != ref.StudentID {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "you can only delete your own achievements",
		})
	}
	
	// Check status (can only delete draft)
	if ref.Status != "draft" {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "can only delete draft achievements",
		})
	}
	
	// Step 1: Soft delete in PostgreSQL (update status to 'deleted')
	if err := s.achievementRepo.DeleteReference(achievementID); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to delete achievement",
		})
	}
	
	// Step 2: Hard delete from MongoDB (optional, bisa juga soft delete)
	// s.achievementRepo.DeleteAchievement(ref.MongoAchievementID)
	
	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement deleted successfully",
	})
}

//
// ==================== FR-004: SUBMIT FOR VERIFICATION (POST /achievements/:id/submit) ======================
// Actor: Mahasiswa
// Precondition: Prestasi berstatus 'draft'
// Flow:
// 1. Mahasiswa submit prestasi
// 2. Update status menjadi 'submitted'
// 3. Create notification untuk dosen wali (skip for now)
// 4. Return updated status
//

func (s *AchievementService) SubmitForVerification(c *fiber.Ctx) error {
	claims := c.Locals("user").(*model.JWTClaims)
	achievementID := c.Params("id")
	
	// Verify user is Mahasiswa
	if claims.Role != "Mahasiswa" {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "only students can submit achievements",
		})
	}
	
	// Get reference
	ref, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}
	
	// Check ownership
	student, err := s.studentRepo.FindByUserID(claims.UserID)
	if err != nil || student.ID != ref.StudentID {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "you can only submit your own achievements",
		})
	}
	
	// Check status (must be draft)
	if ref.Status != "draft" {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement must be in draft status to submit",
		})
	}
	
	// Check if student has advisor
	if student.AdvisorID == nil {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "cannot submit: no advisor assigned",
		})
	}
	
	// Step 1 & 2: Update status to 'submitted'
	now := time.Now()
	ref.Status = "submitted"
	ref.SubmittedAt = &now
	ref.UpdatedAt = now
	
	if err := s.achievementRepo.UpdateReference(ref); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to submit achievement",
		})
	}
	
	// Step 3: Create notification (TODO: implement notification system)
	
	// Step 4: Return response
	achievement, _ := s.achievementRepo.GetAchievementByID(ref.MongoAchievementID)
	response := s.buildAchievementResponse(ref, achievement, student, nil, nil)
	
	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement submitted for verification",
		Data:    response,
	})
}

//
// ==================== FR-007: VERIFY ACHIEVEMENT (POST /achievements/:id/verify) ======================
// Actor: Dosen Wali
// Precondition: Status 'submitted'
// Flow:
// 1. Dosen review prestasi detail
// 2. Dosen approve prestasi
// 3. Update status menjadi 'verified'
// 4. Set verified_by dan verified_at
// 5. Return updated status
//

func (s *AchievementService) VerifyAchievement(c *fiber.Ctx) error {
	claims := c.Locals("user").(*model.JWTClaims)
	achievementID := c.Params("id")
	
	// Verify user is Dosen Wali
	if claims.Role != "Dosen Wali" {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "only lecturers can verify achievements",
		})
	}
	
	// Get reference
	ref, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}
	
	// Check if lecturer is the advisor of this student
	lecturer, err := s.lecturerRepo.FindByUserID(claims.UserID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "lecturer profile not found",
		})
	}
	
	student, err := s.studentRepo.FindByID(ref.StudentID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "student not found",
		})
	}
	
	if student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "you can only verify achievements of your advisees",
		})
	}
	
	// Check status (must be submitted)
	if ref.Status != "submitted" {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement must be in submitted status to verify",
		})
	}
	
	// Step 2 & 3: Update status to 'verified'
	now := time.Now()
	ref.Status = "verified"
	ref.VerifiedAt = &now
	ref.VerifiedBy = &claims.UserID
	ref.UpdatedAt = now
	ref.RejectionNote = nil // Clear rejection note if any
	
	if err := s.achievementRepo.UpdateReference(ref); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to verify achievement",
		})
	}
	
	// Step 5: Return response
	achievement, _ := s.achievementRepo.GetAchievementByID(ref.MongoAchievementID)
	verifiedByUser, _ := s.userRepo.FindByID(claims.UserID)
	response := s.buildAchievementResponse(ref, achievement, student, verifiedByUser, nil)
	
	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement verified successfully",
		Data:    response,
	})
}

//
// ==================== FR-008: REJECT ACHIEVEMENT (POST /achievements/:id/reject) ======================
// Actor: Dosen Wali
// Precondition: Status 'submitted'
// Flow:
// 1. Dosen input rejection note
// 2. Update status menjadi 'rejected'
// 3. Save rejection_note
// 4. Create notification untuk mahasiswa (skip for now)
// 5. Return updated status
//

func (s *AchievementService) RejectAchievement(c *fiber.Ctx) error {
	claims := c.Locals("user").(*model.JWTClaims)
	achievementID := c.Params("id")
	
	// Verify user is Dosen Wali
	if claims.Role != "Dosen Wali" {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "only lecturers can reject achievements",
		})
	}
	
	// Parse request
	req := new(model.RejectAchievementRequest)
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
			Error:  "rejection note is required",
		})
	}
	
	// Get reference
	ref, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}
	
	// Check if lecturer is the advisor of this student
	lecturer, err := s.lecturerRepo.FindByUserID(claims.UserID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "lecturer profile not found",
		})
	}
	
	student, err := s.studentRepo.FindByID(ref.StudentID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "student not found",
		})
	}
	
	if student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "you can only reject achievements of your advisees",
		})
	}
	
	// Check status (must be submitted)
	if ref.Status != "submitted" {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement must be in submitted status to reject",
		})
	}
	
	// Step 1, 2, 3: Update status to 'rejected' with note
	now := time.Now()
	ref.Status = "rejected"
	ref.RejectionNote = &req.RejectionNote
	ref.VerifiedBy = &claims.UserID
	ref.VerifiedAt = &now
	ref.UpdatedAt = now
	
	if err := s.achievementRepo.UpdateReference(ref); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to reject achievement",
		})
	}
	
	// Step 4: Create notification (TODO: implement notification system)
	
	// Step 5: Return response
	achievement, _ := s.achievementRepo.GetAchievementByID(ref.MongoAchievementID)
	verifiedByUser, _ := s.userRepo.FindByID(claims.UserID)
	response := s.buildAchievementResponse(ref, achievement, student, verifiedByUser, nil)
	
	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement rejected",
		Data:    response,
	})
}

//
// ==================== UPLOAD ATTACHMENT (POST /achievements/:id/attachments) ======================
//

func (s *AchievementService) UploadAttachment(c *fiber.Ctx) error {
	claims := c.Locals("user").(*model.JWTClaims)
	achievementID := c.Params("id")
	
	// Get reference
	ref, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}
	
	// Check access
	if err := s.checkAchievementAccess(claims, ref); err != nil {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  err.Error(),
		})
	}
	
	// Parse request
	req := new(model.UploadAttachmentRequest)
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
	
	// Add attachment
	attachment := model.Attachment{
		FileName: req.FileName,
		FileURL:  req.FileURL,
		FileType: req.FileType,
	}
	
	if err := s.achievementRepo.AddAttachment(ref.MongoAchievementID, attachment); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to add attachment",
		})
	}
	
	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "attachment added successfully",
		Data:    attachment,
	})
}

//
// ==================== HELPER FUNCTIONS ======================
//

func (s *AchievementService) checkAchievementAccess(claims *model.JWTClaims, ref *model.AchievementReference) error {
	switch claims.Role {
	case "Admin":
		return nil // Admin can access all
		
	case "Mahasiswa":
		student, err := s.studentRepo.FindByUserID(claims.UserID)
		if err != nil || student.ID != ref.StudentID {
			return fmt.Errorf("you can only access your own achievements")
		}
		return nil
		
	case "Dosen Wali":
		lecturer, err := s.lecturerRepo.FindByUserID(claims.UserID)
		if err != nil {
			return fmt.Errorf("lecturer profile not found")
		}
		
		student, err := s.studentRepo.FindByID(ref.StudentID)
		if err != nil {
			return fmt.Errorf("student not found")
		}
		
		if student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
			return fmt.Errorf("you can only access achievements of your advisees")
		}
		return nil
		
	default:
		return fmt.Errorf("unauthorized access")
	}
}

func (s *AchievementService) buildAchievementResponse(
	ref *model.AchievementReference,
	achievement *model.Achievement,
	student *model.Student,
	verifiedByUser *model.User,
	user *model.User,
) *model.AchievementResponse {
	response := &model.AchievementResponse{
		ID:              ref.ID,
		StudentID:       ref.StudentID,
		AchievementType: achievement.Type,
		Title:           achievement.Title,
		Description:     achievement.Description,
		Details:         achievement.Details,
		Attachments:     achievement.Attachments,
		Tags:            achievement.Tags,
		Points:          achievement.Points,
		Status:          ref.Status,
		CreatedAt:       ref.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       ref.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	
	// Add student info
	if student != nil {
		response.StudentNIM = student.StudentID
		// Get student user info for name
		if studentUser, err := s.userRepo.FindByID(student.UserID); err == nil {
			response.StudentName = studentUser.FullName
		}
	}
	
	// Add submitted_at
	if ref.SubmittedAt != nil {
		submittedStr := ref.SubmittedAt.Format("2006-01-02 15:04:05")
		response.SubmittedAt = &submittedStr
	}
	
	// Add verified info
	if ref.VerifiedAt != nil {
		verifiedStr := ref.VerifiedAt.Format("2006-01-02 15:04:05")
		response.VerifiedAt = &verifiedStr
	}
	
	if ref.VerifiedBy != nil {
		response.VerifiedBy = ref.VerifiedBy
		if verifiedByUser != nil {
			response.VerifiedByName = &verifiedByUser.FullName
		}
	}
	
	// Add rejection note
	if ref.RejectionNote != nil {
		response.RejectionNote = ref.RejectionNote
	}
	
	return response
}