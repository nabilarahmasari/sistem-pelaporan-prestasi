package repository

import (
	"database/sql"
	"project_uas/app/model"
)

type RoleRepository interface {
	GetRoleByID(id string) (*model.Role, error)
}

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetRoleByID(id string) (*model.Role, error) {
	role := &model.Role{}
	
	query := `SELECT id, name, description, created_at FROM roles WHERE id = $1`
	
	err := r.db.QueryRow(query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return role, nil
}