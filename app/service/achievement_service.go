package service

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"project_uas/app/model"
	"project_uas/app/repository"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AchievementService struct {
	achievementRepo repository.AchievementRepository
	studentRepo     repository.StudentRepository
	lecturerRepo    repository.LecturerRepository
	userRepo        repository.UserRepository
	validate        *validator.Validate
}

func NewAchievementService(
	achievementRepo repository.AchievementRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
	userRepo repository.UserRepository,
) *AchievementService {
	return &AchievementService{
		achievementRepo: achievementRepo,
		studentRepo:     studentRepo,
		lecturerRepo:    lecturerRepo,
		userRepo:        userRepo,
		validate:        validator.New(),
	}
}

//
// ==================== CREATE ACHIEVEMENT (POST /achievements) ======================
// FR-003: Mahasiswa dapat menambahkan laporan prestasi
//

func (s *AchievementService) CreateAchievement(c *fiber.Ctx) error {
	// Get user dari context
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Parse request
	req := new(model.AchievementCreateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "invalid request body",
		})
	}

	// Validasi input
	if err := s.validate.Struct(req); err != nil {
		return c.Status(422).JSON(model.APIResponse{
			Status: "error",
			Error:  err.Error(),
		})
	}

	// Get student profile dari user yang login
	student, err := s.studentRepo.FindByUserID(claims.UserID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "student profile not found",
		})
	}

	// Create achievement di MongoDB
	achievement := &model.Achievement{
		StudentID:       student.ID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		Points:          req.Points,
		Attachments:     []model.Attachment{}, // empty initially
	}

	mongoID, err := s.achievementRepo.CreateAchievement(achievement)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to create achievement",
		})
	}

	// Create reference di PostgreSQL
	reference := &model.AchievementReference{
		StudentID:          student.ID,
		MongoAchievementID: mongoID,
		Status:             "draft", // Status awal: draft
	}

	if err := s.achievementRepo.CreateReference(reference); err != nil {
		// Rollback: hapus achievement di MongoDB
		s.achievementRepo.DeleteAchievement(mongoID)
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to create achievement reference",
		})
	}

	// Build response
	response := s.buildAchievementResponse(achievement, reference, mongoID)

	return c.Status(201).JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement created successfully",
		Data:    response,
	})
}

//
// ==================== GET ACHIEVEMENT HISTORY (GET /achievements/:id/history) ======================
// Menampilkan riwayat perubahan status achievement
//

func (s *AchievementService) GetAchievementHistory(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	// Get user dari context
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Get reference dari PostgreSQL
	reference, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}

	// Check authorization
	if claims.Role == "Mahasiswa" {
		student, _ := s.studentRepo.FindByUserID(claims.UserID)
		if student == nil || student.ID != reference.StudentID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden",
			})
		}
	} else if claims.Role == "Dosen Wali" {
		lecturer, _ := s.lecturerRepo.FindByUserID(claims.UserID)
		student, _ := s.studentRepo.FindByID(reference.StudentID)
		if lecturer == nil || student == nil || student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden",
			})
		}
	}
	// Admin dapat akses semua

	// Build history
	history := s.buildAchievementHistory(reference)

	return c.JSON(model.APIResponse{
		Status: "success",
		Data: fiber.Map{
			"achievement_id": achievementID,
			"current_status": reference.Status,
			"history":        history,
		},
	})
}

//
// ==================== HELPER: BUILD ACHIEVEMENT HISTORY ======================
//

type HistoryEntry struct {
	Status    string  `json:"status"`
	Timestamp string  `json:"timestamp"`
	Actor     string  `json:"actor,omitempty"`
	ActorID   *string `json:"actor_id,omitempty"`
	Action    string  `json:"action"`
	Notes     *string `json:"notes,omitempty"`
}

