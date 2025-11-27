package models

import "time"

type Lecturer struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     string    `gorm:"type:uuid;unique"`
	User       User      `gorm:"foreignKey:UserID"`
	LecturerID string    `gorm:"unique;not null"`
	Department string
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
