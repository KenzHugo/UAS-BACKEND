package model

// ===================== STATISTICS RESPONSE ========================
// FR-011: Achievement Statistics
// Output statistik prestasi untuk berbagai actor

type AchievementStatistics struct {
	// Total prestasi per tipe
	TotalByType map[string]int `json:"total_by_type"`
	
	// Total prestasi per periode (tahun-bulan)
	TotalByPeriod []PeriodStats `json:"total_by_period"`
	
	// Top mahasiswa berprestasi
	TopStudents []TopStudent `json:"top_students"`
	
	// Distribusi tingkat kompetisi
	CompetitionLevelDistribution map[string]int `json:"competition_level_distribution"`
	
	// Total keseluruhan
	TotalAchievements int `json:"total_achievements"`
	
	// Breakdown by status
	StatusBreakdown map[string]int `json:"status_breakdown"`
}

type PeriodStats struct {
	Period string `json:"period"` // Format: "2025-01", "2025-02"
	Count  int    `json:"count"`
}

type TopStudent struct {
	StudentID       string `json:"student_id"`
	StudentNIM      string `json:"student_nim"`
	StudentName     string `json:"student_name"`
	ProgramStudy    string `json:"program_study"`
	AchievementCount int   `json:"achievement_count"`
	TotalPoints     int    `json:"total_points"`
}

// ===================== STUDENT REPORT RESPONSE ========================
// GET /api/v1/reports/student/:id
// Detail report untuk satu mahasiswa

type StudentReport struct {
	Student StudentInfo `json:"student"`
	
	// Summary statistics
	Summary StudentSummary `json:"summary"`
	
	// Achievement breakdown
	AchievementsByType map[string]int `json:"achievements_by_type"`
	AchievementsByStatus map[string]int `json:"achievements_by_status"`
	
	// Recent achievements (last 10)
	RecentAchievements []AchievementResponse `json:"recent_achievements"`
	
	// Timeline (monthly)
	Timeline []PeriodStats `json:"timeline"`
}

type StudentInfo struct {
	ID           string  `json:"id"`
	StudentID    string  `json:"student_id"`
	FullName     string  `json:"full_name"`
	Email        string  `json:"email"`
	ProgramStudy string  `json:"program_study"`
	AcademicYear int     `json:"academic_year"`
	AdvisorName  *string `json:"advisor_name,omitempty"`
}

type StudentSummary struct {
	TotalAchievements int `json:"total_achievements"`
	VerifiedCount     int `json:"verified_count"`
	PendingCount      int `json:"pending_count"`
	DraftCount        int `json:"draft_count"`
	RejectedCount     int `json:"rejected_count"`
	TotalPoints       int `json:"total_points"`
}
