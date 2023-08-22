#!/bin/sh

HD_PATH="m/44'/60'/1'/0/11"
#HD_PATH="m/44'/60'/1'/0/13"

HD_PATH="$HD_PATH" docker compose -f validators/honest/docker-compose.yaml -p honest-2 up -d
