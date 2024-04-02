package validator

import (
	"context"
	"errors"

	opsources "github.com/ethereum-optimism/optimism/op-service/sources"
	optxmgr "github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/log"
	"github.com/kroma-network/kroma/kroma-validator"
	"github.com/kroma-network/kroma/kroma-validator/challenge"
	"github.com/kroma-network/kroma/kroma-validator/metrics"
	opservice "github.com/kroma-network/kroma/op-service"
	"github.com/kroma-network/kroma/op-service/dial"
)

func NewMaliciousValidatorConfig(cfg validator.CLIConfig, l log.Logger, m *metrics.Metrics, maliciousBlockNumber uint64, submissionInterval uint64) (*validator.Config, error) {
	l2ooAddress, err := opservice.ParseAddress(cfg.L2OOAddress)
	if err != nil {
		return nil, err
	}

	colosseumAddress, err := opservice.ParseAddress(cfg.ColosseumAddress)
	if err != nil {
		return nil, err
	}

	securityCouncilAddress, err := opservice.ParseAddress(cfg.SecurityCouncilAddress)
	if err != nil {
		return nil, err
	}

	valPoolAddress, err := opservice.ParseAddress(cfg.ValPoolAddress)
	if err != nil {
		return nil, err
	}

	//opCfg := optxmgr.NewCLIConfig(cfg.L1EthRpc, optxmgr.DefaultBatcherFlagValues)
	opCfg := cfg.TxMgrConfig
	simpleTxManager, err := optxmgr.NewSimpleTxManager("validator", l, m, opCfg)
	if err != nil {
		return nil, err
	}
	bfTxManager := &optxmgr.BufferedTxManager{
		SimpleTxManager: *simpleTxManager,
	}

	if !cfg.OutputSubmitterEnabled && !cfg.ChallengerEnabled {
		return nil, errors.New("output submitter and challenger are disabled. either output submitter or challenger must be enabled")
	}

	if cfg.ChallengerEnabled && len(cfg.ProverRPC) == 0 {
		return nil, errors.New("ProverRPC is required when challenger enabled, but given empty")
	}

	var fetcher validator.ProofFetcher
	if cfg.ChallengerEnabled && len(cfg.ProverRPC) > 0 {
		fetcher, err = challenge.NewFetcher(cfg.ProverRPC, cfg.FetchingProofTimeout, l)
		if err != nil {
			return nil, err
		}
	}

	// Connect to L1 and L2 providers. Perform these last since they are the most expensive.
	ctx := context.Background()
	l1Client, err := dial.DialEthClientWithTimeout(ctx, dial.DefaultDialTimeout, l, cfg.L1EthRpc)
	if err != nil {
		return nil, err
	}

	l2Client, err := dial.DialEthClientWithTimeout(ctx, dial.DefaultDialTimeout, l, cfg.L2EthRpc)
	if err != nil {
		return nil, err
	}

	maliciousRollupRPC, err := NewMaliciousRollupRPC(cfg.RollupRpc)
	if err != nil {
		return nil, err
	}
	maliciousRollupRPC.SetCustomFlags(maliciousBlockNumber, submissionInterval)
	rollupClient := opsources.NewRollupClient(maliciousRollupRPC)

	rollupConfig, err := rollupClient.RollupConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &validator.Config{
		L2OutputOracleAddr:              l2ooAddress,
		ColosseumAddr:                   colosseumAddress,
		SecurityCouncilAddr:             securityCouncilAddress,
		ValidatorPoolAddr:               valPoolAddress,
		ChallengerPollInterval:          cfg.ChallengerPollInterval,
		NetworkTimeout:                  cfg.TxMgrConfig.NetworkTimeout,
		TxManager:                       bfTxManager,
		L1Client:                        l1Client,
		L2Client:                        l2Client,
		RollupClient:                    rollupClient,
		RollupConfig:                    rollupConfig,
		AllowNonFinalized:               cfg.AllowNonFinalized,
		OutputSubmitterEnabled:          cfg.OutputSubmitterEnabled,
		OutputSubmitterRetryInterval:    cfg.OutputSubmitterRetryInterval,
		OutputSubmitterRoundBuffer:      cfg.OutputSubmitterRoundBuffer,
		ChallengerEnabled:               cfg.ChallengerEnabled,
		GuardianEnabled:                 cfg.GuardianEnabled,
		ProofFetcher:                    fetcher,
		OutputSubmitterAllowPublicRound: cfg.OutputSubmitterAllowPublicRound,
	}, nil
}
