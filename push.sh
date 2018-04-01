#!/bin/bash
set -e

version=$(date -u '+%Y%m%d%H%M%S')
image="us.gcr.io/neil-164300/webhooksub"

docker build -t "${image}:latest" .
docker tag "${image}:latest" "${image}:${version}"
docker push "${image}:latest"
docker push "${image}:${version}"
echo "pushed ${image}:${version}"
