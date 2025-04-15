package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/services"
)

type UserOperationsHandler struct {
	operationsService *services.UserOperationsService
}

func NewUserOperationsHandler(operationsService *services.UserOperationsService) *UserOperationsHandler {
	return &UserOperationsHandler{
		operationsService: operationsService,
	}
}

// CreateQuoteRequest handles quote request creation from any user type
func (h *UserOperationsHandler) CreateQuoteRequest(c *gin.Context) {
	var request models.UserQuoteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user type from path parameter
	userType := models.UserType(c.Param("user_type"))
	request.RequestedBy = userType

	if err := h.operationsService.CreateQuoteRequest(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, request)
}

// GenerateQuoteResponse handles quote response generation
func (h *UserOperationsHandler) GenerateQuoteResponse(c *gin.Context) {
	var request models.UserQuoteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.operationsService.GenerateQuoteResponse(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ConfirmQuote handles quote confirmation by any user type
func (h *UserOperationsHandler) ConfirmQuote(c *gin.Context) {
	var quote models.UserQuoteResponse
	if err := c.ShouldBindJSON(&quote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userType := models.UserType(c.Param("user_type"))
	confirmation, err := h.operationsService.ConfirmQuote(&quote, userType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, confirmation)
}

// CreateBooking handles booking creation
func (h *UserOperationsHandler) CreateBooking(c *gin.Context) {
	var request struct {
		ConfirmationID string                 `json:"confirmation_id"`
		Booking        models.UserBookingRequest `json:"booking"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userType := models.UserType(c.Param("user_type"))
	request.Booking.RequestedBy = userType

	// Get quote confirmation
	confirmation := &models.UserQuoteConfirmation{
		ID: request.ConfirmationID,
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
func (h *UserOperationsHandler) ConfirmBooking(c *gin.Context) {
	var request models.UserBookingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	confirmation, err := h.operationsService.ConfirmBooking(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, confirmation)
}

// GetQuoteRequest retrieves a quote request
func (h *UserOperationsHandler) GetQuoteRequest(c *gin.Context) {
	requestID := c.Param("id")
	userType := models.UserType(c.Param("user_type"))

	// In a real implementation, this would retrieve from storage
	request := &models.UserQuoteRequest{
		ID:          requestID,
		RequestedBy: userType,
		Status:      "NEW",
	}

	c.JSON(http.StatusOK, request)
}

// GetQuoteResponse retrieves a quote response
func (h *UserOperationsHandler) GetQuoteResponse(c *gin.Context) {
	quoteID := c.Param("id")
	userType := models.UserType(c.Param("user_type"))

	// In a real implementation, this would retrieve from storage
	response := &models.UserQuoteResponse{
		ID:          quoteID,
		RespondedBy: models.FreightForwarder,
		Status:      "DRAFT",
	}

	c.JSON(http.StatusOK, response)
}

// GetBooking retrieves a booking
func (h *UserOperationsHandler) GetBooking(c *gin.Context) {
	bookingID := c.Param("id")
	userType := models.UserType(c.Param("user_type"))

	// In a real implementation, this would retrieve from storage
	booking := &models.UserBookingRequest{
		ID:          bookingID,
		RequestedBy: userType,
		Status:      "PENDING",
	}

	c.JSON(http.StatusOK, booking)
}

// ListQuoteRequests lists quote requests for a user
func (h *UserOperationsHandler) ListQuoteRequests(c *gin.Context) {
	userType := models.UserType(c.Param("user_type"))
	status := c.Query("status") // Optional status filter

	// In a real implementation, this would retrieve from storage
	requests := []models.UserQuoteRequest{
		{
			ID:          "QR-1",
			RequestedBy: userType,
			Status:      status,
		},
	}

	c.JSON(http.StatusOK, requests)
}

// ListQuoteResponses lists quote responses for a user
func (h *UserOperationsHandler) ListQuoteResponses(c *gin.Context) {
	userType := models.UserType(c.Param("user_type"))
	status := c.Query("status") // Optional status filter

	// In a real implementation, this would retrieve from storage
	responses := []models.UserQuoteResponse{
		{
			ID:          "QT-1",
			RespondedBy: models.FreightForwarder,
			Status:      status,
		},
	}

	c.JSON(http.StatusOK, responses)
}

// ListBookings lists bookings for a user
func (h *UserOperationsHandler) ListBookings(c *gin.Context) {
	userType := models.UserType(c.Param("user_type"))
	status := c.Query("status") // Optional status filter

	// In a real implementation, this would retrieve from storage
	bookings := []models.UserBookingRequest{
		{
			ID:          "BK-1",
			RequestedBy: userType,
			Status:      status,
		},
	}

	c.JSON(http.StatusOK, bookings)
}
