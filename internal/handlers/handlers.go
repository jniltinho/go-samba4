package handlers

import (
	"gorm.io/gorm"
	"go-samba4/internal/auth"
	"go-samba4/internal/config"
	"go-samba4/internal/ldap"
)

// AppContext holds dependencies for the handlers
type AppContext struct {
	Config     *config.Config
	DB         *gorm.DB
	LDAPClient *ldap.Client
	SessionMgr *auth.SessionManager
}
