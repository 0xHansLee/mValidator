#!/bin/sh

DUMMY_TX_ACC_PRIV_KEY="0x3bc84101cccc2988dbe3ca7a160bd130e86af12271e25c73506c323c87ceb263" # Ed
# DUMMY_TX_ACC_PRIV_KEY="0xff995c6d5b2851750146977f197369fc641fcbb58c8f084646f49572c7a79879" # Ellie

DUMMY_TX_ACC_PRIV_KEY="${DUMMY_TX_ACC_PRIV_KEY}" docker compose -f tx-generator/docker-compose.yaml -p tx-generator up -d
