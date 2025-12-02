package repository

import (
	"database/sql"
)

type PermissionRepository interface {
	GetPermissionsByRoleID(roleID string) ([]string, error)
}

type permissionRepository struct {
	db *sql.DB
}

func NewPermissionRepository(db *sql.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) GetPermissionsByRoleID(roleID string) ([]string, error) {
	query := `
		SELECT p.name
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`
	
	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var permissions []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}
	
	return permissions, nil
}