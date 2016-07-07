#!/bin/sh

# run from project dir

echo "killing running main process... " 
echo "(will fail to kill ps proc)"
kill $(ps aux | grep main | awk '{print $2}')

echo "rebuilding main into main* from main.go"
go build -o main main.go

echo "starting main with out to main.log"
./main &>"$PWD/main.log" &!

