package mongo_models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attachment struct {
	FileName   string    `bson:"fileName"`
	FileUrl    string    `bson:"fileUrl"`
	FileType   string    `bson:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt"`
}

type Achievement struct {
	ID              primitive.ObjectID      `bson:"_id,omitempty"`
	StudentID       string                  `bson:"studentId"`
	AchievementType string                  `bson:"achievementType"`
	Title           string                  `bson:"title"`
	Description     string                  `bson:"description"`
	Details         map[string]interface{}  `bson:"details"`
	Attachments     []Attachment            `bson:"attachments"`
	Tags            []string                `bson:"tags"`
	Points          int                     `bson:"points"`
	CreatedAt       time.Time               `bson:"createdAt"`
	UpdatedAt       time.Time               `bson:"updatedAt"`
}
