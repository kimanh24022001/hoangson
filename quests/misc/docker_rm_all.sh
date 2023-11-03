#! /usr/bin/bash

# Remove all containers
sudo docker rm -vf $(sudo docker ps -aq)

# Remove all images
sudo docker rmi -f $(sudo docker images -aq)
