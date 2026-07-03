package services

import (
	"context"
	"errors"
	"ppdb-be/models"
	"ppdb-be/models/request"
	"ppdb-be/repositories"
	"ppdb-be/utils"
	"strings"
	"time"
)

type UserService interface {
	Login(ctx context.Context, email, password, ip, userAgent string) (map[string]interface{}, error)
	CreateUser(ctx context.Context, req *request.CreateUserRequest) (*models.User, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, req *request.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	Logout(ctx context.Context, token string) error
}

type userService struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
}

func NewUserService(userRepo repositories.UserRepository, sessionRepo repositories.SessionRepository) UserService {
	return &userService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func parseUserAgent(ua string) (platform, browser, device string) {
	if ua == "" {
		return "Unknown", "Unknown", "Unknown"
	}
	lowerUA := strings.ToLower(ua)

	if strings.Contains(lowerUA, "edg") {
		browser = "Edge"
	} else if strings.Contains(lowerUA, "chrome") {
		browser = "Chrome"
	} else if strings.Contains(lowerUA, "firefox") {
		browser = "Firefox"
	} else if strings.Contains(lowerUA, "safari") {
		browser = "Safari"
	} else {
		browser = "Unknown"
	}

	if strings.Contains(lowerUA, "windows") {
		platform = "Windows"
		device = "PC"
	} else if strings.Contains(lowerUA, "macintosh") || strings.Contains(lowerUA, "mac os") {
		platform = "macOS"
		device = "Mac"
	} else if strings.Contains(lowerUA, "android") {
		platform = "Android"
		device = "Mobile"
	} else if strings.Contains(lowerUA, "iphone") || strings.Contains(lowerUA, "ipad") {
		platform = "iOS"
		if strings.Contains(lowerUA, "ipad") {
			device = "iPad"
		} else {
			device = "iPhone"
		}
	} else if strings.Contains(lowerUA, "linux") {
		platform = "Linux"
		device = "PC"
	} else {
		platform = "Unknown"
		device = "Unknown"
	}

	return platform, browser, device
}

func (s *userService) Login(ctx context.Context, email, password, ip, userAgent string) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("Email atau password salah")
	}

	if err := utils.ComparePassword(user.UserPassword, password); err != nil {
		return nil, errors.New("Email atau password salah")
	}

	accessToken, err := utils.GenerateToken(user.UserID, user.UserName, user.UserEmail, 24*time.Hour)
	if err != nil {
		return nil, errors.New("Gagal membuat access token")
	}

	refreshToken, err := utils.GenerateToken(user.UserID, user.UserName, user.UserEmail, 7*24*time.Hour)
	if err != nil {
		return nil, errors.New("Gagal membuat refresh token")
	}

	platform, browser, device := parseUserAgent(userAgent)
	now := time.Now()
	expiredAt := now.Add(7 * 24 * time.Hour)

	// Revoke old sessions
	_ = s.sessionRepo.RevokeActiveSessions(ctx, user.UserID, now)

	session := &models.UserSession{
		UserID:       user.UserID,
		AccessToken:  accessToken,
		RefreshToken: &refreshToken,
		LoginAt:      now,
		ExpiredAt:    expiredAt,
		IPAddress:    &ip,
		UserAgent:    &userAgent,
		DeviceName:   &device,
		Platform:     &platform,
		Browser:      &browser,
		IsActive:     true,
		IsRevoked:    false,
		CreatedOn:    now,
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, errors.New("Gagal menyimpan sesi login")
	}

	return map[string]interface{}{
		"session_id":    session.SessionID,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, req *request.CreateUserRequest) (*models.User, error) {
	exists, err := s.userRepo.CheckExists(ctx, req.UserName, req.UserEmail)
	if err != nil {
		return nil, errors.New("Gagal memvalidasi data user")
	}
	if exists {
		return nil, errors.New("Username atau Email sudah terdaftar")
	}

	hashedPassword, err := utils.HashPassword(req.UserPassword)
	if err != nil {
		return nil, errors.New("Gagal mengamankan password")
	}

	user := &models.User{
		RoleID:       req.RoleID,
		UserFullName: req.UserFullName,
		UserName:     req.UserName,
		UserPhone:    req.UserPhone,
		UserEmail:    req.UserEmail,
		UserPassword: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.New("Gagal menyimpan data user")
	}

	return user, nil
}

func (s *userService) GetUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.FindAll(ctx)
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id int64, req *request.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("User tidak ditemukan")
	}

	user.RoleID = req.RoleID
	user.UserFullName = req.UserFullName
	user.UserName = req.UserName
	user.UserPhone = req.UserPhone
	user.UserEmail = req.UserEmail

	if req.UserPassword != "" {
		hashed, err := utils.HashPassword(req.UserPassword)
		if err != nil {
			return nil, errors.New("Gagal mengamankan password")
		}
		user.UserPassword = hashed
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, errors.New("Gagal memperbarui data user")
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("User tidak ditemukan")
	}

	user.IsDeleted = true

	if err := s.userRepo.Update(ctx, user); err != nil {
		return errors.New("Gagal menghapus data user")
	}

	return nil
}

func (s *userService) Logout(ctx context.Context, token string) error {
	session, err := s.sessionRepo.FindActiveSessionByToken(ctx, token)
	if err != nil {
		// If not found, just return success since it's already logged out or invalid
		return nil
	}

	now := time.Now()
	session.IsActive = false
	session.IsRevoked = true
	session.LogoutAt = &now

	if err := s.sessionRepo.UpdateSession(ctx, session); err != nil {
		return errors.New("Gagal mengakhiri sesi")
	}

	return nil
}
