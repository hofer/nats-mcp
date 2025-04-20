#!/usr/bin/env bash

set -e

rm nats-mcp
cp ../../nats-mcp .
./nats-mcp tool -c "uv" run app.py