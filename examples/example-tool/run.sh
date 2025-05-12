#!/usr/bin/env bash

set -e

rm -f nats-mcp
cp ../../nats-mcp .
./nats-mcp tool -n "example" -c "uv" run app.py
