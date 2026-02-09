#!/bin/env python

import platform
import os

OS = platform.system()
GO_CODE_PATH = os.path.join('.', 'decision-maker' )

MAKEFILE_CONTENTS = \
f"""GOOSE=go run github.com/pressly/goose/v3/cmd/goose@latest

dev: mup
	@cd {GO_CODE_PATH} &&\\
 	go run .
.PHONY: dev

gen-api:
	@cd {GO_CODE_PATH} &&\\
    go generate ./api/...
.PHONY: gen-api

start-db:
	@docker compose up -d
.PHONY: start-db

stop-db:
	@docker compose down
.PHONY: stop-db

stop-and-erase-db:
	@docker compose down -v
.PHONY: stop-and-erase-db

mup: start-db check-db-connectivity
	@$(GOOSE) up
.PHONY: mup

mdown: start-db check-db-connectivity
	@$(GOOSE) down
.PHONY: mdown

check-db-connectivity:
	@docker compose exec postgres bash /check-db-connection.sh
.PHONY:check-db-connectivity

migration:
	@$(if $(strip $(name)),,$(error Usage: make migration name=name_of_migration_file))
	@$(GOOSE) create $(name) sql
.PHONY: migration
"""

with open('Makefile', 'w') as f:
    f.write(MAKEFILE_CONTENTS)

print('Makefile generated')
