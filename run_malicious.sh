#!/bin/sh

PRIVATE_KEY=""

PRIVATE_KEY="$PRIVATE_KEY" docker compose -f validators/malicious/docker-compose.yaml -p malicious up -d
