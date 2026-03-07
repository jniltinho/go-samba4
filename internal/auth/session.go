package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"gorm.io/gorm"
	"go-samba4/internal/config"
	"go-samba4/internal/models"
)

const SessionCookieName = "samba4_admin_session"

type SessionManager struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewSessionManager(db *gorm.DB, cfg *config.Config) *SessionManager {
	return &SessionManager{db: db, cfg: cfg}
}

func (sm *SessionManager) CreateSession(username, ip, userAgent string, isAdmin bool) (*models.Session, error) {
	tokenBytes := asSecret(64)
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	expiresAt := time.Now().Add(time.Duration(sm.cfg.Session.TimeoutMinutes) * time.Minute)

	session := &models.Session{
		Token:     token,
		Username:  username,
		IsAdmin:   isAdmin,
		IPAddress: ip,
		UserAgent: userAgent,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := sm.db.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (sm *SessionManager) GetSession(token string) (*models.Session, error) {
	var session models.Session
	// Remove expired sessions implicitly or by querying only valid ones
	if err := sm.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (sm *SessionManager) DeleteSession(token string) error {
	return sm.db.Where("token = ?", token).Delete(&models.Session{}).Error
}

func (sm *SessionManager) SetCookie(w http.ResponseWriter, token string) {
	sameSite := http.SameSiteStrictMode
	if sm.cfg.Session.CookieSameSite == "lax" {
		sameSite = http.SameSiteLaxMode
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   sm.cfg.Session.CookieSecure,
		SameSite: sameSite,
		Expires:  time.Now().Add(time.Duration(sm.cfg.Session.TimeoutMinutes) * time.Minute),
	})
}

func (sm *SessionManager) ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   sm.cfg.Session.CookieSecure,
		MaxAge:   -1,
	})
}

func asSecret(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}
