# FinEdu — AI-Powered Investment Research & Education Platform

An AI-powered platform that helps retail investors research companies, understand
financial metrics, manage a portfolio, and learn investing concepts — all in one place,
with transparent reasoning instead of blind "Buy/Sell" calls.

## Tech Stack

| Layer          | Choice                                             |
|----------------|-----------------------------------------------------|
| Frontend       | React 18 + TypeScript + Vite + Tailwind CSS v4      |
| Routing        | React Router                                        |
| Charts         | Recharts                                            |
| Backend        | Go + Gin                                            |
| Database (dev) | PostgreSQL, local via Docker Compose                |
| Database (prod)| Supabase Postgres                                   |
| Auth           | Supabase Auth (used in both dev and prod)           |
| AI             | Gemini API (news summaries); chat/portfolio analysis TBD |
| Market Data    | Alpha Vantage / Finnhub / Twelve Data (TBD by API limits) |
| Deployment     | Vercel (frontend), Render/Railway (backend), Supabase (prod DB) |

**Why this stack:** Go pairs well with TypeScript's static typing for a
multi-person team shipping features in parallel — shared type contracts between
API responses and frontend props cut down on integration bugs. Postgres runs
locally during development (no network dependency, fast iteration) and the team
switches to Supabase's hosted Postgres for production. Supabase Auth is used in
both environments since it's a separate hosted service from the DB itself.

## Project Structure

```
FinEdu/
├── docker-compose.yml       # Local Postgres for dev
├── backend/                 # Go + Gin API
│   ├── cmd/server/          # Entry point (main.go)
│   └── internal/
│       ├── config/          # Env var loading
│       ├── db/              # Postgres connection pool
│       ├── handlers/        # Route handlers, one file per feature area
│       ├── models/          # Shared data structs
│       └── routes/          # Route registration
└── frontend/                # React + TypeScript + Vite
    └── src/
        ├── components/layout/  # Navbar, page layout
        ├── features/            # One folder per feature area (below)
        ├── lib/                 # Supabase client, API fetch wrapper
        ├── pages/                # Top-level routes (Home)
        └── types/                # Shared TS interfaces
```

## Feature Ownership

Matches the assignments from the project proposal:

| Feature area                                          | Owner  | Frontend                          | Backend handlers          |
|--------------------------------------------------------|--------|------------------------------------|----------------------------|
| Company Search & Market Dashboard + Visualizations     | Nhien  | `features/company-dashboard`       | `handlers/company.go`      |
| News Aggregation Summary + Interactive Stock Charts    | Nhi    | `features/news-charts`             | `handlers/news.go`         |
| AI Learning Center                                     | Hiếu   | `features/learning-center`         | `handlers/learning.go`     |
| Portfolio Management + Risk                            | Quang  | `features/portfolio`               | `handlers/portfolio.go`    |

Each handler file currently returns `501 Not Implemented` stubs — replace the
body with real logic as you build out your feature. Add new routes in
`backend/internal/routes/routes.go`.

## Getting Started

### Prerequisites

- [Node.js](https://nodejs.org/) 20.19+ or 22.12+ (older 22.x patch versions hit a
  known npm/Windows bug with newer tooling — check with `node -v`)
- [Go](https://go.dev/) 1.23+
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (for local Postgres)
- A free [Supabase](https://supabase.com/) project (for Auth, and later prod DB)

### 1. Clone and configure environment variables

```bash
git clone <repo-url>
cd FinEdu
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

Fill in `SUPABASE_URL` (backend) and `VITE_SUPABASE_URL` / `VITE_SUPABASE_ANON_KEY`
(frontend) from your Supabase project's Settings → API page. Ask in the team
chat if you don't have access to the shared Supabase project yet.

### 2. Start local Postgres

```bash
docker compose up -d
```

This starts Postgres on `localhost:5432` with the credentials already set in
`backend/.env.example` (`finedu` / `finedu`). No manual DB setup needed.

### 3. Run the backend

```bash
cd backend
go mod tidy
go run ./cmd/server
```

The API starts on `http://localhost:8080`. Check `http://localhost:8080/health`
to confirm it's up.

### 4. Run the frontend

```bash
cd frontend
npm install
npm run dev
```

The app starts on `http://localhost:5173`.

## Scripts

**Frontend** (`cd frontend`)
- `npm run dev` — start dev server
- `npm run build` — typecheck + production build
- `npm run lint` — run ESLint

**Backend** (`cd backend`)
- `go run ./cmd/server` — run the API
- `go build ./...` — build all packages
- `go vet ./...` — static analysis

## Branching

- `main` is protected — no direct commits.
- Branch per feature: `feature/<your-name>-<short-description>` (e.g. `feature/quang-portfolio-crud`).
- Open a PR into `main` and get at least one teammate review before merging.
