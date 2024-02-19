#!/bin/sh

while getopts 'v:t:d:c:i:qV' flag; do
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
      ;;
    q)
      quiet=sim
      ;;
    V)
      curl_v=' -v'
  esac
done

valor=${valor:-1}
tipo=${tipo:-d}
descricao=${descricao:-teste}
cliente_id=${cliente_id:-1}

if [ -z "$quiet" ]; then
  echo "curl$curl_v -X POST --header 'Content-Type: application/json;' localhost:4000/clientes/$cliente_id/transacoes -d \"{\"valor\": $valor, \"tipo\": \"$tipo\", \"descricao\": \"$descricao\"}\" $curl_v"
fi

curl \
  $curl_v \
  -X POST \
  --header 'Content-Type: application/json;' \
  "localhost:4000/clientes/$cliente_id/transacoes" \
  -d "{\"valor\": ${valor}, \"tipo\": \"${tipo}\", \"descricao\": \"${descricao}\"}" \
