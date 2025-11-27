package models

import "time"

type Student struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       string    `gorm:"type:uuid;unique"`
	User         User      `gorm:"foreignKey:UserID"`
	StudentID    string    `gorm:"unique;not null"`
	ProgramStudy string
	AcademicYear string
	AdvisorID    string `gorm:"type:uuid"`
	Advisor      Lecturer `gorm:"foreignKey:AdvisorID"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
