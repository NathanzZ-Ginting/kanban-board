package service

import (
	"errors"
	"time"

	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/auth-service/config"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/auth-service/internal/domain"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/auth-service/internal/dto"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/auth-service/internal/repository"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/auth-service/pkg/jwt"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/auth-service/pkg/validator"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already exists")
	ErrUsernameExists     = errors.New("username already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(req *dto.LogoutRequest) error
	ValidateRegisterRequest(req *dto.RegisterRequest) []validator.ValidationError
	ValidateLoginRequest(req *dto.LoginRequest) []validator.ValidationError
}

type authService struct {
	userRepo   repository.UserRepository
	jwtService *jwt.JWTService
	cfg        *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	jwtService := jwt.NewJWTService(cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshExpiry)
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtService,
		cfg:        cfg,
	}
}

func (s *authService) ValidateRegisterRequest(req *dto.RegisterRequest) []validator.ValidationError {
	var errs []validator.ValidationError

	if !validator.IsNotEmpty(req.Email) {
		errs = append(errs, validator.ValidationError{Field: "email", Message: "Email is required"})
	} else if !validator.IsValidEmail(req.Email) {
		errs = append(errs, validator.ValidationError{Field: "email", Message: "Invalid email format"})
	}

	if !validator.IsNotEmpty(req.Username) {
		errs = append(errs, validator.ValidationError{Field: "username", Message: "Username is required"})
	} else if !validator.IsValidUsername(req.Username) {
		errs = append(errs, validator.ValidationError{Field: "username", Message: "Username must be 3-30 characters, alphanumeric and underscore only"})
	}

	if !validator.IsNotEmpty(req.Password) {
		errs = append(errs, validator.ValidationError{Field: "password", Message: "Password is required"})
	} else if !validator.IsValidPassword(req.Password) {
		errs = append(errs, validator.ValidationError{Field: "password", Message: "Password must be at least 8 characters"})
	}

	return errs
}

func (s *authService) ValidateLoginRequest(req *dto.LoginRequest) []validator.ValidationError {
	var errs []validator.ValidationError

	if !validator.IsNotEmpty(req.Email) {
		errs = append(errs, validator.ValidationError{Field: "email", Message: "Email is required"})
	} else if !validator.IsValidEmail(req.Email) {
		errs = append(errs, validator.ValidationError{Field: "email", Message: "Invalid email format"})
	}

	if !validator.IsNotEmpty(req.Password) {
		errs = append(errs, validator.ValidationError{Field: "password", Message: "Password is required"})
	}

	return errs
}

func (s *authService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if email exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailExists
	}

	// Check if username exists
	existingUser, err = s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUsernameExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
		IsActive:  true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate tokens
	return s.generateAuthResponse(user)
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	return s.generateAuthResponse(user)
}

func (s *authService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Find stored refresh token
	storedToken, err := s.userRepo.FindRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if storedToken == nil {
		return nil, ErrInvalidToken
	}

	// Check if token is expired
	if storedToken.ExpiresAt.Before(time.Now()) {
		s.userRepo.DeleteRefreshToken(req.RefreshToken)
		return nil, ErrInvalidToken
	}

	// Find user
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Delete old refresh token
	s.userRepo.DeleteRefreshToken(req.RefreshToken)

	// Generate new tokens
	return s.generateAuthResponse(user)
}

func (s *authService) Logout(req *dto.LogoutRequest) error {
	// Delete refresh token
	return s.userRepo.DeleteRefreshToken(req.RefreshToken)
}

func (s *authService) generateAuthResponse(user *domain.User) (*dto.AuthResponse, error) {
	// Generate access token
	accessToken, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	refreshTokenEntity := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshExpiry()),
	}
	if err := s.userRepo.SaveRefreshToken(refreshTokenEntity); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Avatar:    user.Avatar,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}
