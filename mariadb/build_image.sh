#!/bin/sh

docker build --no-cache -f ./mariadb/mariadb.containerfile -t ilhammhdd/mariadb-oauth2-go:ubuntu-focal --build-arg MARIADB_ROOT_PASSWORD=passwordformariadbroot --build-arg MARIADB_CLIENT_PASSWORD=passwordformariadb .