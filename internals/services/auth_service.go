package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/repositories"
)

type AuthService interface {
	Login(email, password string) (models.User, error)
	Register(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	Logout(token string, expiresAt time.Time) error
	RequestPasswordReset(email string) (string, error)
	ResetPassword(token, password string) error
}

type authService struct {
	repo         repositories.UserRepository
	tokenService TokenService
}

func NewAuthService(repo repositories.UserRepository, tokenService TokenService) AuthService {
	return &authService{repo: repo, tokenService: tokenService}
}

func (s *authService) Login(email, password string) (models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return models.User{}, err
	}
	if err := user.CheckPassword(password); err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func (s *authService) Register(user *models.User) error {
	var existingUser *models.User
	existingUser, _ = s.repo.GetUserByEmail(user.Email)
	if existingUser != nil {
		return http.ErrBodyNotAllowed
	}
	if err := user.HashPassword(); err != nil {
		return err
	}
	return s.repo.CreateUser(user)
}

func (s *authService) GetUserByID(id string) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *authService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *authService) Logout(token string, expiresAt time.Time) error {
	return s.tokenService.BlacklistToken(token, expiresAt)
}

func (s *authService) RequestPasswordReset(email string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", nil
	}

	resetToken, err := generateSecureToken(32)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(30 * time.Minute)
	if err := s.tokenService.CreatePasswordResetToken(user.ID, resetToken, expiresAt); err != nil {
		return "", err
	}

	return resetToken, nil
}

func (s *authService) ResetPassword(token, password string) error {
	resetToken, err := s.tokenService.GetValidPasswordResetToken(token)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	user, err := s.repo.GetUserByID(strconv.FormatUint(uint64(resetToken.UserID), 10))
	if err != nil {
		return err
	}

	user.Password = password
	if err := user.HashPassword(); err != nil {
		return err
	}

	if err := s.repo.UpdateUser(user); err != nil {
		return err
	}

	if err := s.tokenService.MarkPasswordResetTokenUsed(resetToken); err != nil {
		return err
	}

	return nil
}

func generateSecureToken(size int) (string, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
