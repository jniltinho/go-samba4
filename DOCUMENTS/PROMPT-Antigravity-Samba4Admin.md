# Prompt — Create Samba4 AD Web Admin Panel project in Antigravity

---

## Project Context

Create a complete project named **`go-samba4`** — a web administration panel for **Samba 4 Active Directory**, written in **Go 1.26**. It is a single, self-sufficient binary that embeds all HTML templates and static assets into the executable via `//go:embed`. There is no SPA — everything is SSR with Go's `html/template`.

---

## Mandatory Stack

*   **Go 1.26** — primary language
*   **Echo v4** — HTTP server, routing, middleware
*   **GORM v2** — ORM with SQLite and MySQL/MariaDB support
*   **Cobra + Viper** — CLI and configuration via `config.toml`
*   **go-ldap/ldap v3** (`github.com/go-ldap/ldap/v3`) — LDAP integration with Samba 4 AD
*   **gokrb5 v8** (`github.com/jcmturner/gokrb5/v8`) — Kerberos/SPNEGO authentication
*   **Tailwind CSS 4.2** — styling (compiled to `web/static/css/app.css`)
*   **jQuery 4.0.0** — DOM interactivity and AJAX
*   **Go `embed` + `html/template`** (stdlib) — templates and assets embedded in the binary

---

## Directory Structure

Create exactly this structure:

```
go-samba4/
├── cmd/
│   ├── root.go          # Cobra root + Viper (reads config.toml)
│   ├── serve.go         # go-samba4 serve [--port] [--config]
│   ├── migrate.go       # go-samba4 migrate
│   └── user.go          # go-samba4 user create/list
├── internal/
│   ├── auth/
│   │   ├── ldap.go      # LDAP bind — user authentication
│   │   ├── kerberos.go  # SPNEGO middleware for Windows SSO
│   │   ├── totp.go      # 2FA TOTP (RFC 6238)
│   │   └── session.go   # HTTPOnly Sessions + CSRF token
│   ├── config/
│   │   └── config.go    # Configuration structs + Viper loader
│   ├── handlers/
│   │   ├── dashboard.go
│   │   ├── users.go
│   │   ├── groups.go
│   │   ├── ous.go
│   │   ├── search.go
│   │   └── settings.go
│   ├── ldap/
│   │   ├── client.go    # go-ldap connection pool
│   │   ├── users.go     # AD user CRUD
│   │   ├── groups.go    # AD group CRUD
│   │   ├── ous.go       # Read/navigate OUs
│   │   └── schema.go    # AD attributes ↔ Go structs mapping
│   ├── models/
│   │   ├── session.go   # GORM: web sessions table
│   │   ├── audit.go     # GORM: audit log (append-only)
│   │   └── setting.go   # GORM: persisted local settings
│   └── middleware/
│       ├── auth.go      # Checks valid session on protected routes
│       ├── csrf.go      # Double CSRF token on POST/PUT/DELETE
│       ├── rbac.go      # Checks user role by route
│       └── ratelimit.go # Max 5 login attempts per IP / 5 min
├── web/
│   ├── templates/
│   │   ├── layout/
│   │   │   ├── base.html     # Base: <head>, scripts, meta, nav
│   │   │   └── sidebar.html  # 240px Sidebar with navigation menu
│   │   ├── dashboard.html
│   │   ├── auth/
│   │   │   └── login.html
│   │   ├── users/
│   │   │   ├── list.html
│   │   │   ├── form.html
│   │   │   └── detail.html
│   │   ├── groups/
│   │   │   ├── list.html
│   │   │   └── form.html
│   │   ├── ous/
│   │   │   └── tree.html
│   │   └── audit/
│   │       └── list.html
│   └── static/
│       ├── css/
│       │   └── app.css       # Compiled Tailwind CSS 4.2
│       └── js/
│           └── app.js        # jQuery 4.0.0 + helpers
├── embed.go             # //go:embed — main package, project root
├── main.go
├── config.toml
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

## File `embed.go` (root, main package)

```go
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

