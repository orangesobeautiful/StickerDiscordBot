#!/bin/sh

go mod tidy && \
go generate ./app/ent && \
go build -ldflags="-s" . && 
./backend