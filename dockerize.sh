#!/bin/bash

set -xe

cd server
./dockerize.sh
cd ../client
./dockerize.sh
cd ..
