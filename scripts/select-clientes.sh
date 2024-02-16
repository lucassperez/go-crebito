#!/bin/sh

docker compose exec db psql -U postgres -d rinha-go-crebito --echo-all -c 'SELECT * FROM clientes;'
