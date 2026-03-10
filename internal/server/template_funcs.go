package server

import (
	"fmt"
	"html/template"

	"go-samba4/internal/buildinfo"
	"go-samba4/internal/i18n"
)

// TemplateFuncMap returns the global template function map available to all templates.
func TemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		// i18n: {{ T $.Lang "Key" }}
		"T": func(lang, msgID string) string {
			return i18n.Translate(lang, msgID, nil)
		},
		// i18n with data substitution: {{ TData $.Lang "Key" (dict "Name" .User.Name) }}
		"TData": func(lang, msgID string, data map[string]any) string {
			return i18n.Translate(lang, msgID, data)
		},
		// Application version: {{ version }}
		"version": func() string {
			return buildinfo.Version
		},
		// Render trusted HTML: {{ unescapeHTML .Content }}
		"unescapeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		// Render trusted CSS: {{ safeCSS .Style }}
		"safeCSS": func(s string) template.CSS {
			return template.CSS(s)
		},
		// Build a map from key-value pairs (for TData): {{ dict "Key" .Value "Key2" .Value2 }}
		"dict": func(kvs ...any) (map[string]any, error) {
			if len(kvs)%2 != 0 {
				return nil, fmt.Errorf("dict requires an even number of arguments")
			}
			m := make(map[string]any, len(kvs)/2)
			for i := 0; i < len(kvs); i += 2 {
				key, ok := kvs[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				m[key] = kvs[i+1]
			}
			return m, nil
		},
	}
}
