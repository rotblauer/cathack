#!/bin/sh

# run from project dir

echo "killing running chat process... " 
echo "(will fail to kill ps proc)"
kill $(ps aux | grep chat | awk '{print $2}')

echo "rebuilding chat from go"
go build -o main main.go

echo "starting chat with out to chat.log"
./main &>"$PWD/chat.log" &!

