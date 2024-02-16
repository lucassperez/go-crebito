#!/bin/sh

if [ "$1" = nojq ]; then
  nojq=sim
  shift
fi

# At first I tried echoing directly curl ${@:- },
# but then when $@ was empty, it was left with ugly
# spaces between curl and the rest of the arguments.
# So I figured a good ol' if/else would fix it.
if [ -z "$@" ]; then
  curl_to_print='curl'
else
  curl_to_print="curl $@"
fi

if [ "$nojq" = sim ]; then
  echo "$(tput bold)$curl_to_print localhost:4000/clientes/1/extrato$(tput sgr0)"
  curl $@ localhost:4000/clientes/1/extrato

  echo "$(tput bold)$curl_to_print -X POST localhost:4000/clientes/1/transacoes$(tput sgr0)"
  curl $@ -X POST localhost:4000/clientes/1/transacoes
else
  echo "$(tput bold)$curl_to_print localhost:4000/clientes/1/extrato | jq$(tput sgr0)"
  curl $@ localhost:4000/clientes/1/extrato | jq

  echo "$(tput bold)$curl_to_print -X POST localhost:4000/clientes/1/transacoes | jq$(tput sgr0)"
  curl $@ -X POST localhost:4000/clientes/1/transacoes | jq
fi
