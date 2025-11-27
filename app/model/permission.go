package models

type Permission struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string `gorm:"unique;not null"`
	Resource    string
	Action      string
	Description string
}