func (s *AchievementService) buildAchievementHistory(reference *model.AchievementReference) []HistoryEntry {
	var history []HistoryEntry

	// 1. Draft/Created
	history = append(history, HistoryEntry{
		Status:    "draft",
		Timestamp: reference.CreatedAt.Format("2006-01-02 15:04:05"),
		Action:    "Achievement created",
		Notes:     nil,
	})

	// 2. Submitted (jika ada)
	if reference.SubmittedAt != nil {
		history = append(history, HistoryEntry{
			Status:    "submitted",
			Timestamp: reference.SubmittedAt.Format("2006-01-02 15:04:05"),
			Action:    "Submitted for verification",
			Notes:     nil,
		})
	}

	// 3. Verified (jika ada)
	if reference.Status == "verified" && reference.VerifiedAt != nil {
		var actorName string
		var actorID *string

		if reference.VerifiedBy != nil {
			// Get verifier info
			user, err := s.userRepo.FindByID(*reference.VerifiedBy)
			if err == nil {
				actorName = user.FullName + " (Dosen Wali)"
				actorID = reference.VerifiedBy
			}
		}

		history = append(history, HistoryEntry{
			Status:    "verified",
			Timestamp: reference.VerifiedAt.Format("2006-01-02 15:04:05"),
			Actor:     actorName,
			ActorID:   actorID,
			Action:    "Achievement verified",
			Notes:     nil,
		})
	}

	// 4. Rejected (jika ada)
	if reference.Status == "rejected" {
		var actorName string
		var actorID *string

		if reference.VerifiedBy != nil {
			// Note: VerifiedBy juga dipakai untuk rejection
			user, err := s.userRepo.FindByID(*reference.VerifiedBy)
			if err == nil {
				actorName = user.FullName + " (Dosen Wali)"
				actorID = reference.VerifiedBy
			}
		}

		history = append(history, HistoryEntry{
			Status:    "rejected",
			Timestamp: reference.UpdatedAt.Format("2006-01-02 15:04:05"),
			Actor:     actorName,
			ActorID:   actorID,
			Action:    "Achievement rejected",
			Notes:     reference.RejectionNote,
		})
	}

	// 5. Deleted (jika ada)
	if reference.Status == "deleted" {
		history = append(history, HistoryEntry{
			Status:    "deleted",
			Timestamp: reference.UpdatedAt.Format("2006-01-02 15:04:05"),
			Action:    "Achievement deleted",
			Notes:     nil,
		})
	}

	return history
}


//
// ==================== GET ACHIEVEMENTS (GET /achievements) ======================
// Filtered by role:
// - Mahasiswa: hanya prestasi sendiri
// - Dosen Wali: prestasi mahasiswa bimbingannya
// - Admin: semua prestasi
//

func (s *AchievementService) GetAchievements(c *fiber.Ctx) error {
	// Get user dari context
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Parse query params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	status := c.Query("status", "") // filter by status

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var references []model.AchievementReference
	var total int
	var err error

	// Filter berdasarkan role
	if claims.Role == "Mahasiswa" {
		// Get student profile
		student, err := s.studentRepo.FindByUserID(claims.UserID)
		if err != nil {
			return c.Status(404).JSON(model.APIResponse{
				Status: "error",
				Error:  "student profile not found",
			})
		}

		references, err = s.achievementRepo.GetReferencesByStudentID(student.ID, status, pageSize, offset)
		if err != nil {
			return c.Status(500).JSON(model.APIResponse{
				Status: "error",
				Error:  "failed to fetch achievements",
			})
		}

		total, err = s.achievementRepo.CountReferencesByStudentID(student.ID, status)
		if err != nil {
			return c.Status(500).JSON(model.APIResponse{
				Status: "error",
				Error:  "failed to count achievements",
			})
		}

	} else if claims.Role == "Dosen Wali" {
		// Get lecturer profile
		lecturer, err := s.lecturerRepo.FindByUserID(claims.UserID)
		if err != nil {
			return c.Status(404).JSON(model.APIResponse{
				Status: "error",
				Error:  "lecturer profile not found",
			})
		}

		references, err = s.achievementRepo.GetReferencesByAdvisorID(lecturer.ID, status, pageSize, offset)
		if err != nil {
			return c.Status(500).JSON(model.APIResponse{
				Status: "error",
				Error:  "failed to fetch achievements",
			})
		}

		total, err = s.achievementRepo.CountReferencesByAdvisorID(lecturer.ID, status)
		if err != nil {
			return c.Status(500).JSON(model.APIResponse{
				Status: "error",
				Error:  "failed to count achievements",
			})
		}

	} else if claims.Role == "Admin" {
		// Admin dapat melihat semua
		references, err = s.achievementRepo.GetAllReferences(status, pageSize, offset)
		if err != nil {
			return c.Status(500).JSON(model.APIResponse{
				Status: "error",
				Error:  "failed to fetch achievements",
			})
		}

		total, err = s.achievementRepo.CountAllReferences(status)
		if err != nil {
			return c.Status(500).JSON(model.APIResponse{
				Status: "error",
				Error:  "failed to count achievements",
			})
		}
	} else {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "forbidden",
		})
	}

	// Fetch details dari MongoDB
	var achievements []model.AchievementResponse
	for _, ref := range references {
		achievement, err := s.achievementRepo.GetAchievementByID(ref.MongoAchievementID)
		if err != nil {
			continue // Skip jika tidak ditemukan
		}

		response := s.buildAchievementResponse(achievement, &ref, ref.MongoAchievementID)
		achievements = append(achievements, *response)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	response := model.AchievementListResponse{
		Achievements: achievements,
		Total:        total,
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
	}

	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   response,
	})
}

