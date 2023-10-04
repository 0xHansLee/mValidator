package validator

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/kroma-network/kroma/components/validator"
	"github.com/kroma-network/kroma/components/validator/metrics"
	"github.com/kroma-network/kroma/utils"
	"github.com/kroma-network/kroma/utils/monitoring"
	klog "github.com/kroma-network/kroma/utils/service/log"
	krpc "github.com/kroma-network/kroma/utils/service/rpc"
)

func Main(version string, cliCtx *cli.Context) error {
	cliCfg := validator.NewCLIConfig(cliCtx)
	if err := cliCfg.Check(); err != nil {
		return fmt.Errorf("invalid CLI flags: %w", err)
	}

	// target malicious block number
	maliciousBlockNumber := cliCtx.Uint64(MaliciousBlockNumberFlag.Name)
	outputSubmissionInterval := cliCtx.Uint64(OutputSubmissionIntervalFlag.Name)

	l := klog.NewLogger(cliCfg.LogConfig)
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
	server, err := monitoring.StartRPC(cliCfg.RPCConfig, version, krpc.WithLogger(l))
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

	validator, err := validator.NewValidator(ctx, *validatorCfg, l, m)
	if err != nil {
		return err
	}

	if err := validator.Start(); err != nil {
		l.Error("failed to start validator", "err", err)
		return err
	}
	<-utils.WaitInterrupt()
	if err := validator.Stop(); err != nil {
		l.Error("failed to stop validator", "err", err)
		return err
	}

	return nil
}
