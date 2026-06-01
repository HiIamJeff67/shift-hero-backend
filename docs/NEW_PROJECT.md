# New Project Guide

## 1) Create From Template

1. Create a new repository from this template (GitHub template or manual copy).
2. Clone the new repository locally.
3. Ensure `.git` exists and working tree is clean.

## 2) Bootstrap Project Identity

Run bootstrap once from repo root:

```bash
./scripts/bootstrap.sh \
  --project "My Backend Platform" \
  --module "github.com/acme/my-backend-platform" \
  --service "my-backend-platform" \
  --api-base "api"
```

What it does:

- Replaces module/import path defaults.
- Replaces service naming defaults used by compose/otel/docs/env.
- Updates `.env` / `.env.example` key defaults.
- Optionally updates API base path.
- Runs `go mod tidy`.

## 3) Start The Project

### Local Host Mode

```bash
cp .env.example .env
go run main.go
```

### Docker Compose Mode

```bash
cp .env.example .env
docker compose up -d --build
make migrate-hotreload-db
make seed-hotreload-db
```

## 4) Add A New Module (Checklist)

1. Add schema/input/repository/service/controller/binder for the module under `app/`.
2. Register module wiring in `app/modules/`.
3. Add route entry under `app/routes/developmentroutes/` and call it from `development_routes.go`.
4. Add migrations/seeds if new tables/enums are introduced.
5. Add tests under `test/unit` and/or `test/e2e`.
6. Run:

```bash
go test ./...
go build ./...
```
