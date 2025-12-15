package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"UASBE/app/model"
	"UASBE/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AchievementRepository interface {
	// PostgreSQL operations (achievement_references)
	CreateReference(ref *model.AchievementReference) error
	GetReferenceByID(id string) (*model.AchievementReference, error)
	GetReferenceByMongoID(mongoID string) (*model.AchievementReference, error)
	UpdateReferenceStatus(id string, status string) error
	UpdateReference(ref *model.AchievementReference) error
	DeleteReference(id string) error
	GetReferences(filters *model.AchievementFilters) ([]model.AchievementReference, error)
	CountReferences(filters *model.AchievementFilters) (int, error)
	
	// MongoDB operations (achievements)
	CreateAchievement(achievement *model.Achievement) (string, error)
	GetAchievementByID(id string) (*model.Achievement, error)
	UpdateAchievement(id string, achievement *model.Achievement) error
	DeleteAchievement(id string) error
	GetAchievements(mongoIDs []string) ([]model.Achievement, error)
	
	// Attachments
	AddAttachment(achievementID string, attachment model.Attachment) error
}

type achievementRepository struct {
	db *sql.DB
}

func NewAchievementRepository(db *sql.DB) AchievementRepository {
	return &achievementRepository{db: db}
}

// ===================== POSTGRESQL OPERATIONS ========================

func (r *achievementRepository) CreateReference(ref *model.AchievementReference) error {
	ref.CreatedAt = time.Now()
	ref.UpdatedAt = time.Now()
	
	query := `
		INSERT INTO achievement_references 
		(id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query,
		ref.ID,
		ref.StudentID,
		ref.MongoAchievementID,
		ref.Status,
		ref.CreatedAt,
		ref.UpdatedAt,
	)
	return err
}

func (r *achievementRepository) GetReferenceByID(id string) (*model.AchievementReference, error) {
	ref := &model.AchievementReference{}
	query := `
		SELECT id, student_id, mongo_achievement_id, status, 
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

func (r *achievementRepository) GetReferenceByMongoID(mongoID string) (*model.AchievementReference, error) {
	ref := &model.AchievementReference{}
	query := `
		SELECT id, student_id, mongo_achievement_id, status, 
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE mongo_achievement_id = $1
	`
	err := r.db.QueryRow(query, mongoID).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

func (r *achievementRepository) UpdateReferenceStatus(id string, status string) error {
	query := `
		UPDATE achievement_references
		SET status = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

func (r *achievementRepository) UpdateReference(ref *model.AchievementReference) error {
	ref.UpdatedAt = time.Now()
	query := `
		UPDATE achievement_references
		SET status = $1, submitted_at = $2, verified_at = $3, 
		    verified_by = $4, rejection_note = $5, updated_at = $6
		WHERE id = $7
	`
	_, err := r.db.Exec(query,
		ref.Status,
		ref.SubmittedAt,
		ref.VerifiedAt,
		ref.VerifiedBy,
		ref.RejectionNote,
		ref.UpdatedAt,
		ref.ID,
	)
	return err
}

func (r *achievementRepository) DeleteReference(id string) error {
	query := `
		UPDATE achievement_references
		SET status = 'deleted', updated_at = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

func (r *achievementRepository) GetReferences(filters *model.AchievementFilters) ([]model.AchievementReference, error) {
	query := `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status,
		       ar.submitted_at, ar.verified_at, ar.verified_by, ar.rejection_note,
		       ar.created_at, ar.updated_at
		FROM achievement_references ar
	`
	
	// Add joins and filters
	var args []interface{}
	argCount := 0
	whereClauses := []string{}
	
	if filters.StudentID != "" {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("ar.student_id = $%d", argCount))
		args = append(args, filters.StudentID)
	}
	
	if filters.Status != "" {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("ar.status = $%d", argCount))
		args = append(args, filters.Status)
	}
	
	// FR-006: Filter by advisor (dosen wali only see their advisees)
	if filters.AdvisorID != "" {
		query += ` INNER JOIN students s ON ar.student_id = s.id`
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("s.advisor_id = $%d", argCount))
		args = append(args, filters.AdvisorID)
	}
	
	// Add WHERE clause
	if len(whereClauses) > 0 {
		query += " WHERE " + whereClauses[0]
		for i := 1; i < len(whereClauses); i++ {
			query += " AND " + whereClauses[i]
		}
	}
	
	// Add pagination
	query += " ORDER BY ar.created_at DESC"
	
	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)
		
		if filters.Page > 0 {
			offset := (filters.Page - 1) * filters.PageSize
			argCount++
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var references []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		)
		if err != nil {
			continue
		}
		references = append(references, ref)
	}
	
	return references, nil
}

