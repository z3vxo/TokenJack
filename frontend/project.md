# TokenJack Admin — Frontend

## Scope

- React SPA, served only at `/admin` and `/admin/*`, embedded into the Go binary.
- No other routes are part of this app.
- Design goal: sleek SaaS, clean UI. No emojis, no gradients, no AI slop.

## Design direction

- **Style:** Minimalism & Swiss style — spacious, grid-based, geometric, high contrast, essential elements only.
- **Modes:** light (default) + **dark**, switched via `data-theme` on `<html>` (see Dark mode).
- **Separation:** 1px borders, not shadows (Swiss). No decorative elevation. The only shadow is the dropdown menu (`--menu-shadow`).
- **Voice:** precise, functional, operator-focused. Terse labels, sentence case.

---

## Skills to use (from `.opencode/skills/`)

| Skill | Use it for | Key entry point |
|-------|-----------|-----------------|
| **ui-ux-pro-max** | Primary driver. Generate/refresh the design system, and query focused guidance. | `.opencode/skills/ui-ux-pro-max/scripts/search.py` |
| **design-system** | Three-layer token architecture (primitive → semantic → component), component state specs. | `.opencode/skills/design-system/references/token-architecture.md` |
| **ui-styling** | Tailwind v4 setup, component/a11y patterns, responsive utilities. We use the **Tailwind-only Vite** path (not shadcn). | `.opencode/skills/ui-styling/references/` |
| **brand** | Voice, messaging, logo/asset rules — for copy & future assets. | `.opencode/skills/brand/references/voice-framework.md` |
| **design** | Router to sub-skills (logo/icon/banner/slides) if those are ever needed. | `.opencode/skills/design/references/design-routing.md` |
| **customize-opencode** | Only when editing opencode config itself — never for app code. | built-in |

### ui-ux-pro-max commands

```bash
# Regenerate the full design system (persisted):
python3 .opencode/skills/ui-ux-pro-max/scripts/search.py "saas admin panel minimal swiss clean white" \
  --design-system --persist -p "TokenJack" --output-dir .

# Focused lookups (one dominant intent, 2-5 terms):
python3 .opencode/skills/ui-ux-pro-max/scripts/search.py "minimalism clean saas light" --domain style
python3 .opencode/skills/ui-ux-pro-max/scripts/search.py "saas clean white minimal light" --domain color
python3 .opencode/skills/ui-ux-pro-max/scripts/search.py "saas dashboard minimal professional precise" --domain typography
python3 .opencode/skills/ui-ux-pro-max/scripts/search.py "accessible authentication login form password" --domain ux
python3 .opencode/skills/ui-ux-pro-max/scripts/search.py "navigation hierarchy top tabs sidebar" --domain ux
```

> Note: the generated `design-system/tokenjack/MASTER.md` came back as a dark
> glassmorphism ("Smart Home/IoT") system, which conflicts with the required
> look. **This file + `src/styles/tokens.css` are the source of truth.**
> Re-run the command above only to pull fresh *guidance*, then reconcile back to
> the Swiss tokens below.

---

## Design tokens

Source of truth: `frontend/src/styles/tokens.css` (three-layer), exposed to
Tailwind via `@theme inline` in `frontend/src/index.css`. Light lives in
`:root`; dark overrides only the semantic + component layers under
`[data-theme="dark"]` (primitives are shared).

### Primitive

```
--white #ffffff
--slate-50 #f8fafc   --slate-100 #f1f5f9  --slate-200 #e2e8f0  --slate-300 #cbd5e1
--slate-400 #94a3b8  --slate-500 #64748b  --slate-600 #475569  --slate-700 #334155
--slate-800 #1e293b  --slate-900 #0f172a  --slate-950 #020617
--blue-400 #60a5fa   --blue-500 #3b82f6   --blue-600 #2563eb   --blue-700 #1d4ed8
--green-500 #22c55e  --green-600 #16a34a
--red-50 #fef2f2     --red-200 #fecaca    --red-500 #ef4444    --red-600 #dc2626
--red-900 #7f1d1d    --red-950 #2a1416
```

### Semantic (light / dark)

