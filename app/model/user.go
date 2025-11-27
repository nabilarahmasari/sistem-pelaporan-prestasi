package models

import "time"

type User struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username     string    `gorm:"size:50;unique;not null"`
	Email        string    `gorm:"size:100;unique;not null"`
	PasswordHash string    `gorm:"not null"`
	FullName     string    `gorm:"size:100;not null"`
	RoleID       string    `gorm:"type:uuid"`
	Role         Role      `gorm:"foreignKey:RoleID"`
	IsActive     bool      `gorm:"default:true"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