func loadTemplates() *template.Template {
    return template.Must(
        template.New("").ParseFS(TemplatesFS, "web/templates/**/*.html"),
    )
}

func staticHandler() http.Handler {
    sub, _ := fs.Sub(StaticFS, "web/static")
    return http.FileServer(http.FS(sub))
}
```

---

## File `config.toml`

```toml
[server]
host     = "0.0.0.0"
port     = 8080
tls_cert = ""
tls_key  = ""
dev_mode = false   # true = reads templates from disk (hot-reload in dev)

[ldap]
host            = "dc1.empresa.local"
port            = 636
use_tls         = true
skip_tls_verify = false
base_dn         = "DC=empresa,DC=local"
bind_user       = "CN=samba4admin,CN=Users,DC=empresa,DC=local"
# password via env: SAMBA4_LDAP_PASS

[database]
# SQLite (default)
driver = "sqlite"
path   = "/var/lib/go-samba4/data.db"
# MySQL/MariaDB:
# driver = "mysql"
# dsn    = "user:pass@tcp(localhost:3306)/samba4admin?charset=utf8mb4&parseTime=True&loc=Local"

[session]
secret           = ""   # auto-generated on first start if empty
timeout_minutes  = 30
cookie_secure    = true
cookie_same_site = "strict"

[security]
max_login_attempts = 5
lockout_minutes    = 15
require_totp       = false

[rbac]
admin_group    = "Domain Admins"
operator_group = "SambaWebOperators"
helpdesk_group = "SambaWebHelpdesk"
readonly_group = "SambaWebReadOnly"
```

---

## Design System — Neo-Brutalism (MANDATORY)

The visual style is **Neo-Brutalist / Square-Modern (Brutalist-Lite)**. Strictly apply:

### Absolute rules
*   **ZERO border-radius** on any interactive component (buttons, inputs, cards, modals, badges).
*   **Solid borders** 2–4px `solid #1A1A1A` as the only visual delimiter.
*   **Brutalist offset shadows**: `box-shadow: 4px 4px 0px #1A1A1A` (not diffuse).
*   **Hover via inversion**: background and text swap colors — no fade, no transition.
*   **Monospace typography** (`JetBrains Mono` or `Fira Code`) on all technical data: DNs, SIDs, IPs, timestamps, sAMAccountNames.
*   **No smooth animations** — `transition: none` on all interactive elements.

### Color palette (CSS custom properties)
```css
:root {
  --color-base:    #F5F5F0;  /* application background */
  --color-surface: #FFFFFF;  /* cards, tables, panels */
  --color-ink:     #1A1A1A;  /* main text and borders */
  --color-accent:  #E63946;  /* highlight, destructive actions */
  --color-primary: #2B2D42;  /* sidebar, headers, CTAs */
  --color-success: #2D6A4F;  /* confirmations, enabled */
  --color-warning: #E07B00;  /* warnings */
  --color-muted:   #888888;  /* secondary texts */
}
```

### UI Components

**Buttons:**
```css
.btn {
  border: 2px solid var(--color-ink);
  border-radius: 0;
  box-shadow: 4px 4px 0px var(--color-ink);
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 0.5rem 1.25rem;
  transition: none;
  cursor: pointer;
}
.btn:hover {
  background: var(--color-ink);
  color: #FFFFFF;
  box-shadow: none;
  transform: translate(4px, 4px);
}
.btn-primary { background: var(--color-primary); color: #FFF; border-color: var(--color-primary); box-shadow: 4px 4px 0px var(--color-ink); }
.btn-danger  { border-color: var(--color-accent); box-shadow: 4px 4px 0px var(--color-accent); }
.btn-danger:hover { background: var(--color-accent); color: #FFF; }
```

**Inputs:**
```css
input, select, textarea {
  border: 2px solid var(--color-ink);
  border-radius: 0;
  padding: 0.5rem 0.75rem;
  font-family: inherit;
  transition: none;
}
input:focus { outline: none; box-shadow: 4px 4px 0px var(--color-primary); }
input.error { border-color: var(--color-accent); box-shadow: 4px 4px 0px var(--color-accent); }
```

