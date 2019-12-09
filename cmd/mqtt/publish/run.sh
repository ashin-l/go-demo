#!/bin/bash

command=$1

while IFS='=' read -r key val
do
    export $key=$val;
done < .env

./publish

