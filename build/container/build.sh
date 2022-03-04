#!/bin/bash

docker build -t main-server:5443/builder .
# docker push main-server:5443/builder
docker run -p 8080:8080 main-server:5443/builder