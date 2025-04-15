package stellar

import (
	"fmt"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/protocols/horizon"
)

// TokenManager handles token-related operations
type TokenManager struct {
	accountManager *AccountManager
	tokenCode     string
	issuerAccount string
}

// NewTokenManager creates a new TokenManager instance
func NewTokenManager(accountManager *AccountManager, tokenCode, issuerAccount string) *TokenManager {
	return &TokenManager{
		accountManager: accountManager,
		tokenCode:     tokenCode,
		issuerAccount: issuerAccount,
	}
}

// CreateToken issues a new token on the Stellar network
func (tm *TokenManager) CreateToken(distributorAccount string) error {
	// Create trust line operation
	trustLineOp := &txnbuild.ChangeTrust{
		Line: txnbuild.CreditAsset{
			Code:   tm.tokenCode,
			Issuer: tm.issuerAccount,
		},
		Limit: "100000000000", // 100B tokens
	}

	// Build and submit trust line transaction
	tx, err := tm.accountManager.BuildTransaction(distributorAccount, trustLineOp)
	if err != nil {
		return fmt.Errorf("failed to build trust line transaction: %w", err)
	}

	// Payment operation to issue tokens
	paymentOp := &txnbuild.Payment{
		Destination: distributorAccount,
		Asset: txnbuild.CreditAsset{
			Code:   tm.tokenCode,
			Issuer: tm.issuerAccount,
		},
		Amount: "100000000000",
	}

	// Build and submit payment transaction
	tx, err = tm.accountManager.BuildTransaction(tm.issuerAccount, paymentOp)
	if err != nil {
		return fmt.Errorf("failed to build payment transaction: %w", err)
	}

	return nil
}

// TransferTokens transfers tokens between accounts
func (tm *TokenManager) TransferTokens(fromAccount, toAccount string, amount string) error {
	paymentOp := &txnbuild.Payment{
		Destination: toAccount,
		Asset: txnbuild.CreditAsset{
			Code:   tm.tokenCode,
			Issuer: tm.issuerAccount,
		},
		Amount: amount,
	}

	// Build transaction
	tx, err := tm.accountManager.BuildTransaction(fromAccount, paymentOp)
	if err != nil {
		return fmt.Errorf("failed to build transfer transaction: %w", err)
	}

	return nil
}

// GetTokenBalance retrieves the token balance for an account
func (tm *TokenManager) GetTokenBalance(account string) (string, error) {
	accountDetails, err := tm.accountManager.GetAccountDetails(account)
	if err != nil {
		return "", fmt.Errorf("failed to get account details: %w", err)
	}

	for _, balance := range accountDetails.Balances {
		if balance.Asset.Code == tm.tokenCode && balance.Asset.Issuer == tm.issuerAccount {
			return balance.Balance, nil
		}
	}

	return "0", nil
}

// EstablishTrustLine creates a trust line for the token
func (tm *TokenManager) EstablishTrustLine(account string) error {
	trustLineOp := &txnbuild.ChangeTrust{
		Line: txnbuild.CreditAsset{
			Code:   tm.tokenCode,
			Issuer: tm.issuerAccount,
		},
		Limit: "100000000000", // Maximum trust line limit
	}

	// Build and submit transaction
	tx, err := tm.accountManager.BuildTransaction(account, trustLineOp)
	if err != nil {
		return fmt.Errorf("failed to build trust line transaction: %w", err)
	}

	return nil
}

// GetTokenTransactions retrieves token transfer history for an account
func (tm *TokenManager) GetTokenTransactions(account string) ([]horizon.Transaction, error) {
	txRequest := horizonclient.TransactionRequest{
		ForAccount: account,
		Limit:      50,
	}

	txs, err := tm.accountManager.client.Transactions(txRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return txs.Embedded.Records, nil
}

// LockTokens implements token locking for escrow during booking process
func (tm *TokenManager) LockTokens(account string, amount string, duration uint64) error {
	// Create a claimable balance with time-based clawback
	claimant := txnbuild.NewClaimant(account, &txnbuild.UnixTimePredicate{TimePoint: duration})
	
	createClaimableBalance := &txnbuild.CreateClaimableBalance{
		Amount: amount,
		Asset: txnbuild.CreditAsset{
			Code:   tm.tokenCode,
			Issuer: tm.issuerAccount,
		},
		Claimants: []txnbuild.Claimant{claimant},
	}

	// Build and submit transaction
	tx, err := tm.accountManager.BuildTransaction(account, createClaimableBalance)
	if err != nil {
		return fmt.Errorf("failed to build lock transaction: %w", err)
	}

	return nil
}

// UnlockTokens releases locked tokens after service completion
func (tm *TokenManager) UnlockTokens(claimableBalanceID string, account string) error {
	claimBalance := &txnbuild.ClaimClaimableBalance{
		BalanceID: claimableBalanceID,
	}

	// Build and submit transaction
	tx, err := tm.accountManager.BuildTransaction(account, claimBalance)
	if err != nil {
		return fmt.Errorf("failed to build unlock transaction: %w", err)
	}

	return nil
}
