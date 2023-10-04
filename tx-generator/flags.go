package tx_generator

import (
	"time"

	"github.com/urfave/cli/v2"
)

var (
	DummyTransactionTypeFlag = &cli.Uint64Flag{
		Name:     "dummy-tx-type",
		Usage:    "generate dummy transaction of given type repeatedly",
		EnvVars:  []string{"DUMMY_TX_TYPE"},
		Required: true,
	}
	DummyTransactionAccPrivateKeyFlag = &cli.StringFlag{
		Name:     "dummy-tx-account",
		Usage:    "EOA to generate dummy transaction",
		EnvVars:  []string{"DUMMY_TX_ACC_PRIV_KEY"},
		Required: true,
	}
	DummyTransactionSendIntervalFlag = &cli.DurationFlag{
		Name:     "dummy-tx-send-interval",
		Usage:    "Interval second to send dummy transaction",
		EnvVars:  []string{"DUMMY_TX_SEND_INTERVAL"},
		Required: false,
		Value:    1 * time.Second,
	}
	ChainIDFlag = &cli.Uint64Flag{
		Name:     "chain-id",
		Usage:    "chain id to send transaction",
		EnvVars:  []string{"CHAIN_ID"},
		Required: true,
	}
)
