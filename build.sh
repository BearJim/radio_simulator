#!/usr/bin/env bash

echo "Start build Radio Simulator...."
go build -o bin/simulator -x src/simulator.go
