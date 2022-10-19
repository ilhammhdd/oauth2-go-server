ARG ARCH=linux/amd64

FROM --platform=${ARCH} ubuntu:focal

RUN apt-get update
RUN apt-get install -y software-properties-common gnupg

RUN apt-key adv --fetch-keys 'https://mariadb.org/mariadb_release_signing_key.asc'
RUN add-apt-repository 'deb [arch=amd64] https://archive.mariadb.org/mariadb-10.6.9/repo/ubuntu focal main'
RUN apt-get update
RUN apt-get install -y mariadb-server mariadb-client libmariadb3 mariadb-common mariadb-plugin-cracklib-password-check mariadb-plugin-connect