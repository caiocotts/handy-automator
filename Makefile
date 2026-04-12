GOOSE_OPTS=-dir /db/migrations postgres
TOOL=docker compose run --rm tool

help: ##@Miscellaneous Show this help message
	@uv run ./scripts/help.py
.PHONY: help

dm: ##@Development Start the decision-maker in development mode
	@cd ./decision-maker && go tool reflex --decoration=none --config=./reflex.conf
.PHONY: dm

pu: ##@Development Start the processing-unit in development mode
	@cd ./processing-unit && uv run main.py
.PHONY: pu

test-pu: ##@Development Run the processing-unit test suite
	@cd ./processing-unit && uv run python ./test_suite/test.py
.PHONY: test-pu

gen-api: lint-api ##@Development Generate API code from the OpenAPI spec
	@$(TOOL) go generate ./api/...
.PHONY: gen-api

lint-api: ##@Development Lint OpenAPI spec
	@docker compose run --rm --workdir /api tool redocly lint /api/openapi.yml # can't use $(TOOL) cause I need to specify the workdir
.PHONY: lint-api

start-db: ##@Database Instantiate a database container
	@docker compose up -d postgres
.PHONY: start-db

stop-db: ##@Database Stop the database container
	@docker compose stop postgres
.PHONY: stop-db

destroy-db: ##@Database Stop the database container and erase all its data
	@docker compose down postgres -v
.PHONY: destroy-db

seed-db: start-db check-db-status ##@Database Seed the database with sample data
	@$(TOOL) goose -dir /db/seed postgres -no-versioning up
.PHONY: seed-db

mup: start-db check-db-status ##@Database Apply all migrations
	@$(TOOL) goose $(GOOSE_OPTS) up
.PHONY: mup

mdown: start-db check-db-status ##@Database Undo the last migration
	@$(TOOL) goose $(GOOSE_OPTS) down
.PHONY: mdown

migration: ##@Database Create a new migration SQL file in db/migrations
	@$(if $(strip $(name)),,$(error Usage: make migration name=name_of_migration_file))
	@$(TOOL) goose $(GOOSE_OPTS) create $(name) sql
.PHONY: migration

check-db-status: ##@Miscellaneous Check if the database is ready to accept connections (internal)
	@docker compose exec postgres bash /check-db-connection.sh
.PHONY:check-db-status
