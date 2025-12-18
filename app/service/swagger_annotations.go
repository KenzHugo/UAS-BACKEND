package service

// File ini berisi Swagger annotations untuk semua service endpoints
// Letakkan di folder app/service/

// ==================== AUTH SERVICE ANNOTATIONS ======================

// Login godoc
// @Summary Login to the system
// @Description Authenticate user with username/email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.APIResponse{data=model.LoginResponse} "Login successful"
// @Failure 400 {object} model.APIResponse "Invalid request body"
// @Failure 401 {object} model.APIResponse "Invalid username or password"
// @Router /auth/login [post]
func (s *AuthService) LoginSwagger() {}

// Refresh godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} model.APIResponse{data=model.LoginResponse} "Token refreshed"
// @Failure 401 {object} model.APIResponse "Invalid refresh token"
// @Failure 404 {object} model.APIResponse "User not found"
// @Router /auth/refresh [post]
func (s *AuthService) RefreshSwagger() {}

// Profile godoc
// @Summary Get current user profile
// @Description Get profile of currently authenticated user
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.APIResponse{data=model.UserResponse} "User profile"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 404 {object} model.APIResponse "User not found"
// @Router /auth/profile [get]
func (s *AuthService) ProfileSwagger() {}

// Logout godoc
// @Summary Logout from system
// @Description Logout current user (invalidate token on client side)
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.APIResponse "Logout successful"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Router /auth/logout [post]
func (s *AuthService) LogoutSwagger() {}

// ==================== USER SERVICE ANNOTATIONS ======================

// CreateUser godoc
// @Summary Create new user (Admin only)
// @Description Create new user with role and profile. Supports creating Mahasiswa with student profile or Dosen Wali with lecturer profile.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.UserCreateRequest true "User data with optional student_profile or lecturer_profile"
// @Success 201 {object} model.APIResponse{data=model.UserResponse} "User created successfully"
// @Failure 400 {object} model.APIResponse "Invalid request body"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Admin only"
// @Failure 404 {object} model.APIResponse "Role not found"
// @Failure 409 {object} model.APIResponse "Username/Email/StudentID/LecturerID already exists"
// @Failure 422 {object} model.APIResponse "Validation error"
// @Router /users [post]
func (s *UserService) CreateUserSwagger() {}

// GetUsers godoc
// @Summary Get all users (Admin only)
// @Description Get list of all users with pagination and optional role filter. Includes student or lecturer profile if applicable.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size (max 100)" default(10)
// @Param role query string false "Filter by role name" Enums(Admin, Mahasiswa, Dosen Wali)
// @Success 200 {object} model.APIResponse{data=model.UserListResponse} "List of users"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Admin only"
// @Router /users [get]
func (s *UserService) GetUsersSwagger() {}

// GetUserByID godoc
// @Summary Get user by ID (Admin only)
// @Description Get detailed user information by ID including profile (student/lecturer)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} model.APIResponse{data=model.UserResponse} "User details"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Admin only"
// @Failure 404 {object} model.APIResponse "User not found"
// @Router /users/{id} [get]
func (s *UserService) GetUserByIDSwagger() {}

// UpdateUser godoc
// @Summary Update user (Admin only)
// @Description Update user information (email, full_name, is_active). Cannot update username or password.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Param request body model.UserUpdateRequest true "Update data"
// @Success 200 {object} model.APIResponse{data=model.UserResponse} "User updated successfully"
// @Failure 400 {object} model.APIResponse "Invalid request body"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Admin only"
// @Failure 404 {object} model.APIResponse "User not found"
// @Failure 409 {object} model.APIResponse "Email already used by another user"
// @Router /users/{id} [put]
func (s *UserService) UpdateUserSwagger() {}

// DeleteUser godoc
// @Summary Delete user (Admin only)
// @Description Delete user and associated profile (student/lecturer). This will also cascade delete related data.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} model.APIResponse "User deleted successfully"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Admin only"
// @Failure 404 {object} model.APIResponse "User not found"
// @Router /users/{id} [delete]
func (s *UserService) DeleteUserSwagger() {}

