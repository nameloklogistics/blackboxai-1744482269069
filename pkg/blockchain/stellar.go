package blockchain

import (
    "fmt"
    "github.com/stellar/go/clients/horizonclient"
    "github.com/stellar/go/keypair"
    "github.com/stellar/go/network"
    "github.com/stellar/go/txnbuild"
    "logistics-marketplace/pkg/config"
)

type StellarClient struct {
    HorizonClient *horizonclient.Client
    NetworkPassphrase string
}

// NewStellarClient creates a new instance of the Stellar blockchain client
func NewStellarClient(cfg *config.Config) (*StellarClient, error) {
    client := horizonclient.NewClient(horizonclient.HistoryURLFlagName(cfg.Network.HorizonURL))
    
    return &StellarClient{
        HorizonClient: client,
        NetworkPassphrase: cfg.Network.NetworkPassphrase,
    }, nil
}

// CreateAccount generates a new Stellar account
func (sc *StellarClient) CreateAccount() (*keypair.Full, error) {
    kp, err := keypair.Random()
    if err != nil {
        return nil, fmt.Errorf("failed to create keypair: %v", err)
    }
    
    return kp, nil
}

// IssueToken creates and issues the logistics service token
func (sc *StellarClient) IssueToken(issuerKP *keypair.Full, cfg *config.Config) error {
    // Create asset
    asset := txnbuild.CreditAsset{
        Code:   cfg.Token.TokenCode,
        Issuer: issuerKP.Address(),
    }

    // Get account details
    account, err := sc.HorizonClient.AccountDetail(horizonclient.AccountRequest{
        AccountID: issuerKP.Address(),
    })
    if err != nil {
        return fmt.Errorf("failed to get account details: %v", err)
    }

    // Create transaction to set trust line
    tx, err := txnbuild.NewTransaction(
        txnbuild.TransactionParams{
            SourceAccount:        &account,
            IncrementSequenceNum: true,
            Operations: []txnbuild.Operation{
                &txnbuild.ChangeTrust{
                    Line:  asset.ToChangeTrustAsset(),
                    Limit: cfg.Token.MaxSupply,
                },
            },
            BaseFee:    txnbuild.MinBaseFee,
            Timebounds: txnbuild.NewTimeout(300),
        },
    )
    if err != nil {
        return fmt.Errorf("failed to build transaction: %v", err)
    }

    // Sign transaction
    tx, err = tx.Sign(sc.NetworkPassphrase, issuerKP)
    if err != nil {
        return fmt.Errorf("failed to sign transaction: %v", err)
    }

    // Submit transaction
    _, err = sc.HorizonClient.SubmitTransaction(tx)
    if err != nil {
        return fmt.Errorf("failed to submit transaction: %v", err)
    }

    return nil
}

// CreateTrustline establishes a trustline for the token
func (sc *StellarClient) CreateTrustline(accountKP *keypair.Full, issuerAddress string, assetCode string) error {
    asset := txnbuild.CreditAsset{
        Code:   assetCode,
        Issuer: issuerAddress,
    }

    account, err := sc.HorizonClient.AccountDetail(horizonclient.AccountRequest{
        AccountID: accountKP.Address(),
    })
    if err != nil {
        return fmt.Errorf("failed to get account details: %v", err)
    }

    tx, err := txnbuild.NewTransaction(
        txnbuild.TransactionParams{
            SourceAccount:        &account,
            IncrementSequenceNum: true,
            Operations: []txnbuild.Operation{
                &txnbuild.ChangeTrust{
                    Line:  asset.ToChangeTrustAsset(),
                    Limit: "100000000000", // Maximum trust line limit
                },
            },
            BaseFee:    txnbuild.MinBaseFee,
            Timebounds: txnbuild.NewTimeout(300),
        },
    )
    if err != nil {
        return fmt.Errorf("failed to build transaction: %v", err)
    }

    tx, err = tx.Sign(sc.NetworkPassphrase, accountKP)
    if err != nil {
        return fmt.Errorf("failed to sign transaction: %v", err)
    }

    _, err = sc.HorizonClient.SubmitTransaction(tx)
    if err != nil {
        return fmt.Errorf("failed to submit transaction: %v", err)
    }

    return nil
}

// TransferTokens sends tokens from one account to another
func (sc *StellarClient) TransferTokens(
    fromKP *keypair.Full,
    toAddress string,
    amount string,
    assetCode string,
    issuerAddress string,
) error {
    asset := txnbuild.CreditAsset{
        Code:   assetCode,
        Issuer: issuerAddress,
    }

    account, err := sc.HorizonClient.AccountDetail(horizonclient.AccountRequest{
        AccountID: fromKP.Address(),
    })
    if err != nil {
        return fmt.Errorf("failed to get account details: %v", err)
    }

    tx, err := txnbuild.NewTransaction(
        txnbuild.TransactionParams{
            SourceAccount:        &account,
            IncrementSequenceNum: true,
            Operations: []txnbuild.Operation{
                &txnbuild.Payment{
                    Destination: toAddress,
                    Asset:       asset,
                    Amount:      amount,
                },
            },
            BaseFee:    txnbuild.MinBaseFee,
            Timebounds: txnbuild.NewTimeout(300),
        },
    )
    if err != nil {
        return fmt.Errorf("failed to build transaction: %v", err)
    }

    tx, err = tx.Sign(sc.NetworkPassphrase, fromKP)
    if err != nil {
        return fmt.Errorf("failed to sign transaction: %v", err)
    }

    _, err = sc.HorizonClient.SubmitTransaction(tx)
    if err != nil {
        return fmt.Errorf("failed to submit transaction: %v", err)
    }

    return nil
}
