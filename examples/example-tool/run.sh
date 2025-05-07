#!/usr/bin/env bash

set -e

rm -f nats-mcp
cp ../../nats-mcp .
./nats-mcp tool -c "uv" run app.py
