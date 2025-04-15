package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/services"
)

type ForwarderOperationsHandler struct {
	operationsService *services.ForwarderOperationsService
}

func NewForwarderOperationsHandler(operationsService *services.ForwarderOperationsService) *ForwarderOperationsHandler {
	return &ForwarderOperationsHandler{
		operationsService: operationsService,
	}
}

// CreateQuoteRequest handles creation of new freight quote requests
func (h *ForwarderOperationsHandler) CreateQuoteRequest(c *gin.Context) {
	var request models.QuoteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get forwarder ID from authenticated user
	forwarderID := c.GetString("user_id")
	request.ForwarderID = forwarderID

	if err := h.operationsService.CreateQuoteRequest(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, request)
}

// GenerateFreightQuote handles generation of detailed freight quotes
func (h *ForwarderOperationsHandler) GenerateFreightQuote(c *gin.Context) {
	var request models.QuoteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get forwarder ID from authenticated user
	forwarderID := c.GetString("user_id")
	request.ForwarderID = forwarderID

	quote, err := h.operationsService.GenerateFreightQuote(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quote)
}

// ConfirmRate handles rate confirmation
func (h *ForwarderOperationsHandler) ConfirmRate(c *gin.Context) {
	var quote models.FreightQuote
	if err := c.ShouldBindJSON(&quote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get forwarder ID from authenticated user
	forwarderID := c.GetString("user_id")
	quote.ForwarderID = forwarderID

	confirmation, err := h.operationsService.ConfirmRate(&quote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, confirmation)
}

// CreateBooking handles creation of new shipment bookings
func (h *ForwarderOperationsHandler) CreateBooking(c *gin.Context) {
	var request struct {
		ConfirmationID string                `json:"confirmation_id"`
		Booking        models.ShipmentBooking `json:"booking"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get forwarder ID from authenticated user
	forwarderID := c.GetString("user_id")
	request.Booking.ForwarderID = forwarderID

	// Get rate confirmation
	confirmation := &models.RateConfirmation{
		ID:          request.ConfirmationID,
		ForwarderID: forwarderID,
		// In a real implementation, this would be retrieved from storage
	}

	booking, err := h.operationsService.CreateBooking(confirmation, &request.Booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// ConfirmBooking handles booking confirmation
func (h *ForwarderOperationsHandler) ConfirmBooking(c *gin.Context) {
	var booking models.ShipmentBooking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get forwarder ID from authenticated user
	forwarderID := c.GetString("user_id")
	booking.ForwarderID = forwarderID

	confirmation, err := h.operationsService.ConfirmBooking(&booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, confirmation)
}

// GetQuoteRequest retrieves a quote request by ID
func (h *ForwarderOperationsHandler) GetQuoteRequest(c *gin.Context) {
	requestID := c.Param("id")
	forwarderID := c.GetString("user_id")

	// In a real implementation, this would retrieve from storage
	request := &models.QuoteRequest{
		ID:          requestID,
		ForwarderID: forwarderID,
		Status:      "NEW",
	}

	c.JSON(http.StatusOK, request)
}

// GetFreightQuote retrieves a freight quote by ID
func (h *ForwarderOperationsHandler) GetFreightQuote(c *gin.Context) {
	quoteID := c.Param("id")
	forwarderID := c.GetString("user_id")

	// In a real implementation, this would retrieve from storage
	quote := &models.FreightQuote{
		ID:          quoteID,
		ForwarderID: forwarderID,
		Status:      "DRAFT",
	}

	c.JSON(http.StatusOK, quote)
}

// GetBooking retrieves a booking by ID
func (h *ForwarderOperationsHandler) GetBooking(c *gin.Context) {
	bookingID := c.Param("id")
	forwarderID := c.GetString("user_id")

	// In a real implementation, this would retrieve from storage
	booking := &models.ShipmentBooking{
		ID:          bookingID,
		ForwarderID: forwarderID,
		Status:      "PENDING",
	}

	c.JSON(http.StatusOK, booking)
}

// ListQuoteRequests lists all quote requests for a forwarder
func (h *ForwarderOperationsHandler) ListQuoteRequests(c *gin.Context) {
	forwarderID := c.GetString("user_id")
	status := c.Query("status") // Optional status filter

	// In a real implementation, this would retrieve from storage
	requests := []models.QuoteRequest{
		{
			ID:          "QR-1",
			ForwarderID: forwarderID,
			Status:      status,
		},
	}

	c.JSON(http.StatusOK, requests)
}

// ListFreightQuotes lists all freight quotes for a forwarder
func (h *ForwarderOperationsHandler) ListFreightQuotes(c *gin.Context) {
	forwarderID := c.GetString("user_id")
	status := c.Query("status") // Optional status filter

	// In a real implementation, this would retrieve from storage
	quotes := []models.FreightQuote{
		{
			ID:          "FQ-1",
			ForwarderID: forwarderID,
			Status:      status,
		},
	}

	c.JSON(http.StatusOK, quotes)
}

// ListBookings lists all bookings for a forwarder
func (h *ForwarderOperationsHandler) ListBookings(c *gin.Context) {
	forwarderID := c.GetString("user_id")
	status := c.Query("status") // Optional status filter

	// In a real implementation, this would retrieve from storage
	bookings := []models.ShipmentBooking{
		{
			ID:          "BK-1",
			ForwarderID: forwarderID,
			Status:      status,
		},
	}

	c.JSON(http.StatusOK, bookings)
}
