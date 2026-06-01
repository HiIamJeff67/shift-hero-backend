#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage:
  ./scripts/bootstrap.sh --project <display-name> --module <go-module-path> --service <service-name> [--api-base <base-path>]

Required:
  --project   Human-readable project name for docs/readme
  --module    Go module path (example: github.com/acme/my-backend)
  --service   Service name used for compose/container/otel/logical naming

Optional:
  --api-base  API base path segment (default keeps current value, usually: api)
USAGE
}

PROJECT_NAME=""
MODULE_PATH=""
SERVICE_NAME=""
API_BASE_PATH=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --project)
      PROJECT_NAME="${2:-}"
      shift 2
      ;;
    --module)
      MODULE_PATH="${2:-}"
      shift 2
      ;;
    --service)
      SERVICE_NAME="${2:-}"
      shift 2
      ;;
    --api-base)
      API_BASE_PATH="${2:-}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$PROJECT_NAME" || -z "$MODULE_PATH" || -z "$SERVICE_NAME" ]]; then
  echo "Missing required arguments." >&2
  usage
  exit 1
fi

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "Error: bootstrap must run inside a git repository." >&2
  echo "Hint: run 'git init' first if this folder has no .git directory." >&2
  exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
  echo "Error: git working tree is not clean. Commit or stash changes before bootstrap." >&2
  exit 1
fi

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

CURRENT_MODULE="$(awk '/^module / { print $2; exit }' go.mod 2>/dev/null || true)"
CURRENT_SERVICE="$(rg -n '^APP_NAME=' .env 2>/dev/null | head -n1 | sed 's/^.*APP_NAME=//' || true)"
CURRENT_PROJECT="$(head -n1 README.md 2>/dev/null | sed 's/^# //')"

replace_all() {
  local search="$1"
  local replace="$2"

  while IFS= read -r file; do
    SEARCH="$search" REPLACE="$replace" perl -0pi -e 's/\Q$ENV{SEARCH}\E/$ENV{REPLACE}/g' "$file"
  done < <(rg -l --hidden --glob '!.git' --fixed-strings "$search" . || true)
}

update_env_key() {
  local file="$1"
  local key="$2"
  local value="$3"

  [[ -f "$file" ]] || return 0

  if rg -q "^${key}=" "$file"; then
    KEY="$key" VALUE="$value" perl -0pi -e 's/^$ENV{KEY}=.*$/$ENV{KEY}=$ENV{VALUE}/mg' "$file"
  else
    printf "\n%s=%s\n" "$key" "$value" >> "$file"
  fi
}

# 1) module/import replacements
if [[ -n "$CURRENT_MODULE" ]]; then
  replace_all "$CURRENT_MODULE" "$MODULE_PATH"
fi
replace_all "github.com/HiIamJeff67/shift-hero-backend" "$MODULE_PATH"
replace_all "github.com/HiIamJeff67/shift-hero-backend" "$MODULE_PATH"

# 2) service naming replacements
if [[ -n "$CURRENT_SERVICE" ]]; then
  replace_all "$CURRENT_SERVICE" "$SERVICE_NAME"
fi
replace_all "shift-hero" "$SERVICE_NAME"
replace_all "shift-hero" "$SERVICE_NAME"

# 3) project display name replacements in docs
if [[ -n "$CURRENT_PROJECT" ]]; then
  replace_all "$CURRENT_PROJECT" "$PROJECT_NAME"
fi
replace_all "shift-hero-backend" "$PROJECT_NAME"
replace_all "shift-hero-backend" "$PROJECT_NAME"

# 4) env alignment
for env_file in .env .env.example; do
  update_env_key "$env_file" "APP_NAME" "$SERVICE_NAME"
  update_env_key "$env_file" "WORKSPACE_DIR" "/$SERVICE_NAME"
  update_env_key "$env_file" "ENTRYPOINT_CMD" "./$SERVICE_NAME"

  update_env_key "$env_file" "DB_HOST" "$SERVICE_NAME-db"
  update_env_key "$env_file" "DB_NAME" "$SERVICE_NAME-db"
  update_env_key "$env_file" "REDIS_HOST" "$SERVICE_NAME-redis"

  update_env_key "$env_file" "DOCKER_REDIS_SERVICE_NAME" "$SERVICE_NAME-redis"
  update_env_key "$env_file" "DOCKER_DB_SERVICE_NAME" "$SERVICE_NAME-db"
  update_env_key "$env_file" "DOCKER_API_SERVICE_NAME" "$SERVICE_NAME-api"
  update_env_key "$env_file" "DOCKER_NGINX_SERVICE_NAME" "$SERVICE_NAME-nginx"
  update_env_key "$env_file" "DOCKER_LOKI_SERVICE_NAME" "$SERVICE_NAME-loki"
  update_env_key "$env_file" "DOCKER_TEMPO_SERVICE_NAME" "$SERVICE_NAME-tempo"
  update_env_key "$env_file" "DOCKER_MIMIR_SERVICE_NAME" "$SERVICE_NAME-mimir"
  update_env_key "$env_file" "DOCKER_OTEL_COLLECTOR_SERVICE_NAME" "$SERVICE_NAME-otel-collector"
  update_env_key "$env_file" "DOCKER_GRAFANA_SERVICE_NAME" "$SERVICE_NAME-grafana"

done

if [[ -n "$API_BASE_PATH" ]]; then
  replace_all "API_BASE_PATH=api" "API_BASE_PATH=$API_BASE_PATH"
  replace_all "api" "$API_BASE_PATH"
  update_env_key ".env" "API_BASE_PATH" "$API_BASE_PATH"
  update_env_key ".env.example" "API_BASE_PATH" "$API_BASE_PATH"
fi

# tidy dependencies after module path replacement
GOFLAGS="${GOFLAGS:-}" go mod tidy

cat <<NEXT
Bootstrap completed.

Next steps:
  1. cp .env.example .env
  2. go run main.go

Or with Docker Compose:
  1. cp .env.example .env
  2. docker compose up -d --build
  3. make migrate-hotreload-db
  4. make seed-hotreload-db
NEXT
