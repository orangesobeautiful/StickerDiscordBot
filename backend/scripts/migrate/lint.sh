#!/bin/bash

# latest argument (default to 1)

latest=${1:-1}

# Generate migration file

atlas migrate lint \
    --dir "file://migrations?format=goose" \
    --dev-url "docker://postgres/16/test?search_path=public" \
    --latest $latest