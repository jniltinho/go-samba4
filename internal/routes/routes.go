package routes

import (
	"github.com/labstack/echo/v5"
	"go-samba4/internal/auth"
	"go-samba4/internal/handlers"
	"go-samba4/internal/middleware"
)

// RegisterRoutes sets up all the HTTP routes for the application
func RegisterRoutes(e *echo.Echo, appCtx *handlers.AppContext, sm *auth.SessionManager) {
	e.GET("/", func(c *echo.Context) error { return c.Redirect(302, "/dashboard") })

	// Public Auth Routes
	authGrp := e.Group("/auth", middleware.CSRF())
	authGrp.GET("/login", appCtx.AuthLoginGET)
	authGrp.POST("/login", appCtx.AuthLoginPOST, middleware.RateLimit())
	authGrp.GET("/logout", appCtx.AuthLogout)

	// Protected Routes
	p := e.Group("", middleware.RequireAuth(sm), middleware.CSRF())
	p.GET("/dashboard", appCtx.DashboardGET)
	p.GET("/users", appCtx.UsersListGET)
	p.GET("/users/new", appCtx.UsersFormGET, middleware.RequireAdmin())
	p.POST("/users/new", appCtx.UsersCreatePOST, middleware.RequireAdmin())
	p.GET("/users/:id", appCtx.UsersDetailGET)
	p.GET("/users/:id/edit", appCtx.UsersFormGET, middleware.RequireAdminOrSelf())
	p.POST("/users/:id/edit", appCtx.UsersUpdatePOST, middleware.RequireAdminOrSelf())
	p.POST("/users/:id/delete", appCtx.UsersDeletePOST, middleware.RequireAdmin())

	p.GET("/groups", appCtx.GroupsListGET)
	p.GET("/groups/new", appCtx.GroupsFormGET)

	p.GET("/ous", appCtx.OUsTreeGET)

	p.GET("/search", appCtx.SearchGET)

	p.GET("/audit", appCtx.AuditListGET)
	p.GET("/settings", appCtx.SettingsGET)
	p.POST("/settings", appCtx.SettingsPOST)
}
