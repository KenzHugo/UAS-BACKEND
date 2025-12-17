package model

import "time"

// ===================== STUDENT ENTITY ========================
// Representasi tabel "students" di database
// PENTING: id di students = user_id (one-to-one dengan users table)
// Tidak ada kolom user_id terpisah di migration!

type Student struct {
	ID           string    `json:"id" db:"id"`                         // Primary key & Foreign key ke users.id
	StudentID    string    `json:"student_id" db:"student_id"`         // NIM
	ProgramStudy string    `json:"program_study" db:"program_study"`   
	AcademicYear int       `json:"academic_year" db:"academic_year"`   // INT di database
	AdvisorID    *string   `json:"advisor_id" db:"advisor_id"`         // Nullable, FK ke lecturers.id
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ===================== STUDENT PROFILE DTO ====================
// Dipakai saat create/update user dengan role mahasiswa

type StudentProfileRequest struct {
	StudentID    string  `json:"student_id" validate:"required"`
	ProgramStudy string  `json:"program_study" validate:"required"`
	AcademicYear int     `json:"academic_year" validate:"required"`    // INT untuk konsistensi
	AdvisorID    *string `json:"advisor_id,omitempty"`                 // Optional saat create
}

// ===================== STUDENT RESPONSE =======================
// Response untuk menampilkan data student

type StudentResponse struct {
	ID           string  `json:"id"`
	StudentID    string  `json:"student_id"`
	ProgramStudy string  `json:"program_study"`
	AcademicYear int     `json:"academic_year"`                        // INT
	AdvisorID    *string `json:"advisor_id,omitempty"`
	CreatedAt    string  `json:"created_at"`
}

// ===================== SET ADVISOR REQUEST ====================
// Untuk endpoint PUT /students/:id/advisor

type SetAdvisorRequest struct {
	AdvisorID string `json:"advisor_id" validate:"required"`
}