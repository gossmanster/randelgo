#!/usr/bin/env bash
set -xe 	

PUSHIMAGE=0

while getopts ":p" option
do
    case "$option" in
    
        p)
            PUSHIMAGE=1
            ;;
    esac
done

# Statically link all packages so program is completely standalone and
# works in a minimal Docker image
CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo .

if [ $PUSHIMAGE -ne "0" ]; then
    sudo docker build -t gossmanster/randelgo .
    sudo docker push gossmanster/randelgo
else
    sudo docker build -t randelgo .
fi