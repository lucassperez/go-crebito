#!/bin/sh

# ORDER BY (timestamp) DESC coloca as mais recentes primeiro
docker compose exec db psql -U postgres -d rinha-go-crebito --echo-all -c 'SELECT * FROM transacoes ORDER BY realizada_em DESC;'
