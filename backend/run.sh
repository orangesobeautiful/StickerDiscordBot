#!/bin/sh

go mod tidy && \
go run -mod=mod entgo.io/ent/cmd/ent generate --template app/ent/extemplates --feature sql/upsert ./app/ent/schema && \
go build -ldflags="-s" . && 
./backend