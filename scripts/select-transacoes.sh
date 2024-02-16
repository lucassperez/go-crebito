#!/bin/sh

docker compose exec db psql -U postgres --echo-all -c 'SELECT * FROM transacoes;'
