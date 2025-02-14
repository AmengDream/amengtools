@echo off
go env -w GOOS=linux
go build -ldflags "-s -w"