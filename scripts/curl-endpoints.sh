#!/bin/sh

# At first I tried echoing directly curl ${@:- },
# but then when $@ was empty, it was left with ugly
# spaces between curl and the rest of the arguments.
# So I figured a good ol' if/else would fix it.
if [ -z "$@" ]; then
  curl_to_print='curl'
else
  curl_to_print="curl $@"
fi

echo ">> $curl_to_print localhost:4000/clientes/1/extrato | jq"
curl $@ localhost:4000/clientes/1/extrato | jq

echo ">> $curl_to_print -X POST localhost:4000/clientes/1/transacoes | jq"
curl $@ -X POST localhost:4000/clientes/1/transacoes | jq
