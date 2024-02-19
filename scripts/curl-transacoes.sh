#!/bin/sh

while getopts 'v:t:d:c:i:' flag; do
  case $flag in
    v)
      valor="$OPTARG"
      ;;
    t)
      tipo="$OPTARG"
      ;;
    d)
      descricao="$OPTARG"
      ;;
    c|i)
      cliente_id="$OPTARG"
  esac
done

valor=${valor:-1}
tipo=${tipo:-d}
descricao=${descricao:-teste}
cliente_id=${cliente_id:-1}

curl \
  -X POST \
  --header 'Content-Type: application/json;' \
  "localhost:4000/clientes/$cliente_id/transacoes" \
  -d "{\"valor\": ${valor}, \"tipo\": \"${tipo}\", \"descricao\": \"${descricao}\"}" \
