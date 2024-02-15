help:
# Shows any commands below with ## after the make target name as a list of available commands
	@echo "Available Commands:"
	@grep -E '^[a-zA-Z0-9 -]+:.*#' Makefile | while read -r l; do printf "$$(tput setaf 2)$$(tput bold)$$(echo $$l | cut -f 1 -d':')$$(tput sgr0):$$(echo $$l | cut -f 2- -d'#')\n"; done

DB_COMMAND := docker compose exec db psql -U postgres --echo-all

bash:
	docker compose exec db bash

server: # starts the db service and run air
	docker compose up db -d
	air

up: # executes docker compose up
	docker compose up -d

down: # executes docker compose down
	docker compose down

psql: # start psql inside the db container
	docker compose exec db psql -U postgres

drop: # drop database postgres
	$(DB_COMMAND) -d template1 -c 'DROP DATABASE postgres;'

create: # create postgres database and executes the init.sql file
	$(DB_COMMAND) -d template1 -c 'CREATE DATABASE postgres;'
	$(DB_COMMAND) -d postgres --file /docker-entrypoint-initdb.d/init.sql

seed: # executes the seed.sql
	$(DB_COMMAND) --file /seed.sql

reset: drop create seed # database drop, create and seed