**Tables:**
*   Header: `background: var(--color-primary)`, white text.
*   Alternating rows: `#FFFFFF` / `#F5F5F0`.
*   Row hover: `border-left: 4px solid var(--color-accent)` + `background: #FFF5F5`.
*   Solid borders on all sides.

**Modals:**
*   Overlay: `background: rgba(0,0,0,0.85)` — no blur.
*   Container: `border: 3px solid var(--color-ink)` + `box-shadow: 10px 10px 0px var(--color-ink)`.
*   No `transition` or `animation` — appears instantly.
*   Close: button `[✕]` as pure text, no SVG.

**Sidebar:**
*   Fixed width: `240px`.
*   Background: `var(--color-primary)` (`#2B2D42`).
*   Active item: `background: var(--color-accent)` + left border `4px solid #FFF`.
*   Hover: text/background inversion.

---

## Features to Implement (by phase)

### Phase 1 — Foundation (implement first)
1. `go.mod` with all dependencies declared.
2. Cobra CLI with `serve`, `migrate`, `user` commands.
3. Viper loading `config.toml` with all sections above.
4. GORM with auto-migration for `Session`, `AuditLog`, `Setting` models.
5. SQLite (default) and MySQL/MariaDB support via `config.toml`.
6. LDAP client (`internal/ldap/client.go`) with connection pool, auto-reconnect, and TLS.
7. LDAP bind authentication (`internal/auth/ldap.go`).
8. Session system with HTTPOnly cookies + per-form CSRF token.
9. Authentication middleware to protect routes.
10. Rate limiting: 5 login attempts per IP / 5 minutes.
11. Login page (`web/templates/auth/login.html`) with Neo-Brutalist design.
12. Base layout with sidebar + header (`web/templates/layout/`).
13. Dashboard with counters: total users, groups, computers, OUs; DC status; latest audit actions.
14. `embed.go` at the root loading the entire `web/` into the binary.

### Phase 2 — Core AD
15. User list with pagination, filters (enabled/disabled/expired), and search.
16. Create user: form with `displayName`, `givenName`, `sn`, `sAMAccountName`, `userPrincipalName`, `mail`, `telephoneNumber`, `title`, `department`, OU selection.
17. Edit user: all attributes + account options (must change password, never expires, account expires).
18. Disable/enable user.
19. Password reset with validation against domain policy.
20. Move user to another OU.
21. Group list with type (security/distribution) and member count.
22. Create/edit groups, add/remove members.
23. View OUs in a hierarchical tree with object count per OU.
24. Advanced LDAP search with custom filters by attribute.
25. Audit log: every action writes `{timestamp, admin_user, ip, action, object_dn, details}` to the database.

### Phase 3 — Production
26. RBAC via AD groups: Admin / Operator / Helpdesk / Read-Only.
27. Fine-Grained Password Policies: list and display details.
28. TOTP 2FA: QR code for setup, validation on login.
29. Exportable reports: inactive users (90+ days), passwords expiring in 7/14/30 days.
30. Security headers: CSP, X-Frame-Options, HSTS, X-Content-Type-Options.
31. Multi-stage Dockerfile (builder + minimal runtime).
32. `docker-compose.yml` with the `go-samba4` service.

---

## Implementation Rules

### Go
*   `main` package at the root — `embed.go` and `main.go` stay together at the root.
*   All business logic in `internal/` — handlers only orchestrate.
*   No hardcoded credentials — only via `config.toml` or environment variable.
*   `go vet` and `golangci-lint` must pass without warnings.
*   LDAP errors always logged with sufficient context for debugging.
*   LDAP connections with configurable timeout and automatic retry.
*   Passwords **never** logged — not even in debug mode.

