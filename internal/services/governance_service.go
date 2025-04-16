package services

import (
	"context"
	"fmt"
	"time"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"

	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/stellar"
)

type GovernanceService struct {
	stellarClient    *horizonclient.Client
	governanceContract *stellar.Contract
	tokenContract    *stellar.Contract
}

func NewGovernanceService(stellarClient *horizonclient.Client, governanceContract, tokenContract *stellar.Contract) *GovernanceService {
	return &GovernanceService{
		stellarClient:    stellarClient,
		governanceContract: governanceContract,
		tokenContract:    tokenContract,
	}
}

// CreateProposal creates a new governance proposal
func (s *GovernanceService) CreateProposal(ctx context.Context, req models.ProposalCreateRequest, creatorAddress string) (*models.ProposalResponse, error) {
	// Check if creator has enough tokens to create proposal
	balance, err := s.tokenContract.GetBalance(creatorAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get token balance: %w", err)
	}

	minThreshold := s.governanceContract.GetParameter("minProposalThreshold").(uint64)
	if balance < minThreshold {
		return nil, fmt.Errorf("insufficient tokens to create proposal: required %d, got %d", minThreshold, balance)
	}

	// Create proposal on blockchain
	proposalID, err := s.governanceContract.CreateProposal(
		creatorAddress,
		req.Title,
		req.Description,
		string(req.ProposalType),
		req.ProposalData,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create proposal on blockchain: %w", err)
	}

	// Create proposal in database
	proposal := &models.GovernanceProposal{
		Title:           req.Title,
		Description:     req.Description,
		ProposalType:    req.ProposalType,
		Creator:         creatorAddress,
		Status:          models.ProposalStatusActive,
		StartTime:       time.Now(),
		EndTime:         time.Now().Add(7 * 24 * time.Hour), // 7 days voting period
		ProposalData:    req.ProposalData,
		ContractAddress: s.governanceContract.Address(),
	}

	// Calculate vote summary
	voteSummary := struct {
		TotalVotes       uint64
		QuorumReached    bool
		ApprovalRate     float64
		ParticipationRate float64
	}{
		TotalVotes:       0,
		QuorumReached:    false,
		ApprovalRate:     0,
		ParticipationRate: 0,
	}

	return &models.ProposalResponse{
		Proposal:    *proposal,
		VoteSummary: voteSummary,
	}, nil
}

// CastVote casts a vote on a proposal
func (s *GovernanceService) CastVote(ctx context.Context, req models.VoteCastRequest, voterAddress string) (*models.VoteResponse, error) {
	// Get proposal
	proposal, err := s.GetProposal(ctx, req.ProposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal: %w", err)
	}

	// Check if proposal is active
	if proposal.Proposal.Status != models.ProposalStatusActive {
		return nil, fmt.Errorf("proposal is not active")
	}

	// Check if voting period is still open
	if time.Now().After(proposal.Proposal.EndTime) {
		return nil, fmt.Errorf("voting period has ended")
	}

	// Get voter's token balance at proposal start time
	balance, err := s.tokenContract.GetBalanceAt(voterAddress, proposal.Proposal.StartTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get token balance: %w", err)
	}

	// Cast vote on blockchain
	err = s.governanceContract.CastVote(req.ProposalID, voterAddress, string(req.VoteType))
	if err != nil {
		return nil, fmt.Errorf("failed to cast vote on blockchain: %w", err)
	}

	// Create vote record
	vote := &models.GovernanceVote{
		ProposalID: req.ProposalID,
		Voter:      voterAddress,
		VoteType:   req.VoteType,
		VotePower:  balance,
		VoteTime:   time.Now(),
	}

	return &models.VoteResponse{
		Vote:     *vote,
		Proposal: proposal.Proposal,
	}, nil
}

