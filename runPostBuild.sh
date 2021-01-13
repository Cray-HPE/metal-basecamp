#!/bin/bash

set -ex

SANS_TIME_TAG="$(cat .version)-$(git rev-parse --short HEAD)"
docker pull "$IMAGE_VERSIONED"
docker tag "$IMAGE_VERSIONED" "${IMAGE_VERSIONED%:*}:$SANS_TIME_TAG"
docker push "${IMAGE_VERSIONED%:*}:$SANS_TIME_TAG" || echo "$SANS_TIME_TAG already exists because the RPM beat the docker image"
