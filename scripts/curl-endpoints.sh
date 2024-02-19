#!/bin/sh

curl \
  --header 'Content-Type: application/json;' \
  localhost:4000/clientes/1/extrato

curl \
  -X POST \
  --header 'Content-Type: application/json;' \
  -d '{"valor": 100, "tipo": "c", "descricao": "teste-curl"}' \
  localhost:4000/clientes/1/transacoes
