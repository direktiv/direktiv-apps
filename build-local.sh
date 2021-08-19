#!/bin/bash

app=$1
if [[ $app =~ .*\/$ ]]; then
  app=${app::-1}
fi

echo "building $app"

docker build -t localhost:5000/$app:latest $1
docker push localhost:5000/$app:latest

if [[ $2 == "run" ]]; then
  docker run -it -p 8080:8080 localhost:5000/$app:latest
fi
