#!/bin/bash

set -xe

docker build -t ngpitt/blinkt-animation-controller-client:v1 .
docker push ngpitt/blinkt-animation-controller-client
