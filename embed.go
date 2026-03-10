package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed all:web/templates
var TemplatesFS embed.FS

//go:embed all:web/static
var StaticFS embed.FS

//go:embed all:locales
var LocalesFS embed.FS

func loadTemplates() *template.Template {
	return template.Must(
		template.New("").ParseFS(TemplatesFS, "web/templates/**/*.html", "web/templates/*.html"),
	)
}

func staticHandler() http.Handler {
	sub, _ := fs.Sub(StaticFS, "web/static")
	return http.FileServer(http.FS(sub))
}
