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
	OutputSubmissionIntervalFlag = cli.Uint64Flag{
		Name:     "output-submission-interval",
		Usage:    "Output submission interval in block number",
		EnvVar:   "VALIDATOR_OUTPUT_SUBMISSION_INTERVAL",
		Required: true,
	}
)
