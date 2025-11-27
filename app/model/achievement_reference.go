package models

import "time"

type AchievementReference struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	StudentID  string    `gorm:"type:uuid"`
	Student    Student   `gorm:"foreignKey:StudentID"`
	MongoID    string    `gorm:"column:mongo_achievement_id"`
	Status     string    `gorm:"type:enum('draft','submitted','verified','rejected')"`
	SubmittedAt *time.Time
	VerifiedAt  *time.Time
	VerifiedBy  *string   `gorm:"type:uuid"`
	RejectionNote string
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
