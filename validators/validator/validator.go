package validator

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	oprpc "github.com/ethereum-optimism/optimism/op-service/rpc"
	"github.com/kroma-network/kroma/kroma-validator"
	"github.com/kroma-network/kroma/kroma-validator/metrics"
	klog "github.com/kroma-network/kroma/op-service/log"
	"github.com/kroma-network/kroma/op-service/monitoring"
	"github.com/kroma-network/kroma/op-service/opio"
)

func Main(version string, cliCtx *cli.Context) error {
	cliCfg := validator.NewConfig(cliCtx)
	if err := cliCfg.Check(); err != nil {
		return fmt.Errorf("invalid CLI flags: %w", err)
	}

	// target malicious block number
	maliciousBlockNumber := cliCtx.Uint64(MaliciousBlockNumberFlag.Name)
	outputSubmissionInterval := cliCtx.Uint64(OutputSubmissionIntervalFlag.Name)

	l := klog.NewLogger(klog.AppOut(cliCtx), klog.DefaultCLIConfig())
	m := metrics.NewMetrics("default")
	l.Info("initializing Validator")

	validatorCfg, err := NewMaliciousValidatorConfig(cliCfg, l, m, maliciousBlockNumber, outputSubmissionInterval)
	if err != nil {
		l.Error("Unable to create validator config", "err", err)
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	monitoring.MaybeStartPprof(ctx, cliCfg.PprofConfig, l)
	monitoring.MaybeStartMetrics(ctx, cliCfg.MetricsConfig, l, m, validatorCfg.L1Client, validatorCfg.TxManager.From())
	server, err := monitoring.StartRPC(cliCfg.RPCConfig, version, oprpc.WithLogger(l))
	if err != nil {
		return err
	}
	defer func() {
		if err = server.Stop(); err != nil {
			l.Error("Error shutting down http server: %w", err)
		}
	}()

	m.RecordInfo(version)
	m.RecordUp()

	validator, err := validator.NewValidator(*validatorCfg, l, m)
	if err != nil {
		return err
	}

	if err := validator.Start(); err != nil {
		l.Error("failed to start validator", "err", err)
		return err
	}
	opio.BlockOnInterrupts()
	if err := validator.Stop(); err != nil {
		l.Error("failed to stop validator", "err", err)
		return err
	}

	return nil
}
