#!/bin/sh
go build -ldflags "-X main.buildDate=$(date -Iseconds) -X main.commit=$(git rev-list HEAD --count)"
