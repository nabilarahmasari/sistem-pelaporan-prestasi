package service

import (
	

	"github.com/gofiber/fiber/v2"

	"project_uas/app/model"
	"project_uas/app/repository"
)

type ReportService struct {
	achievementRepo repository.AchievementRepository
	studentRepo     repository.StudentRepository
	lecturerRepo    repository.LecturerRepository
	userRepo        repository.UserRepository
}

func NewReportService(
	achievementRepo repository.AchievementRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
	userRepo repository.UserRepository,
) *ReportService {
	return &ReportService{
		achievementRepo: achievementRepo,
		studentRepo:     studentRepo,
		lecturerRepo:    lecturerRepo,
		userRepo:        userRepo,
	}
}

//
// ==================== GET STATISTICS (GET /reports/statistics) ======================
// SRS FR-011: Achievement Statistics
// Actor: Mahasiswa (own), Dosen Wali (advisees), Admin (all)
// Output:
// • Total prestasi per tipe
// • Total prestasi per periode
// • Top mahasiswa berprestasi
// • Distribusi tingkat kompetisi
//

func (s *ReportService) GetStatistics(c *fiber.Ctx) error {
	// Get user dari context
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	var references []model.AchievementReference
	var err error

	// Filter berdasarkan role (sesuai FR-011)
	if claims.Role == "Mahasiswa" {
		// Mahasiswa: hanya statistik prestasi sendiri
		student, err := s.studentRepo.FindByUserID(claims.UserID)
		if err != nil {
			return c.Status(404).JSON(model.APIResponse{
				Status: "error",
				Error:  "student profile not found",
			})
		}

		references, err = s.achievementRepo.GetReferencesByStudentID(student.ID, "", 10000, 0)
		if err != nil {
			return c.Status(500).JSON(model.APIResponse{
				Status: "error",
				Error:  "failed to fetch achievements",
			})
		}

	} else if claims.Role == "Dosen Wali" {
		// Dosen Wali: statistik mahasiswa bimbingannya
		lecturer, err := s.lecturerRepo.FindByUserID(claims.UserID)
		if err != nil {
			return c.Status(404).JSON(model.APIResponse{
				Status: "error",
				Error:  "lecturer profile not found",
			})
		}

		references, err = s.achievementRepo.GetReferencesByAdvisorID(lecturer.ID, "", 10000, 0)
		if err != nil {
			return c.Status(500).JSON(model.APIResponse{
				Status: "error",
				Error:  "failed to fetch achievements",
			})
		}

	} else if claims.Role == "Admin" {
		// Admin: statistik seluruh sistem
		references, err = s.achievementRepo.GetAllReferences("", 10000, 0)
		if err != nil {
			return c.Status(500).JSON(model.APIResponse{
				Status: "error",
				Error:  "failed to fetch achievements",
			})
		}
	} else {
		return c.Status(403).JSON(model.APIResponse{
			Status: "error",
			Error:  "forbidden",
		})
	}

	// Aggregate statistics
	stats := s.aggregateStatistics(references)

	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   stats,
	})
}

//
// ==================== GET STUDENT REPORT (GET /reports/student/:id) ======================
// SRS Section 5.8: GET /api/v1/reports/student/:id
// Authorization: Mahasiswa (own), Dosen Wali (advisees), Admin (all)
//

