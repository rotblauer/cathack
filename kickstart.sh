#!/bin/sh

# run from project dir

echo "rebuilding chat from go"
go build -o chat chat.go

echo "starting chat with out to chat.log"
./chat &>"$PWD/chat.log" &!

