GOOSE_OPTS=-dir /db/migrations postgres
RUN=docker compose run --rm development

help: ##@Miscellaneous Show this help message
	@python3 ./scripts/help.py
.PHONY: help

dev: ##@Development Start handy-automator in development mode
	@docker compose up development
.PHONY: dev

gen-api: ##@Development Generate API code from the OpenAPI spec
	@$(RUN) go generate ./api/...
.PHONY: gen-api

lint-api: ##@Development Lint OpenAPI spec
	@docker compose run --rm --workdir /api development redocly lint /api/openapi.yml # can't use $(RUN) cause I need to specify the workdir
.PHONY: lint-api

start-db: ##@Database Instantiate a database container
	@docker compose up -d postgres
.PHONY: start-db

stop-db: ##@Database Stop the database container
	@docker compose down postgres
.PHONY: stop-db

destroy-db: ##@Database Stop the database container and erase all its data
	@docker compose down postgres -v
.PHONY: destroy-db

mup: start-db check-db-status ##@Database Apply all migrations
	@$(RUN) goose $(GOOSE_OPTS) up
.PHONY: mup

mdown: start-db check-db-status ##@Database Undo the last migration
	@$(RUN) goose $(GOOSE_OPTS) down
.PHONY: mdown

migration: ##@Database Create a new migration SQL file in db/migrations
	@$(if $(strip $(name)),,$(error Usage: make migration name=name_of_migration_file))
	@$(RUN) goose $(GOOSE_OPTS) create $(name) sql
.PHONY: migration

check-db-status: ##@Miscellaneous Check if the database is ready to accept connections (internal)
	@docker compose exec postgres bash /check-db-connection.sh
.PHONY:check-db-status

run-py: ##@Development Run the main.py file and start the program
	@cd ./python-processing && uv run main.py
.PHONY: run-py