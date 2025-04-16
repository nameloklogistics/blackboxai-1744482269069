package marketplace

import (
	"github.com/stellar/soroban-sdk/go/soroban"
	"time"
)

const (
	MinProposalThreshold = 1_000_000     // Minimum tokens required to create proposal
	VotingPeriod        = 7 * 24 * 3600  // 7 days in seconds
	ExecutionDelay      = 2 * 24 * 3600  // 2 days delay after voting ends
	QuorumPercentage    = 4              // 4% of total supply must vote
	MajorityPercentage  = 51             // 51% of votes must be in favor
)

type ProposalType int

const (
	ParameterChange ProposalType = iota
	ContractUpgrade
	FundsAllocation
	ServiceUpdate
)

type VoteType int

const (
	Against VoteType = iota
	For
	Abstain
)

type Proposal struct {
	ID          string
	Creator     string
	Title       string
	Description string
	Type        ProposalType
	StartTime   int64
	EndTime     int64
	Status      string
	ForVotes    uint64
	AgainstVotes uint64
	AbstainVotes uint64
	Executed    bool
	Data        []byte // Encoded proposal-specific data
}

type Vote struct {
	Voter     string
	VoteType  VoteType
	Amount    uint64
	Timestamp int64
}

type GovernanceContract struct {
	soroban.Contract
	token          *TokenContract
	proposals      map[string]Proposal
	votes         map[string]map[string]Vote // proposalID -> voter -> Vote
	parameters    map[string]interface{}
}

func (c *GovernanceContract) Initialize(env soroban.Env, tokenAddress string) {
	if c.token != nil {
		panic("Contract already initialized")
	}

	// Initialize the governance contract with token contract
	c.token = &TokenContract{}
	c.proposals = make(map[string]Proposal)
	c.votes = make(map[string]map[string]Vote)
	c.parameters = make(map[string]interface{})

	// Set initial governance parameters
	c.parameters["minProposalThreshold"] = MinProposalThreshold
	c.parameters["votingPeriod"] = VotingPeriod
	c.parameters["executionDelay"] = ExecutionDelay
	c.parameters["quorumPercentage"] = QuorumPercentage
	c.parameters["majorityPercentage"] = MajorityPercentage
}

func (c *GovernanceContract) CreateProposal(env soroban.Env, creator string, title string, description string, proposalType ProposalType, data []byte) string {
	// Check if creator has enough tokens
	creatorBalance := c.token.BalanceOf(creator)
	if creatorBalance < MinProposalThreshold {
		panic("Insufficient tokens to create proposal")
	}

	// Generate proposal ID
	proposalID := env.Crypto().RandomBytes(32).String()

	// Create new proposal
	proposal := Proposal{
		ID:          proposalID,
		Creator:     creator,
		Title:       title,
		Description: description,
		Type:        proposalType,
		StartTime:   env.Ledger().Timestamp(),
		EndTime:     env.Ledger().Timestamp() + VotingPeriod,
		Status:      "ACTIVE",
		Data:        data,
	}

	c.proposals[proposalID] = proposal

	// Initialize votes mapping for this proposal
	c.votes[proposalID] = make(map[string]Vote)

	// Emit proposal creation event
	env.Events().Publish("proposal_created", map[string]interface{}{
		"proposal_id": proposalID,
		"creator":     creator,
		"title":       title,
	})

	return proposalID
}

