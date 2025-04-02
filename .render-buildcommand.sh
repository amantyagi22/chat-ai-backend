#!/bin/bash
go mod download
go build -tags netgo -ldflags '-s -w' -o app 