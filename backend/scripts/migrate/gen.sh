#!/bin/bash

# migrate name argument

if [ -z "$1" ]; then
    echo "Please provide a name for the migration"
    exit 1
fi

migrate_name=$1

# Generate migration file

atlas migrate diff $migrate_name \
    --dir "file://migrations?format=goose" \
    --to "ent://app/ent/schema" \
    --dev-url "docker://postgres/16/test?search_path=public" \
    --format '{{ sql . "  " }}'