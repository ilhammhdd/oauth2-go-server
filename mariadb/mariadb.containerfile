FROM ilhammhdd/mariadb-base-oauth2-go:ubuntu-focal 

ARG MARIADB_ROOT_PASSWORD=
ARG MARIADB_CLIENT_PASSWORD=

ENV MARIADB_ROOT_PASSWORD=$MARIADB_ROOT_PASSWORD
ENV MARIADB_CLIENT_PASSWORD=$MARIADB_CLIENT_PASSWORD

VOLUME /var/lib/mariadb-data

COPY ./mariadb/70-mariadb.cnf /etc/mysql/mariadb.conf.d/
RUN chmod 400 /etc/mysql/mariadb.conf.d/70-mariadb.cnf
COPY ./mariadb/entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh
RUN mkdir /home/mariadb_migrations
COPY ./mariadb/migrations/ /home/mariadb_migrations/

EXPOSE 3306

ENTRYPOINT ["entrypoint.sh"]
CMD ["mariadbd-safe"]