func (s *ReportService) GetStudentReport(c *fiber.Ctx) error {
	studentID := c.Params("id")

	// Get user dari context
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Get student
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "student not found",
		})
	}

	// Authorization check
	if claims.Role == "Mahasiswa" {
		// Mahasiswa hanya bisa akses report sendiri
		currentStudent, _ := s.studentRepo.FindByUserID(claims.UserID)
		if currentStudent == nil || currentStudent.ID != studentID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden",
			})
		}
	} else if claims.Role == "Dosen Wali" {
		// Dosen hanya bisa akses report mahasiswa bimbingannya
		lecturer, _ := s.lecturerRepo.FindByUserID(claims.UserID)
		if lecturer == nil || student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden: not your advisee",
			})
		}
	}
	// Admin dapat akses semua

	// Get student user info
	user, _ := s.userRepo.FindByID(student.UserID)

	// Get all achievements
	references, err := s.achievementRepo.GetReferencesByStudentID(studentID, "", 10000, 0)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch achievements",
		})
	}

	// Aggregate statistics
	stats := s.aggregateStatistics(references)

	// Get recent achievements (last 5)
	recentRefs, _ := s.achievementRepo.GetReferencesByStudentID(studentID, "", 5, 0)
	var recentAchievements []map[string]interface{}
	for _, ref := range recentRefs {
		achievement, err := s.achievementRepo.GetAchievementByID(ref.MongoAchievementID)
		if err != nil {
			continue
		}

		recentAchievements = append(recentAchievements, map[string]interface{}{
			"id":     ref.ID,
			"title":  achievement.Title,
			"type":   achievement.AchievementType,
			"status": ref.Status,
			"points": achievement.Points,
			"date":   achievement.CreatedAt.Format("2006-01-02"),
		})
	}

	// Build response
	response := fiber.Map{
		"student": fiber.Map{
			"id":            student.ID,
			"student_id":    student.StudentID,
			"full_name":     user.FullName,
			"email":         user.Email,
			"program_study": student.ProgramStudy,
			"academic_year": student.AcademicYear,
		},
		"summary": fiber.Map{
			"total_achievements":    stats["total_achievements"],
			"total_points":          stats["total_points"],
			"verified_achievements": stats["by_status"].(map[string]int)["verified"],
			"pending_achievements":  stats["by_status"].(map[string]int)["submitted"],
			"rejected_achievements": stats["by_status"].(map[string]int)["rejected"],
		},
		"by_type":            stats["by_type"],
		"by_period":          stats["by_period"],
		"competition_levels": stats["competition_levels"],
		"recent_achievements": recentAchievements,
	}

	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   response,
	})
}

//
// ==================== HELPER: AGGREGATE STATISTICS ======================
//

func (s *ReportService) aggregateStatistics(references []model.AchievementReference) map[string]interface{} {
	// Initialize counters
	byType := make(map[string]int)
	byStatus := make(map[string]int)
	byPeriod := make(map[string]int)
	competitionLevels := make(map[string]int)
	studentPoints := make(map[string]int)      // student_id -> total_points
	studentCounts := make(map[string]int)      // student_id -> achievement_count
	studentNames := make(map[string]string)    // student_id -> full_name

	totalPoints := 0
	totalAchievements := len(references)

	// Process each achievement
	for _, ref := range references {
		// Get achievement details dari MongoDB
		achievement, err := s.achievementRepo.GetAchievementByID(ref.MongoAchievementID)
		if err != nil {
			continue
		}

		// Count by type
		byType[achievement.AchievementType]++

		// Count by status
		byStatus[ref.Status]++

		// Count by period (year-month)
		period := achievement.CreatedAt.Format("2006-01")
		byPeriod[period]++

		// Count competition levels
		if achievement.AchievementType == "competition" {
			if details, ok := achievement.Details["competitionLevel"].(string); ok {
				competitionLevels[details]++
			}
		}

		// Accumulate points per student
		studentPoints[achievement.StudentID] += achievement.Points
		studentCounts[achievement.StudentID]++
		totalPoints += achievement.Points

		// Get student name (cache untuk performa)
		if _, exists := studentNames[achievement.StudentID]; !exists {
			student, err := s.studentRepo.FindByID(achievement.StudentID)
			if err == nil {
				user, err := s.userRepo.FindByID(student.UserID)
				if err == nil {
					studentNames[achievement.StudentID] = user.FullName
				}
			}
		}
	}

	// Build top students list
	type StudentRank struct {
		StudentID         string `json:"student_id"`
		StudentNIM        string `json:"student_nim"`
		FullName          string `json:"full_name"`
		TotalAchievements int    `json:"total_achievements"`
		TotalPoints       int    `json:"total_points"`
	}

	var topStudents []StudentRank
	for studentID, points := range studentPoints {
		student, _ := s.studentRepo.FindByID(studentID)
		
		topStudents = append(topStudents, StudentRank{
			StudentID:         studentID,
			StudentNIM:        student.StudentID,
			FullName:          studentNames[studentID],
			TotalAchievements: studentCounts[studentID],
			TotalPoints:       points,
		})
	}

	// Sort top students by points (simple bubble sort for top 10)
	for i := 0; i < len(topStudents); i++ {
		for j := i + 1; j < len(topStudents); j++ {
			if topStudents[j].TotalPoints > topStudents[i].TotalPoints {
				topStudents[i], topStudents[j] = topStudents[j], topStudents[i]
			}
		}
	}

	// Take only top 10
	if len(topStudents) > 10 {
		topStudents = topStudents[:10]
	}

	// Build response
	return map[string]interface{}{
		"total_achievements":  totalAchievements,
		"total_points":        totalPoints,
		"by_type":             byType,
		"by_status":           byStatus,
		"by_period":           byPeriod,
		"competition_levels":  competitionLevels,
		"top_students":        topStudents,
	}
}
