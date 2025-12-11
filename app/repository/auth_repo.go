package repository

import (
	"database/sql"
	"errors"
	"UASBE/model"
)

type AuthRepository interface {
	FindUserByUsername(username string) (*model.User, error)
	FindUserByEmail(email string) (*model.User, error)
	FindUserByID(id string) (*model.User, error)
	GetUserRole(userID string) (*model.Role, error)
	GetUserPermissions(roleID string) ([]string, error)
	GetUserWithRole(userID string) (*model.UserWithRole, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindUserByID(id string) (*model.User, error) {
	var user model.User
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserRole(userID string) (*model.Role, error) {
	var role model.Role
	query := `
		SELECT r.id, r.name, r.description, r.created_at
		FROM roles r
		INNER JOIN users u ON u.role_id = r.id
		WHERE u.id = $1
	`
	err := r.db.QueryRow(query, userID).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("role not found")
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *authRepository) GetUserPermissions(roleID string) ([]string, error) {
	query := `
		SELECT p.name
		FROM permissions p
		INNER JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
	`
	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (r *authRepository) GetUserWithRole(userID string) (*model.UserWithRole, error) {
	var userWithRole model.UserWithRole
	query := `
		SELECT 
			u.id, u.username, u.email, u.password_hash, u.full_name, 
			u.role_id, u.is_active, u.created_at, u.updated_at,
			COALESCE(r.name, '') as role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`
	err := r.db.QueryRow(query, userID).Scan(
		&userWithRole.ID,
		&userWithRole.Username,
		&userWithRole.Email,
		&userWithRole.PasswordHash,
		&userWithRole.FullName,
		&userWithRole.RoleID,
		&userWithRole.IsActive,
		&userWithRole.CreatedAt,
		&userWithRole.UpdatedAt,
		&userWithRole.RoleName,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	// Get permissions
	if userWithRole.RoleID != nil {
		permissions, err := r.GetUserPermissions(*userWithRole.RoleID)
		if err != nil {
			return nil, err
		}
		userWithRole.Permissions = permissions
	}

	return &userWithRole, nil
}