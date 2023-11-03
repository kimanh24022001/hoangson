#! /usr/bin/bash

# sudo docker pull postgres
sudo docker run --name postgres -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres
