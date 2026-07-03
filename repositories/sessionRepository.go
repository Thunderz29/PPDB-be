package repositories

import (
	"context"
	"ppdb-be/models"
	"time"

	"gorm.io/gorm"
)

type SessionRepository interface {
	RevokeActiveSessions(ctx context.Context, userID int64, now time.Time) error
	CreateSession(ctx context.Context, session *models.UserSession) error
	FindActiveSessionByToken(ctx context.Context, token string) (*models.UserSession, error)
	UpdateSession(ctx context.Context, session *models.UserSession) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) RevokeActiveSessions(ctx context.Context, userID int64, now time.Time) error {
	return r.db.WithContext(ctx).Model(&models.UserSession{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Updates(map[string]interface{}{
			"is_active":  false,
			"is_revoked": true,
			"logout_at":  &now,
		}).Error
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *models.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepository) FindActiveSessionByToken(ctx context.Context, token string) (*models.UserSession, error) {
	var session models.UserSession
	if err := r.db.WithContext(ctx).Where("access_token = ? AND is_active = ?", token, true).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) UpdateSession(ctx context.Context, session *models.UserSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}
