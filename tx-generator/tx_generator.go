package tx_generator

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	"github.com/kroma-network/kroma/components/validator"
	"github.com/kroma-network/kroma/components/validator/metrics"
	"github.com/kroma-network/kroma/utils"
	kcrypto "github.com/kroma-network/kroma/utils/service/crypto"
	klog "github.com/kroma-network/kroma/utils/service/log"
	"github.com/kroma-network/kroma/utils/service/txmgr"
)

func Main(cliCtx *cli.Context) error {
	cliCfg := validator.NewCLIConfig(cliCtx)
	if err := cliCfg.Check(); err != nil {
		return fmt.Errorf("invalid CLI flags: %w", err)
	}

	dummyTxType := cliCtx.Uint64(DummyTransactionTypeFlag.Name)
	dummyTxAccPrivKey := cliCtx.String(DummyTransactionAccPrivateKeyFlag.Name)
	dummyTxSendInterval := cliCtx.Duration(DummyTransactionSendIntervalFlag.Name)
	chainID := cliCtx.Uint64(ChainIDFlag.Name)

	l := klog.NewLogger(cliCfg.LogConfig)
	l.Info("initializing TxGenerator")

	txGeneratorCfg, err := validator.NewValidatorConfig(cliCfg, l, metrics.NoopMetrics)
	if err != nil {
		l.Error("Unable to create tx generator config", "err", err)
		return err
	}

	generator, err := NewTxGenerator(l, *txGeneratorCfg, dummyTxType, dummyTxAccPrivKey, dummyTxSendInterval, chainID)
	if err != nil {
		l.Error("failed to create tx generator")
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	generator.Start(ctx)
	<-utils.WaitInterrupt()
	generator.Stop()

	return nil
}

type TxGenerator struct {
	ctx        context.Context
	cancel     context.CancelFunc
	txType     uint64
	txInterval time.Duration
	l          log.Logger
	Backend    txmgr.ETHBackend
	chainID    *big.Int
	from       common.Address
	signer     kcrypto.SignerFn
}

func NewTxGenerator(l log.Logger, cfg validator.Config, dummyTxType uint64, dummyTxAccPrivKey string, dummyTxSendInterval time.Duration, chainID uint64) (*TxGenerator, error) {
	privKeyBz, err := crypto.HexToECDSA(strings.TrimPrefix(dummyTxAccPrivKey, "0x"))
	if err != nil {
		return nil, err
	}

	from := crypto.PubkeyToAddress(privKeyBz.PublicKey)
	signer := func(chainID *big.Int) kcrypto.SignerFn {
		s := kcrypto.PrivateKeySignerFn(privKeyBz, chainID)
		return func(_ context.Context, addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return s(addr, tx)
		}
	}

	intChainID := new(big.Int).SetUint64(chainID)

	return &TxGenerator{
		txType:     dummyTxType,
		txInterval: dummyTxSendInterval,
		l:          l,
		Backend:    cfg.L2Client,
		chainID:    intChainID,
		from:       from,
		signer:     signer(intChainID),
	}, nil
}

func (g TxGenerator) Start(ctx context.Context) {
	g.ctx, g.cancel = context.WithCancel(ctx)

	ticker := time.NewTicker(g.txInterval)
	defer ticker.Stop()

	g.l.Info("tx generator started", "address", g.from, "txType", g.txType)

	txType := [3]uint64{0, 1, 2}
	num := 0
	for ; ; <-ticker.C {
		select {
		case <-g.ctx.Done():
			g.l.Info("stopping tx generator")
			return
		default:
			if err := g.generateTx(txType[num%3]); err != nil {
				g.l.Error("failed to generate dummy tx", "err", err)
			}
			num += 1
		}
	}
}

func (g TxGenerator) Stop() {
	g.cancel()
}

func (g TxGenerator) generateTx(txType uint64) error {
	g.l.Info("generating dummy tx...", "txType", txType)

	tx, err := g.getTx(txType)
	if err != nil {
		return fmt.Errorf("failed to get dummy tx: %w", err)
	}

	signedTx, err := g.signer(g.ctx, g.from, tx)
	if err != nil {
		return fmt.Errorf("failed to sign tx: %w", err)
	}

	if err = g.Backend.SendTransaction(g.ctx, signedTx); err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}

	return nil
}

func (g TxGenerator) getTx(txType uint64) (*types.Transaction, error) {
	zeroAddr := common.Address{0}
	nonce, err := g.Backend.PendingNonceAt(g.ctx, g.from)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	accessList := types.AccessList{
		types.AccessTuple{
			Address:     zeroAddr,
			StorageKeys: []common.Hash{common.BigToHash(common.Big0)},
		},
	}

	var txData types.TxData
	switch txType {
	case 0:
		txData = &types.LegacyTx{
			Nonce:    nonce,
			GasPrice: hexutil.MustDecodeBig("0x600f1f"),
			Gas:      21000,
			To:       &zeroAddr,
			Value:    common.Big1,
			Data:     hexutil.MustDecode("0x"),
		}
	case 1:
		txData = &types.AccessListTx{
			ChainID:    g.chainID,
			Nonce:      nonce,
			GasPrice:   hexutil.MustDecodeBig("0x600f1f"),
			Gas:        30000,
			To:         &zeroAddr,
			Value:      common.Big1,
			Data:       hexutil.MustDecode("0x"),
			AccessList: accessList,
		}
	case 2:
		txData = &types.DynamicFeeTx{
			ChainID:    g.chainID,
			Nonce:      nonce,
			Gas:        30000,
			GasFeeCap:  new(big.Int).SetInt64(1000000005),
			GasTipCap:  new(big.Int).SetInt64(1000000005),
			To:         &zeroAddr,
			Value:      common.Big1,
			Data:       hexutil.MustDecode("0x"),
			AccessList: accessList,
		}
	}

	return types.NewTx(txData), nil
}
