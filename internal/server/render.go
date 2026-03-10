package server

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/labstack/echo/v5"
	"go-samba4/internal/config"
	"go-samba4/internal/handlers"
)

// TemplateRegistry is a custom html/template renderer for Echo framework
type TemplateRegistry struct {
	Templates map[string]*template.Template
}

// NewTemplateRegistry initializes and creates the map of all templates
func NewTemplateRegistry(cfg *config.Config, tplFS embed.FS) (*TemplateRegistry, error) {
	t := &TemplateRegistry{
		Templates: make(map[string]*template.Template),
	}

	var fileSystem fs.FS
	if cfg.Server.DevMode {
		fileSystem = os.DirFS(".") // we will walk "web/templates" relatively from project root
	} else {
		fileSystem = tplFS // embedded filesystem also has the native "web/templates" folder
	}

	layout := "web/templates/layout/base.html"
	sidebar := "web/templates/layout/sidebar.html"

	err := fs.WalkDir(fileSystem, "web/templates", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".html") {
			return nil
		}

		// skip layouts as main targets
		if filePath == layout || filePath == sidebar {
			return nil
		}

		// Calculate mapping key (e.g. "web/templates/auth/login.html" -> "auth/login")
		tmplKey := strings.TrimPrefix(filePath, "web/templates/")
		tmplKey = strings.TrimSuffix(tmplKey, ".html")

		// Each page template gets parsed together with the layout foundations
		// Funcs must be registered before parsing
		tmpl := template.New(path.Base(filePath)).Funcs(TemplateFuncMap())
		tmpl, parseErr := tmpl.ParseFS(fileSystem, layout, sidebar, filePath)

		if parseErr != nil {
			return fmt.Errorf("failed to parse template %s: %w", filePath, parseErr)
		}

		t.Templates[tmplKey] = tmpl
		return nil
	})

	if err != nil {
		return nil, err
	}

	return t, nil
}

// Render implements echo.Renderer interface
func (t *TemplateRegistry) Render(c *echo.Context, w io.Writer, name string, data any) error {
	tmpl, ok := t.Templates[name]
	if !ok {
		return fmt.Errorf("Template not found: %s", name)
	}

	// Add global data here like CSRF token
	viewData := map[string]interface{}{}
	if data != nil {
		if d, ok := data.(map[string]interface{}); ok {
			viewData = d
		}
	}
	viewData["CSRFToken"] = c.Get("csrf")
	viewData["Username"] = c.Get("username")
	viewData["Lang"] = handlers.LangFromRequest(c)
	viewData["CurrentPath"] = c.Request().URL.Path

	// Render the primary 'base' template layout which intrinsically yields the corresponding sub-content
	return tmpl.ExecuteTemplate(w, "base", viewData)
}
