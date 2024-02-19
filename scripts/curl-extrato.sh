#!/bin/sh

while getopts 'q' flag; do
  case $flag in
    q)
      quiet=sim
      shift
      ;;
  esac
done

if [ "$#" -gt 0 ]; then
  args=" $@"
fi

if [ -z "$quiet" ]; then
  echo "$(tput bold)curl$args localhost:4000/clientes/1/extrato$(tput sgr0)"
fi
curl $@ localhost:4000/clientes/1/extrato
