#!/bin/sh

# ORDER BY (timestamp) DESC coloca as mais recentes primeiro
# Quero o contrário, então ASC
docker compose exec db psql -U rinheiro -d rinha-go-crebito --echo-all -c 'SELECT * FROM transacoes ORDER BY realizada_em ASC;'