// AssignRole godoc
// @Summary Assign role to user (Admin only)
// @Description Change user's role. Note: Changing role does not automatically create/delete profiles.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Param request body model.AssignRoleRequest true "Role name (Admin, Mahasiswa, or Dosen Wali)"
// @Success 200 {object} model.APIResponse{data=model.UserResponse} "Role assigned successfully"
// @Failure 400 {object} model.APIResponse "Invalid request body"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Admin only"
// @Failure 404 {object} model.APIResponse "User or role not found"
// @Router /users/{id}/role [put]
func (s *UserService) AssignRoleSwagger() {}

// ==================== STUDENT SERVICE ANNOTATIONS ======================

// GetAllStudents godoc
// @Summary Get all students
// @Description Get list of all students with pagination. Admin and Mahasiswa can see all, Dosen Wali only sees their advisees.
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size (max 100)" default(10)
// @Success 200 {object} model.APIResponse{data=object} "List of students with user info and advisor details"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 404 {object} model.APIResponse "Lecturer profile not found (for Dosen Wali)"
// @Router /students [get]
func (s *StudentService) GetAllStudentsSwagger() {}

// GetStudentByID godoc
// @Summary Get student by ID
// @Description Get detailed student information including advisor. Mahasiswa can only view own profile, Dosen Wali can view advisees, Admin can view all.
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID (UUID)"
// @Success 200 {object} model.APIResponse{data=object} "Student details with user and advisor info"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Not authorized to view this student"
// @Failure 404 {object} model.APIResponse "Student not found"
// @Router /students/{id} [get]
func (s *StudentService) GetStudentByIDSwagger() {}

// GetStudentAchievements godoc
// @Summary Get student achievements
// @Description Get all achievements of a student with pagination and status filter. Authorization checks apply based on role.
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID (UUID)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size (max 100)" default(10)
// @Param status query string false "Filter by status" Enums(draft, submitted, verified, rejected)
// @Success 200 {object} model.APIResponse{data=object} "List of achievements with student info"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Not authorized to view this student's achievements"
// @Failure 404 {object} model.APIResponse "Student not found"
// @Router /students/{id}/achievements [get]
func (s *StudentService) GetStudentAchievementsSwagger() {}

// SetAdvisor godoc
// @Summary Set student advisor (Admin only)
// @Description Assign or change student's advisor (Dosen Wali)
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID (UUID)"
// @Param request body model.SetAdvisorRequest true "Advisor (Lecturer) ID"
// @Success 200 {object} model.APIResponse{data=object} "Advisor set successfully with updated student info"
// @Failure 400 {object} model.APIResponse "Invalid request body"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Admin only"
// @Failure 404 {object} model.APIResponse "Student or advisor (lecturer) not found"
// @Failure 422 {object} model.APIResponse "Validation error"
// @Router /students/{id}/advisor [put]
func (s *StudentService) SetAdvisorSwagger() {}

// ==================== LECTURER SERVICE ANNOTATIONS ======================

// GetAllLecturers godoc
// @Summary Get all lecturers
// @Description Get list of all lecturers with pagination and their user details
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size (max 100)" default(10)
// @Success 200 {object} model.APIResponse{data=object} "List of lecturers with user info"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Router /lecturers [get]
func (s *LecturerService) GetAllLecturersSwagger() {}

// GetLecturerAdvisees godoc
// @Summary Get lecturer's advisees
// @Description Get all students advised by this lecturer. Dosen Wali can only view own advisees, Admin can view all.
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lecturer ID (UUID)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size (max 100)" default(10)
// @Param include_achievements query boolean false "Include achievements summary" default(false)
// @Success 200 {object} model.APIResponse{data=object} "List of advisees with optional achievement summary"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Dosen Wali can only view own advisees"
// @Failure 404 {object} model.APIResponse "Lecturer not found"
// @Router /lecturers/{id}/advisees [get]
func (s *LecturerService) GetLecturerAdviseesSwagger() {}

