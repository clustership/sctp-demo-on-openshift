#!/bin/bash

IMAGE=go-sctp
REPO=quay.io/xymox

podman build --rm -t $IMAGE -f Dockerfile.golang .
podman tag $IMAGE $REPO/${IMAGE}:latest
podman push $IMAGE $REPO/${IMAGE}:latest
