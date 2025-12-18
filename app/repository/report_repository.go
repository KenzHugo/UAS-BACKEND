package repository

import (
	"context"
	"database/sql"
	"UASBE/app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportRepository interface {
	// Statistics methods
	GetTotalByType(studentID *string, advisorID *string) (map[string]int, error)
	GetTotalByPeriod(studentID *string, advisorID *string) ([]model.PeriodStats, error)
	GetTopStudents(limit int, advisorID *string) ([]model.TopStudent, error)
	GetCompetitionLevelDistribution(studentID *string, advisorID *string) (map[string]int, error)
	GetStatusBreakdown(studentID *string, advisorID *string) (map[string]int, error)
	
	// Student report methods
	GetStudentSummary(studentID string) (*model.StudentSummary, error)
	GetStudentAchievementsByType(studentID string) (map[string]int, error)
	GetStudentAchievementsByStatus(studentID string) (map[string]int, error)
	GetStudentTimeline(studentID string) ([]model.PeriodStats, error)
}

type reportRepository struct {
	pgDB    *sql.DB
	mongoDB *mongo.Database
}

func NewReportRepository(pgDB *sql.DB, mongoDB *mongo.Database) ReportRepository {
	return &reportRepository{
		pgDB:    pgDB,
		mongoDB: mongoDB,
	}
}

//
// ==================== STATISTICS METHODS ======================
//

