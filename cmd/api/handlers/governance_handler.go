package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/services"
)

type GovernanceHandler struct {
	governanceService *services.GovernanceService
}

func NewGovernanceHandler(governanceService *services.GovernanceService) *GovernanceHandler {
	return &GovernanceHandler{
		governanceService: governanceService,
	}
}

// RegisterRoutes registers the governance routes
func (h *GovernanceHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/governance/proposals", h.CreateProposal).Methods("POST")
	router.HandleFunc("/api/v1/governance/proposals", h.ListProposals).Methods("GET")
	router.HandleFunc("/api/v1/governance/proposals/{id}", h.GetProposal).Methods("GET")
	router.HandleFunc("/api/v1/governance/proposals/{id}/vote", h.CastVote).Methods("POST")
	router.HandleFunc("/api/v1/governance/proposals/{id}/execute", h.ExecuteProposal).Methods("POST")
	router.HandleFunc("/api/v1/governance/proposals/{id}/votes/{voter}", h.GetVote).Methods("GET")
	router.HandleFunc("/api/v1/governance/parameters/{name}", h.GetParameter).Methods("GET")
}

// CreateProposal handles proposal creation requests
func (h *GovernanceHandler) CreateProposal(w http.ResponseWriter, r *http.Request) {
	var req models.ProposalCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get creator address from authenticated user
	creatorAddress := r.Context().Value("user_address").(string)

	response, err := h.governanceService.CreateProposal(r.Context(), req, creatorAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

// ListProposals handles listing proposals with optional filters
func (h *GovernanceHandler) ListProposals(w http.ResponseWriter, r *http.Request) {
	status := models.ProposalStatus(r.URL.Query().Get("status"))
	proposalType := models.ProposalType(r.URL.Query().Get("type"))

	proposals, err := h.governanceService.ListProposals(r.Context(), status, proposalType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(proposals)
}

// GetProposal handles getting a specific proposal
func (h *GovernanceHandler) GetProposal(w http.ResponseWriter, r *http.Request) {
	proposalID := mux.Vars(r)["id"]

	proposal, err := h.governanceService.GetProposal(r.Context(), proposalID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(proposal)
}

// CastVote handles vote casting on proposals
func (h *GovernanceHandler) CastVote(w http.ResponseWriter, r *http.Request) {
	proposalID := mux.Vars(r)["id"]

	var req models.VoteCastRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.ProposalID = proposalID

	// Get voter address from authenticated user
	voterAddress := r.Context().Value("user_address").(string)

	response, err := h.governanceService.CastVote(r.Context(), req, voterAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

// ExecuteProposal handles proposal execution requests
func (h *GovernanceHandler) ExecuteProposal(w http.ResponseWriter, r *http.Request) {
	proposalID := mux.Vars(r)["id"]

	req := models.ProposalExecuteRequest{
		ProposalID: proposalID,
	}

	response, err := h.governanceService.ExecuteProposal(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

// GetVote handles getting a specific vote
func (h *GovernanceHandler) GetVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proposalID := vars["id"]
	voterAddress := vars["voter"]

	vote, err := h.governanceService.GetVote(r.Context(), proposalID, voterAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(vote)
}

// GetParameter handles getting governance parameters
func (h *GovernanceHandler) GetParameter(w http.ResponseWriter, r *http.Request) {
	paramName := mux.Vars(r)["name"]

	parameter, err := h.governanceService.GetParameter(r.Context(), paramName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(parameter)
}