// ==================== ACHIEVEMENT SERVICE ANNOTATIONS ======================

// CreateAchievement godoc
// @Summary Create achievement (Mahasiswa only)
// @Description Create new achievement draft. Status will be set to 'draft' initially.
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.AchievementCreateRequest true "Achievement data with dynamic details field"
// @Success 201 {object} model.APIResponse{data=model.AchievementResponse} "Achievement created successfully"
// @Failure 400 {object} model.APIResponse "Invalid request body"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Mahasiswa only"
// @Failure 404 {object} model.APIResponse "Student profile not found"
// @Failure 422 {object} model.APIResponse "Validation error"
// @Router /achievements [post]
func (s *AchievementService) CreateAchievementSwagger() {}

// GetAchievements godoc
// @Summary Get achievements (filtered by role)
// @Description Get achievements list with role-based filtering: Mahasiswa sees own achievements, Dosen Wali sees advisees' achievements, Admin sees all.
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size (max 100)" default(10)
// @Param status query string false "Filter by status" Enums(draft, submitted, verified, rejected)
// @Success 200 {object} model.APIResponse{data=model.AchievementListResponse} "List of achievements"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden"
// @Failure 404 {object} model.APIResponse "Profile not found (student/lecturer)"
// @Router /achievements [get]
func (s *AchievementService) GetAchievementsSwagger() {}

// GetAchievementByID godoc
// @Summary Get achievement by ID
// @Description Get detailed achievement information with authorization check based on role
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement Reference ID (UUID from PostgreSQL)"
// @Success 200 {object} model.APIResponse{data=model.AchievementResponse} "Achievement details"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Not authorized to view this achievement"
// @Failure 404 {object} model.APIResponse "Achievement not found"
// @Router /achievements/{id} [get]
func (s *AchievementService) GetAchievementByIDSwagger() {}

// UpdateAchievement godoc
// @Summary Update achievement (Mahasiswa only, draft status)
// @Description Update achievement data. Can only update if status is 'draft' and you are the owner.
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement Reference ID (UUID)"
// @Param request body model.AchievementUpdateRequest true "Update data (all fields optional)"
// @Success 200 {object} model.APIResponse{data=model.AchievementResponse} "Achievement updated successfully"
// @Failure 400 {object} model.APIResponse "Can only update achievements with status 'draft'"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Not your achievement"
// @Failure 404 {object} model.APIResponse "Achievement not found"
// @Router /achievements/{id} [put]
func (s *AchievementService) UpdateAchievementSwagger() {}

// DeleteAchievement godoc
// @Summary Delete achievement (Mahasiswa only, draft status)
// @Description Soft delete achievement. Can only delete if status is 'draft' and you are the owner.
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement Reference ID (UUID)"
// @Success 200 {object} model.APIResponse "Achievement deleted successfully"
// @Failure 400 {object} model.APIResponse "Can only delete achievements with status 'draft'"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Not your achievement"
// @Failure 404 {object} model.APIResponse "Achievement not found"
// @Router /achievements/{id} [delete]
func (s *AchievementService) DeleteAchievementSwagger() {}

// SubmitForVerification godoc
// @Summary Submit achievement for verification (Mahasiswa only)
// @Description Submit draft achievement to advisor for verification. Changes status from 'draft' to 'submitted'.
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement Reference ID (UUID)"
// @Success 200 {object} model.APIResponse{data=object} "Achievement submitted with new status and timestamp"
// @Failure 400 {object} model.APIResponse "Achievement already submitted or processed"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Not your achievement"
// @Failure 404 {object} model.APIResponse "Achievement not found"
// @Router /achievements/{id}/submit [post]
func (s *AchievementService) SubmitForVerificationSwagger() {}

