package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/services"
)

type TransportHandler struct {
	marketplaceService *services.MarketplaceService
}

func NewTransportHandler(marketplaceService *services.MarketplaceService) *TransportHandler {
	return &TransportHandler{
		marketplaceService: marketplaceService,
	}
}

// CreateSeaService handles creation of sea transport services
func (h *TransportHandler) CreateSeaService(c *gin.Context) {
	var service models.SeaService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set transport mode
	service.Mode = models.SeaTransport

	// Get provider ID from authenticated user
	providerID := c.GetString("user_id")
	service.Provider.ID = providerID

	if err := h.marketplaceService.CreateTransportService(&service.TransportService); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// CreateAirService handles creation of air transport services
func (h *TransportHandler) CreateAirService(c *gin.Context) {
	var service models.AirService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set transport mode
	service.Mode = models.AirTransport

	// Get provider ID from authenticated user
	providerID := c.GetString("user_id")
	service.Provider.ID = providerID

	if err := h.marketplaceService.CreateTransportService(&service.TransportService); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// CreateRailService handles creation of rail transport services
func (h *TransportHandler) CreateRailService(c *gin.Context) {
	var service models.RailService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set transport mode
	service.Mode = models.RailTransport

	// Get provider ID from authenticated user
	providerID := c.GetString("user_id")
	service.Provider.ID = providerID

	if err := h.marketplaceService.CreateTransportService(&service.TransportService); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// CreateRoadService handles creation of road transport services
func (h *TransportHandler) CreateRoadService(c *gin.Context) {
	var service models.RoadService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set transport mode
	service.Mode = models.RoadTransport

	// Get provider ID from authenticated user
	providerID := c.GetString("user_id")
	service.Provider.ID = providerID

	if err := h.marketplaceService.CreateTransportService(&service.TransportService); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// ListServicesByMode handles listing services by transport mode
func (h *TransportHandler) ListServicesByMode(c *gin.Context) {
	mode := models.TransportMode(c.Param("mode"))
	category := c.Query("category") // Optional category filter

	services, err := h.marketplaceService.GetServicesByMode(mode, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// GetServiceTemplates returns service templates for a transport mode
func (h *TransportHandler) GetServiceTemplates(c *gin.Context) {
	mode := models.TransportMode(c.Param("mode"))
	
	templates := models.GetServiceTemplatesByMode(mode)
	if templates == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transport mode"})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// GetEquipmentTypes returns available equipment types for a transport mode
func (h *TransportHandler) GetEquipmentTypes(c *gin.Context) {
	mode := models.TransportMode(c.Param("mode"))
	
	equipment := models.GetEquipmentByMode(mode)
	if equipment == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transport mode"})
		return
	}

	c.JSON(http.StatusOK, equipment)
}

// UpdateService handles updates to transport services
func (h *TransportHandler) UpdateService(c *gin.Context) {
	serviceID := c.Param("id")
	mode := models.TransportMode(c.Param("mode"))

	// Validate transport mode
	switch mode {
	case models.SeaTransport, models.AirTransport, models.RailTransport, models.RoadTransport:
		// Valid mode
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transport mode"})
		return
	}

	var service models.TransportService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify provider ID matches authenticated user
	providerID := c.GetString("user_id")
	if service.Provider.ID != providerID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	service.ID = serviceID
	service.Mode = mode

	if err := h.marketplaceService.UpdateTransportService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

// GetServiceQuotation handles quotation requests for transport services
func (h *TransportHandler) GetServiceQuotation(c *gin.Context) {
	var request struct {
		ServiceID    string            `json:"service_id"`
		Mode         models.TransportMode `json:"mode"`
		CargoDetails models.CargoDetails `json:"cargo_details"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quotation, err := h.marketplaceService.GetTransportQuotation(
		request.ServiceID,
		request.Mode,
		&request.CargoDetails,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quotation)
}