| Role | Token | Light | Dark |
|------|-------|-------|------|
| Background | `--background` | `slate-50 #f8fafc` | `slate-950 #020617` |
| Foreground | `--foreground` | `slate-900 #0f172a` | `slate-100 #f1f5f9` |
| Card | `--card` | `white #ffffff` | `slate-900 #0f172a` |
| Card fg | `--card-foreground` | `slate-900` | `slate-100` |
| Muted | `--muted` | `slate-100 #f1f5f9` | `slate-800 #1e293b` |
| Muted fg | `--muted-foreground` | `slate-500 #64748b` | `slate-400 #94a3b8` |
| Border | `--border` | `slate-200 #e2e8f0` | `slate-800 #1e293b` |
| Input border | `--input` | `slate-300 #cbd5e1` | `slate-700 #334155` |
| Placeholder | `--placeholder` | `slate-400 #94a3b8` | `slate-500 #64748b` |
| Primary / CTA | `--primary` | `blue-600 #2563eb` | `blue-500 #3b82f6` |
| Primary hover | `--primary-hover` | `blue-700 #1d4ed8` | `blue-400 #60a5fa` |
| On primary | `--on-primary` | `white` | `white` |
| Ring | `--ring` | `blue-600 #2563eb` | `blue-500 #3b82f6` |
| Destructive | `--destructive` | `red-600 #dc2626` | `red-500 #ef4444` |
| Destructive surface | `--destructive-surface` | `red-50 #fef2f2` | `red-950 #2a1416` |
| Destructive border | `--destructive-border` | `red-200 #fecaca` | `red-900 #7f1d1d` |
| Success | `--success` | `green-600 #16a34a` | `green-500 #22c55e` |

### Component

| Token | Light | Dark | Usage |
|-------|-------|------|-------|
| `--control-height` | `2.75rem` (44px) | same | Inputs & buttons |
| `--radius-control` | `0.5rem` (8px) | same | Buttons, inputs, alerts, menu |
| `--radius-mark` | `0.625rem` (10px) | same | Brand mark |
| `--focus-ring-width` | `3px` | same | Visible focus ring |
| `--focus-ring-color` | `color-mix(blue-600 28%, transparent)` | `color-mix(blue-500 45%, transparent)` | Focus ring fill |
| `--menu-shadow` | `0 8px 24px rgba(15,23,42,0.08)` | `0 10px 30px rgba(0,0,0,0.45)` | Dropdown menu only |

Tailwind class mapping: `bg-background`, `text-foreground`, `border-border`,
`bg-card`, `text-muted-foreground`, `bg-primary`, `text-on-primary`,
`rounded-control`, `placeholder:text-placeholder`, `shadow-menu`, etc.

---

## Typography

- **Font:** Inter (single family, weights 300–700). Matches "Minimal Swiss" / admin-panel guidance.
- **Mono:** Fira Code, only for token/code-like data.
- Loaded in `frontend/index.html`:

```html
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=Fira+Code:wght@400;500;600&display=swap" rel="stylesheet" />
```

- Headings: semibold, tight tracking (`tracking-[-0.02em]`).
- Body/labels: `text-sm` (14px), `text-muted-foreground` for secondary.

---

## Dark mode

- Mechanism: `data-theme="light|dark"` on `<html>`; tokens override in CSS. No `dark:` utility classes needed.
- Persistence: `localStorage["tj_theme"]` (`src/lib/theme.js`). Falls back to `prefers-color-scheme`.
- No flash: an inline script in `index.html` sets the attribute before first paint.
- `color-scheme` is set per theme (`light` / `dark`) so native scrollbars, autofill, etc. follow.
- Toggle: `ThemeToggle` (Sun/Moon, Phosphor) in the app header and on the login page.

---

## Style rules (enforce on every component)

- **No emojis as icons.** Use vector icons only: **Phosphor** (`@phosphor-icons/react`) as the default set; Heroicons as fallback. Match stroke weight across a layer.
- **Icons:** decorative → `aria-hidden="true"`; standalone controls → accessible name (`aria-label`).
- **No layout-shifting hovers.** Change color/opacity/border, never size or position.
- **Transitions** 150–300ms, color-only.
- **Focus:** always visible — `focus-visible:ring-[3px] ring-[color:var(--focus-ring-color)]`.
- **`cursor-pointer`** on every clickable element.
- **Contrast ≥ 4.5:1** for body text, in **both** light and dark.
- **`prefers-reduced-motion`** honored globally (set in `index.css`).
- **Responsive:** verify at 375 / 768 / 1024 / 1440px. No horizontal scroll.
- Separation via borders, not shadows.

