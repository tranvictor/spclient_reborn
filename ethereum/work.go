// Package ethereum contains all necessary components to plug into smartpool
// to work with ethereum blockchain. Such as: Contract, Network Client,
// Share receiver...
// This package also provides interfaces for different ethereum clients to
// be able to work with smartpool.
package ethereum

import (
	"../"
	"./ethash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

// Work represents Ethereum pow work
type Work struct {
	blockHeader     *types.Header
	powHash         string
	seedHash        string
	shareDifficulty *big.Int
}

func (w *Work) ID() string {
	return w.powHash
}

func (w *Work) AcceptSolution(sol smartpool.Solution) smartpool.Share {
	solution := sol.(*Solution)
	s := &Share{
		blockHeader:     w.blockHeader,
		nonce:           solution.Nonce,
		mixDigest:       solution.MixDigest,
		shareDifficulty: w.ShareDifficulty(),
	}
	s.SolutionState = ethash.Instance.SolutionState(s, w.ShareDifficulty())
	return s
}

func (w *Work) PoWHash() common.Hash {
	return common.HexToHash(w.powHash)
}

func (w Work) SeedHash() string {
	return w.seedHash
}

func (w Work) ShareDifficulty() *big.Int {
	return w.shareDifficulty
}

func (w Work) BlockHeader() *types.Header {
	return w.blockHeader
}

func NewWork(h *types.Header, ph string, sh string, diff *big.Int) *Work {
	return &Work{h, ph, sh, diff}
}
