package mocks

import (
	"UASBE/app/model"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository
type MockUserRepository struct{ mock.Mock }
func (m *MockUserRepository) FindByUsername(un string) (*model.User, error) {
	args := m.Called(un)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockUserRepository) FindByID(id string) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockUserRepository) GetRoleName(rid string) (string, error) {
	args := m.Called(rid)
	return args.String(0), args.Error(1)
}
func (m *MockUserRepository) GetPermissions(role string) ([]string, error) {
	args := m.Called(role)
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockUserRepository) Create(u *model.User) error { return m.Called(u).Error(0) }
func (m *MockUserRepository) Update(u *model.User) error { return m.Called(u).Error(0) }
func (m *MockUserRepository) Delete(id string) error { return m.Called(id).Error(0) }
func (m *MockUserRepository) FindByEmail(e string) (*model.User, error) {
	args := m.Called(e)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockUserRepository) GetAll(l, o int, r string) ([]model.User, error) {
	args := m.Called(l, o, r)
	return args.Get(0).([]model.User), args.Error(1)
}
func (m *MockUserRepository) CountAll(r string) (int, error) {
	args := m.Called(r)
	return args.Int(0), args.Error(1)
}
func (m *MockUserRepository) UpdateRole(uid, rid string) error { return m.Called(uid, rid).Error(0) }

// MockRoleRepository
type MockRoleRepository struct{ mock.Mock }
func (m *MockRoleRepository) GetRoleByID(id string) (*model.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Role), args.Error(1)
}
func (m *MockRoleRepository) GetRoleByName(n string) (*model.Role, error) {
	args := m.Called(n)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Role), args.Error(1)
}

// MockPermissionRepository
type MockPermissionRepository struct{ mock.Mock }
func (m *MockPermissionRepository) GetPermissionsByRoleID(id string) ([]string, error) {
	args := m.Called(id)
	return args.Get(0).([]string), args.Error(1)
}

// MockStudentRepository
type MockStudentRepository struct{ mock.Mock }
func (m *MockStudentRepository) FindByUserID(uid string) (*model.Student, error) {
	args := m.Called(uid)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Student), args.Error(1)
}
func (m *MockStudentRepository) FindByID(id string) (*model.Student, error) {
	args := m.Called(id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Student), args.Error(1)
}
func (m *MockStudentRepository) Create(s *model.Student) error { return m.Called(s).Error(0) }
func (m *MockStudentRepository) Update(s *model.Student) error { return m.Called(s).Error(0) }
func (m *MockStudentRepository) Delete(id string) error { return m.Called(id).Error(0) }
func (m *MockStudentRepository) FindByStudentID(sid string) (*model.Student, error) {
	args := m.Called(sid)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Student), args.Error(1)
}
func (m *MockStudentRepository) SetAdvisor(sid, aid string) error { return m.Called(sid, aid).Error(0) }
func (m *MockStudentRepository) GetAll(l, o int) ([]model.Student, error) {
	args := m.Called(l, o)
	return args.Get(0).([]model.Student), args.Error(1)
}
func (m *MockStudentRepository) CountAll() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

// MockLecturerRepository
type MockLecturerRepository struct{ mock.Mock }
func (m *MockLecturerRepository) FindByUserID(uid string) (*model.Lecturer, error) {
	args := m.Called(uid)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Lecturer), args.Error(1)
}
func (m *MockLecturerRepository) FindByID(id string) (*model.Lecturer, error) {
	args := m.Called(id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Lecturer), args.Error(1)
}
func (m *MockLecturerRepository) Create(l *model.Lecturer) error { return m.Called(l).Error(0) }
func (m *MockLecturerRepository) FindByLecturerID(lid string) (*model.Lecturer, error) {
	args := m.Called(lid)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Lecturer), args.Error(1)
}
func (m *MockLecturerRepository) Update(l *model.Lecturer) error { return m.Called(l).Error(0) }
func (m *MockLecturerRepository) Delete(id string) error { return m.Called(id).Error(0) }
func (m *MockLecturerRepository) GetAll(l, o int) ([]model.Lecturer, error) {
	args := m.Called(l, o)
	return args.Get(0).([]model.Lecturer), args.Error(1)
}
func (m *MockLecturerRepository) CountAll() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

