#!/bin/bash

set -xe

cd server
./build.sh
cd ../client
./build.sh
cd ..
