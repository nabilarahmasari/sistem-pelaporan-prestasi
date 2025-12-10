package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ===================== ACHIEVEMENT ENTITY (MONGODB) ========================
// Collection: achievements
// Berisi data prestasi dinamis dengan field yang fleksibel

type Achievement struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	StudentID       string                 `bson:"studentId" json:"student_id"`
	AchievementType string                 `bson:"achievementType" json:"achievement_type"` // 'academic', 'competition', 'organization', 'publication', 'certification', 'other'
	Title           string                 `bson:"title" json:"title"`
	Description     string                 `bson:"description" json:"description"`
	Details         map[string]interface{} `bson:"details" json:"details"` // Field dinamis
	Attachments     []Attachment           `bson:"attachments" json:"attachments"`
	Tags            []string               `bson:"tags" json:"tags"`
	Points          int                    `bson:"points" json:"points"`
	CreatedAt       time.Time              `bson:"createdAt" json:"created_at"`
	UpdatedAt       time.Time              `bson:"updatedAt" json:"updated_at"`
}

type Attachment struct {
	FileName   string    `bson:"fileName" json:"file_name"`
	FileURL    string    `bson:"fileUrl" json:"file_url"`
	FileType   string    `bson:"fileType" json:"file_type"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploaded_at"`
}

// ===================== ACHIEVEMENT REFERENCE (POSTGRESQL) ========================
// Tabel: achievement_references
// Link antara student dan achievement di MongoDB + status workflow

type AchievementReference struct {
	ID                 string     `json:"id" db:"id"`
	StudentID          string     `json:"student_id" db:"student_id"`
	MongoAchievementID string     `json:"mongo_achievement_id" db:"mongo_achievement_id"`
	Status             string     `json:"status" db:"status"` // 'draft', 'submitted', 'verified', 'rejected', 'deleted'
	SubmittedAt        *time.Time `json:"submitted_at,omitempty" db:"submitted_at"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	VerifiedBy         *string    `json:"verified_by,omitempty" db:"verified_by"`
	RejectionNote      *string    `json:"rejection_note,omitempty" db:"rejection_note"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// ===================== CREATE ACHIEVEMENT REQUEST ========================

type AchievementCreateRequest struct {
	AchievementType string                 `json:"achievement_type" validate:"required"`
	Title           string                 `json:"title" validate:"required"`
	Description     string                 `json:"description" validate:"required"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags"`
	Points          int                    `json:"points"`
}

// ===================== UPDATE ACHIEVEMENT REQUEST ========================

type AchievementUpdateRequest struct {
	AchievementType string                 `json:"achievement_type,omitempty"`
	Title           string                 `json:"title,omitempty"`
	Description     string                 `json:"description,omitempty"`
	Details         map[string]interface{} `json:"details,omitempty"`
	Tags            []string               `json:"tags,omitempty"`
	Points          int                    `json:"points,omitempty"`
}

// ===================== VERIFY/REJECT REQUEST ========================

type VerifyAchievementRequest struct {
	// Kosong, cukup POST saja
}

type RejectAchievementRequest struct {
	RejectionNote string `json:"rejection_note" validate:"required"`
}

// ===================== ACHIEVEMENT RESPONSE ========================

type AchievementResponse struct {
	ID              string                 `json:"id"`
	StudentID       string                 `json:"student_id"`
	AchievementType string                 `json:"achievement_type"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Attachments     []Attachment           `json:"attachments"`
	Tags            []string               `json:"tags"`
	Points          int                    `json:"points"`
	Status          string                 `json:"status"`
	SubmittedAt     *string                `json:"submitted_at,omitempty"`
	VerifiedAt      *string                `json:"verified_at,omitempty"`
	VerifiedBy      *string                `json:"verified_by,omitempty"`
	RejectionNote   *string                `json:"rejection_note,omitempty"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
}

// ===================== ACHIEVEMENT LIST RESPONSE ========================

type AchievementListResponse struct {
	Achievements []AchievementResponse `json:"achievements"`
	Total        int                   `json:"total"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
	TotalPages   int                   `json:"total_pages"`
}

// ===================== UPLOAD ATTACHMENT REQUEST ========================

type UploadAttachmentRequest struct {
	FileName string `json:"file_name" validate:"required"`
	FileURL  string `json:"file_url" validate:"required"`
	FileType string `json:"file_type" validate:"required"`
}
