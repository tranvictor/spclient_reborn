package ethminer

import (
	"../"
	"../../protocol"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

var SmartPool *protocol.SmartPool

type SmartPoolService struct{}

func (SmartPoolService) GetWork() ([3]string, error) {
	var res [3]string
	w := SmartPool.GetWork().(*ethereum.Work)
	res[0] = w.PoWHash().Hex()
	res[1] = w.SeedHash()
	n := big.NewInt(1)
	n.Lsh(n, 255)
	n.Div(n, w.ShareDifficulty())
	n.Lsh(n, 1)
	res[2] = common.BytesToHash(n.Bytes()).Hex()
	return res, nil
}

func (SmartPoolService) SubmitHashrate(hashrate hexutil.Uint64, id common.Hash) bool {
	nc := SmartPool.NetworkClient.(*ethereum.NetworkClient)
	return nc.SubmitHashrate(hashrate, id)
}

func (SmartPoolService) SubmitWork(nonce types.BlockNonce, hash, mixDigest common.Hash) bool {
	// Because it's time critical when miner found a full block so just broadcast
	// everything miner submitted
	nc := SmartPool.NetworkClient.(*ethereum.NetworkClient)
	sol := &ethereum.Solution{
		Nonce:     nonce,
		Hash:      hash,
		MixDigest: mixDigest,
	}
	nc.SubmitSolution(sol)
	return SmartPool.AcceptSolution(sol)
}
