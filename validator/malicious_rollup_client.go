package validator

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/kroma-network/kroma/components/node/client"
)

type MaliciousRollupRPC struct {
	rpc               client.RPC
	targetBlockNumber *hexutil.Uint64
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

func (r *MaliciousRollupRPC) SetTargetBlockNum(blockNum uint64) {
	r.targetBlockNumber = new(hexutil.Uint64)
	*r.targetBlockNumber = hexutil.Uint64(blockNum)
}

func (r *MaliciousRollupRPC) Close() {
	r.rpc.Close()
}

func (r *MaliciousRollupRPC) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	if method == "kroma_outputAtBlock" || method == "kroma_outputWithProofAtBlock" {
		blockNumber := args[0].(hexutil.Uint64)

		if r.targetBlockNumber != nil && blockNumber >= *r.targetBlockNumber {
			err := r.rpc.CallContext(ctx, &result, method, blockNumber-1)
			if err != nil {
				return err
			}
			return nil
		}
		err := r.rpc.CallContext(ctx, &result, method, blockNumber)
		if err != nil {
			return err
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
