package validator

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/log"
	"github.com/kroma-network/kroma/components/node/sources"
	"github.com/kroma-network/kroma/components/validator"
	"github.com/kroma-network/kroma/components/validator/metrics"
	"github.com/kroma-network/kroma/utils"
	"github.com/kroma-network/kroma/utils/service/txmgr"
)

func NewMaliciousValidatorConfig(cfg validator.CLIConfig, l log.Logger, m *metrics.Metrics, maliciousBlockNumber uint64) (*validator.Config, error) {
	l2ooAddress, err := utils.ParseAddress(cfg.L2OOAddress)
	if err != nil {
		return nil, err
	}

	colosseumAddress, err := utils.ParseAddress(cfg.ColosseumAddress)
	if err != nil {
		return nil, err
	}

	securityCouncilAddress, err := utils.ParseAddress(cfg.SecurityCouncilAddress)
	if err != nil {
		return nil, err
	}

	valPoolAddress, err := utils.ParseAddress(cfg.ValPoolAddress)
	if err != nil {
		return nil, err
	}

	txManager, err := txmgr.NewSimpleTxManager("validator", l, m, cfg.TxMgrConfig)
	if err != nil {
		return nil, err
	}

	if cfg.OutputSubmitterDisabled && cfg.ChallengerDisabled {
		return nil, errors.New("output submitter and challenger are disabled. either output submitter or challenger must be enabled")
	}

	if !cfg.ChallengerDisabled && len(cfg.ProverGrpc) == 0 {
		return nil, errors.New("ProverGrpc is required but given empty")
	}

	// mock fetcher
	fetcher := NewFetcher(l, "/app/validator/proof")

	// Connect to L1 and L2 providers. Perform these last since they are the most expensive.
	ctx := context.Background()
	l1Client, err := utils.DialEthClientWithTimeout(ctx, cfg.L1EthRpc)
	if err != nil {
		return nil, err
	}

	maliciousRollupRPC, err := NewMaliciousRollupRPC(cfg.RollupRpc)
	if err != nil {
		return nil, err
	}
	maliciousRollupRPC.SetTargetBlockNum(maliciousBlockNumber)
	rollupClient := sources.NewRollupClient(maliciousRollupRPC)

	rollupConfig, err := rollupClient.RollupConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &validator.Config{
		L2OutputOracleAddr:           l2ooAddress,
		ColosseumAddr:                colosseumAddress,
		SecurityCouncilAddr:          securityCouncilAddress,
		ValidatorPoolAddr:            valPoolAddress,
		ChallengerPollInterval:       cfg.ChallengerPollInterval,
		NetworkTimeout:               cfg.TxMgrConfig.NetworkTimeout,
		TxManager:                    txManager,
		L1Client:                     l1Client,
		RollupClient:                 rollupClient,
		RollupConfig:                 rollupConfig,
		AllowNonFinalized:            cfg.AllowNonFinalized,
		OutputSubmitterDisabled:      cfg.OutputSubmitterDisabled,
		OutputSubmitterBondAmount:    cfg.OutputSubmitterBondAmount,
		OutputSubmitterRetryInterval: cfg.OutputSubmitterRetryInterval,
		OutputSubmitterRoundBuffer:   cfg.OutputSubmitterRoundBuffer,
		ChallengerDisabled:           cfg.ChallengerDisabled,
		GuardianEnabled:              cfg.GuardianEnabled,
		ProofFetcher:                 fetcher,
	}, nil
}
