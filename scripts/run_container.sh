#!/bin/bash
set -e

DIR=$(cd $( dirname ${BASH_SOURCE[0]})/.. && pwd)

# run test.
cd $DIR/src && go test .

# build.
DOCKER_BUILDKIT=1 docker build $DIR/src/ \
  -f $DIR/docker/ratelimit.Dockerfile \
  -t local/ratelimit

# run.
docker run \
  --rm \
  -p 8080:8080 \
  -e RL_LIMIT=60 \
  -e RL_WINDOW=60 \
  local/ratelimit