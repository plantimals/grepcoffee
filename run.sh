#!/bin/bash

if [ -a coffee ]
  then
    echo "removing old coffee"
    rm coffee
fi
go build models/models.go
go build coffee.go
./coffee
