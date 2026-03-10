package i18n

import (
	"bytes"
	"embed"
	"io/fs"
	"log/slog"
	"strings"
	"text/template"

	"github.com/leonelquinteros/gotext"
)

var locales map[string]*gotext.Po

// Init loads all PO files from locales/*/default.po in the given embedded FS.
func Init(fsys embed.FS) {
	locales = make(map[string]*gotext.Po)

	entries, err := fs.ReadDir(fsys, "locales")
	if err != nil {
		slog.Warn("i18n: could not read locales directory", "err", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		lang := entry.Name()
		poPath := "locales/" + lang + "/default.po"

		data, err := fsys.ReadFile(poPath)
		if err != nil {
			slog.Warn("i18n: could not read PO file", "path", poPath, "err", err)
			continue
		}

		po := gotext.NewPo()
		po.Parse(data)
		locales[lang] = po
		slog.Debug("i18n: loaded locale", "lang", lang, "path", poPath)
	}
}

// Translate returns the translated string for the given language and message ID.
// If templateData is provided, {{.Key}} placeholders in the translated string are substituted.
// Falls back to msgID if translation is not found.
func Translate(lang, msgID string, data map[string]any) string {
	normalized := normalizeLang(lang)

	po, ok := locales[normalized]
	if !ok {
		// try "en" as fallback
		po, ok = locales["en"]
		if !ok {
			return msgID
		}
	}

	translated := po.Get(msgID)
	if translated == "" || translated == msgID {
		// try English fallback before giving up
		if normalized != "en" {
			if enPo, ok := locales["en"]; ok {
				translated = enPo.Get(msgID)
			}
		}
		if translated == "" {
			return msgID
		}
	}

	if data == nil {
		return translated
	}

	// Apply Go template substitution for {{.Var}} placeholders
	tmpl, err := template.New("").Parse(translated)
	if err != nil {
		return translated
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return translated
	}
	return buf.String()
}

// normalizeLang maps incoming language codes to PO directory names.
func normalizeLang(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	switch {
	case lang == "pt" || lang == "pt-br" || lang == "pt_br":
		return "pt_BR"
	case strings.HasPrefix(lang, "pt"):
		return "pt_BR"
	case strings.HasPrefix(lang, "es"):
		return "es"
	default:
		return "en"
	}
}
