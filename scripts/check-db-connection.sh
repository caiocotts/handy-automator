#!/usr/bin/env bash

# exit immediately if database is ready
if [[ "$(pg_isready)" =~ "accepting connections" ]]; then exit 0; fi

echo "initializing database"
while ! [[ "$(pg_isready)" =~ "accepting connections" ]]; do
    sleep 1
    echo -n "."
done

printf "\ndatabase system is ready to accept connections\n"
exit 0