### Templates
*   All in `web/templates/` — embedded via `//go:embed all:web/templates`.
*   Use `{{template "base" .}}` and `{{define "content"}}` for composition.
*   Technical data (DNs, SIDs, IPs) always in `<code>` with `font-family: monospace`.
*   Forms always with CSRF token: `<input type="hidden" name="_csrf" value="{{.CSRFToken}}">`.
*   Labels above fields — never placeholder as the only label.
*   Inline error messages below the field with `.field-error` class.

### Security
*   Sanitize **all** input before assembling LDAP filters (prevent LDAP injection).
*   Cookies: `HttpOnly=true`, `Secure=true` (in production), `SameSite=Strict`.
*   Mandatory CSRF on all POST/PUT/DELETE forms.
*   Destructive actions (delete, disable) require confirmation modal with re-typing the object name.
*   Audit log for **every** write operation in the AD.

### Database
*   GORM with auto-migrate on startup (`go-samba4 migrate` or automatic in `serve`).
*   `AuditLog` is append-only — never UPDATE or DELETE in this table.
*   Indexes on: `audit_logs.created_at`, `audit_logs.admin_user`, `sessions.token`, `sessions.expires_at`.

---

## Initial `go.mod`

```
module github.com/youruser/go-samba4

go 1.26

require (
    github.com/labstack/echo/v4          v4.13.0
    github.com/spf13/cobra               v1.9.1
    github.com/spf13/viper               v1.20.1
    gorm.io/gorm                         v1.25.12
    gorm.io/driver/sqlite                v1.5.7
    gorm.io/driver/mysql                 v1.5.7
    github.com/go-ldap/ldap/v3           v3.4.10
    github.com/jcmturner/gokrb5/v8       v8.4.4
    github.com/pquerna/otp               v1.4.0
    golang.org/x/crypto                  v0.36.0
)
```

---

## Existing Code References

Consult these projects for implementation references:
*   **go-samba4**: https://github.com/jniltinho/go-samba4 — Samba-specific AD operations.
*   **BLAZAM**: https://github.com/Blazam-App/BLAZAM — comprehensive UI/UX reference for AD panels.
*   **SWAT2**: https://github.com/rnapoles/swat2 — reference for Samba management via web.
*   **Samba Wiki**: https://wiki.samba.org/index.php/Main_Page — official documentation of Samba 4 LDAP attributes.

---

## File Creation Order

Create in this sequence to ensure incremental compilation:

1. `go.mod` + `go.sum`
2. `internal/config/config.go`
3. `internal/models/session.go`, `audit.go`, `setting.go`
4. `internal/ldap/schema.go`, `client.go`, `users.go`, `groups.go`, `ous.go`
5. `internal/auth/ldap.go`, `session.go`, `kerberos.go`, `totp.go`
6. `internal/middleware/auth.go`, `csrf.go`, `rbac.go`, `ratelimit.go`
7. `internal/handlers/dashboard.go`, `users.go`, `groups.go`, `ous.go`, `search.go`, `settings.go`
8. `web/static/css/app.css` (Tailwind + Neo-Brutalist variables)
9. `web/static/js/app.js` (jQuery 4 + modals + AJAX helpers)
10. `web/templates/layout/base.html`, `sidebar.html`
11. `web/templates/auth/login.html`
12. `web/templates/dashboard.html`
13. `web/templates/users/list.html`, `form.html`, `detail.html`
14. `web/templates/groups/list.html`, `form.html`
15. `web/templates/ous/tree.html`
16. `web/templates/audit/list.html`
17. `embed.go` (root, main package)
18. `cmd/root.go`, `serve.go`, `migrate.go`, `user.go`
19. `main.go`
20. `config.toml`
21. `Dockerfile` (multi-stage)
22. `docker-compose.yml`
23. `README.md`

---

## Expected Outcome

At the end, the command below must work:

```bash
# Build
go build -o go-samba4 .

# Initialize database and start server
./go-samba4 serve --config config.toml

# Access
# http://localhost:8080 → redirects to /login
# Login with AD user → dashboard
# Single binary, no external template or static files
```
