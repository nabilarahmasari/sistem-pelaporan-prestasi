package repository

import (
	"database/sql"
	"project_uas/app/model"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	// Auth methods (sudah ada)
	FindByUsername(username string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	GetRoleName(roleID string) (string, error)
	GetPermissions(role string) ([]string, error)

	// User management methods (BARU)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id string) error
	FindByEmail(email string) (*model.User, error)
	GetAll(limit, offset int, roleName string) ([]model.User, error)
	CountAll(roleName string) (int, error)
	UpdateRole(userID string, roleID string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

//
// ==================== AUTH METHODS (EXISTING) ======================
//

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	user := model.User{}
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE username = $1 LIMIT 1
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
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	user := model.User{}
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 LIMIT 1
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
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetRoleName(roleID string) (string, error) {
	var role string
	err := r.db.QueryRow("SELECT name FROM roles WHERE id=$1", roleID).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}

func (r *userRepository) GetPermissions(role string) ([]string, error) {
	rows, err := r.db.Query(`SELECT permission FROM role_permissions WHERE role=$1`, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var p string
		rows.Scan(&p)
		permissions = append(permissions, p)
	}
	return permissions, nil
}

//
// ==================== USER MANAGEMENT METHODS (NEW) ======================
//

// Create - Insert user baru
func (r *userRepository) Create(user *model.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.RoleID,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

// Update - Update user data
func (r *userRepository) Update(user *model.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET email = $1, full_name = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(query,
		user.Email,
		user.FullName,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)
	return err
}

// Delete - Hapus user (hard delete)
func (r *userRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// FindByEmail - Cari user berdasarkan email
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
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
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAll - Ambil semua users dengan pagination dan filter role
func (r *userRepository) GetAll(limit, offset int, roleName string) ([]model.User, error) {
	var query string
	var rows *sql.Rows
	var err error

	if roleName != "" {
		query = `
			SELECT u.id, u.username, u.email, u.password_hash, u.full_name, u.role_id, u.is_active, u.created_at, u.updated_at
			FROM users u
			JOIN roles r ON u.role_id = r.id
			WHERE r.name = $1
			ORDER BY u.created_at DESC
			LIMIT $2 OFFSET $3
		`
		rows, err = r.db.Query(query, roleName, limit, offset)
	} else {
		query = `
			SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
			FROM users
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		rows, err = r.db.Query(query, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.PasswordHash,
			&u.FullName,
			&u.RoleID,
			&u.IsActive,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

// CountAll - Hitung total users dengan filter role
func (r *userRepository) CountAll(roleName string) (int, error) {
	var count int
	var query string

	if roleName != "" {
		query = `
			SELECT COUNT(*)
			FROM users u
			JOIN roles r ON u.role_id = r.id
			WHERE r.name = $1
		`
		err := r.db.QueryRow(query, roleName).Scan(&count)
		return count, err
	}

	query = `SELECT COUNT(*) FROM users`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// UpdateRole - Update role user
func (r *userRepository) UpdateRole(userID string, roleID string) error {
	query := `
		UPDATE users
		SET role_id = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(query, roleID, time.Now(), userID)
	return err
}