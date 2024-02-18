.PHONY: *

help:
# Shows with # after the make target name as a list of available commands. Also show the aliases
	@echo "Available Commands:"
	@grep -E '^[a-zA-Z0-9 -]+:.*#' Makefile | while read -r l; do printf "$$(tput setaf 2)$$(tput bold)$$(echo $$l | cut -f 1 -d':')$$(tput sgr0):$$(echo $$l | cut -f 2- -d'#')\n"; done
	@echo --- Aliases ---
	@grep '#[ a-zA-Z]*[aA]liases' Makefile -A 50 | tail -n +2 | sort | while read -r l; do printf "$$(tput setaf 3)$$(tput bold)$$(echo $$l | cut -d ':' -f 1)$$(tput sgr0):$$(echo $$l | cut -d ':' -f 2-)\n"; done

PSQL_COMMAND := docker compose exec db psql -U postgres -d rinha-go-crebito
DB_COMMAND := $(PSQL_COMMAND) --echo-all

bash: api # starts bash in api container
	docker compose exec api bash

server: db # starts the db service on the background and run api service with --rm option
	./scripts/test-db.sh && docker compose run --rm --service-ports api || docker compose stop db

up: # executes docker compose up
	docker compose up -d

down: # executes docker compose down
	docker compose down

db: # executes docker compose up db -d
	docker compose up db -d

api: db # executes docker compose up api -d
	docker compose up api -d

psql: # start psql inside the db container
	$(PSQL_COMMAND)

drop: # drop database rinha-go-crebito
	$(DB_COMMAND) -d template1 -c 'DROP DATABASE "rinha-go-crebito";'

create: # create rinha-go-crebito database and executes the init.sql file
	$(DB_COMMAND) -d template1 -c 'CREATE DATABASE "rinha-go-crebito";'
	$(DB_COMMAND) --file /docker-entrypoint-initdb.d/init.sql

seed: # executes the seed.sql
	$(DB_COMMAND) --file /seed.sql

reset: drop create seed # database drop, create and seed

# Aliases
b: bash
p: psql
s: server
u: up
d: down
r: reset