//
// ==================== GET ACHIEVEMENT BY ID (GET /achievements/:id) ======================
//

func (s *AchievementService) GetAchievementByID(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	// Get user dari context
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Get reference dari PostgreSQL
	reference, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}

	// Check authorization
	if claims.Role == "Mahasiswa" {
		student, _ := s.studentRepo.FindByUserID(claims.UserID)
		if student == nil || student.ID != reference.StudentID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden",
			})
		}
	} else if claims.Role == "Dosen Wali" {
		lecturer, _ := s.lecturerRepo.FindByUserID(claims.UserID)
		student, _ := s.studentRepo.FindByID(reference.StudentID)
		if lecturer == nil || student == nil || student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden",
			})
		}
	}
	// Admin dapat akses semua

	// Get detail dari MongoDB
	achievement, err := s.achievementRepo.GetAchievementByID(reference.MongoAchievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement detail not found",
		})
	}

	response := s.buildAchievementResponse(achievement, reference, reference.MongoAchievementID)

	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   response,
	})
}

//
// ==================== UPDATE ACHIEVEMENT (PUT /achievements/:id) ======================
// Hanya mahasiswa pemilik yang bisa update, dan hanya jika status = draft
//

func (s *AchievementService) UpdateAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	// Get user dari context
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Get reference
	reference, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}

	// Check authorization (hanya mahasiswa pemilik)
	student, _ := s.studentRepo.FindByUserID(claims.UserID)
	if student == nil || student.ID != reference.StudentID {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "forbidden",
		})
	}

	// Hanya bisa update jika status = draft
	if reference.Status != "draft" {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "can only update achievement with status 'draft'",
		})
	}

	// Parse request
	req := new(model.AchievementUpdateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "invalid request body",
		})
	}

	// Get existing achievement dari MongoDB
	achievement, err := s.achievementRepo.GetAchievementByID(reference.MongoAchievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement detail not found",
		})
	}

	// Update fields (hanya yang diisi)
	if req.AchievementType != "" {
		achievement.AchievementType = req.AchievementType
	}
	if req.Title != "" {
		achievement.Title = req.Title
	}
	if req.Description != "" {
		achievement.Description = req.Description
	}
	if req.Details != nil {
		achievement.Details = req.Details
	}
	if req.Tags != nil {
		achievement.Tags = req.Tags
	}
	if req.Points > 0 {
		achievement.Points = req.Points
	}

	// Update di MongoDB
	if err := s.achievementRepo.UpdateAchievement(reference.MongoAchievementID, achievement); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to update achievement",
		})
	}

	response := s.buildAchievementResponse(achievement, reference, reference.MongoAchievementID)

	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement updated successfully",
		Data:    response,
	})
}

//
// ==================== DELETE ACHIEVEMENT (DELETE /achievements/:id) ======================
// FR-005: Mahasiswa dapat menghapus prestasi draft
// Flow SRS:
// 1. Soft delete data di MongoDB
// 2. Update reference di PostgreSQL
// 3. Return success message
//

