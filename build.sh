#!/bin/bash

set -e

export GOOS=linux
export GOARCH=arm
export GOARM=5
go build -o remote
