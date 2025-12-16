package mocks

import (
	"project_uas/app/model"
	"github.com/stretchr/testify/mock"
)

// ==================== MOCK ACHIEVEMENT REPOSITORY ====================

type MockAchievementRepository struct {
	mock.Mock
}

func (m *MockAchievementRepository) CreateReference(ref *model.AchievementReference) error {
	args := m.Called(ref)
	return args.Error(0)
}

func (m *MockAchievementRepository) UpdateReference(ref *model.AchievementReference) error {
	args := m.Called(ref)
	return args.Error(0)
}

func (m *MockAchievementRepository) GetReferenceByID(id string) (*model.AchievementReference, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AchievementReference), args.Error(1)
}

func (m *MockAchievementRepository) GetReferenceByMongoID(mongoID string) (*model.AchievementReference, error) {
	args := m.Called(mongoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AchievementReference), args.Error(1)
}

func (m *MockAchievementRepository) GetReferencesByStudentID(studentID string, status string, limit, offset int) ([]model.AchievementReference, error) {
	args := m.Called(studentID, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AchievementReference), args.Error(1)
}

func (m *MockAchievementRepository) CountReferencesByStudentID(studentID string, status string) (int, error) {
	args := m.Called(studentID, status)
	return args.Int(0), args.Error(1)
}

func (m *MockAchievementRepository) GetReferencesByAdvisorID(advisorID string, status string, limit, offset int) ([]model.AchievementReference, error) {
	args := m.Called(advisorID, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AchievementReference), args.Error(1)
}

func (m *MockAchievementRepository) CountReferencesByAdvisorID(advisorID string, status string) (int, error) {
	args := m.Called(advisorID, status)
	return args.Int(0), args.Error(1)
}

func (m *MockAchievementRepository) GetAllReferences(status string, limit, offset int) ([]model.AchievementReference, error) {
	args := m.Called(status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AchievementReference), args.Error(1)
}

func (m *MockAchievementRepository) CountAllReferences(status string) (int, error) {
	args := m.Called(status)
	return args.Int(0), args.Error(1)
}

func (m *MockAchievementRepository) CreateAchievement(achievement *model.Achievement) (string, error) {
	args := m.Called(achievement)
	return args.String(0), args.Error(1)
}

func (m *MockAchievementRepository) UpdateAchievement(id string, achievement *model.Achievement) error {
	args := m.Called(id, achievement)
	return args.Error(0)
}

func (m *MockAchievementRepository) GetAchievementByID(id string) (*model.Achievement, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Achievement), args.Error(1)
}

func (m *MockAchievementRepository) DeleteAchievement(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAchievementRepository) AddAttachment(achievementID string, attachment model.Attachment) error {
	args := m.Called(achievementID, attachment)
	return args.Error(0)
}

// ==================== MOCK STUDENT REPOSITORY ====================

type MockStudentRepository struct {
	mock.Mock
}

func (m *MockStudentRepository) FindByUserID(userID string) (*model.Student, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockStudentRepository) FindByID(id string) (*model.Student, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockStudentRepository) FindByStudentID(studentID string) (*model.Student, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockStudentRepository) Create(student *model.Student) error {
	args := m.Called(student)
	return args.Error(0)
}

func (m *MockStudentRepository) Update(student *model.Student) error {
	args := m.Called(student)
	return args.Error(0)
}

func (m *MockStudentRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStudentRepository) SetAdvisor(studentID string, advisorID string) error {
	args := m.Called(studentID, advisorID)
	return args.Error(0)
}

func (m *MockStudentRepository) GetAll(limit, offset int) ([]model.Student, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Student), args.Error(1)
}

func (m *MockStudentRepository) CountAll() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

// ==================== MOCK LECTURER REPOSITORY ====================

type MockLecturerRepository struct {
	mock.Mock
}

func (m *MockLecturerRepository) FindByUserID(userID string) (*model.Lecturer, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Lecturer), args.Error(1)
}

func (m *MockLecturerRepository) FindByID(id string) (*model.Lecturer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Lecturer), args.Error(1)
}

func (m *MockLecturerRepository) FindByLecturerID(lecturerID string) (*model.Lecturer, error) {
	args := m.Called(lecturerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Lecturer), args.Error(1)
}

func (m *MockLecturerRepository) Create(lecturer *model.Lecturer) error {
	args := m.Called(lecturer)
	return args.Error(0)
}

func (m *MockLecturerRepository) Update(lecturer *model.Lecturer) error {
	args := m.Called(lecturer)
	return args.Error(0)
}

func (m *MockLecturerRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLecturerRepository) GetAll(limit, offset int) ([]model.Lecturer, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Lecturer), args.Error(1)
}

func (m *MockLecturerRepository) CountAll() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

// ==================== MOCK USER REPOSITORY ====================

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(id string) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll(limit, offset int, roleName string) ([]model.User, error) {
	args := m.Called(limit, offset, roleName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) CountAll(roleName string) (int, error) {
	args := m.Called(roleName)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) GetRoleName(roleID string) (string, error) {
	args := m.Called(roleID)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) GetPermissions(role string) ([]string, error) {
	args := m.Called(role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockUserRepository) UpdateRole(userID string, roleID string) error {
	args := m.Called(userID, roleID)
	return args.Error(0)
}

// ==================== MOCK ROLE REPOSITORY ====================

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetRoleByID(id string) (*model.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

func (m *MockRoleRepository) GetRoleByName(name string) (*model.Role, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

// ==================== MOCK PERMISSION REPOSITORY ====================

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) GetPermissionsByRoleID(roleID string) ([]string, error) {
	args := m.Called(roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}