func (s *AchievementService) DeleteAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	// Get user dari context
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Get reference
	reference, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}

	// Check authorization (hanya mahasiswa pemilik)
	student, _ := s.studentRepo.FindByUserID(claims.UserID)
	if student == nil || student.ID != reference.StudentID {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "forbidden",
		})
	}

	// Precondition: Hanya bisa delete jika status = draft
	if reference.Status != "draft" {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "can only delete achievement with status 'draft'",
		})
	}

	// FR-005: Soft delete sesuai SRS
	// 1. Soft delete data di MongoDB
	if err := s.achievementRepo.DeleteAchievement(reference.MongoAchievementID); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to delete achievement from MongoDB",
		})
	}

	// 2. Update reference di PostgreSQL dengan status 'deleted'
	reference.Status = "deleted"
	if err := s.achievementRepo.UpdateReference(reference); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to update reference status",
		})
	}

	// 3. Return success message
	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement deleted successfully",
	})
}

//
// ==================== SUBMIT FOR VERIFICATION (POST /achievements/:id/submit) ======================
// FR-004: Mahasiswa submit prestasi draft untuk diverifikasi
//

func (s *AchievementService) SubmitForVerification(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	// Get user dari context
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Get reference
	reference, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}

	// Check authorization
	student, _ := s.studentRepo.FindByUserID(claims.UserID)
	if student == nil || student.ID != reference.StudentID {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "forbidden",
		})
	}

	// Hanya bisa submit jika status = draft
	if reference.Status != "draft" {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement already submitted or processed",
		})
	}

	// Update status menjadi 'submitted'
	now := time.Now()
	reference.Status = "submitted"
	reference.SubmittedAt = &now

	if err := s.achievementRepo.UpdateReference(reference); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to submit achievement",
		})
	}

	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement submitted for verification",
		Data: fiber.Map{
			"status":       reference.Status,
			"submitted_at": reference.SubmittedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

//
// ==================== VERIFY ACHIEVEMENT (POST /achievements/:id/verify) ======================
// FR-007: Dosen wali memverifikasi prestasi mahasiswa
//

func (s *AchievementService) VerifyAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	// Get user dari context
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Get reference
	reference, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}

	// Check authorization (hanya dosen wali dari mahasiswa tersebut)
	lecturer, _ := s.lecturerRepo.FindByUserID(claims.UserID)
	student, _ := s.studentRepo.FindByID(reference.StudentID)

	if lecturer == nil || student == nil || student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "forbidden: you are not the advisor of this student",
		})
	}

	// Hanya bisa verify jika status = submitted
	if reference.Status != "submitted" {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement must be in 'submitted' status",
		})
	}

	// Update status menjadi 'verified'
	now := time.Now()
	reference.Status = "verified"
	reference.VerifiedAt = &now
	reference.VerifiedBy = &claims.UserID

	if err := s.achievementRepo.UpdateReference(reference); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to verify achievement",
		})
	}

	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement verified successfully",
		Data: fiber.Map{
			"status":      reference.Status,
			"verified_at": reference.VerifiedAt.Format("2006-01-02 15:04:05"),
			"verified_by": reference.VerifiedBy,
		},
	})
}

//
// ==================== REJECT ACHIEVEMENT (POST /achievements/:id/reject) ======================
// FR-008: Dosen wali menolak prestasi dengan catatan
//

func (s *AchievementService) RejectAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	// Get user dari context
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Parse request
	req := new(model.RejectAchievementRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "invalid request body",
		})
	}

	// Validasi
	if err := s.validate.Struct(req); err != nil {
		return c.Status(422).JSON(model.APIResponse{
			Status: "error",
			Error:  err.Error(),
		})
	}

	// Get reference
	reference, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}

	// Check authorization
	lecturer, _ := s.lecturerRepo.FindByUserID(claims.UserID)
	student, _ := s.studentRepo.FindByID(reference.StudentID)

	if lecturer == nil || student == nil || student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "forbidden: you are not the advisor of this student",
		})
	}

	// Hanya bisa reject jika status = submitted
	if reference.Status != "submitted" {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement must be in 'submitted' status",
		})
	}

	// Update status menjadi 'rejected'
	reference.Status = "rejected"
	reference.RejectionNote = &req.RejectionNote

	if err := s.achievementRepo.UpdateReference(reference); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to reject achievement",
		})
	}

	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "achievement rejected",
		Data: fiber.Map{
			"status":         reference.Status,
			"rejection_note": reference.RejectionNote,
		},
	})
}

