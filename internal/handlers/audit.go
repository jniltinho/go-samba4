package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"go-samba4/internal/models"
)

// AuditListGET renders the full audit log page.
func (app *AppContext) AuditListGET(c *echo.Context) error {
	var logs []models.AuditLog
	app.DB.Order("created_at desc").Find(&logs)
	return c.Render(http.StatusOK, "audit/list", map[string]interface{}{
		"AuditLogs": logs,
	})
}
