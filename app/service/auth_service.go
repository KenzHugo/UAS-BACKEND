package service

import (
	"errors"
	"UASBE/app/repository"
	"UASBE/model"
	"UASBE/utils"
	"strings"
)

type AuthService interface {
	Login(req model.LoginRequest) (*model.LoginResponse, error)
	RefreshToken(refreshToken string) (*model.LoginResponse, error)
	GetProfile(userID string) (*model.UserProfile, error)
}

type authService struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}

// FR-001: Login
// Flow:
// 1. User mengirim kredensial
// 2. Sistem memvalidasi kredensial
// 3. Sistem mengecek status aktif user
// 4. Sistem generate JWT token dengan role dan permissions
// 5. Return token dan user profile
func (s *authService) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	// Step 1: Find user by username or email
	var user *model.User
	var err error

	// Check if input is email or username
	if strings.Contains(req.Username, "@") {
		user, err = s.authRepo.FindUserByEmail(req.Username)
	} else {
		user, err = s.authRepo.FindUserByUsername(req.Username)
	}

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Step 2: Validate password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Step 3: Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Get user with role and permissions
	userWithRole, err := s.authRepo.GetUserWithRole(user.ID)
	if err != nil {
		return nil, err
	}

	// Step 4: Generate JWT tokens
	token, err := utils.GenerateToken(user.ID, userWithRole.RoleName, userWithRole.Permissions)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Step 5: Return response
	response := &model.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: model.UserProfile{
			ID:          user.ID,
			Username:    user.Username,
			FullName:    user.FullName,
			Email:       user.Email,
			Role:        userWithRole.RoleName,
			Permissions: userWithRole.Permissions,
		},
	}

	return response, nil
}

func (s *authService) RefreshToken(refreshToken string) (*model.LoginResponse, error) {
	// Validate refresh token
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Get user
	userWithRole, err := s.authRepo.GetUserWithRole(claims.UserID)
	if err != nil {
		return nil, err
	}

	// Check if user is active
	if !userWithRole.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Generate new tokens
	newToken, err := utils.GenerateToken(userWithRole.ID, userWithRole.RoleName, userWithRole.Permissions)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := utils.GenerateRefreshToken(userWithRole.ID)
	if err != nil {
		return nil, err
	}

	response := &model.LoginResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
		User: model.UserProfile{
			ID:          userWithRole.ID,
			Username:    userWithRole.Username,
			FullName:    userWithRole.FullName,
			Email:       userWithRole.Email,
			Role:        userWithRole.RoleName,
			Permissions: userWithRole.Permissions,
		},
	}

	return response, nil
}

func (s *authService) GetProfile(userID string) (*model.UserProfile, error) {
	userWithRole, err := s.authRepo.GetUserWithRole(userID)
	if err != nil {
		return nil, err
	}

	profile := &model.UserProfile{
		ID:          userWithRole.ID,
		Username:    userWithRole.Username,
		FullName:    userWithRole.FullName,
		Email:       userWithRole.Email,
		Role:        userWithRole.RoleName,
		Permissions: userWithRole.Permissions,
	}

	return profile, nil
}