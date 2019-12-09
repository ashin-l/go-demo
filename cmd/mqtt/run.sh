#!/bin/bash

command=$1

while IFS='=' read -r key val
do
    export $key=$val;
done < .env

echo $MY_ADDR

function publish() {
    go run ./publish/pub.go
}

function subscribe() {
    go run ./subscribe/sub.go
}

case $command in
    ('pub')
    publish
    ;;
    ('sub')
    subscribe
    ;;
    (*)
    echo 'error'
esac