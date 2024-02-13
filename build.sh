#!/bin/sh
go build -ldflags "-X main.buildDate=$(date -I) -X main.commit=$(git rev-list HEAD --count)"
