package repository

import (
	"database/sql"
	"project_uas/app/model"
)

type UserRepository interface {
	FindByUsername(username string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	GetRoleName(roleID string) (string, error)
	GetPermissions(role string) ([]string, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

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
