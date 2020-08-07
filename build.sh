#!/bin/bash

set -e

env GOOS=linux GOARCH=amd64 go build -a -v -o executable .
docker build -t brunoanjos/trainers-test:latest .
docker push brunoanjos/trainers-test:latest
