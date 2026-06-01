# shift-hero-backend

This repository is a reusable Go monolithic backend template developed by **Notezy**.

Ownership note:

- The architecture and original implementation are fully owned by Notezy.
- Notezy grants open-source usage rights for this repository under **Apache-2.0**.

Technical stack:

- Go `1.26`
- Gin + GORM + Cobra CLI
- REST + GraphQL
- PostgreSQL + Redis
- OpenTelemetry + Loki/Tempo/Mimir/Grafana

Current state: business modules are intentionally kept to preserve runtime stability while template/bootstrap infrastructure is in place.

## Quick Start

1. Copy environment file:

```bash
cp .env.example .env
```

2. Run locally:

```bash
go run main.go
```

Or run with Docker Compose:

```bash
docker compose up -d --build
make migrate-hotreload-db
make seed-hotreload-db
```

## Bootstrap For New Project

Use this once after creating a new repository from this template:

```bash
./scripts/bootstrap.sh \
  --project "My Backend Platform" \
  --module "github.com/acme/my-backend-platform" \
  --service "my-backend-platform" \
  --api-base "api"
```

Parameter reference:

- `--project`:
  - Project display name used in README/docs wording.
  - Example: `"My Backend Platform"`.
- `--module`:
  - Go module path used in `go.mod` and internal import paths.
  - Example: `"github.com/acme/my-backend-platform"` or `"my-backend-platform"`.
- `--service`:
  - Runtime/service identity used for docker service/container names, env defaults, and service naming constants.
  - Example: `"my-backend-platform"`.
- `--api-base` (optional):
  - API base segment for `API_BASE_PATH`.
  - Example: `"api"` -> `/api/development/v1/...`.

Bootstrap behavior:

- Replaces default module/import naming (`github.com/HiIamJeff67/shift-hero-backend`).
- Replaces service naming defaults (`shift-hero`) in compose/env/docs/constants.
- Updates env defaults (`APP_NAME`, docker service names, db/redis hosts, etc.).
- Optionally updates `API_BASE_PATH`.
- Runs `go mod tidy`.

## Common Commands

```bash
go test ./...
go build ./...
make migrate-hotreload-db
make seed-hotreload-db
```

## Docs

- New project flow: `docs/NEW_PROJECT.md`
- Migration scan notes: `docs/TEMPLATE_MIGRATION_NOTES.md`

## Repository Structure

```text
app/
infra/
shared/
test/
docs/
scripts/
```

## License

- `LICENSE` (Apache-2.0, granted by Notezy)
- `.github/LICENSE.zh.md` (Chinese reference translation)
- `.github/THIRD_PARTY_NOTICES.md` (consolidated third-party notices)

## Community

- Contributing: `.github/CONTRIBUTING.md` / `.github/CONTRIBUTING.zh.md`
- Donate/Sponsor: `.github/DONATE.md` / `.github/DONATE.zh.md`
- Funding config: `.github/FUNDING.yml`
- Code of Conduct: `.github/CODE_OF_CONDUCT.md` / `.github/CODE_OF_CONDUCT.zh.md`
- Security policy: `.github/SECURITY.md` / `.github/SECURITY.zh.md`
- Support: `.github/SUPPORT.md` / `.github/SUPPORT.zh.md`
