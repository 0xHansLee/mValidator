#!/bin/sh

HD_PATH="m/44'/60'/1'/0/12"
#HD_PATH="m/44'/60'/1'/0/14"

HD_PATH="$HD_PATH" docker compose -f validators/malicious/docker-compose.yaml -p malicious-2 up -d
