package repositories

import (
	"time"

	"github.com/tacheraSasi/go-api-starter/internals/models"
	"gorm.io/gorm"
)

type TokenRepository interface {
	BlacklistToken(token *models.BlacklistedToken) error
	IsTokenBlacklisted(token string) (bool, error)
	CreatePasswordResetToken(token *models.PasswordResetToken) error
	GetValidPasswordResetToken(token string) (*models.PasswordResetToken, error)
	MarkPasswordResetTokenUsed(token *models.PasswordResetToken) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) BlacklistToken(token *models.BlacklistedToken) error {
	return r.db.Create(token).Error
}

func (r *tokenRepository) IsTokenBlacklisted(token string) (bool, error) {
	var count int64
	err := r.db.Model(&models.BlacklistedToken{}).Where("token = ?", token).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *tokenRepository) CreatePasswordResetToken(token *models.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *tokenRepository) GetValidPasswordResetToken(token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.db.Where("token = ? AND used_at IS NULL AND expires_at > NOW()", token).First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

func (r *tokenRepository) MarkPasswordResetTokenUsed(token *models.PasswordResetToken) error {
	now := time.Now()
	token.UsedAt = &now
	return r.db.Save(token).Error
}