// GetTotalByType - Hitung total prestasi per tipe
func (r *reportRepository) GetTotalByType(studentID *string, advisorID *string) (map[string]int, error) {
	collection := r.mongoDB.Collection("achievements")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build filter based on scope
	var refQuery string
	var args []interface{}
	
	if studentID != nil {
		// Specific student
		refQuery = `SELECT mongo_achievement_id FROM achievement_references WHERE status = 'verified' AND student_id = $1`
		args = append(args, *studentID)
	} else if advisorID != nil {
		// All students of advisor
		refQuery = `
			SELECT ar.mongo_achievement_id 
			FROM achievement_references ar
			JOIN students s ON ar.student_id = s.id
			WHERE ar.status = 'verified' AND s.advisor_id = $1
		`
		args = append(args, *advisorID)
	} else {
		// All achievements
		refQuery = `SELECT mongo_achievement_id FROM achievement_references WHERE status = 'verified'`
	}

	rows, err := r.pgDB.Query(refQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Count by type from MongoDB
	result := make(map[string]int)
	
	for rows.Next() {
		var mongoID string
		rows.Scan(&mongoID)
		
		// Convert hex string to ObjectID
		objectID, err := primitive.ObjectIDFromHex(mongoID)
		if err != nil {
			// Log error jika ObjectID invalid
			println("ERROR: Invalid ObjectID:", mongoID, err.Error())
			continue
		}
		
		var achievement model.Achievement
		filter := bson.M{"_id": objectID}
		err = collection.FindOne(ctx, filter).Decode(&achievement)
		if err == nil {
			result[achievement.AchievementType]++
			// Debug log
			println("Found achievement:", achievement.Title, "Points:", achievement.Points)
		} else {
			// Log error jika achievement tidak ditemukan
			println("ERROR: Achievement not found in MongoDB:", mongoID, err.Error())
		}
	}

	return result, nil
}

// GetTotalByPeriod - Hitung total prestasi per periode (bulan)
func (r *reportRepository) GetTotalByPeriod(studentID *string, advisorID *string) ([]model.PeriodStats, error) {
	var query string
	var args []interface{}

	if studentID != nil {
		query = `
			SELECT 
				TO_CHAR(verified_at, 'YYYY-MM') as period,
				COUNT(*) as count
			FROM achievement_references
			WHERE status = 'verified' AND student_id = $1 AND verified_at IS NOT NULL
			GROUP BY TO_CHAR(verified_at, 'YYYY-MM')
			ORDER BY period DESC
			LIMIT 12
		`
		args = append(args, *studentID)
	} else if advisorID != nil {
		query = `
			SELECT 
				TO_CHAR(ar.verified_at, 'YYYY-MM') as period,
				COUNT(*) as count
			FROM achievement_references ar
			JOIN students s ON ar.student_id = s.id
			WHERE ar.status = 'verified' AND s.advisor_id = $1 AND ar.verified_at IS NOT NULL
			GROUP BY TO_CHAR(ar.verified_at, 'YYYY-MM')
			ORDER BY period DESC
			LIMIT 12
		`
		args = append(args, *advisorID)
	} else {
		query = `
			SELECT 
				TO_CHAR(verified_at, 'YYYY-MM') as period,
				COUNT(*) as count
			FROM achievement_references
			WHERE status = 'verified' AND verified_at IS NOT NULL
			GROUP BY TO_CHAR(verified_at, 'YYYY-MM')
			ORDER BY period DESC
			LIMIT 12
		`
	}

	rows, err := r.pgDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []model.PeriodStats
	for rows.Next() {
		var period sql.NullString
		var count int
		
		rows.Scan(&period, &count)
		
		if period.Valid {
			stats = append(stats, model.PeriodStats{
				Period: period.String,
				Count:  count,
			})
		}
	}

	return stats, nil
}

// GetTopStudents - Dapatkan top mahasiswa berprestasi
func (r *reportRepository) GetTopStudents(limit int, advisorID *string) ([]model.TopStudent, error) {
	var query string
	var args []interface{}

	if advisorID != nil {
		query = `
			SELECT 
				s.id,
				s.student_id,
				u.full_name,
				s.program_study,
				COUNT(ar.id) as achievement_count
			FROM students s
			JOIN users u ON s.id = u.id
			JOIN achievement_references ar ON s.id = ar.student_id
			WHERE ar.status = 'verified' AND s.advisor_id = $1
			GROUP BY s.id, s.student_id, u.full_name, s.program_study
			ORDER BY achievement_count DESC
			LIMIT $2
		`
		args = append(args, *advisorID, limit)
	} else {
		query = `
			SELECT 
				s.id,
				s.student_id,
				u.full_name,
				s.program_study,
				COUNT(ar.id) as achievement_count
			FROM students s
			JOIN users u ON s.id = u.id
			JOIN achievement_references ar ON s.id = ar.student_id
			WHERE ar.status = 'verified'
			GROUP BY s.id, s.student_id, u.full_name, s.program_study
			ORDER BY achievement_count DESC
			LIMIT $1
		`
		args = append(args, limit)
	}

	rows, err := r.pgDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topStudents []model.TopStudent
	collection := r.mongoDB.Collection("achievements")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for rows.Next() {
		var student model.TopStudent
		
		rows.Scan(
			&student.StudentID,
			&student.StudentNIM,
			&student.StudentName,
			&student.ProgramStudy,
			&student.AchievementCount,
		)

		// Calculate total points from MongoDB
		achievementQuery := `SELECT mongo_achievement_id FROM achievement_references WHERE student_id = $1 AND status = 'verified'`
		achRows, _ := r.pgDB.Query(achievementQuery, student.StudentID)
		
		totalPoints := 0
		for achRows.Next() {
			var mongoID string
			achRows.Scan(&mongoID)
			
			objectID, err := primitive.ObjectIDFromHex(mongoID)
			if err != nil {
				continue
			}
			
			var achievement model.Achievement
			filter := bson.M{"_id": objectID}
			err = collection.FindOne(ctx, filter).Decode(&achievement)
			if err == nil {
				totalPoints += achievement.Points
			}
		}
		achRows.Close()
		
		student.TotalPoints = totalPoints
		topStudents = append(topStudents, student)
	}

	return topStudents, nil
}

// GetCompetitionLevelDistribution - Distribusi tingkat kompetisi
func (r *reportRepository) GetCompetitionLevelDistribution(studentID *string, advisorID *string) (map[string]int, error) {
	collection := r.mongoDB.Collection("achievements")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get verified achievement mongo IDs
	var query string
	var args []interface{}

	if studentID != nil {
		query = `SELECT mongo_achievement_id FROM achievement_references WHERE status = 'verified' AND student_id = $1`
		args = append(args, *studentID)
	} else if advisorID != nil {
		query = `
			SELECT ar.mongo_achievement_id 
			FROM achievement_references ar
			JOIN students s ON ar.student_id = s.id
			WHERE ar.status = 'verified' AND s.advisor_id = $1
		`
		args = append(args, *advisorID)
	} else {
		query = `SELECT mongo_achievement_id FROM achievement_references WHERE status = 'verified'`
	}

	rows, err := r.pgDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	distribution := make(map[string]int)
	
	for rows.Next() {
		var mongoID string
		rows.Scan(&mongoID)
		
		objectID, err := primitive.ObjectIDFromHex(mongoID)
		if err != nil {
			continue
		}
		
		var achievement model.Achievement
		filter := bson.M{"_id": objectID}
		err = collection.FindOne(ctx, filter).Decode(&achievement)
		
		if err == nil && achievement.AchievementType == "competition" {
			// Get competition level from details
			if level, ok := achievement.Details["competitionLevel"].(string); ok {
				distribution[level]++
			}
		}
	}

	return distribution, nil
}

// GetStatusBreakdown - Breakdown by status
func (r *reportRepository) GetStatusBreakdown(studentID *string, advisorID *string) (map[string]int, error) {
	var query string
	var args []interface{}

	if studentID != nil {
		query = `
			SELECT status, COUNT(*) as count
			FROM achievement_references
			WHERE student_id = $1 AND status != 'deleted'
			GROUP BY status
		`
		args = append(args, *studentID)
	} else if advisorID != nil {
		query = `
			SELECT ar.status, COUNT(*) as count
			FROM achievement_references ar
			JOIN students s ON ar.student_id = s.id
			WHERE s.advisor_id = $1 AND ar.status != 'deleted'
			GROUP BY ar.status
		`
		args = append(args, *advisorID)
	} else {
		query = `
			SELECT status, COUNT(*) as count
			FROM achievement_references
			WHERE status != 'deleted'
			GROUP BY status
		`
	}

	rows, err := r.pgDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	breakdown := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		rows.Scan(&status, &count)
		breakdown[status] = count
	}

	return breakdown, nil
}