// MockAchievementRepository
type MockAchievementRepository struct{ mock.Mock }
func (m *MockAchievementRepository) CreateAchievement(a *model.Achievement) (string, error) {
	args := m.Called(a)
	return args.String(0), args.Error(1)
}
func (m *MockAchievementRepository) CreateReference(r *model.AchievementReference) error {
	return m.Called(r).Error(0)
}
func (m *MockAchievementRepository) GetReferenceByID(id string) (*model.AchievementReference, error) {
	args := m.Called(id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.AchievementReference), args.Error(1)
}
func (m *MockAchievementRepository) UpdateReference(r *model.AchievementReference) error {
	return m.Called(r).Error(0)
}
func (m *MockAchievementRepository) GetAchievementByID(id string) (*model.Achievement, error) {
	args := m.Called(id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Achievement), args.Error(1)
}
func (m *MockAchievementRepository) DeleteAchievement(id string) error { return m.Called(id).Error(0) }
func (m *MockAchievementRepository) GetReferenceByMongoID(mid string) (*model.AchievementReference, error) {
	args := m.Called(mid)
	return args.Get(0).(*model.AchievementReference), args.Error(1)
}
func (m *MockAchievementRepository) GetReferencesByStudentID(sid, s string, l, o int) ([]model.AchievementReference, error) {
	args := m.Called(sid, s, l, o)
	return args.Get(0).([]model.AchievementReference), args.Error(1)
}
func (m *MockAchievementRepository) CountReferencesByStudentID(sid, s string) (int, error) {
	args := m.Called(sid, s)
	return args.Int(0), args.Error(1)
}
func (m *MockAchievementRepository) GetReferencesByAdvisorID(aid, s string, l, o int) ([]model.AchievementReference, error) {
	args := m.Called(aid, s, l, o)
	return args.Get(0).([]model.AchievementReference), args.Error(1)
}
func (m *MockAchievementRepository) CountReferencesByAdvisorID(aid, s string) (int, error) {
	args := m.Called(aid, s)
	return args.Int(0), args.Error(1)
}
func (m *MockAchievementRepository) GetAllReferences(s string, l, o int) ([]model.AchievementReference, error) {
	args := m.Called(s, l, o)
	return args.Get(0).([]model.AchievementReference), args.Error(1)
}
func (m *MockAchievementRepository) CountAllReferences(s string) (int, error) {
	args := m.Called(s)
	return args.Int(0), args.Error(1)
}
func (m *MockAchievementRepository) UpdateAchievement(id string, a *model.Achievement) error { return m.Called(id, a).Error(0) }
func (m *MockAchievementRepository) AddAttachment(id string, at model.Attachment) error { return m.Called(id, at).Error(0) }

// MockReportRepository
type MockReportRepository struct{ mock.Mock }

func (m *MockReportRepository) GetTotalByType(sid, aid *string) (map[string]int, error) {
	args := m.Called(sid, aid)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockReportRepository) GetTotalByPeriod(sid, aid *string) ([]model.PeriodStats, error) {
	args := m.Called(sid, aid)
	return args.Get(0).([]model.PeriodStats), args.Error(1)
}

func (m *MockReportRepository) GetTopStudents(limit int, aid *string) ([]model.TopStudent, error) {
	args := m.Called(limit, aid)
	return args.Get(0).([]model.TopStudent), args.Error(1)
}

func (m *MockReportRepository) GetCompetitionLevelDistribution(sid, aid *string) (map[string]int, error) {
	args := m.Called(sid, aid)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockReportRepository) GetStatusBreakdown(sid, aid *string) (map[string]int, error) {
	args := m.Called(sid, aid)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockReportRepository) GetStudentSummary(sid string) (*model.StudentSummary, error) {
	args := m.Called(sid)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.StudentSummary), args.Error(1)
}

func (m *MockReportRepository) GetStudentAchievementsByType(sid string) (map[string]int, error) {
	args := m.Called(sid)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockReportRepository) GetStudentAchievementsByStatus(sid string) (map[string]int, error) {
	args := m.Called(sid)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockReportRepository) GetStudentTimeline(sid string) ([]model.PeriodStats, error) {
	args := m.Called(sid)
	return args.Get(0).([]model.PeriodStats), args.Error(1)
}