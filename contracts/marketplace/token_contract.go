package marketplace

import (
	"github.com/stellar/soroban-sdk/go/soroban"
)

const (
	MaxSupply = 100_000_000_000 // 100 billion tokens
	TokenName = "Logistics Marketplace Token"
	TokenSymbol = "LMT"
	TokenDecimals = 7
)

type TokenContract struct {
	soroban.Contract
	token    soroban.Token
	balances map[string]uint64
}

func (c *TokenContract) Initialize(env soroban.Env) {
	if c.token != nil {
		panic("Contract already initialized")
	}

	// Create the token with initial supply
	c.token = env.Token()
	c.balances = make(map[string]uint64)
	
	// Mint initial supply to contract creator
	admin := env.Current().Contract().Address()
	c.balances[admin.String()] = MaxSupply
}

func (c *TokenContract) Name() string {
	return TokenName
}

func (c *TokenContract) Symbol() string {
	return TokenSymbol
}

func (c *TokenContract) Decimals() uint32 {
	return TokenDecimals
}

func (c *TokenContract) TotalSupply() uint64 {
	return MaxSupply
}

func (c *TokenContract) BalanceOf(owner string) uint64 {
	balance, exists := c.balances[owner]
	if !exists {
		return 0
	}
	return balance
}

func (c *TokenContract) Transfer(env soroban.Env, from string, to string, amount uint64) bool {
	if amount == 0 {
		return false
	}

	fromBalance := c.BalanceOf(from)
	if fromBalance < amount {
		return false
	}

	c.balances[from] = fromBalance - amount
	c.balances[to] += amount

	// Emit transfer event
	env.Events().Publish("transfer", map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	})

	return true
}

func (c *TokenContract) Approve(env soroban.Env, owner string, spender string, amount uint64) bool {
	// Implementation for token approval
	// This would allow marketplace contract to transfer tokens on behalf of users
	return true
}