//
// ==================== STUDENT REPORT METHODS ======================
//

// GetStudentSummary - Summary statistics untuk satu student
func (r *reportRepository) GetStudentSummary(studentID string) (*model.StudentSummary, error) {
	summary := &model.StudentSummary{}

	// Count by status
	query := `
		SELECT 
			COUNT(*) FILTER (WHERE status = 'verified') as verified,
			COUNT(*) FILTER (WHERE status = 'submitted') as pending,
			COUNT(*) FILTER (WHERE status = 'draft') as draft,
			COUNT(*) FILTER (WHERE status = 'rejected') as rejected,
			COUNT(*) as total
		FROM achievement_references
		WHERE student_id = $1 AND status != 'deleted'
	`

	err := r.pgDB.QueryRow(query, studentID).Scan(
		&summary.VerifiedCount,
		&summary.PendingCount,
		&summary.DraftCount,
		&summary.RejectedCount,
		&summary.TotalAchievements,
	)
	if err != nil {
		return nil, err
	}

	// Calculate total points from MongoDB
	collection := r.mongoDB.Collection("achievements")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoQuery := `SELECT mongo_achievement_id FROM achievement_references WHERE student_id = $1 AND status = 'verified'`
	rows, err := r.pgDB.Query(mongoQuery, studentID)
	if err != nil {
		return summary, nil // Return summary tanpa points
	}
	defer rows.Close()

	totalPoints := 0
	for rows.Next() {
		var mongoID string
		rows.Scan(&mongoID)
		
		objectID, err := primitive.ObjectIDFromHex(mongoID)
		if err != nil {
			continue
		}
		
		var achievement model.Achievement
		filter := bson.M{"_id": objectID}
		err = collection.FindOne(ctx, filter).Decode(&achievement)
		if err == nil {
			totalPoints += achievement.Points
		}
	}

	summary.TotalPoints = totalPoints

	return summary, nil
}

// GetStudentAchievementsByType - Count by type
func (r *reportRepository) GetStudentAchievementsByType(studentID string) (map[string]int, error) {
	return r.GetTotalByType(&studentID, nil)
}

// GetStudentAchievementsByStatus - Count by status
func (r *reportRepository) GetStudentAchievementsByStatus(studentID string) (map[string]int, error) {
	return r.GetStatusBreakdown(&studentID, nil)
}

// GetStudentTimeline - Monthly timeline
func (r *reportRepository) GetStudentTimeline(studentID string) ([]model.PeriodStats, error) {
	return r.GetTotalByPeriod(&studentID, nil)
}