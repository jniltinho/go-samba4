package handlers

import (
	"gorm.io/gorm"
	"samba4-admin/internal/auth"
	"samba4-admin/internal/config"
	"samba4-admin/internal/ldap"
)

// AppContext holds dependencies for the handlers
type AppContext struct {
	Config     *config.Config
	DB         *gorm.DB
	LDAPClient *ldap.Client
	SessionMgr *auth.SessionManager
}
