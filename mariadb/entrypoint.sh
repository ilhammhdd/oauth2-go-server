#!/bin/bash

if [ -z "$(ls -A /var/lib/mariadb-data)" ]; then
  mariadb-install-db --user=mysql
fi

echo "starting mariadb service..."
service mariadb start

echo "Waiting for MariaDB to be ready (accepting connections)"
until mariadb -e "SELECT 1"; do sleep 1; done
echo "MariaDB ready for accepting connection"

mariadb-tzinfo-to-sql /usr/share/zoneinfo | mariadb -u root mysql

mariadb -e "UPDATE mysql.global_priv SET priv=json_set(priv, '$.plugin', 'mysql_native_password', '$.authentication_string', PASSWORD('${MARIADB_ROOT_PASSWORD}')) WHERE User='root'; \
  DELETE FROM mysql.global_priv WHERE User=''; \
  DELETE FROM mysql.global_priv WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1'); \
  DROP DATABASE IF EXISTS test; \
  DELETE FROM mysql.db WHERE Db='test' OR Db='test\\_%'; \
  FLUSH PRIVILEGES;"
  
echo "starting mariadb service..."
service mariadb stop
service mariadb start

echo "Waiting for MariaDB to be ready (accepting connections)"
until mariadb -e "SELECT 1"; do sleep 1; done
echo "MariaDB ready for accepting connection"

mariadb --user=root --password=$MARIADB_ROOT_PASSWORD -e "CREATE USER IF NOT EXISTS \
  'client'@'%' IDENTIFIED BY '${MARIADB_CLIENT_PASSWORD}'; \
  CREATE DATABASE IF NOT EXISTS oauth2_go CHARACTER SET utf8; \
  GRANT ALL PRIVILEGES ON oauth2_go.* TO 'client'@'%';"

for filename in /home/mariadb_migrations/*.sql; do
  echo "migrating $filename"
  mariadb --user=root --password=$MARIADB_ROOT_PASSWORD oauth2_go < $filename
done

service mariadb stop

exec "$@"