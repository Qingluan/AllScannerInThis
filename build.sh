#!/bin/bash
rm -rf asset
go get  github.com/jteeuwen/go-bindata/...
go-bindata -o=./asset/asset.go -pkg=asset Res/...
go build -ldflags="-s -w" -trimpath