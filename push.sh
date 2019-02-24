#!/bin/bash
set -e

version=$(date -u '+%Y%m%d%H%M%S')
image="us.gcr.io/neil-164300/webhooksub"

docker build -t "${image}:latest" .
docker tag "${image}:latest" "${image}:${version}"
docker push "${image}:latest"
docker push "${image}:${version}"
echo "pushed ${image}:${version}"

echo "to update, run: gcloud beta compute instances update-container irbridge --container-image ${image}:${version}`