func (r *achievementRepository) CountReferences(filters *model.AchievementFilters) (int, error) {
	query := `SELECT COUNT(*) FROM achievement_references ar`
	
	var args []interface{}
	argCount := 0
	whereClauses := []string{}
	
	if filters.StudentID != "" {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("ar.student_id = $%d", argCount))
		args = append(args, filters.StudentID)
	}
	
	if filters.Status != "" {
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("ar.status = $%d", argCount))
		args = append(args, filters.Status)
	}
	
	if filters.AdvisorID != "" {
		query += ` INNER JOIN students s ON ar.student_id = s.id`
		argCount++
		whereClauses = append(whereClauses, fmt.Sprintf("s.advisor_id = $%d", argCount))
		args = append(args, filters.AdvisorID)
	}
	
	if len(whereClauses) > 0 {
		query += " WHERE " + whereClauses[0]
		for i := 1; i < len(whereClauses); i++ {
			query += " AND " + whereClauses[i]
		}
	}
	
	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

// ===================== MONGODB OPERATIONS ========================

func (r *achievementRepository) CreateAchievement(achievement *model.Achievement) (string, error) {
	ctx := context.Background()
	collection := database.GetCollection("achievements")
	
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()
	
	result, err := collection.InsertOne(ctx, achievement)
	if err != nil {
		return "", err
	}
	
	objectID := result.InsertedID.(primitive.ObjectID)
	return objectID.Hex(), nil
}

func (r *achievementRepository) GetAchievementByID(id string) (*model.Achievement, error) {
	ctx := context.Background()
	collection := database.GetCollection("achievements")
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	var achievement model.Achievement
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&achievement)
	if err != nil {
		return nil, err
	}
	
	achievement.ID = objectID.Hex()
	return &achievement, nil
}

func (r *achievementRepository) UpdateAchievement(id string, achievement *model.Achievement) error {
	ctx := context.Background()
	collection := database.GetCollection("achievements")
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	achievement.UpdatedAt = time.Now()
	
	update := bson.M{
		"$set": bson.M{
			"title":       achievement.Title,
			"description": achievement.Description,
			"details":     achievement.Details,
			"tags":        achievement.Tags,
			"points":      achievement.Points,
			"updated_at":  achievement.UpdatedAt,
		},
	}
	
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *achievementRepository) DeleteAchievement(id string) error {
	ctx := context.Background()
	collection := database.GetCollection("achievements")
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *achievementRepository) GetAchievements(mongoIDs []string) ([]model.Achievement, error) {
	ctx := context.Background()
	collection := database.GetCollection("achievements")
	
	// Convert string IDs to ObjectIDs
	objectIDs := make([]primitive.ObjectID, 0, len(mongoIDs))
	for _, id := range mongoIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}
	
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var achievements []model.Achievement
	for cursor.Next(ctx) {
		var achievement model.Achievement
		if err := cursor.Decode(&achievement); err != nil {
			continue
		}
		
		// Set ID from ObjectID
		if oid, ok := cursor.Current.Lookup("_id").ObjectIDOK(); ok {
			achievement.ID = oid.Hex()
		}
		
		achievements = append(achievements, achievement)
	}
	
	return achievements, nil
}

func (r *achievementRepository) AddAttachment(achievementID string, attachment model.Attachment) error {
	ctx := context.Background()
	collection := database.GetCollection("achievements")
	
	objectID, err := primitive.ObjectIDFromHex(achievementID)
	if err != nil {
		return err
	}
	
	attachment.UploadedAt = time.Now()
	
	update := bson.M{
		"$push": bson.M{"attachments": attachment},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update, options.Update().SetUpsert(false))
	return err
}