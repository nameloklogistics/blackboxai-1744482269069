package stellar

import (
	"fmt"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

// AccountManager handles Stellar account operations
type AccountManager struct {
	client       *horizonclient.Client
	networkPassphrase string
}

// NewAccountManager creates a new AccountManager instance
func NewAccountManager(isTestnet bool) *AccountManager {
	var client *horizonclient.Client
	var networkPassphrase string

	if isTestnet {
		client = horizonclient.DefaultTestNetClient
		networkPassphrase = network.TestNetworkPassphrase
	} else {
		client = horizonclient.DefaultPublicNetClient
		networkPassphrase = network.PublicNetworkPassphrase
	}

	return &AccountManager{
		client:           client,
		networkPassphrase: networkPassphrase,
	}
}

// CreateAccount generates a new Stellar account
func (am *AccountManager) CreateAccount() (*keypair.Full, error) {
	// Generate new keypair
	kp, err := keypair.Random()
	if err != nil {
		return nil, fmt.Errorf("failed to generate keypair: %w", err)
	}

	// For testnet, we can fund the account using friendbot
	if am.networkPassphrase == network.TestNetworkPassphrase {
		resp, err := http.Get("https://friendbot.stellar.org/?addr=" + kp.Address())
		if err != nil {
			return nil, fmt.Errorf("failed to fund account: %w", err)
		}
		defer resp.Body.Close()
	}

	return kp, nil
}

// GetAccountDetails retrieves account information from the network
func (am *AccountManager) GetAccountDetails(address string) (*horizonclient.Account, error) {
	account, err := am.client.AccountDetail(horizonclient.AccountRequest{
		AccountID: address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get account details: %w", err)
	}

	return &account, nil
}

// BuildTransaction creates a new transaction with the given operations
func (am *AccountManager) BuildTransaction(sourceAccount string, operations ...txnbuild.Operation) (*txnbuild.Transaction, error) {
	account, err := am.GetAccountDetails(sourceAccount)
	if err != nil {
		return nil, err
	}

	params := txnbuild.TransactionParams{
		SourceAccount:        &txnbuild.SimpleAccount{AccountID: sourceAccount, Sequence: account.Sequence},
		Operations:          operations,
		BaseFee:            txnbuild.MinBaseFee,
		Timebounds:         txnbuild.NewTimeout(300),
		NetworkPassphrase:  am.networkPassphrase,
	}

	tx, err := txnbuild.NewTransaction(params)
	if err != nil {
		return nil, fmt.Errorf("failed to build transaction: %w", err)
	}

	return tx, nil
}

// SignTransaction signs a transaction with the given secret key
func (am *AccountManager) SignTransaction(tx *txnbuild.Transaction, secretKey string) (*txnbuild.Transaction, error) {
	kp, err := keypair.Parse(secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse secret key: %w", err)
	}

	signedTx, err := tx.Sign(am.networkPassphrase, kp.(*keypair.Full))
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return signedTx, nil
}

// SubmitTransaction submits a signed transaction to the network
func (am *AccountManager) SubmitTransaction(tx *txnbuild.Transaction) (*horizonclient.Transaction, error) {
	result, err := am.client.SubmitTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to submit transaction: %w", err)
	}

	return &result, nil
}