---

## Navigation

- **Tabs, centered top nav** (not a sidebar): Campaigns / Automation / Codes / Settings.
- Implemented in `src/components/AppNav.jsx` with `NavLink`; active item = `text-foreground` + 2px `border-primary` underline, inactive = `text-muted-foreground` → `text-foreground` on hover. `mx-auto w-max` keeps them centered while still scrolling on narrow screens.
- Revisit: switch to a static slim sidebar only past ~6–7 sections or with nested pages.

---

## Components

- **`BrandMark`** (`src/components/BrandMark.jsx`) — geometric keyhole SVG in a primary rounded square. No emoji, no raster assets.
- **`AppLayout`** (`src/components/AppLayout.jsx`) — authenticated shell: header (brand, campaign dropdown, theme toggle, sign out) + `AppNav` + `<Outlet/>`.
- **`AppNav`** (`src/components/AppNav.jsx`) — the four section tabs.
- **`CampaignDropdown`** (`src/components/CampaignDropdown.jsx`) — upper-right dropdown listing campaign names; selecting one stores the current campaign context. `role="listbox"`/`option`, `aria-expanded`, closes on outside-click and `Escape`.
- **`ThemeToggle`** (`src/components/ThemeToggle.jsx`) — light/dark switch.
- **`Login`** (`src/pages/Login.jsx`) — accessible sign-in form:
  - Labels associated with inputs (`htmlFor`/`id`); no placeholder-only fields.
  - `autocomplete="username"` / `autocomplete="current-password"`; paste allowed.
  - Password show/hide toggle (`aria-pressed`).
  - Inline error banner with `role="alert"`, `aria-invalid`, `aria-describedby`.
  - Submit feedback: loading spinner → success/error.
- **`Placeholder`** (`src/pages/Placeholder.jsx`) — stub body for the tab routes (real pages pending).

### State / context

- **`CampaignProvider`** (`src/context/CampaignContext.jsx`) — loads campaigns on mount, exposes `{ campaigns, selected, select, status }` via `useCampaign()`.
  - On 401 it clears tokens and redirects to `/login`.
  - Selected campaign is persisted to `localStorage["tj_campaign"]` as `{ name, tenant_id, uniqueURI }` (`src/lib/campaigns.js`); cleared if the saved one no longer exists in the list.
- **`src/lib/api.js`** — `login()`, `getCampaigns()` (bearer auth from `src/lib/auth.js`).
- **`src/lib/auth.js`** — token storage (`tj_access`, `tj_refresh`).

---

## Build & embedding

- Vite `base: "/admin/"`, `build.outDir: ../backend/internal/web/dist`.
- `backend/internal/web/web.go` embeds `dist` and serves it for `/admin*` with SPA fallback to `/`.
- `backend/Makefile` runs the frontend build before `go build`, so `make server` produces a single binary with the SPA inside.
- Generated `dist/` and `node_modules/` are gitignored.

Dev loop: `npm --prefix frontend run dev` for the SPA; backend serves the built output.

---

## Routes / API

Backend (all except login require `Authorization: Bearer <access>`):

- `POST /rest/user/login` — JSON `{"username": "...", "password": "..."}`; returns `{"access": "...", "refresh": "..."}` on 200.
- `GET /rest/user/me` — on 401, redirect to `/login`.
- `GET /rest/campaigns/list` — returns `{"total": int, "campaigns": [{"name": str, "tenant_id": str, "uniqueURI": str}, ...]}`. Backed by `db.GetCampaigns()` / `server.HandleListCampaigns`.
- `POST /rest/campaign/new` — create a campaign; see `backend/internal/server/types.go`. Returns 200 / 400 / 500.

Client routing (React Router, `basename="/admin"`):

- `/login` — public.
- `/` (protected) → redirects to `/campaigns`.
- `/campaigns`, `/automation`, `/codes`, `/settings` — protected tabs (placeholder bodies).
- `*` → `/`.
