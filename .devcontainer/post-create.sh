#!/usr/bin/env sh

go install github.com/spf13/cobra-cli@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

sudo apt update

sudo apt install --yes flac
