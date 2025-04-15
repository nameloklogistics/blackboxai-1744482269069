package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/services"
)

type FreightForwarderHandler struct {
	forwarderService *services.FreightForwarderService
}

func NewFreightForwarderHandler(forwarderService *services.FreightForwarderService) *FreightForwarderHandler {
	return &FreightForwarderHandler{
		forwarderService: forwarderService,
	}
}

// CreateRateRequest handles creation of new rate requests
func (h *FreightForwarderHandler) CreateRateRequest(c *gin.Context) {
	var request models.RateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get forwarder ID from authenticated user
	forwarderID := c.GetString("user_id")
	request.ForwarderID = forwarderID

	if err := h.forwarderService.CreateRateRequest(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, request)
}

// CreateRateQuotation handles creation of rate quotations
func (h *FreightForwarderHandler) CreateRateQuotation(c *gin.Context) {
	var quotation models.RateQuotation
	if err := c.ShouldBindJSON(&quotation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get forwarder ID from authenticated user
	forwarderID := c.GetString("user_id")
	quotation.ForwarderID = forwarderID

	if err := h.forwarderService.CreateRateQuotation(&quotation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, quotation)
}

// ConfirmRate handles rate confirmation
func (h *FreightForwarderHandler) ConfirmRate(c *gin.Context) {
	var confirmation models.RateConfirmation
	if err := c.ShouldBindJSON(&confirmation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get forwarder ID from authenticated user
	forwarderID := c.GetString("user_id")
	confirmation.ForwarderID = forwarderID

	if err := h.forwarderService.ConfirmRate(&confirmation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, confirmation)
}

// CreateBooking handles creation of shipment bookings
func (h *FreightForwarderHandler) CreateBooking(c *gin.Context) {
	var booking models.ShipmentBooking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get forwarder ID from authenticated user
	forwarderID := c.GetString("user_id")
	booking.ForwarderID = forwarderID

	if err := h.forwarderService.CreateBooking(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// ConfirmBooking handles booking confirmation
func (h *FreightForwarderHandler) ConfirmBooking(c *gin.Context) {
	var confirmation models.BookingConfirmation
	if err := c.ShouldBindJSON(&confirmation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.forwarderService.ConfirmBooking(&confirmation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, confirmation)
}

// CalculateRates handles rate calculation requests
func (h *FreightForwarderHandler) CalculateRates(c *gin.Context) {
	var request models.RateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get forwarder ID from authenticated user
	forwarderID := c.GetString("user_id")
	request.ForwarderID = forwarderID

	calculation, err := h.forwarderService.CalculateRates(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, calculation)
}

// AddBookingInstruction handles adding instructions to bookings
func (h *FreightForwarderHandler) AddBookingInstruction(c *gin.Context) {
	var instruction models.BookingInstruction
	if err := c.ShouldBindJSON(&instruction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from authenticated user
	userID := c.GetString("user_id")
	instruction.UpdatedBy = userID

	if err := h.forwarderService.AddBookingInstruction(&instruction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, instruction)
}

// UpdateDocumentRequirement handles updating document requirements
func (h *FreightForwarderHandler) UpdateDocumentRequirement(c *gin.Context) {
	var requirement models.DocumentRequirement
	if err := c.ShouldBindJSON(&requirement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.forwarderService.UpdateDocumentRequirement(&requirement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requirement)
}

// GetRateQuotation handles retrieving rate quotations
func (h *FreightForwarderHandler) GetRateQuotation(c *gin.Context) {
	quotationID := c.Param("id")
	forwarderID := c.GetString("user_id")

	// In a real implementation, this would retrieve the quotation from storage
	quotation := &models.RateQuotation{
		ID:          quotationID,
		ForwarderID: forwarderID,
		// Other fields would be populated from storage
	}

	c.JSON(http.StatusOK, quotation)
}

// GetBooking handles retrieving booking details
func (h *FreightForwarderHandler) GetBooking(c *gin.Context) {
	bookingID := c.Param("id")
	forwarderID := c.GetString("user_id")

	// In a real implementation, this would retrieve the booking from storage
	booking := &models.ShipmentBooking{
		ID:          bookingID,
		ForwarderID: forwarderID,
		// Other fields would be populated from storage
	}

	c.JSON(http.StatusOK, booking)
}

// GetBookingInstructions handles retrieving booking instructions
func (h *FreightForwarderHandler) GetBookingInstructions(c *gin.Context) {
	bookingID := c.Param("booking_id")

	// In a real implementation, this would retrieve instructions from storage
	instructions := []models.BookingInstruction{
		{
			BookingID:       bookingID,
			InstructionType: "PICKUP",
			// Other fields would be populated from storage
		},
	}

	c.JSON(http.StatusOK, instructions)
}

// GetDocumentRequirements handles retrieving document requirements
func (h *FreightForwarderHandler) GetDocumentRequirements(c *gin.Context) {
	bookingID := c.Param("booking_id")

	// In a real implementation, this would retrieve requirements from storage
	requirements := []models.DocumentRequirement{
		{
			BookingID:    bookingID,
			DocumentType: "COMMERCIAL_INVOICE",
			Required:     true,
			// Other fields would be populated from storage
		},
	}

	c.JSON(http.StatusOK, requirements)
}
