#!/bin/env python

import platform
import os

GO_CODE_PATH = os.path.join('.', 'decision-maker' )
HELP_SCRIPT_PATH = os.path.join('.', 'scripts', 'help.py')
MIGRATIONS_PATH = os.path.join('.', 'db', 'migrations')
API_SPEC_PATH = os.path.join('.', 'api')

MAKEFILE_CONTENTS = \
f"""GOOSE=go tool goose
GOOSE_OPTS=-dir {MIGRATIONS_PATH} postgres

dev: ##@Development Start handy-automator in development mode
	@cd {GO_CODE_PATH} &&\\
 	go tool reflex -c reflex.conf
.PHONY: dev

gen-api: ##@Development Generate API code from the OpenAPI spec
	@cd {GO_CODE_PATH} &&\\
    go generate ./api/...
.PHONY: gen-api

lint-api: ##@Development Lint OpenAPI spec
	@docker run --rm -tv {API_SPEC_PATH}:/spec redocly/cli lint openapi.yml
.PHONY: lint-api

start-db: ##@Database Instantiate a database container
	@docker compose up -d
.PHONY: start-db

stop-db: ##@Database Stop the database container
	@docker compose down
.PHONY: stop-db

destroy-db: ##@Database Stop the database container and erase all its data
	@docker compose down -v
.PHONY: destroy-db

mup: start-db check-db-status ##@Database Apply all migrations
	@$(GOOSE) $(GOOSE_OPTS) up
.PHONY: mup

mdown: start-db check-db-status ##@Database Undo the last migration
	@$(GOOSE) $(GOOSE_OPTS) down
.PHONY: mdown

migration: ##@Database Create a new migration SQL file in db/migrations
	@$(if $(strip $(name)),,$(error Usage: make migration name=name_of_migration_file))
	@$(GOOSE) $(GOOSE_OPTS) create $(name) sql
.PHONY: migration

check-db-status: ##@Miscellaneous Check if the database is ready to accept connections (internal)
	@docker compose exec postgres bash /check-db-connection.sh
.PHONY:check-db-status

help: ##@Miscellaneous Show this help message
	@python3 {HELP_SCRIPT_PATH}
.PHONY: help
"""

with open('Makefile', 'w') as f:
    f.write(MAKEFILE_CONTENTS)

print('Makefile generated')
