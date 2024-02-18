#!/bin/sh

retries=''

echo "`date '+[%Y-%m-%d/%H:%M:%S]'` Checking if database is already up and running"
printf "`date '+[%Y-%m-%d/%H:%M:%S]'` Retries: "

while true; do
  if docker compose exec db pg_isready --quiet; then
    printf "\n`date '+[%Y-%m-%d/%H:%M:%S]'` Postgres is already accepting connections\n"
    exit 0
  fi
  retries="$retries."
  if [ $retries = '..........' ]; then
    echo
    echo "`date '+[%Y-%m-%d/%H:%M:%S]'` `tput bold``tput setaf 3`[`basename $0`]`tput sgr0` Exceeded retry limit"
    exit 1
  fi
  printf .
  sleep 2
done
