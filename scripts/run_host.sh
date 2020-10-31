#!/bin/bash
set -e

DIR=$(cd $( dirname ${BASH_SOURCE[0]})/.. && pwd)

# run test.
cd $DIR/src && go test .

# run.
cd $DIR/src/server && RL_LIMIT=60 RL_WINDOW=60 go run .