package model

import "time"

// ===================== ACHIEVEMENT ENTITY (PostgreSQL) ========================
// Tabel achievement_references - hanya menyimpan reference dan status

type AchievementReference struct {
	ID                 string    `json:"id" db:"id"`
	StudentID          string    `json:"student_id" db:"student_id"`
	MongoAchievementID string    `json:"mongo_achievement_id" db:"mongo_achievement_id"`
	Status             string    `json:"status" db:"status"` // draft, submitted, verified, rejected, deleted
	SubmittedAt        *time.Time `json:"submitted_at,omitempty" db:"submitted_at"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	VerifiedBy         *string    `json:"verified_by,omitempty" db:"verified_by"`
	RejectionNote      *string    `json:"rejection_note,omitempty" db:"rejection_note"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// ===================== ACHIEVEMENT ENTITY (MongoDB) ========================
// Collection achievements - menyimpan data prestasi lengkap dengan field dinamis

type Achievement struct {
	ID          string                 `json:"_id,omitempty" bson:"_id,omitempty"`
	StudentID   string                 `json:"student_id" bson:"student_id"`
	Type        string                 `json:"achievement_type" bson:"achievement_type"` // academic, competition, organization, publication, certification, other
	Title       string                 `json:"title" bson:"title"`
	Description string                 `json:"description" bson:"description"`
	Details     map[string]interface{} `json:"details" bson:"details"` // Field dinamis
	Attachments []Attachment           `json:"attachments,omitempty" bson:"attachments,omitempty"`
	Tags        []string               `json:"tags,omitempty" bson:"tags,omitempty"`
	Points      int                    `json:"points,omitempty" bson:"points,omitempty"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
}

type Attachment struct {
	FileName   string    `json:"file_name" bson:"file_name"`
	FileURL    string    `json:"file_url" bson:"file_url"`
	FileType   string    `json:"file_type" bson:"file_type"`
	UploadedAt time.Time `json:"uploaded_at" bson:"uploaded_at"`
}

// ===================== CREATE ACHIEVEMENT REQUEST ========================
// FR-003: Submit Prestasi

type CreateAchievementRequest struct {
	AchievementType string                 `json:"achievement_type" validate:"required"`
	Title           string                 `json:"title" validate:"required"`
	Description     string                 `json:"description" validate:"required"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags,omitempty"`
	Points          int                    `json:"points,omitempty"`
}

// ===================== UPDATE ACHIEVEMENT REQUEST ========================

type UpdateAchievementRequest struct {
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Points      int                    `json:"points,omitempty"`
}

// ===================== SUBMIT FOR VERIFICATION REQUEST ========================
// FR-004: Submit untuk Verifikasi

type SubmitForVerificationRequest struct {
	// No body needed, just change status
}

// ===================== VERIFY ACHIEVEMENT REQUEST ========================
// FR-007: Verify Prestasi

type VerifyAchievementRequest struct {
	// No additional data needed
}

// ===================== REJECT ACHIEVEMENT REQUEST ========================
// FR-008: Reject Prestasi

type RejectAchievementRequest struct {
	RejectionNote string `json:"rejection_note" validate:"required"`
}

// ===================== ACHIEVEMENT RESPONSE ========================
// Menggabungkan data dari PostgreSQL dan MongoDB

type AchievementResponse struct {
	ID                 string                 `json:"id"`
	StudentID          string                 `json:"student_id"`
	StudentName        string                 `json:"student_name,omitempty"`
	StudentNIM         string                 `json:"student_nim,omitempty"`
	AchievementType    string                 `json:"achievement_type"`
	Title              string                 `json:"title"`
	Description        string                 `json:"description"`
	Details            map[string]interface{} `json:"details"`
	Attachments        []Attachment           `json:"attachments,omitempty"`
	Tags               []string               `json:"tags,omitempty"`
	Points             int                    `json:"points"`
	Status             string                 `json:"status"`
	SubmittedAt        *string                `json:"submitted_at,omitempty"`
	VerifiedAt         *string                `json:"verified_at,omitempty"`
	VerifiedBy         *string                `json:"verified_by,omitempty"`
	VerifiedByName     *string                `json:"verified_by_name,omitempty"`
	RejectionNote      *string                `json:"rejection_note,omitempty"`
	CreatedAt          string                 `json:"created_at"`
	UpdatedAt          string                 `json:"updated_at"`
}

// ===================== ACHIEVEMENT LIST RESPONSE ========================

type AchievementListResponse struct {
	Achievements []AchievementResponse `json:"achievements"`
	Total        int                   `json:"total"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
	TotalPages   int                   `json:"total_pages"`
}

// ===================== ACHIEVEMENT FILTERS ========================

type AchievementFilters struct {
	StudentID       string   // Filter by student
	Status          string   // Filter by status
	AchievementType string   // Filter by type
	AdvisorID       string   // For dosen wali - filter mahasiswa bimbingan
	Page            int
	PageSize        int
}

// ===================== UPLOAD ATTACHMENT REQUEST ========================

type UploadAttachmentRequest struct {
	FileName string `json:"file_name" validate:"required"`
	FileURL  string `json:"file_url" validate:"required"`
	FileType string `json:"file_type" validate:"required"`
}