// VerifyAchievement godoc
// @Summary Verify achievement (Dosen Wali only)
// @Description Approve submitted achievement. Can only verify if you are the advisor of the student and status is 'submitted'.
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement Reference ID (UUID)"
// @Success 200 {object} model.APIResponse{data=object} "Achievement verified with timestamp and verifier info"
// @Failure 400 {object} model.APIResponse "Achievement must be in 'submitted' status"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Not advisor of this student"
// @Failure 404 {object} model.APIResponse "Achievement not found"
// @Router /achievements/{id}/verify [post]
func (s *AchievementService) VerifyAchievementSwagger() {}

// RejectAchievement godoc
// @Summary Reject achievement (Dosen Wali only)
// @Description Reject submitted achievement with mandatory rejection note. Can only reject if you are the advisor and status is 'submitted'.
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement Reference ID (UUID)"
// @Param request body model.RejectAchievementRequest true "Rejection note (required)"
// @Success 200 {object} model.APIResponse{data=object} "Achievement rejected with note"
// @Failure 400 {object} model.APIResponse "Achievement must be in 'submitted' status"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Not advisor of this student"
// @Failure 404 {object} model.APIResponse "Achievement not found"
// @Failure 422 {object} model.APIResponse "Validation error - rejection note required"
// @Router /achievements/{id}/reject [post]
func (s *AchievementService) RejectAchievementSwagger() {}

// UploadAttachment godoc
// @Summary Upload attachment file (Mahasiswa only)
// @Description Upload file attachment to achievement. Accepts PDF, JPG, PNG with max size 5MB. Can upload for 'draft' or 'submitted' status.
// @Tags Achievements
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement Reference ID (UUID)"
// @Param file formData file true "File to upload (PDF, JPG, PNG, max 5MB)"
// @Success 201 {object} model.APIResponse{data=model.Attachment} "Attachment uploaded successfully"
// @Failure 400 {object} model.APIResponse "Invalid file, file too large, or file type not allowed"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Not your achievement"
// @Failure 404 {object} model.APIResponse "Achievement not found"
// @Router /achievements/{id}/attachments [post]
func (s *AchievementService) UploadAttachmentSwagger() {}

// GetAchievementHistory godoc
// @Summary Get achievement status history
// @Description Get timeline of achievement status changes (draft → submitted → verified/rejected). Includes actor info and notes.
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement Reference ID (UUID)"
// @Success 200 {object} model.APIResponse{data=object} "Achievement history timeline"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Not authorized to view this achievement"
// @Failure 404 {object} model.APIResponse "Achievement not found"
// @Router /achievements/{id}/history [get]
func (s *AchievementService) GetAchievementHistorySwagger() {}

// ==================== REPORT SERVICE ANNOTATIONS ======================

// GetStatistics godoc
// @Summary Get achievement statistics
// @Description Get comprehensive statistics based on role: Mahasiswa gets own stats, Dosen Wali gets advisees' stats, Admin gets all stats. Includes breakdown by type, period, status, and competition level.
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.APIResponse{data=model.AchievementStatistics} "Achievement statistics"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden"
// @Failure 404 {object} model.APIResponse "Profile not found (student/lecturer)"
// @Router /reports/statistics [get]
func (s *ReportService) GetStatisticsSwagger() {}

// GetStudentReport godoc
// @Summary Get student achievement report
// @Description Get comprehensive achievement report for a specific student including summary, type breakdown, recent achievements, and timeline. Mahasiswa can only view own report, Dosen Wali can view advisees' reports, Admin can view all.
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID (UUID)"
// @Success 200 {object} model.APIResponse{data=model.StudentReport} "Student report with all details"
// @Failure 401 {object} model.APIResponse "Unauthorized"
// @Failure 403 {object} model.APIResponse "Forbidden - Not authorized for this student"
// @Failure 404 {object} model.APIResponse "Student not found"
// @Router /reports/student/{id} [get]
func (s *ReportService) GetStudentReportSwagger() {}