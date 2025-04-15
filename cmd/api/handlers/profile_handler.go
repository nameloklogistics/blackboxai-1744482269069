package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/services"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// CreateServiceProvider handles creation of new service provider profiles
func (h *ProfileHandler) CreateServiceProvider(c *gin.Context) {
	var provider models.ServiceProvider
	if err := c.ShouldBindJSON(&provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.profileService.CreateServiceProvider(&provider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, provider)
}

// CreateServiceBuyer handles creation of new service buyer profiles
func (h *ProfileHandler) CreateServiceBuyer(c *gin.Context) {
	var buyer models.ServiceBuyer
	if err := c.ShouldBindJSON(&buyer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.profileService.CreateServiceBuyer(&buyer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, buyer)
}

// UpdateServiceProvider handles updates to service provider profiles
func (h *ProfileHandler) UpdateServiceProvider(c *gin.Context) {
	id := c.Param("id")
	var provider models.ServiceProvider
	if err := c.ShouldBindJSON(&provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify provider ID matches authenticated user
	if id != c.GetString("user_id") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	provider.ID = id
	if err := h.profileService.UpdateServiceProvider(&provider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, provider)
}

// UpdateServiceBuyer handles updates to service buyer profiles
func (h *ProfileHandler) UpdateServiceBuyer(c *gin.Context) {
	id := c.Param("id")
	var buyer models.ServiceBuyer
	if err := c.ShouldBindJSON(&buyer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify buyer ID matches authenticated user
	if id != c.GetString("user_id") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	buyer.ID = id
	if err := h.profileService.UpdateServiceBuyer(&buyer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, buyer)
}

// AddPort handles adding ports to service provider infrastructure
func (h *ProfileHandler) AddPort(c *gin.Context) {
	providerID := c.Param("id")
	var port models.Port
	if err := c.ShouldBindJSON(&port); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify provider ID matches authenticated user
	if providerID != c.GetString("user_id") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.profileService.AddPort(providerID, port); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "port added successfully"})
}

// AddAirport handles adding airports to service provider infrastructure
func (h *ProfileHandler) AddAirport(c *gin.Context) {
	providerID := c.Param("id")
	var airport models.Airport
	if err := c.ShouldBindJSON(&airport); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify provider ID matches authenticated user
	if providerID != c.GetString("user_id") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.profileService.AddAirport(providerID, airport); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "airport added successfully"})
}

// AddTerminal handles adding container terminals to service provider infrastructure
func (h *ProfileHandler) AddTerminal(c *gin.Context) {
	providerID := c.Param("id")
	var terminal models.Terminal
	if err := c.ShouldBindJSON(&terminal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify provider ID matches authenticated user
	if providerID != c.GetString("user_id") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.profileService.AddTerminal(providerID, terminal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "terminal added successfully"})
}