func (c *GovernanceContract) CastVote(env soroban.Env, voter string, proposalID string, voteType VoteType) bool {
	proposal, exists := c.proposals[proposalID]
	if !exists {
		panic("Proposal does not exist")
	}

	// Check if proposal is active
	currentTime := env.Ledger().Timestamp()
	if currentTime > proposal.EndTime || currentTime < proposal.StartTime {
		panic("Proposal is not active")
	}

	// Get voter's token balance at proposal start
	voterBalance := c.token.BalanceOf(voter)
	if voterBalance == 0 {
		panic("No voting power")
	}

	// Record the vote
	vote := Vote{
		Voter:     voter,
		VoteType:  voteType,
		Amount:    voterBalance,
		Timestamp: currentTime,
	}

	// Update vote counts
	if oldVote, voted := c.votes[proposalID][voter]; voted {
		// Remove old vote counts
		switch oldVote.VoteType {
		case Against:
			proposal.AgainstVotes -= oldVote.Amount
		case For:
			proposal.ForVotes -= oldVote.Amount
		case Abstain:
			proposal.AbstainVotes -= oldVote.Amount
		}
	}

	// Add new vote counts
	switch voteType {
	case Against:
		proposal.AgainstVotes += voterBalance
	case For:
		proposal.ForVotes += voterBalance
	case Abstain:
		proposal.AbstainVotes += voterBalance
	}

	c.votes[proposalID][voter] = vote
	c.proposals[proposalID] = proposal

	// Emit vote cast event
	env.Events().Publish("vote_cast", map[string]interface{}{
		"proposal_id": proposalID,
		"voter":      voter,
		"vote_type":  voteType,
		"amount":     voterBalance,
	})

	return true
}

func (c *GovernanceContract) ExecuteProposal(env soroban.Env, proposalID string) bool {
	proposal, exists := c.proposals[proposalID]
	if !exists {
		panic("Proposal does not exist")
	}

	// Check if proposal can be executed
	currentTime := env.Ledger().Timestamp()
	if currentTime < proposal.EndTime + ExecutionDelay {
		panic("Execution delay not met")
	}

	if proposal.Executed {
		panic("Proposal already executed")
	}

	// Calculate total votes
	totalVotes := proposal.ForVotes + proposal.AgainstVotes + proposal.AbstainVotes
	
	// Check quorum
	quorum := (totalVotes * 100) / c.token.TotalSupply()
	if quorum < QuorumPercentage {
		proposal.Status = "FAILED_QUORUM"
		c.proposals[proposalID] = proposal
		return false
	}

	// Check majority
	forPercentage := (proposal.ForVotes * 100) / totalVotes
	if forPercentage < MajorityPercentage {
		proposal.Status = "REJECTED"
		c.proposals[proposalID] = proposal
		return false
	}

	// Execute proposal based on type
	switch proposal.Type {
	case ParameterChange:
		c.executeParameterChange(env, proposal)
	case ContractUpgrade:
		c.executeContractUpgrade(env, proposal)
	case FundsAllocation:
		c.executeFundsAllocation(env, proposal)
	case ServiceUpdate:
		c.executeServiceUpdate(env, proposal)
	}

	// Mark proposal as executed
	proposal.Executed = true
	proposal.Status = "EXECUTED"
	c.proposals[proposalID] = proposal

	// Emit proposal execution event
	env.Events().Publish("proposal_executed", map[string]interface{}{
		"proposal_id": proposalID,
		"status":     proposal.Status,
	})

	return true
}

func (c *GovernanceContract) executeParameterChange(env soroban.Env, proposal Proposal) {
	// Implementation for parameter change execution
}

func (c *GovernanceContract) executeContractUpgrade(env soroban.Env, proposal Proposal) {
	// Implementation for contract upgrade execution
}

func (c *GovernanceContract) executeFundsAllocation(env soroban.Env, proposal Proposal) {
	// Implementation for funds allocation execution
}

func (c *GovernanceContract) executeServiceUpdate(env soroban.Env, proposal Proposal) {
	// Implementation for service update execution
}

func (c *GovernanceContract) GetProposal(proposalID string) Proposal {
	proposal, exists := c.proposals[proposalID]
	if !exists {
		panic("Proposal does not exist")
	}
	return proposal
}

func (c *GovernanceContract) GetVote(proposalID string, voter string) Vote {
	votes, exists := c.votes[proposalID]
	if !exists {
		panic("Proposal does not exist")
	}
	vote, voted := votes[voter]
	if !voted {
		panic("Vote does not exist")
	}
	return vote
}

func (c *GovernanceContract) GetParameter(name string) interface{} {
	value, exists := c.parameters[name]
	if !exists {
		panic("Parameter does not exist")
	}
	return value
}
