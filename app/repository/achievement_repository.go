package repository

import (
	"context"
	"database/sql"
	"project_uas/app/model"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AchievementRepository interface {
	// PostgreSQL - Achievement References
	CreateReference(ref *model.AchievementReference) error
	UpdateReference(ref *model.AchievementReference) error
	GetReferenceByID(id string) (*model.AchievementReference, error)
	GetReferenceByMongoID(mongoID string) (*model.AchievementReference, error)
	GetReferencesByStudentID(studentID string, status string, limit, offset int) ([]model.AchievementReference, error)
	CountReferencesByStudentID(studentID string, status string) (int, error)
	GetReferencesByAdvisorID(advisorID string, status string, limit, offset int) ([]model.AchievementReference, error)
	CountReferencesByAdvisorID(advisorID string, status string) (int, error)
	GetAllReferences(status string, limit, offset int) ([]model.AchievementReference, error)
	CountAllReferences(status string) (int, error)

	// MongoDB - Achievements
	CreateAchievement(achievement *model.Achievement) (string, error)
	UpdateAchievement(id string, achievement *model.Achievement) error
	GetAchievementByID(id string) (*model.Achievement, error)
	DeleteAchievement(id string) error
	AddAttachment(achievementID string, attachment model.Attachment) error
}

type achievementRepository struct {
	pgDB    *sql.DB
	mongoDB *mongo.Database
}

func NewAchievementRepository(pgDB *sql.DB, mongoDB *mongo.Database) AchievementRepository {
	return &achievementRepository{
		pgDB:    pgDB,
		mongoDB: mongoDB,
	}
}

//
// ==================== POSTGRESQL METHODS (REFERENCES) ======================
//

// CreateReference - Insert reference baru
func (r *achievementRepository) CreateReference(ref *model.AchievementReference) error {
	ref.ID = uuid.New().String()
	ref.CreatedAt = time.Now()
	ref.UpdatedAt = time.Now()

	query := `
		INSERT INTO achievement_references 
		(id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pgDB.Exec(query,
		ref.ID,
		ref.StudentID,
		ref.MongoAchievementID,
		ref.Status,
		ref.SubmittedAt,
		ref.VerifiedAt,
		ref.VerifiedBy,
		ref.RejectionNote,
		ref.CreatedAt,
		ref.UpdatedAt,
	)
	return err
}

// UpdateReference - Update reference
func (r *achievementRepository) UpdateReference(ref *model.AchievementReference) error {
	ref.UpdatedAt = time.Now()

	query := `
		UPDATE achievement_references
		SET status = $1, submitted_at = $2, verified_at = $3, verified_by = $4, rejection_note = $5, updated_at = $6
		WHERE id = $7
	`
	_, err := r.pgDB.Exec(query,
		ref.Status,
		ref.SubmittedAt,
		ref.VerifiedAt,
		ref.VerifiedBy,
		ref.RejectionNote,
		ref.UpdatedAt,
		ref.ID,
	)
	return err
}

// GetReferenceByID - Get reference by ID
func (r *achievementRepository) GetReferenceByID(id string) (*model.AchievementReference, error) {
	ref := &model.AchievementReference{}
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE id = $1
	`
	err := r.pgDB.QueryRow(query, id).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

// GetReferenceByMongoID - Get reference by mongo_achievement_id
func (r *achievementRepository) GetReferenceByMongoID(mongoID string) (*model.AchievementReference, error) {
	ref := &model.AchievementReference{}
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE mongo_achievement_id = $1
	`
	err := r.pgDB.QueryRow(query, mongoID).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

// GetReferencesByStudentID - Get references by student_id dengan filter status
func (r *achievementRepository) GetReferencesByStudentID(studentID string, status string, limit, offset int) ([]model.AchievementReference, error) {
	var query string
	var rows *sql.Rows
	var err error

	if status != "" {
		query = `
			SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
			FROM achievement_references
			WHERE student_id = $1 AND status = $2 AND status != 'deleted'
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
		`
		rows, err = r.pgDB.Query(query, studentID, status, limit, offset)
	} else {
		query = `
			SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
			FROM achievement_references
			WHERE student_id = $1 AND status != 'deleted'
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		rows, err = r.pgDB.Query(query, studentID, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReferences(rows)
}

// CountReferencesByStudentID - Count references by student_id
func (r *achievementRepository) CountReferencesByStudentID(studentID string, status string) (int, error) {
	var count int
	var query string

	if status != "" {
		query = `SELECT COUNT(*) FROM achievement_references WHERE student_id = $1 AND status = $2 AND status != 'deleted'`
		err := r.pgDB.QueryRow(query, studentID, status).Scan(&count)
		return count, err
	}

	query = `SELECT COUNT(*) FROM achievement_references WHERE student_id = $1 AND status != 'deleted'`
	err := r.pgDB.QueryRow(query, studentID).Scan(&count)
	return count, err
}

// GetReferencesByAdvisorID - Get achievements dari mahasiswa bimbingan (untuk Dosen Wali)
func (r *achievementRepository) GetReferencesByAdvisorID(advisorID string, status string, limit, offset int) ([]model.AchievementReference, error) {
	var query string
	var rows *sql.Rows
	var err error

	if status != "" {
		query = `
			SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, ar.submitted_at, ar.verified_at, ar.verified_by, ar.rejection_note, ar.created_at, ar.updated_at
			FROM achievement_references ar
			JOIN students s ON ar.student_id = s.id
			WHERE s.advisor_id = $1 AND ar.status = $2 AND ar.status != 'deleted'
			ORDER BY ar.created_at DESC
			LIMIT $3 OFFSET $4
		`
		rows, err = r.pgDB.Query(query, advisorID, status, limit, offset)
	} else {
		query = `
			SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, ar.submitted_at, ar.verified_at, ar.verified_by, ar.rejection_note, ar.created_at, ar.updated_at
			FROM achievement_references ar
			JOIN students s ON ar.student_id = s.id
			WHERE s.advisor_id = $1 AND ar.status != 'deleted'
			ORDER BY ar.created_at DESC
			LIMIT $2 OFFSET $3
		`
		rows, err = r.pgDB.Query(query, advisorID, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReferences(rows)
}

// CountReferencesByAdvisorID - Count achievements dari mahasiswa bimbingan
func (r *achievementRepository) CountReferencesByAdvisorID(advisorID string, status string) (int, error) {
	var count int
	var query string

	if status != "" {
		query = `
			SELECT COUNT(*)
			FROM achievement_references ar
			JOIN students s ON ar.student_id = s.id
			WHERE s.advisor_id = $1 AND ar.status = $2 AND ar.status != 'deleted'
		`
		err := r.pgDB.QueryRow(query, advisorID, status).Scan(&count)
		return count, err
	}

	query = `
		SELECT COUNT(*)
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		WHERE s.advisor_id = $1 AND ar.status != 'deleted'
	`
	err := r.pgDB.QueryRow(query, advisorID).Scan(&count)
	return count, err
}

// GetAllReferences - Get all references (untuk Admin)
func (r *achievementRepository) GetAllReferences(status string, limit, offset int) ([]model.AchievementReference, error) {
	var query string
	var rows *sql.Rows
	var err error

	if status != "" {
		query = `
			SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
			FROM achievement_references
			WHERE status = $1 AND status != 'deleted'
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		rows, err = r.pgDB.Query(query, status, limit, offset)
	} else {
		query = `
			SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
			FROM achievement_references
			WHERE status != 'deleted'
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		rows, err = r.pgDB.Query(query, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReferences(rows)
}

// CountAllReferences - Count all references
func (r *achievementRepository) CountAllReferences(status string) (int, error) {
	var count int
	var query string

	if status != "" {
		query = `SELECT COUNT(*) FROM achievement_references WHERE status = $1 AND status != 'deleted'`
		err := r.pgDB.QueryRow(query, status).Scan(&count)
		return count, err
	}

	query = `SELECT COUNT(*) FROM achievement_references WHERE status != 'deleted'`
	err := r.pgDB.QueryRow(query).Scan(&count)
	return count, err
}

// Helper: scanReferences
func (r *achievementRepository) scanReferences(rows *sql.Rows) ([]model.AchievementReference, error) {
	var refs []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		)
		if err != nil {
			continue
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

//
// ==================== MONGODB METHODS (ACHIEVEMENTS) ======================
//

// CreateAchievement - Insert achievement ke MongoDB
func (r *achievementRepository) CreateAchievement(achievement *model.Achievement) (string, error) {
	collection := r.mongoDB.Collection("achievements")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	result, err := collection.InsertOne(ctx, achievement)
	if err != nil {
		return "", err
	}

	// Return ObjectID as string
	objectID := result.InsertedID.(primitive.ObjectID)
	return objectID.Hex(), nil
}

// UpdateAchievement - Update achievement di MongoDB
func (r *achievementRepository) UpdateAchievement(id string, achievement *model.Achievement) error {
	collection := r.mongoDB.Collection("achievements")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	achievement.UpdatedAt = time.Now()

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": achievement}

	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// GetAchievementByID - Get achievement dari MongoDB
func (r *achievementRepository) GetAchievementByID(id string) (*model.Achievement, error) {
	collection := r.mongoDB.Collection("achievements")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var achievement model.Achievement
	filter := bson.M{"_id": objectID}
	err = collection.FindOne(ctx, filter).Decode(&achievement)
	if err != nil {
		return nil, err
	}

	return &achievement, nil
}

// DeleteAchievement - Soft delete (optional, bisa juga hard delete)
func (r *achievementRepository) DeleteAchievement(id string) error {
	collection := r.mongoDB.Collection("achievements")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = collection.DeleteOne(ctx, filter)
	return err
}

// AddAttachment - Tambah attachment ke achievement
func (r *achievementRepository) AddAttachment(achievementID string, attachment model.Attachment) error {
	collection := r.mongoDB.Collection("achievements")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(achievementID)
	if err != nil {
		return err
	}

	attachment.UploadedAt = time.Now()

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$push": bson.M{"attachments": attachment},
		"$set":  bson.M{"updatedAt": time.Now()},
	}

	_, err = collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
	return err
}
