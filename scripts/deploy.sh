#!/bin/bash

git pull origin master
docker-compose down
docker-compose up -d