package models

import "time"

type Role struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"unique;not null"`
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
