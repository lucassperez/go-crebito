#!/bin/sh

docker compose exec db psql -U rinheiro -d rinha-go-crebito --echo-all -c 'SELECT * FROM clientes;'
