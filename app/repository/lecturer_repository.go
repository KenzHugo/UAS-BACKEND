package repository

import (
	"database/sql"
	"UASBE/app/model"
	"time"
)

type LecturerRepository interface {
	Create(lecturer *model.Lecturer) error
	FindByUserID(userID string) (*model.Lecturer, error)
	FindByID(id string) (*model.Lecturer, error)
	FindByLecturerID(lecturerID string) (*model.Lecturer, error)
	Update(lecturer *model.Lecturer) error
	Delete(id string) error
	GetAll(limit, offset int) ([]model.Lecturer, error)
	CountAll() (int, error)
}

type lecturerRepository struct {
	db *sql.DB
}

func NewLecturerRepository(db *sql.DB) LecturerRepository {
	return &lecturerRepository{db}
}

// Create - Insert lecturer baru
// PENTING: id = user_id (sesuai migration: id REFERENCES users(id))
// Tidak ada kolom user_id terpisah!
func (r *lecturerRepository) Create(lecturer *model.Lecturer) error {
	// ID sudah di-set dari luar (= user.id saat create user)
	lecturer.CreatedAt = time.Now()

	query := `
		INSERT INTO lecturers (id, lecturer_id, department, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query,
		lecturer.ID,         // id = user_id
		lecturer.LecturerID, // lecturer_id (NIP)
		lecturer.Department,
		lecturer.CreatedAt,
	)
	return err
}

// FindByUserID - Cari lecturer berdasarkan user_id
// Karena id = user_id, maka sama dengan FindByID
func (r *lecturerRepository) FindByUserID(userID string) (*model.Lecturer, error) {
	return r.FindByID(userID)
}

// FindByID - Cari lecturer berdasarkan ID
func (r *lecturerRepository) FindByID(id string) (*model.Lecturer, error) {
	lecturer := &model.Lecturer{}
	query := `
		SELECT id, lecturer_id, department, created_at
		FROM lecturers
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&lecturer.ID,
		&lecturer.LecturerID,
		&lecturer.Department,
		&lecturer.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return lecturer, nil
}

// FindByLecturerID - Cari lecturer berdasarkan lecturer_id (NIP)
func (r *lecturerRepository) FindByLecturerID(lecturerID string) (*model.Lecturer, error) {
	lecturer := &model.Lecturer{}
	query := `
		SELECT id, lecturer_id, department, created_at
		FROM lecturers
		WHERE lecturer_id = $1
	`
	err := r.db.QueryRow(query, lecturerID).Scan(
		&lecturer.ID,
		&lecturer.LecturerID,
		&lecturer.Department,
		&lecturer.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return lecturer, nil
}

// Update - Update data lecturer
func (r *lecturerRepository) Update(lecturer *model.Lecturer) error {
	query := `
		UPDATE lecturers
		SET department = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(query, lecturer.Department, lecturer.ID)
	return err
}

// Delete - Hapus lecturer
func (r *lecturerRepository) Delete(id string) error {
	query := `DELETE FROM lecturers WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// GetAll - Ambil semua lecturers dengan pagination
func (r *lecturerRepository) GetAll(limit, offset int) ([]model.Lecturer, error) {
	query := `
		SELECT id, lecturer_id, department, created_at
		FROM lecturers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lecturers []model.Lecturer
	for rows.Next() {
		var l model.Lecturer
		err := rows.Scan(
			&l.ID,
			&l.LecturerID,
			&l.Department,
			&l.CreatedAt,
		)
		if err != nil {
			continue
		}
		lecturers = append(lecturers, l)
	}
	return lecturers, nil
}

// CountAll - Hitung total lecturers
func (r *lecturerRepository) CountAll() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM lecturers`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}