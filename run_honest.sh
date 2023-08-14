#!/bin/sh

HD_PATH="m/44'/60'/1'/0/7"
#HD_PATH="m/44'/60'/1'/0/9"

HD_PATH="$HD_PATH" docker compose -f honest/docker-compose.yaml -p honest up -d
