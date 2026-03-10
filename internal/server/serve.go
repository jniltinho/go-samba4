package server

import (
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"go-samba4/internal/auth"
	"go-samba4/internal/config"
	"go-samba4/internal/handlers"
	"go-samba4/internal/i18n"
	"go-samba4/internal/ldap"
	"go-samba4/internal/models"
	"go-samba4/internal/routes"
)

// Serve initializes Echo and its dependencies before starting the server
func Serve(globalCfg *config.Config, tplFS embed.FS, statFS embed.FS, localesFS embed.FS, debug bool) {
	// 1. Database Connection
	var dialector gorm.Dialector
	if globalCfg.Database.Driver == "mysql" {
		dialector = mysql.Open(globalCfg.Database.DSN)
	} else {
		dialector = sqlite.Open(globalCfg.Database.Path)
	}

	gormCfg := &gorm.Config{}
	if debug {
		gormCfg.Logger = gormLogger.Default.LogMode(gormLogger.Info)
	}
	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		slog.Error("Failed to connect to database", "err", err)
		os.Exit(1)
	}

	// Auto-migrate tables on start
	if err := db.AutoMigrate(&models.Session{}, &models.AuditLog{}, &models.Setting{}); err != nil {
		slog.Error("Failed to auto-migrate database", "err", err)
		os.Exit(1)
	}

	// 2. LDAP Connection
	ldapClient, err := ldap.NewClient(&globalCfg.LDAP)
	if err != nil {
		slog.Error("Failed to establish initial LDAP connection. Will continue starting...", "err", err)
	} else {
		slog.Info("Connected to AD LDAP successfully.")
		defer ldapClient.Close()
	}

	// 3. Setup Session Manager
	sm := auth.NewSessionManager(db, globalCfg)

	// 4. Init i18n
	i18n.Init(localesFS)

	// 5. Echo Instance setup
	if debug {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
		slog.Info("Debug mode enabled")
	}
	e := echo.New()

	e.Use(echoMiddleware.RequestLogger())
	e.Use(echoMiddleware.Recover())

	// Templates rendering mapping
	tmplRegistry, err := NewTemplateRegistry(globalCfg, tplFS)
	if err != nil {
		slog.Error("Failed to initialize template registry", "err", err)
		os.Exit(1)
	}
	e.Renderer = tmplRegistry

	// 5. Static files routing
	if globalCfg.Server.DevMode {
		e.Static("/static", "web/static")
	} else {
		e.StaticFS("/static", echo.MustSubFS(statFS, "web/static"))
	}

	// 6. Context wiring
	appCtx := &handlers.AppContext{
		Config:     globalCfg,
		DB:         db,
		LDAPClient: ldapClient,
		SessionMgr: sm,
	}

	// 7. Routes Definitions
	routes.RegisterRoutes(e, appCtx, sm)

	// 8. Start server
	bindAddr := fmt.Sprintf("%s:%d", globalCfg.Server.Host, globalCfg.Server.Port)
	if globalCfg.Server.TLSCert != "" && globalCfg.Server.TLSKey != "" {
		slog.Info(fmt.Sprintf("Server starting on https://%s", bindAddr))
		if err := e.Start(bindAddr); err != nil { // In some v5 versionsStart also handles TLS or you use http.Server
			slog.Error("Server failed", "err", err)
			os.Exit(1)
		}
	} else {
		slog.Info(fmt.Sprintf("Server starting on http://%s", bindAddr))
		if err := e.Start(bindAddr); err != nil {
			slog.Error("Server failed", "err", err)
			os.Exit(1)
		}
	}
}
