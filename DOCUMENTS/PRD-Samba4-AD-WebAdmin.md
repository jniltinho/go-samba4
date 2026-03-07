# PRD — Samba4 AD Web Administration Panel

> **Version:** 1.1.0 — Revised Draft
> **Date:** March 2026
> **Platform:** IDE Antigravity
> **Status:** 🔴 Under Development

---

## Summary

1.  [Product Overview](#1-product-overview)
2.  [Technology Stack](#2-technology-stack)
3.  [Design System — Neo-Brutalism](#3-design-system--neo-brutalism--square-modern)
4.  [System Architecture](#4-system-architecture)
5.  [Product Features](#5-product-features)
6.  [Security Requirements](#6-security-requirements)
7.  [Interface and User Experience](#7-interface-and-user-experience)
8.  [Integrations and References](#8-integrations-and-references)
9.  [Development Roadmap](#9-development-roadmap)
10. [Acceptance Criteria (DoD)](#10-acceptance-criteria-dod)
11. [Technical Glossary](#11-technical-glossary)

---

## 1. Product Overview

The **Samba4 AD Web Admin Panel** is a modern web administration tool to manage **Samba 4 Active Directory** environments. The project replaces legacy interfaces (such as SWAT) with a contemporary, secure, and extensible solution, built with high-performance modern technologies.

Inspired by projects like [BLAZAM](https://github.com/Blazam-App/BLAZAM) and [go-samba4](https://github.com/jniltinho/go-samba4), the panel offers a **Neo-Brutalist** user experience — direct, functional, and visually distinct — without the excesses of generic SaaS interfaces.

### 1.1 Product Objectives

*   Provide a centralized web interface for comprehensive Samba 4 AD administration.
*   Replace legacy command-line tools and GUIs (`SWAT`, `samba-tool`) with a modern web panel.
*   Implement secure LDAP/Kerberos authentication integrated directly with the AD.
*   Offer granular management of users, groups, OUs, GPOs, and password policies.
*   Be lightweight, fast, and self-hosted, without dependencies on external services.
*   Serve as an extensible foundation for future integrations (DNS, DHCP, Samba shares).

### 1.2 Target Audience

*   Linux systems administrators managing Samba 4 domains.
*   IT teams in small to medium-sized businesses lacking Microsoft infrastructure.
*   Professionals looking for open-source alternatives to RSAT (Remote Server Administration Tools).

---

## 2. Technology Stack

### 2.1 Core Technologies

| Technology | Version | Purpose |
| :--- | :--- | :--- |
| **Golang** | 1.26 | Main language — backend and asset compilation |
| **Echo Framework** | v4+ | HTTP server, middleware, RESTful routing |
| **GORM** | v2 | ORM for local database (SQLite, MySQL, MariaDB) |
| **Cobra + Viper** | latest | Management CLI and TOML configuration |
| **Tailwind CSS** | 4.2 | Utility CSS framework — Neo-Brutalist style |
| **jQuery** | 4.0.0 | DOM interactivity, AJAX calls, UI components |
| **Go `embed` + `html/template`** | stdlib | HTML templates and static assets embedded in the binary |
| **go-ldap/ldap** | v3 | LDAP Protocol — Samba 4 AD integration |
| **gokrb5** | v8+ | Kerberos 5 Authentication — AD SSO |
| **SQLite / MySQL / MariaDB** | — | Local storage for sessions, logs, and settings |

### 2.2 Database — Driver Strategy

GORM abstracts the two supported drivers. The choice is made via `config.toml` without altering application code.

| Criterion | SQLite | MySQL / MariaDB |
| :--- | :--- | :--- |
| **Setup** | Zero configuration — single file | Requires a separate database server |
| **Ideal Use Case** | Single instance, small domains | High concurrency, multiple simultaneous admins |
| **Deploy** | Copy a `.db` file for backup | Backup via `mysqldump` |
| **Practical Limit** | ~100 write req/s | Thousands of req/s |
| **Go Driver (GORM)**| `gorm.io/driver/sqlite` | `gorm.io/driver/mysql` |
| **Supported Versions**| SQLite 3.35+ | MySQL 8.0+ / MariaDB 10.6+ |

> **Recommendation:** Use **SQLite** as the default for simple installations and **MySQL/MariaDB** when there are multiple administrators accessing simultaneously or a need for replication.

### 2.3 LDAP / Active Directory Libraries for Go

After analyzing the Go ecosystem, the following libraries are recommended:

#### 2.3.1 `go-ldap/ldap` — Primary Recommendation

```
github.com/go-ldap/ldap/v3
```

*   The most mature and widely used LDAP library in the Go ecosystem.
*   Full support for LDAP v3, TLS/LDAPS, SASL, result pagination.
*   Complete CRUD operations: `Add`, `Modify`, `Delete`, `Search`, `ModifyDN`.
*   Actively maintained — used by HashiCorp Vault, Grafana, and other large-scale projects.

#### 2.3.2 `go-samba4` — Implementation Reference

```
github.com/jniltinho/go-samba4
```

*   Specific wrapper for Samba 4 AD with high-level helpers.
*   Useful as a reference for AD-specific operations and Samba attributes.
*   Can be used in conjunction with `go-ldap/ldap`.

#### 2.3.3 `gokrb5` — Kerberos Authentication

```
github.com/jcmturner/gokrb5/v8
```

*   Pure Go Kerberos 5 implementation, no system dependencies.
*   Support for SPNEGO/Negotiate authentication for SSO with Windows clients.
*   Integration with `net/http` for transparent authentication middleware.

---

## 3. Design System — Neo-Brutalism / Square-Modern

### 3.1 Visual Philosophy

The **Neo-Brutalist (Brutalist-Lite)** theme prioritizes functional clarity over ornamentation. Every visual element serves an explicit purpose. The interface is square, direct, and unambiguous — suitable for a critical infrastructure administration panel.

### 3.2 Design Principles

1.  **No rounded borders** — `border-radius: 0` on all interactive components.
2.  **Solid, thick borders** — 2–4px as explicit element delimiters.
3.  **Monospace typography** for technical data (IPs, DNs, SIDs, timestamps).
4.  **High contrast** — light background with dark text, or black background with white text.
5.  **Offset shadows** instead of diffuse drop-shadows — `box-shadow: 4px 4px 0px #000`.
6.  **Restricted palette** — maximum 3 highlight colors + black + white + grays.
7.  **Hover via color inversion** — background/text swap on interactive elements.
8.  **Simple linear icons** — no gradient fills.

### 3.3 Color Palette

| CSS Token | HEX Value | Usage |
| :--- | :--- | :--- |
| `--color-base` | `#F5F5F0` | Main application background |
| `--color-surface` | `#FFFFFF` | Cards, tables, panels |
| `--color-ink` | `#1A1A1A` | Main text and borders |
| `--color-accent` | `#E63946` | Highlight, destructive actions, alerts |
| `--color-primary` | `#2B2D42` | Sidebar, headers, primary CTAs |
| `--color-success` | `#2D6A4F` | Confirmations, enabled objects |
| `--color-warning` | `#E07B00` | Warnings, items requiring attention |
| `--color-muted` | `#888888` | Secondary texts, metadata |

### 3.4 Spacing and Typography Tokens

```css
/* Typography */
--font-sans: 'Inter', system-ui, sans-serif;
--font-mono: 'JetBrains Mono', 'Fira Code', monospace;

/* Borders */
--border-thin:   2px solid var(--color-ink);
--border-medium: 3px solid var(--color-ink);
--border-thick:  4px solid var(--color-ink);

/* Brutalist shadows */
--shadow-sm: 3px 3px 0px var(--color-ink);
--shadow-md: 5px 5px 0px var(--color-ink);
--shadow-lg: 8px 8px 0px var(--color-ink);
```

---

## 4. System Architecture

### 4.1 Overview

The application follows a **modular monolithic architecture** — a single self-sufficient Go binary serving both a REST API and server-side rendered (SSR) HTML pages. All HTML templates and static assets (CSS, JS) are **embedded directly into the binary** via `//go:embed`, eliminating external file dependencies in production. The UI is built with Go templates + jQuery for progressive interactivity.

```
┌─────────────────────────────────────────────┐
│               Admin Browser                  │
│         (HTML + Tailwind + jQuery)           │
└────────────────────┬────────────────────────┘
                     │ HTTP/HTTPS
┌────────────────────▼────────────────────────┐
│           Echo HTTP Server (Go)              │
│  ┌──────────┐ ┌──────────┐ ┌─────────────┐ │
│  │Middleware│ │ Handlers │ │  Templates  │ │
│  │Auth/CSRF │ │ REST/SSR │ │  html/tmpl  │ │
│  └──────────┘ └────┬─────┘ └─────────────┘ │
│                    │                         │
│  ┌─────────────────▼──────────────────────┐ │
│  │         Service Layer (Go)              │ │
│  │  ┌──────────┐  ┌──────────────────┐   │ │
│  │  │  LDAP    │  │   GORM / SQLite  │   │ │
│  │  │ Service  │  │  (sessions/logs) │   │ │
│  │  └────┬─────┘  └──────────────────┘   │ │
│  └───────┼─────────────────────────────── ┘ │
└──────────┼──────────────────────────────────┘
           │ LDAP/LDAPS (389/636)
┌──────────▼──────────────────────────────────┐
│           Samba 4 Active Directory           │
│        (Linux Domain Controller)             │
└─────────────────────────────────────────────┘
```

### 4.2 Directory Structure

```
samba4-admin/
├── cmd/
│   ├── root.go              # Cobra root command + Viper setup
│   ├── serve.go             # Command: samba4-admin serve
│   ├── migrate.go           # Command: samba4-admin migrate
│   └── user.go              # Command: samba4-admin user create/list
├── internal/
│   ├── auth/
│   │   ├── ldap.go          # LDAP bind authentication
│   │   ├── kerberos.go      # Kerberos/SPNEGO authentication
│   │   ├── totp.go          # 2FA with TOTP (RFC 6238)
│   │   └── session.go       # Session management
│   ├── config/
│   │   └── config.go        # Viper loader + config structs
│   ├── handlers/
│   │   ├── dashboard.go
│   │   ├── users.go
│   │   ├── groups.go
│   │   ├── ous.go
│   │   ├── search.go
│   │   └── settings.go
│   ├── ldap/
│   │   ├── client.go        # go-ldap wrapper with connection pool
│   │   ├── users.go         # LDAP user operations
│   │   ├── groups.go        # LDAP group operations
│   │   ├── ous.go           # LDAP OU operations
│   │   └── schema.go        # AD attribute mapping
│   ├── models/
│   │   ├── session.go       # GORM model: web sessions
│   │   ├── audit.go         # GORM model: audit log
│   │   └── setting.go       # GORM model: local settings
│   └── middleware/
│       ├── auth.go          # Echo auth middleware
│       ├── csrf.go          # CSRF protection
│       ├── rbac.go          # Role-based access control
│       └── ratelimit.go     # IP rate limiting
├── web/                     # ← compiled and embedded in binary via go:embed
│   ├── templates/
│   │   ├── layout/
│   │   │   ├── base.html    # Base template (head, scripts, meta)
│   │   │   └── sidebar.html # Side navigation
│   │   ├── dashboard.html
│   │   ├── users/
│   │   │   ├── list.html
│   │   │   ├── form.html
│   │   │   └── detail.html
│   │   ├── groups/
│   │   ├── ous/
│   │   ├── audit/
│   │   └── auth/
│   │       └── login.html
│   └── static/
│       ├── css/
│       │   └── app.css      # Compiled Tailwind CSS 4.2
│       └── js/
│           └── app.js       # jQuery 4 + custom helpers
├── embed.go                 # //go:embed — embedded assets (root, package main)
├── config.toml              # Main configuration
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── main.go
```

> **Note on `embed.go`:** Located in the project root within the `main` package, which allows referencing `web/` directly with `//go:embed web/...` without complex relative paths. The final production binary relies on no external files — everything (HTML, CSS, JS) is compiled into the executable.

### 4.3 Template and Asset Embed (`embed.go`)

All files in the `web/` folder are embedded into the binary at compile time using the Go standard library's `//go:embed` directive (available since Go 1.16). The file resides in the project root within the `main` package.

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

// Templates returns the compiled template set from the embedded FS.
func Templates() *template.Template {
    tmpl := template.Must(
        template.New("").ParseFS(TemplatesFS, "web/templates/**/*.html"),
    )
    return tmpl
}

// StaticHandler returns an http.FileServer serving embedded static assets.
func StaticHandler() http.Handler {
    sub, _ := fs.Sub(StaticFS, "web/static")
    return http.FileServer(http.FS(sub))
}
```

> **Production advantage:** Deployment is a single binary copied to the server. No need to sync `templates/` or `static/` folders — everything travels together in the executable. In development, `--dev` mode can override the embedded FS to read directly from disk for hot-reload.

### 4.4 Configuration (`config.toml`)

```toml
[server]
host     = "0.0.0.0"
port     = 8080
tls_cert = "/etc/ssl/certs/admin.crt"
tls_key  = "/etc/ssl/private/admin.key"
debug    = false
# In dev mode, templates are loaded from disk (hot-reload) instead of embed
dev_mode = false

[ldap]
host            = "dc1.empresa.local"
port            = 636
use_tls         = true
skip_tls_verify = false
base_dn         = "DC=empresa,DC=local"
bind_user       = "CN=samba4admin,CN=Users,DC=empresa,DC=local"
# bind_pass is read from environment variable: SAMBA4_LDAP_PASS

[database]
# Option 1 — SQLite (default, zero config, ideal for single instances)
driver = "sqlite"
path   = "/var/lib/samba4-admin/data.db"

# Option 2 — MySQL or MariaDB (recommended for high-concurrency environments)
# driver = "mysql"
# dsn    = "samba4admin:password@tcp(localhost:3306)/samba4admin?charset=utf8mb4&parseTime=True&loc=Local"

[session]
secret           = ""        # Auto-generated on first start
timeout_minutes  = 30
cookie_secure    = true
cookie_same_site = "strict"

[security]
max_login_attempts = 5
lockout_minutes    = 15
require_totp       = false
```

---

## 5. Product Features

### 5.1 Modules and Priorities

| Module / Feature | Priority | Description |
| :--- | :--- | :--- |
| Main Dashboard | 🔴 High | Overview: active users, DCs, alerts, recent activity |
| Web Authentication | 🔴 High | LDAP login with TOTP 2FA support |
| User Management | 🔴 High | Full CRUD: create, edit, disable, reset password, move OU |
| Group Management | 🔴 High | Create/edit security and distribution groups, manage members |
| OU Management | 🔴 High | Organizational Units tree with hierarchical navigation |
| Advanced LDAP Search | 🔴 High | Search AD attributes with custom LDAP filters |
| Password Policies | 🔴 High | View and edit Fine-Grained Password Policies |
| Audit / Logs | 🔴 High | Log all actions with user, IP, timestamp, and affected object |
| GPO Management | 🟡 Medium | List and edit GPOs linked to OUs |
| Integrated DNS | 🟡 Medium | Manage AD DNS records (Zones, A, CNAME, PTR, MX) |
| Computer Management | 🟡 Medium | List/move/disable computer objects in AD |
| Reports | 🟡 Medium | Export: inactive users, expired passwords, group members |
| Multi-Domain | 🟢 Low | Support for multiple AD forests/domains in one interface |
| REST API | 🟢 Low | Documented public API (Swagger/OpenAPI) for automation |
| Integrated DHCP | 🟢 Low | View DHCP leases integrated with AD DNS |

### 5.2 Breakdown — Users Module

#### Listing

*   Paginated table with columns: Name, `sAMAccountName`, Email, Department, Status, Last Login.
*   Quick filters: enabled/disabled, expired account, no login for 90+ days.
*   Search by name, email, username, and custom attributes.
*   CSV and JSON export.

#### Creation and Editing

*   Comprehensive form with all relevant AD attributes.
*   Required fields: `displayName`, `givenName`, `sn`, `sAMAccountName`, `userPrincipalName`.
*   Optional fields: `mail`, `telephoneNumber`, `title`, `department`, `manager`, `description`.
*   Target OU selection via hierarchical tree dropdown.
*   Password setting with real-time validation against domain policy.
*   Account options: must change at next logon, password never expires, account expires.

#### Bulk Actions

*   Multiple selection → move OU, enable/disable, add to group.
*   Bulk CSV import with interactive column mapping.
*   Export results of any search/filter.

### 5.3 Breakdown — Dashboard

```
┌─────────────────────────────────────────────────────┐
│  DOMAIN: empresa.local          [DC: dc1 ● Online]  │
├──────────┬──────────┬──────────┬────────────────────┤
│  👤 847  │  👥 124  │  💻 312 │  ⚠️  3 Alerts      │
│ Users    │  Groups  │  Comps  │                     │
├──────────┴──────────┴──────────┴────────────────────┤
│  Recent Activity                Passwords Expiring   │
│  • jsmith created by admin      • bsantos — 2 days  │
│  • TI-Group modified            • rpereira — 5 days │
│  • CN=PC042 disabled            • amaral — 7 days   │
└─────────────────────────────────────────────────────┘
```

---

## 6. Security Requirements

### 6.1 Authentication and Authorization

*   Primary authentication via **LDAP bind** against Samba 4 AD.
*   **Kerberos SPNEGO** support for SSO with Windows domain clients.
*   Fallback local authentication (emergency account with `bcrypt` + salt).
*   **TOTP** support (RFC 6238) as a mandatory or optional second factor.
*   Sessions via `HTTPOnly` + `Secure` cookies with per-form CSRF tokens.
*   Configurable session timeout (default: 30 minutes of inactivity).

### 6.2 Access Control (RBAC)

| Role | Permissions |
| :--- | :--- |
| **Admin** | Full access — all CRUD operations, settings, RBAC |
| **Operator** | CRUD for users/groups/OUs, no access to system settings |
| **Helpdesk** | Password resets, account enable/disable, read audit logs |
| **Read Only**| View all objects, no modifications |

*   Roles mapped directly to AD groups (configurable in `config.toml`).
*   Granular permissions per module — e.g., Helpdesk can reset passwords but not delete users.
*   Destructive actions require explicit confirmation by re-typing the object name.

### 6.3 Application Security

*   Mandatory TLS in production (no pure HTTP accepted by default).
*   IP rate limiting for login attempts (max 5 attempts / 5 mins).
*   CSRF protection on all POST/PUT/DELETE forms via double-submit tokens.
*   Sanitization of LDAP inputs to prevent LDAP injection.
*   Mandatory security headers: `CSP`, `X-Frame-Options`, `HSTS`, `X-Content-Type-Options`.
*   Immutable audit logs — append-only in the local database.

---

## 7. Interface and User Experience

### 7.1 Main Layout

```
┌─────────────────────────────────────────────────────────┐
│ HEADER: [breadcrumb]           [domain]  [user ▼]      │
├─────────┬───────────────────────────────────────────────┤
│         │                                               │
│ SIDEBAR │              CONTENT AREA                     │
│ 240px   │              (fluid)                          │
│ dark    │                                               │
│         │                                               │
│ ▪ Dashboard                                            │
│ ▪ Users                                                │
│ ▪ Groups                                               │
│ ▪ OUs                                                  │
│ ▪ GPOs                                                 │
│ ▪ DNS                                                  │
│ ▪ Audit                                                │
│ ▪ Settings                                             │
│         │                                               │
└─────────┴───────────────────────────────────────────────┘
```

### 7.2 Component Specifications

#### Buttons

```css
.btn {
  border: 2px solid #1A1A1A;
  border-radius: 0;
  box-shadow: 4px 4px 0px #1A1A1A;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  transition: none; /* brutalism: no smooth transitions */
}
.btn:hover {
  background: #1A1A1A;
  color: #FFFFFF;
  box-shadow: none;
  transform: translate(4px, 4px);
}
.btn-danger { border-color: #E63946; box-shadow: 4px 4px 0px #E63946; }
```

#### Tables

*   Solid borders on all sides — no minimal borderless styles.
*   Header with `#2B2D42` background and white text.
*   Alternating rows: `#FFFFFF` and `#F5F5F0`.
*   Row hover: left border `4px solid #E63946` + `#FFF5F5` background.
*   Action columns right-aligned with compact buttons.

#### Inputs and Forms

*   Labels always visible above the field — never placeholder-only.
*   `border: 2px solid #1A1A1A` — no border-radius.
*   Focus: `box-shadow: 4px 4px 0px #2B2D42`.
*   Error states: border `#E63946` + red message below the field.
*   Grouped fieldsets with `border: 2px solid #1A1A1A` and prominent `<legend>`.

#### Modals

*   Overlay `rgba(0,0,0,0.85)` — no blur.
*   Container with `border: 3px solid #1A1A1A` + `box-shadow: 10px 10px 0px #1A1A1A`.
*   No fade animations — instant display (pure brutalism).
*   Close button `[✕]` in the top right corner as text, not SVG icon.

### 7.3 Responsiveness

The panel is designed primarily for **desktop (min. 1280px)**. On tablets (768px+), the sidebar collapses into a hamburger menu. Mobile is not a primary requirement for v1, but the Tailwind structure should allow future evolution without a rewrite.

---

## 8. Integrations and References

### 8.1 Reference Projects

| Project | Link | Relevance |
| :--- | :--- | :--- |
| **BLAZAM** | [github.com/Blazam-App/BLAZAM](https://github.com/Blazam-App/BLAZAM) | Comprehensive UI/UX reference for AD panels. Inspiration for helpdesk features and granular RBAC. |
| **go-samba4** | [github.com/jniltinho/go-samba4](https://github.com/jniltinho/go-samba4) | Go wrapper for Samba 4 AD — foundation for specific LDAP operations and Samba attributes. |
| **SWAT2** | [github.com/rnapoles/swat2](https://github.com/rnapoles/swat2) | Reference for Samba share management and web-based smb.conf settings. |
| **Samba Wiki** | [wiki.samba.org](https://wiki.samba.org/index.php/Main_Page) | Official documentation — AD schema, LDAP attributes, and samba-tool equivalents. |
| **Samba GitLab** | [gitlab.com/samba-team/samba](https://gitlab.com/samba-team/samba) | Official Samba source code — reference for internal AD behaviors. |
| **ui-ux-pro-max-skill** | [github.com/nextlevelbuilder/ui-ux-pro-max-skill](https://github.com/nextlevelbuilder/ui-ux-pro-max-skill) | UX quality skill for complex functional interfaces. |

### 8.2 Reference Videos

*   [Samba 4 AD Administration Overview](https://www.youtube.com/watch?v=9_lVXZmbBoQ)
*   [Playlist: Samba 4 Active Directory](https://www.youtube.com/playlist?list=PLozhsZB1lLUP8vTzrTTfQ1YWsR0hUAzFu)

---

## 9. Development Roadmap

### Phase 1 — Foundation *(Weeks 1–3)*

*   [ ] Go project setup with Cobra/Viper, Echo, and GORM
*   [ ] Build pipeline configuration (Tailwind CSS 4.2 + assets)
*   [ ] LDAP connection with `go-ldap/ldap` — authentication and basic search
*   [ ] Web authentication system (login, session, CSRF, logout)
*   [ ] Neo-Brutalist base layout with Tailwind — sidebar + header
*   [ ] Dashboard with basic domain metrics (counters, DC status)

### Phase 2 — Core AD *(Weeks 4–7)*

*   [ ] Full User CRUD (create, edit, disable, reset password)
*   [ ] Full Group CRUD (security + distribution, members)
*   [ ] View and navigate OUs in a hierarchical tree
*   [ ] Move objects between OUs via interface
*   [ ] Advanced LDAP search system with dynamic filters
*   [ ] Action auditing — append-only log in the database

### Phase 3 — Production *(Weeks 8–10)*

*   [ ] Full RBAC mapped to AD groups
*   [ ] Password reset with validation against domain policy
*   [ ] View and edit Fine-Grained Password Policies
*   [ ] 2FA via TOTP (Google Authenticator / Authy)
*   [ ] Basic exportable reports (CSV/JSON)
*   [ ] Security hardening — headers, rate limiting, penetration testing
*   [ ] Dockerfile + docker-compose.yml + README

### Phase 4 — Extensions *(Post v1.0)*

*   [ ] GPO Management (list, edit links)
*   [ ] AD DNS Integration (Zones, A/CNAME/PTR/MX records)
*   [ ] Computer Object Management
*   [ ] Multi-domain / Multi-forest
*   [ ] Documented REST API (Swagger/OpenAPI 3.0)
*   [ ] SAML 2.0 Authentication (external SSO)

---

## 10. Acceptance Criteria (DoD)

### Per Feature

*   Code compiles without warnings using `go vet` and `golangci-lint`.
*   Unit tests with a minimum **70%** coverage in logic layers (services).
*   Integration tests against Samba 4 AD in a Docker environment.
*   Interface functions on Firefox ESR, Chrome/Chromium, and Edge (current versions).
*   All actions correctly audited in the log with user, IP, and timestamp.
*   Inline documentation (`godoc`) for all public functions and types.
*   No hardcoded credentials — everything via `config.toml` or environment variables.

### For v1.0 Release

*   Self-sufficient single binary — templates and assets embedded via `//go:embed`, no external files needed in production.
*   `dev_mode = true` in `config.toml` allows hot-reloading templates from disk during development.
*   Dockerfile and `docker-compose.yml` functional and documented.
*   `README.md` with complete installation, configuration, and first-access instructions.
*   Samba 4 AD configuration guide (required LDAP permissions, bind user setup).
*   `CHANGELOG.md` following semantic versioning (SemVer).
*   Response time < 300ms for 95% of read operations (p95).
*   No critical or high vulnerabilities in `govulncheck`.

---

## 11. Technical Glossary

| Term | Definition |
| :--- | :--- |
| **Active Directory** | Microsoft's directory service based on LDAP/Kerberos, implemented in Samba 4 for Linux. |
| **DN** *(Distinguished Name)* | Unique identifier of an object in LDAP — e.g., `CN=John,OU=IT,DC=company,DC=local`. |
| **sAMAccountName** | AD attribute representing the user's NetBIOS login — e.g., `jsmith`. |
| **userPrincipalName** | Login in UPN format — e.g., `jsmith@company.local`. |
| **OU** *(Organizational Unit)* | LDAP container to organize AD objects in an administrative hierarchy. |
| **GPO** *(Group Policy Object)* | Configuration policy applied to users/computers via Group Policy. |
| **LDAP** | *Lightweight Directory Access Protocol* — standard protocol for directory access. |
| **Kerberos** | Ticket-based network authentication protocol, used by AD for SSO. |
| **SID** *(Security Identifier)* | Unique security identifier assigned to each object in AD. |
| **Fine-Grained Password Policy** | Password policies applied to specific groups, superseding the default domain policy. |
| **RBAC** | *Role-Based Access Control* — access control based on user roles/profiles. |
| **SPNEGO** | Authentication negotiation mechanism used for Kerberos SSO over HTTP. |
| **TOTP** | *Time-based One-Time Password* (RFC 6238) — standard for two-factor authentication. |
| **SSR** | *Server-Side Rendering* — rendering HTML pages on the Go server. |
| **SWAT** | *Samba Web Administration Tool* — legacy Samba web interface, discontinued. |

---

*PRD v1.1 — Samba4 AD Web Admin Panel — IDE Antigravity — March 2026*
