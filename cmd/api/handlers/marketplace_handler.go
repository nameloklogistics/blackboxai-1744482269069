package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/services"
)

type MarketplaceHandler struct {
	marketplaceService *services.MarketplaceService
	customsService    *services.CustomsService
}

func NewMarketplaceHandler(
	marketplaceService *services.MarketplaceService,
	customsService *services.CustomsService,
) *MarketplaceHandler {
	return &MarketplaceHandler{
		marketplaceService: marketplaceService,
		customsService:    customsService,
	}
}

// ListServicesByCategory handles listing services by main category
func (h *MarketplaceHandler) ListServicesByCategory(c *gin.Context) {
	category := models.ServiceCategory(c.Param("category"))
	
	// Validate category
	switch category {
	case models.ImportService, models.ExportService, models.TransitService, models.TransshipService:
		// Valid category
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service category"})
		return
	}

	services, err := h.marketplaceService.GetServicesByCategory(string(category))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// ListSubCategories handles listing subcategories for a main category
func (h *MarketplaceHandler) ListSubCategories(c *gin.Context) {
	category := models.ServiceCategory(c.Param("category"))
	
	var subcategories []models.ServiceSubCategory
	switch category {
	case models.ImportService:
		subcategories = models.ImportSubCategories
	case models.ExportService:
		subcategories = models.ExportSubCategories
	case models.TransitService:
		subcategories = models.TransitSubCategories
	case models.TransshipService:
		subcategories = models.TransshipSubCategories
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service category"})
		return
	}

	c.JSON(http.StatusOK, subcategories)
}

// ListServiceItems handles listing service items for a subcategory
func (h *MarketplaceHandler) ListServiceItems(c *gin.Context) {
	subcategoryID := c.Param("subcategory")
	
	items, err := h.marketplaceService.GetServiceItems(subcategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// CreateServiceListing handles creation of new service listings
func (h *MarketplaceHandler) CreateServiceListing(c *gin.Context) {
	var service models.LogisticsService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get provider ID from authenticated user
	providerID := c.GetString("user_id")
	service.Provider.ID = providerID

	if err := h.marketplaceService.CreateServiceListing(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// GetQuotation handles rate quotation requests
func (h *MarketplaceHandler) GetQuotation(c *gin.Context) {
	var request struct {
		ServiceID    string            `json:"service_id"`
		CargoDetails models.CargoDetails `json:"cargo_details"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rate, err := h.marketplaceService.GetQuotation(request.ServiceID, &request.CargoDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rate)
}

// CreateBooking handles new booking creation
func (h *MarketplaceHandler) CreateBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get customer ID from authenticated user
	customerID := c.GetString("user_id")
	booking.CustomerID = customerID

	if err := h.marketplaceService.CreateBooking(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// ProcessPayment handles booking payment
func (h *MarketplaceHandler) ProcessPayment(c *gin.Context) {
	var request struct {
		BookingID string  `json:"booking_id"`
		Amount    float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerID := c.GetString("user_id")
	if err := h.marketplaceService.ProcessPayment(request.BookingID, customerID, request.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "payment processed"})
}

// UpdateShipmentStatus handles shipment status updates
func (h *MarketplaceHandler) UpdateShipmentStatus(c *gin.Context) {
	var event models.TrackingEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.marketplaceService.UpdateShipmentStatus(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

// InitiateCustomsClearance handles customs clearance initiation
func (h *MarketplaceHandler) InitiateCustomsClearance(c *gin.Context) {
	var clearance models.CustomsClearance
	if err := c.ShouldBindJSON(&clearance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.customsService.InitiateCustomsClearance(&clearance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, clearance)
}

// SubmitCustomsDeclaration handles customs declaration submission
func (h *MarketplaceHandler) SubmitCustomsDeclaration(c *gin.Context) {
	var request struct {
		ClearanceID     string   `json:"clearance_id"`
		DeclarationType string   `json:"declaration_type"`
		Documents       []string `json:"documents"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.customsService.SubmitDeclaration(
		request.ClearanceID,
		request.DeclarationType,
		request.Documents,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "declaration submitted"})
}

// GetRequiredDocuments returns required documents for customs clearance
func (h *MarketplaceHandler) GetRequiredDocuments(c *gin.Context) {
	declarationType := c.Query("declaration_type")
	if declarationType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "declaration_type is required"})
		return
	}

	documents := h.customsService.GetRequiredDocuments(declarationType)
	c.JSON(http.StatusOK, gin.H{"documents": documents})
}
