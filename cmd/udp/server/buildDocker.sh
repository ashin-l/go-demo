#!/bin/bash

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main

docker build . -t godemo/udpserver:1.0
