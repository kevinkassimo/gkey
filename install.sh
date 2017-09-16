#!/bin/bash

go build gokey.go

if [ -e "$GOPATH" ]
then
    rm -f "$GOPATH/bin/gokey"
    mv gokey "$GOPATH/bin"
fi