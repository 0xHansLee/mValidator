package validator

import (
	"github.com/urfave/cli"
)

var (
	MaliciousBlockNumberFlag = cli.Uint64Flag{
		Name:     "malicious-block-number",
		Usage:    "Target block number of invalid block for malicious validator",
		EnvVar:   "VALIDATOR_MALICIOUS_BLOCK_NUMBER",
		Required: true,
	}
)