// ExecuteProposal executes a passed proposal
func (s *GovernanceService) ExecuteProposal(ctx context.Context, req models.ProposalExecuteRequest) (*models.ProposalResponse, error) {
	// Get proposal
	proposal, err := s.GetProposal(ctx, req.ProposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal: %w", err)
	}

	// Check if proposal can be executed
	if proposal.Proposal.Status != models.ProposalStatusPassed {
		return nil, fmt.Errorf("proposal is not in passed status")
	}

	// Check execution delay
	executionDelay := s.governanceContract.GetParameter("executionDelay").(int64)
	if time.Now().Before(proposal.Proposal.EndTime.Add(time.Duration(executionDelay) * time.Second)) {
		return nil, fmt.Errorf("execution delay not met")
	}

	// Execute proposal on blockchain
	err = s.governanceContract.ExecuteProposal(req.ProposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute proposal on blockchain: %w", err)
	}

	// Update proposal status
	proposal.Proposal.Status = models.ProposalStatusExecuted
	proposal.Proposal.ExecutionTime = new(time.Time)
	*proposal.Proposal.ExecutionTime = time.Now()

	return proposal, nil
}

// GetProposal gets a proposal by ID
func (s *GovernanceService) GetProposal(ctx context.Context, proposalID string) (*models.ProposalResponse, error) {
	// Get proposal from blockchain
	proposal, err := s.governanceContract.GetProposal(proposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal from blockchain: %w", err)
	}

	// Calculate vote summary
	totalSupply := s.tokenContract.GetTotalSupply()
	totalVotes := proposal.ForVotes + proposal.AgainstVotes + proposal.AbstainVotes
	
	quorumPercentage := s.governanceContract.GetParameter("quorumPercentage").(uint64)
	quorumReached := (totalVotes * 100 / totalSupply) >= quorumPercentage

	var approvalRate float64
	if totalVotes > 0 {
		approvalRate = float64(proposal.ForVotes) / float64(totalVotes) * 100
	}

	participationRate := float64(totalVotes) / float64(totalSupply) * 100

	voteSummary := struct {
		TotalVotes       uint64
		QuorumReached    bool
		ApprovalRate     float64
		ParticipationRate float64
	}{
		TotalVotes:       totalVotes,
		QuorumReached:    quorumReached,
		ApprovalRate:     approvalRate,
		ParticipationRate: participationRate,
	}

	return &models.ProposalResponse{
		Proposal:    proposal,
		VoteSummary: voteSummary,
	}, nil
}

// ListProposals lists all proposals with optional filters
func (s *GovernanceService) ListProposals(ctx context.Context, status models.ProposalStatus, proposalType models.ProposalType) ([]models.ProposalResponse, error) {
	// Get proposals from blockchain
	proposals, err := s.governanceContract.ListProposals()
	if err != nil {
		return nil, fmt.Errorf("failed to list proposals from blockchain: %w", err)
	}

	var response []models.ProposalResponse
	for _, proposal := range proposals {
		if (status == "" || proposal.Status == status) &&
			(proposalType == "" || proposal.ProposalType == proposalType) {
			
			// Get vote summary for each proposal
			proposalResp, err := s.GetProposal(ctx, proposal.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get proposal details: %w", err)
			}
			
			response = append(response, *proposalResp)
		}
	}

	return response, nil
}

// GetVote gets a vote by proposal ID and voter address
func (s *GovernanceService) GetVote(ctx context.Context, proposalID string, voterAddress string) (*models.VoteResponse, error) {
	// Get vote from blockchain
	vote, err := s.governanceContract.GetVote(proposalID, voterAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get vote from blockchain: %w", err)
	}

	// Get proposal details
	proposal, err := s.GetProposal(ctx, proposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal details: %w", err)
	}

	return &models.VoteResponse{
		Vote:     vote,
		Proposal: proposal.Proposal,
	}, nil
}

// GetParameter gets a governance parameter by name
func (s *GovernanceService) GetParameter(ctx context.Context, name string) (*models.ParameterResponse, error) {
	// Get parameter from blockchain
	value := s.governanceContract.GetParameter(name)
	if value == nil {
		return nil, fmt.Errorf("parameter not found: %s", name)
	}

	parameter := &models.GovernanceParameter{
		Name:        name,
		Value:       value,
		LastUpdated: time.Now(), // This should come from blockchain
	}

	// Get parameter history (if implemented)
	history := []struct {
		Value      interface{} `json:"value"`
		UpdatedAt  time.Time   `json:"updated_at"`
		UpdatedBy  string      `json:"updated_by"`
		ProposalID string      `json:"proposal_id"`
	}{}

	return &models.ParameterResponse{
		Parameter: *parameter,
		History:   history,
	}, nil
}
