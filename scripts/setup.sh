#!/usr/bin/env bash

go mod tidy 
go mod vendor
go get -u github.com/goware/modvendor
modvendor -copy="**/*.c **/*.h **/*.proto" -v
go mod tidy 

