package geth

import (
	"../"
	"../../"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
)

type GethContractClient struct {
	// the contract implementation that holds all underlying
	// communication with Ethereum Contract
	pool       *TestPool
	transactor *bind.TransactOpts
	node       ethereum.RPCClient
	sender     common.Address
}

func (cc *GethContractClient) Version() string {
	v, err := cc.pool.Version(nil)
	if err != nil {
		return ""
	}
	return v
}

func (cc *GethContractClient) IsRegistered() bool {
	ok, err := cc.pool.IsRegistered(nil, cc.sender)
	if err != nil {
		return false
	}
	return ok
}

func (cc *GethContractClient) CanRegister() bool {
	ok, err := cc.pool.CanRegister(nil, cc.sender)
	if err != nil {
		return false
	}
	return ok
}

func (cc *GethContractClient) Register(paymentAddress common.Address) error {
	tx, err := cc.pool.Register(cc.transactor, paymentAddress)
	smartpool.Output.Printf("Registering address %s to SmartPool contract by tx: %s\n", paymentAddress.Hex(), tx.Hash())
	if err != nil {
		smartpool.Output.Printf("Registering address %s failed. Error: %s\n", err)
		return err
	}
	NewTxWatcher(tx, cc.node).Wait()
	smartpool.Output.Printf("Registered address %s to SmartPool contract. Tx %s is confirmed\n", paymentAddress.Hex(), tx.Hash())
	return nil
}

func (cc *GethContractClient) GetClaimSeed() *big.Int {
	seed, err := cc.pool.GetClaimSeed(nil, cc.sender)
	if err != nil {
		smartpool.Output.Printf("Getting claim seed failed. Error: %s\n", err)
		return big.NewInt(0)
	}
	return seed
}

func (cc *GethContractClient) SubmitClaim(
	numShares *big.Int,
	difficulty *big.Int,
	min *big.Int,
	max *big.Int,
	augMerkle *big.Int) error {
	tx, err := cc.pool.SubmitClaim(cc.transactor,
		numShares, difficulty, min, max, augMerkle)
	if err != nil {
		smartpool.Output.Printf("Submitting claim failed. Error: %s\n", err)
		return err
	}
	NewTxWatcher(tx, cc.node).Wait()
	return nil
}

func (cc *GethContractClient) VerifyClaim(
	rlpHeader []byte,
	nonce *big.Int,
	shareIndex *big.Int,
	dataSetLookup []*big.Int,
	witnessForLookup []*big.Int,
	augCountersBranch []*big.Int,
	augHashesBranch []*big.Int) error {
	tx, err := cc.pool.VerifyClaim(cc.transactor,
		rlpHeader, nonce, shareIndex, dataSetLookup,
		witnessForLookup, augCountersBranch, augHashesBranch)
	if err != nil {
		smartpool.Output.Printf("Verifying claim failed. Error: %s\n", err)
		return err
	}
	NewTxWatcher(tx, cc.node).Wait()
	return nil
}

func (cc *GethContractClient) SetEpochData(merkleRoot []*big.Int, fullSizeIn128Resolution []uint64, branchDepth []uint64, epoch []*big.Int) error {
	tx, err := cc.pool.SetEpochData(cc.transactor,
		merkleRoot, fullSizeIn128Resolution, branchDepth, epoch)
	if err != nil {
		smartpool.Output.Printf("Setting epoch data. Error: %s\n", err)
		return err
	}
	NewTxWatcher(tx, cc.node).Wait()
	return nil
}

func getClient(rpc string) (*ethclient.Client, error) {
	return ethclient.Dial(rpc)
}

func NewGethContractClient(
	contractAddr common.Address, node ethereum.RPCClient, miner common.Address,
	ipc, keystorePath, passphrase string) (*GethContractClient, error) {
	client, err := getClient(ipc)
	if err != nil {
		smartpool.Output.Printf("Couldn't connect to Geth via IPC file. Error: %s\n", err)
		return nil, err
	}
	pool, err := NewTestPool(contractAddr, client)
	if err != nil {
		smartpool.Output.Printf("Couldn't get SmartPool information from Ethereum Blockchain. Error: %s\n", err)
		return nil, err
	}
	account := GetAccount(keystorePath, miner, passphrase)
	if account == nil {
		smartpool.Output.Printf("Couldn't get any account from key store.\n")
		return nil, err
	}
	keyio, err := os.Open(account.KeyFile())
	if err != nil {
		smartpool.Output.Printf("Failed to open key file: %s\n", err)
		return nil, err
	}
	smartpool.Output.Printf("Unlocking account...")
	auth, err := bind.NewTransactor(keyio, account.PassPhrase())
	if err != nil {
		smartpool.Output.Printf("Failed to create authorized transactor: %s\n", err)
		return nil, err
	}
	// TODO: make gas price one command line flag
	auth.GasPrice = big.NewInt(10000000000)
	smartpool.Output.Printf("Done.\n")
	return &GethContractClient{pool, auth, node, miner}, nil
}
