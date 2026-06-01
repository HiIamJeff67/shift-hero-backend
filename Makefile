# ============================== Database Shortcut Commands ============================== #
view-hotreload-dbs:
	docker compose exec -T go-start-monolithic-kit-api go run main.go viewDatabases

view-hotreload-enums:
	docker compose exec -T go-start-monolithic-kit-api go run main.go viewAllEnums

psql:
	docker exec -it go-start-monolithic-kit-db psql -U jeff -d go-start-monolithic-kit-db

# ============================== Migration Commands ============================== #
migrate-build-db:
	docker compose exec -T go-start-monolithic-kit-api ./go-start-monolithic-kit migrateDB
migrate-hotreload-db:
	docker compose exec -T go-start-monolithic-kit-api go run main.go migrateDB

clear-build-db:
	docker exec -i go-start-monolithic-kit-db psql -U jeff -d go-start-monolithic-kit-db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
clear-hotreload-db: # the same as the build version of db
	docker exec -i go-start-monolithic-kit-db psql -U jeff -d go-start-monolithic-kit-db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

remigrate-build-db:
	make clear-build-db
	make migrate-build-db

remigrate-hotreload-db:
	make clear-hotreload-db
	make migrate-hotreload-db

# ============================== Seeding Commands ============================== #
seed-build-db:
	docker compose exec -T go-start-monolithic-kit-api ./go-start-monolithic-kit seedDB
seed-hotreload-db:
	docker compose exec -T go-start-monolithic-kit-api go run main.go seedDB

clear-go-cache:
	go clean -modcache
	go mod download

test-auth-e2e:
	docker compose exec -T go-start-monolithic-kit-api go test ./test/e2e/auth

# ============================== GraphQL Shortcut Commands ============================== #
gql-generate: # update before generate
	go get github.com/99designs/gqlgen@v0.17.76
	go run github.com/99designs/gqlgen generate --config infra/graphql/gqlgen.yaml

gql-clean:
ifeq ($(OS),Windows_NT)
	@if exist app\graphql\generated\*.* del /q /s app\graphql\generated\*.*
	@if exist app\graphql\models\*.* del /q /s app\graphql\models\*.*
else
	rm -rf app/graphql/generated/*
	rm -rf app/graphql/models/*
endif

gql-regenerate:
	make gql-clean
	make gql-generate
