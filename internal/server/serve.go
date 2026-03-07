package server

import (
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go-samba4/internal/auth"
	"go-samba4/internal/config"
	"go-samba4/internal/handlers"
	"go-samba4/internal/ldap"
	"go-samba4/internal/middleware"
	"go-samba4/internal/models"
	"go-samba4/internal/routes"
)

// Serve initializes Echo and its dependencies before starting the server
func Serve(globalCfg *config.Config, tplFS embed.FS, statFS embed.FS) {
	// 1. Database Connection
	var dialector gorm.Dialector
	if globalCfg.Database.Driver == "mysql" {
		dialector = mysql.Open(globalCfg.Database.DSN)
	} else {
		dialector = sqlite.Open(globalCfg.Database.Path)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
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

	// 4. Echo Instance setup
	e := echo.New()
	e.HideBanner = true

	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(middleware.RateLimit())

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
		e.Logger.Fatal(e.StartTLS(bindAddr, globalCfg.Server.TLSCert, globalCfg.Server.TLSKey))
	} else {
		slog.Info(fmt.Sprintf("Server starting on http://%s", bindAddr))
		e.Logger.Fatal(e.Start(bindAddr))
	}
}
