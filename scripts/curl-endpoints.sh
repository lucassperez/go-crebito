#!/bin/sh

for arg in "$@"; do
  case $arg in
    e|extrato)
      extrato=sim
      ;;
    t|transacoes)
      transacoes=sim
      ;;
    v|-v)
      verbose='-v'
      shift
      ;;
    *)
      echo "Make requests to the server's endpoints."
      echo
      echo -e '\tGET /clientes/1/extrato'
      echo -e '\tPOST /clientes/1/transacoes {"valor": 100, "tipo": "c", "descricao": "teste-curl"}'
      echo
      echo 'If the argument `e` or `extrato` is passed, it will perform the GET to extrato endpoint.'
      echo 'If the argument `t` or `transacoes` is passed, it will perform the POST to transacoes endpoint.'
      echo 'If none of those is passed, it will perform both requests.'
      echo 'If argument `v` or `-v` is passed, curl will be run with verbose flag, -v.'
      echo 'Other arguments show this help message.'
      exit
      ;;
  esac
done

if [ "$#" -eq 0 ]; then
  extrato=${extrato:-sim}
  transacoes=${transacoes:-sim}
fi

if [ -n "$extrato" ]; then
  curl \
    $verbose \
    --header 'Content-Type: application/json;' \
    localhost:4000/clientes/1/extrato
fi

if [ -n "$transacoes" ]; then
  curl \
    $verbose \
    -X POST \
    --header 'Content-Type: application/json;' \
    -d '{"valor": 100, "tipo": "c", "descricao": "teste-curl"}' \
    localhost:4000/clientes/1/transacoes
fi
