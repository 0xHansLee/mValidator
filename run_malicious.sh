#!/bin/sh

HD_PATH="m/44'/60'/1'/0/8"
#HD_PATH="m/44'/60'/1'/0/10"

HD_PATH="$HD_PATH" docker compose -f validators/malicious/docker-compose.yaml -p malicious up -d
