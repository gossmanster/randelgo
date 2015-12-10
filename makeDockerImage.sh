#!/usr/bin/env bash
set -xe 	

# Statically link all packages so program is completely standalone and
# works in a minimal Docker image
CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo .
sudo docker build -t randelgo .
