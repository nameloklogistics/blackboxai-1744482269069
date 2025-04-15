package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/services"
)

type TrackingHandler struct {
	trackingService *services.TrackingService
}

func NewTrackingHandler(trackingService *services.TrackingService) *TrackingHandler {
	return &TrackingHandler{
		trackingService: trackingService,
	}
}

// AddTrackingEvent handles adding new tracking events
func (h *TrackingHandler) AddTrackingEvent(c *gin.Context) {
	var event models.TrackingEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify provider authorization
	providerID := c.GetString("user_id")
	if !h.isAuthorizedProvider(providerID, event.BookingID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized provider"})
		return
	}

	if err := h.trackingService.AddTrackingEvent(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// AddTransshipmentPoint handles adding new transshipment points
func (h *TrackingHandler) AddTransshipmentPoint(c *gin.Context) {
	var point models.TransshipmentPoint
	if err := c.ShouldBindJSON(&point); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify provider authorization
	providerID := c.GetString("user_id")
	if !h.isAuthorizedProvider(providerID, point.BookingID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized provider"})
		return
	}

	if err := h.trackingService.AddTransshipmentPoint(&point); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, point)
}

// UpdateRouting handles updating shipment routing
func (h *TrackingHandler) UpdateRouting(c *gin.Context) {
	bookingID := c.Param("booking_id")
	var points []models.TransshipmentPoint
	if err := c.ShouldBindJSON(&points); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify provider authorization
	providerID := c.GetString("user_id")
	if !h.isAuthorizedProvider(providerID, bookingID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized provider"})
		return
	}

	if err := h.trackingService.UpdateRouting(bookingID, points); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "routing updated"})
}

// GetShipmentTracking handles retrieving tracking history
func (h *TrackingHandler) GetShipmentTracking(c *gin.Context) {
	bookingID := c.Param("booking_id")

	// Verify user authorization (both customer and provider should be able to track)
	userID := c.GetString("user_id")
	if !h.isAuthorizedUser(userID, bookingID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
		return
	}

	events, err := h.trackingService.GetShipmentTracking(bookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetOptimalRoute handles calculating optimal route
func (h *TrackingHandler) GetOptimalRoute(c *gin.Context) {
	var request struct {
		Origin        string            `json:"origin"`
		Destination   string            `json:"destination"`
		ShipmentMode  string            `json:"shipment_mode"`
		CargoDetails  models.CargoDetails `json:"cargo_details"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := h.trackingService.GetOptimalRoute(
		request.Origin,
		request.Destination,
		request.ShipmentMode,
		request.CargoDetails,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate ETA
	eta := h.trackingService.CalculateETA(route, request.ShipmentMode)

	response := gin.H{
		"route": route,
		"eta":   eta,
	}

	c.JSON(http.StatusOK, response)
}

// Helper functions

func (h *TrackingHandler) isAuthorizedProvider(providerID, bookingID string) bool {
	// In a real implementation, this would check if the provider is associated with the booking
	// This is a simplified example
	return true
}

func (h *TrackingHandler) isAuthorizedUser(userID, bookingID string) bool {
	// In a real implementation, this would check if the user is either the customer or provider
	// This is a simplified example
	return true
}