// Ganti fungsi UploadAttachment yang lama dengan ini:

//
// ==================== UPLOAD ATTACHMENT (POST /achievements/:id/attachments) ======================
// Handle REAL file upload dengan multipart/form-data
//

func (s *AchievementService) UploadAttachment(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	// Get user dari context
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Get reference
	reference, err := s.achievementRepo.GetReferenceByID(achievementID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "achievement not found",
		})
	}

	// Check authorization (hanya mahasiswa pemilik)
	student, _ := s.studentRepo.FindByUserID(claims.UserID)
	if student == nil || student.ID != reference.StudentID {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "forbidden",
		})
	}

	// Hanya bisa upload jika status = draft atau submitted
	if reference.Status != "draft" && reference.Status != "submitted" {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "can only upload attachments for draft or submitted achievements",
		})
	}

	// Parse multipart file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "file is required",
		})
	}

	// Validasi ukuran file (max 5MB)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if file.Size > maxSize {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "file size exceeds 5MB limit",
		})
	}

	// Validasi tipe file (hanya PDF, JPG, PNG, JPEG)
	allowedTypes := map[string]bool{
		"application/pdf":  true,
		"image/jpeg":       true,
		"image/jpg":        true,
		"image/png":        true,
	}

	// Get MIME type from header
	fileHeader, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to read file",
		})
	}
	defer fileHeader.Close()

	// Read first 512 bytes untuk detect MIME type
	buffer := make([]byte, 512)
	_, err = fileHeader.Read(buffer)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to read file content",
		})
	}

	// Detect MIME type
	contentType := http.DetectContentType(buffer)
	if !allowedTypes[contentType] {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "file type not allowed. Only PDF, JPG, PNG are accepted",
		})
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	randomString := uuid.New().String()[:8]
	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%s_%d_%s%s", achievementID, timestamp, randomString, ext)

	// Create uploads directory jika belum ada
	uploadsDir := "./uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			return c.Status(500).JSON(model.APIResponse{
				Status: "error",
				Error:  "failed to create uploads directory",
			})
		}
	}

	// Simpan file
	filePath := filepath.Join(uploadsDir, newFilename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to save file",
		})
	}

	// Create attachment object
	attachment := model.Attachment{
		FileName:   file.Filename, // Original filename
		FileURL:    fmt.Sprintf("/uploads/%s", newFilename), // Relative path
		FileType:   contentType,
		UploadedAt: time.Now(),
	}

	// Add attachment ke MongoDB
	if err := s.achievementRepo.AddAttachment(reference.MongoAchievementID, attachment); err != nil {
		// Rollback: hapus file yang sudah diupload
		os.Remove(filePath)
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to save attachment metadata",
		})
	}

	return c.Status(201).JSON(model.APIResponse{
		Status:  "success",
		Message: "attachment uploaded successfully",
		Data:    attachment,
	})
}

//
// ==================== HELPER: BUILD ACHIEVEMENT RESPONSE ======================
//

func (s *AchievementService) buildAchievementResponse(
	achievement *model.Achievement,
	reference *model.AchievementReference,
	mongoID string,
) *model.AchievementResponse {
	response := &model.AchievementResponse{
		ID:              reference.ID,
		StudentID:       achievement.StudentID,
		AchievementType: achievement.AchievementType,
		Title:           achievement.Title,
		Description:     achievement.Description,
		Details:         achievement.Details,
		Attachments:     achievement.Attachments,
		Tags:            achievement.Tags,
		Points:          achievement.Points,
		Status:          reference.Status,
		CreatedAt:       achievement.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       achievement.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if reference.SubmittedAt != nil {
		submittedAt := reference.SubmittedAt.Format("2006-01-02 15:04:05")
		response.SubmittedAt = &submittedAt
	}

	if reference.VerifiedAt != nil {
		verifiedAt := reference.VerifiedAt.Format("2006-01-02 15:04:05")
		response.VerifiedAt = &verifiedAt
	}

	response.VerifiedBy = reference.VerifiedBy
	response.RejectionNote = reference.RejectionNote

	return response
}