package validator

import (
	"context"
	"math/rand"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/testutils"
)

type MaliciousRollupRPC struct {
	rpc                      client.RPC
	targetBlockNumber        *hexutil.Uint64
	outputSubmissionInterval *hexutil.Uint64
}

func NewMaliciousRollupRPC(url string) (*MaliciousRollupRPC, error) {
	rpcCl, err := rpc.DialContext(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return &MaliciousRollupRPC{
		rpc: client.NewBaseRPCClient(rpcCl),
	}, nil
}

func (r *MaliciousRollupRPC) SetCustomFlags(blockNum uint64, submissionInterval uint64) {
	r.targetBlockNumber = new(hexutil.Uint64)
	*r.targetBlockNumber = hexutil.Uint64(blockNum)
	r.outputSubmissionInterval = new(hexutil.Uint64)
	*r.outputSubmissionInterval = hexutil.Uint64(submissionInterval)
}

func (r *MaliciousRollupRPC) Close() {
	r.rpc.Close()
}

func (r *MaliciousRollupRPC) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	if method == "optimism_outputAtBlock" || method == "kroma_outputWithProofAtBlock" {
		blockNumber := args[0].(hexutil.Uint64)

		err := r.rpc.CallContext(ctx, &result, method, blockNumber)
		if err != nil {
			return err
		}
		if r.isMaliciousBlock(blockNumber) {
			rng := rand.New(rand.NewSource(int64(blockNumber)))

			s := result.(**eth.OutputResponse)
			(*s).OutputRoot = eth.Bytes32(testutils.RandomHash(rng))
			(*s).WithdrawalStorageRoot = testutils.RandomHash(rng)
			(*s).StateRoot = testutils.RandomHash(rng)

			return nil
		}
	}

	return r.rpc.CallContext(ctx, result, method, args...)
}

func (r *MaliciousRollupRPC) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return r.rpc.BatchCallContext(ctx, b)
}

func (r *MaliciousRollupRPC) EthSubscribe(ctx context.Context, channel interface{}, args ...interface{}) (ethereum.Subscription, error) {
	return r.rpc.EthSubscribe(ctx, channel, args...)
}

func (r *MaliciousRollupRPC) isMaliciousBlock(blockNumber hexutil.Uint64) bool {
	return r.targetBlockNumber != nil && r.outputSubmissionInterval != nil && blockNumber >= *r.targetBlockNumber && blockNumber < *r.targetBlockNumber+*r.outputSubmissionInterval
}
