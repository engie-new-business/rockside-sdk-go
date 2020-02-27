package rockside

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	_ bind.ContractTransactor = (*Transactor)(nil)
)

func TransactOpts() *bind.TransactOpts {
	return &bind.TransactOpts{Signer: noopSigner()}
}

func noopSigner() bind.SignerFn {
	return func(s types.Signer, c common.Address, t *types.Transaction) (*types.Transaction, error) {
		return t, nil
	}
}

type Transactor struct {
	client           *Client
	rocksideIdentity common.Address
	mu               sync.RWMutex
	transactions     map[common.Hash]string
}

func NewTransactor(rocksideIdentity common.Address, client *Client) *Transactor {
	return &Transactor{
		client:           client,
		rocksideIdentity: rocksideIdentity,
		transactions:     make(map[common.Hash]string),
	}
}

func (t *Transactor) ReturnRocksideTransactionHash(hash common.Hash) string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if txhash, ok := t.transactions[hash]; ok {
		delete(t.transactions, hash)
		return txhash
	}
	return ""
}

func (t *Transactor) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return t.client.RPCClient.PendingCodeAt(ctx, account)
}
func (t *Transactor) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return 0, nil // Rockside manage the nonce
}
func (t *Transactor) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return t.client.RPCClient.SuggestGasPrice(ctx)
}

func (t *Transactor) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	return 0, nil // Rockside manage the gas
}

func (t *Transactor) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	resp, err := t.client.Transaction.Send(Transaction{
		From:     t.rocksideIdentity.String(),
		To:       tx.To().String(),
		Value:    hexutil.EncodeBig(tx.Value()),
		Data:     hexutil.Encode(tx.Data()),
		Gas:      hexutil.EncodeUint64(tx.Gas()),
		GasPrice: hexutil.EncodeBig(tx.GasPrice()),
	})
	if err == nil {
		t.mu.Lock()
		defer t.mu.Unlock()
		t.transactions[tx.Hash()] = resp.TransactionHash
	}
	return err
}
