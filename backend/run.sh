#!/bin/sh

go mod tidy && \
go run -mod=mod entgo.io/ent/cmd/ent generate --template app/ent/extemplates ./app/ent/schema && \
go build -ldflags="-s" . && 
./backend