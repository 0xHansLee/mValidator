#!/bin/sh

PRIVATE_KEY=""

PRIVATE_KEY="$PRIVATE_KEY" docker compose -f validators/honest/docker-compose.yaml -p honest up -d
