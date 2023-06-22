package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/kroma-network/kroma-malicious-validator/validator"
	"github.com/kroma-network/kroma/components/validator/cmd/balance"
	kflags "github.com/kroma-network/kroma/components/validator/flags"
	klog "github.com/kroma-network/kroma/utils/service/log"
	"github.com/urfave/cli"
)

var (
	Version   = ""
	GitCommit = ""
	GitDate   = ""
)

func main() {
	klog.SetupDefaults()

	flags := append(kflags.Flags, validator.MaliciousBlockNumberFlag)
	app := cli.NewApp()
	app.Flags = flags
	app.Version = fmt.Sprintf("%s-%s-%s", Version, GitCommit, GitDate)
	app.Name = "kroma-malicious-validator"
	app.Usage = "Malicious validator (L2 output submitter and/or challenger)"
	app.Commands = []cli.Command{
		{
			Name:  "deposit",
			Usage: "Deposit ETH into ValidatorPool to be used as bond",
			Flags: []cli.Flag{
				cli.Uint64Flag{
					Name:     "amount",
					Usage:    "Amount to deposit into ValidatorPool (in wei)",
					Required: true,
				},
			},
			Action: balance.Deposit,
		},
		{
			Name:  "withdraw",
			Usage: "Withdraw ETH from ValidatorPool",
			Flags: []cli.Flag{
				cli.Uint64Flag{
					Name:     "amount",
					Usage:    "Amount to withdraw from ValidatorPool (in wei)",
					Required: true,
				},
			},
			Action: balance.Withdraw,
		},
		{
			Name:   "unbond",
			Usage:  "Attempt to unbond in ValidatorPool",
			Action: balance.Unbond,
		},
	}

	app.Action = curryMain(Version)
	err := app.Run(os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}

// curryMain transforms the maliciousValidator.Main function into an app.Action
// This is done to capture the Version of the validator.
func curryMain(version string) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		return validator.Main(version, ctx)
	}
}
