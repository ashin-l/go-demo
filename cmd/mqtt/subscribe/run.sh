#!/bin/bash

while IFS='=' read -r key val
do
    export $key=$val;
done < .env

echo $MY_ADDR

go run sub.go