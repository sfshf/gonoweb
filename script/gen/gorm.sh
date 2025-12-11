#!/usr/bin/env bash

# generate models from local db
gentool -dsn 'root:dev123@tcp(127.0.0.1:3306)/gonoweb?charset=utf8mb4&parseTime=True&loc=Local' -onlyModel -outPath './internal/model'

# format go codes
go fmt ./...