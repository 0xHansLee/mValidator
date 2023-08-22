package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli"

	txGenerator "github.com/kroma-network/kroma-malicious-validator/tx-generator"
	kflags "github.com/kroma-network/kroma/components/validator/flags"
	klog "github.com/kroma-network/kroma/utils/service/log"
)

var (
	Version   = ""
	GitCommit = ""
	GitDate   = ""
)

func main() {
	klog.SetupDefaults()

	app := cli.NewApp()
	txGeneratorFlags := []cli.Flag{
		txGenerator.DummyTransactionTypeFlag,
		txGenerator.DummyTransactionAccPrivateKeyFlag,
		txGenerator.DummyTransactionSendIntervalFlag,
		txGenerator.ChainIDFlag,
	}
	flags := append(kflags.Flags, txGeneratorFlags...)
	app.Flags = flags
	app.Version = fmt.Sprintf("%s-%s-%s", Version, GitCommit, GitDate)
	app.Name = "kroma-tx-generator"
	app.Usage = "Dummy tx generator"

	app.Action = curryMain()
	err := app.Run(os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}

// curryMain transforms the txGenerator.Main function into an app.Action
// This is done to capture the Version of the tx generator.
func curryMain() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		return txGenerator.Main(ctx)
	}
}
