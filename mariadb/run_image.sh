#!/bin/sh

docker run --detach --rm --name mariadb_oauth2_go -v /home/ilhammhdd/oauth-go/mariadb-data:/var/lib/mariadb-data -p 4512:3306 ilhammhdd/mariadb-oauth2-go:ubuntu-focal