#!/bin/bash

while IFS='=' read -r key val
do
    export "$key"="$val"
done < .env

./util