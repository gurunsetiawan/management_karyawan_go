package service

import (
	"errors"
	"log"

	"karyawan-app/internal/auth"
	"karyawan-app/internal/models"
	"karyawan-app/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Login memvalidasi kredensial user dan mengembalikan token JWT jika valid
func (s *AuthService) Login(username, password string) (string, *models.UserResponse, error) {
	// Dapatkan user dari database
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid username or password")
	}

	// Periksa apakah user aktif
	if !user.IsActive {
		return "", nil, errors.New("account is not active")
	}

	// Verifikasi password
	err = auth.CheckPassword(password, user.PasswordHash)
	if err != nil {
		return "", nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role.Name)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	// Update last login
	err = s.userRepo.UpdateLastLogin(user.ID)
	if err != nil {
		// Log error tapi jangan return error ke client
		// Karena login sebenarnya sudah berhasil
		// Hanya gagal update last login
		log.Printf("Warning: Failed to update last login for user %s: %v", user.Username, err)
	}

	// Siapkan response
	userResponse := &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.Name,
		IsActive:  user.IsActive,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
	}

	return token, userResponse, nil
}

// GetUserProfile mengambil profil user berdasarkan ID
func (s *AuthService) GetUserProfile(userID int) (*models.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.Name,
		IsActive:  user.IsActive,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
	}, nil
}

// CreateUser membuat user baru
func (s *AuthService) CreateUser(user *models.User, password string) (*models.UserResponse, error) {
	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Set password hash
	user.PasswordHash = hashedPassword

	// Set default role jika tidak diset
	if user.RoleID == 0 {
		user.RoleID = 3 // Default role: karyawan
	}

	// Set is_active default true
	user.IsActive = true

	// Simpan ke database
	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// Dapatkan data user yang baru dibuat
	createdUser, err := s.userRepo.GetUserByID(user.ID)
	if err != nil {
		return nil, err
	}

	// Siapkan response
	return &models.UserResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		Role:      createdUser.Role.Name,
		IsActive:  createdUser.IsActive,
		LastLogin: createdUser.LastLogin,
		CreatedAt: createdUser.CreatedAt,
	}, nil
}
