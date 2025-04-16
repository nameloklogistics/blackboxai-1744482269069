package models

import (
	"time"
)

// ProposalType represents different types of governance proposals
type ProposalType string

const (
	ParameterChange  ProposalType = "PARAMETER_CHANGE"
	ContractUpgrade  ProposalType = "CONTRACT_UPGRADE"
	FundsAllocation  ProposalType = "FUNDS_ALLOCATION"
	ServiceUpdate    ProposalType = "SERVICE_UPDATE"
)

// ProposalStatus represents the status of a governance proposal
type ProposalStatus string

const (
	ProposalStatusActive        ProposalStatus = "ACTIVE"
	ProposalStatusPassed        ProposalStatus = "PASSED"
	ProposalStatusExecuted      ProposalStatus = "EXECUTED"
	ProposalStatusRejected      ProposalStatus = "REJECTED"
	ProposalStatusFailedQuorum  ProposalStatus = "FAILED_QUORUM"
)

// VoteType represents the type of vote cast on a proposal
type VoteType string

const (
	VoteTypeFor     VoteType = "FOR"
	VoteTypeAgainst VoteType = "AGAINST"
	VoteTypeAbstain VoteType = "ABSTAIN"
)

// GovernanceProposal represents a governance proposal in the system
type GovernanceProposal struct {
	BaseModel
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	ProposalType    ProposalType  `json:"proposal_type"`
	Creator         string        `json:"creator"`
	Status          ProposalStatus `json:"status"`
	StartTime       time.Time     `json:"start_time"`
	EndTime         time.Time     `json:"end_time"`
	ForVotes        uint64        `json:"for_votes"`
	AgainstVotes    uint64        `json:"against_votes"`
	AbstainVotes    uint64        `json:"abstain_votes"`
	ExecutionTime   *time.Time    `json:"execution_time,omitempty"`
	ProposalData    []byte        `json:"proposal_data"`
	ContractAddress string        `json:"contract_address"`
}

// GovernanceVote represents a vote cast on a governance proposal
type GovernanceVote struct {
	BaseModel
	ProposalID  string    `json:"proposal_id"`
	Voter       string    `json:"voter"`
	VoteType    VoteType  `json:"vote_type"`
	VotePower   uint64    `json:"vote_power"`
	VoteTime    time.Time `json:"vote_time"`
}

// GovernanceParameter represents a governance parameter in the system
type GovernanceParameter struct {
	BaseModel
	Name         string      `json:"name"`
	Value        interface{} `json:"value"`
	Description  string      `json:"description"`
	LastUpdated  time.Time   `json:"last_updated"`
	UpdatedBy    string      `json:"updated_by"`
}

// ProposalCreateRequest represents the request to create a new proposal
type ProposalCreateRequest struct {
	Title        string       `json:"title" validate:"required"`
	Description  string       `json:"description" validate:"required"`
	ProposalType ProposalType `json:"proposal_type" validate:"required"`
	ProposalData []byte       `json:"proposal_data"`
}

// VoteCastRequest represents the request to cast a vote on a proposal
type VoteCastRequest struct {
	ProposalID string   `json:"proposal_id" validate:"required"`
	VoteType   VoteType `json:"vote_type" validate:"required"`
}

// ProposalExecuteRequest represents the request to execute a proposal
type ProposalExecuteRequest struct {
	ProposalID string `json:"proposal_id" validate:"required"`
}

// ProposalResponse represents the response for proposal-related operations
type ProposalResponse struct {
	Proposal    GovernanceProposal `json:"proposal"`
	VoteSummary struct {
		TotalVotes     uint64  `json:"total_votes"`
		QuorumReached  bool    `json:"quorum_reached"`
		ApprovalRate   float64 `json:"approval_rate"`
		ParticipationRate float64 `json:"participation_rate"`
	} `json:"vote_summary"`
}

// VoteResponse represents the response for vote-related operations
type VoteResponse struct {
	Vote      GovernanceVote `json:"vote"`
	Proposal  GovernanceProposal `json:"proposal"`
}

// ParameterResponse represents the response for parameter-related operations
type ParameterResponse struct {
	Parameter GovernanceParameter `json:"parameter"`
	History   []struct {
		Value       interface{} `json:"value"`
		UpdatedAt   time.Time   `json:"updated_at"`
		UpdatedBy   string      `json:"updated_by"`
		ProposalID  string      `json:"proposal_id"`
	} `json:"history"`
}
