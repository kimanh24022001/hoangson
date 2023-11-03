#! /usr/bin/bash

set -e

if ! [ -x "$(command -v psql)" ]; then
    sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
    apt-get update
    apt-get -y install postgresql

    # TODO: Copy config stuffs
    pushd /tmp
    sudo -u postgres psql -c "create database monolith;"
    sudo -u postgres psql -c "create user monolith_admin with encrypted password 'monolith_admin_arst'";
    sudo -u postgres psql -c "grant all privileges on database monolith to monolith_admin";
    sudo -u postgres psql -d "monolith" -c "grant all on schema public to monolith_admin";
    popd
fi

echo "===> complete installing POSTGRES"
