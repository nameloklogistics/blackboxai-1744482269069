package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/services"
)

type InfrastructureHandler struct {
	infrastructureService *services.InfrastructureService
}

func NewInfrastructureHandler(infrastructureService *services.InfrastructureService) *InfrastructureHandler {
	return &InfrastructureHandler{
		infrastructureService: infrastructureService,
	}
}

// CreateAirport handles creation of new international airports
func (h *InfrastructureHandler) CreateAirport(c *gin.Context) {
	var airport models.Airport
	if err := c.ShouldBindJSON(&airport); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.infrastructureService.CreateAirport(&airport); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, airport)
}

// CreateSeaport handles creation of new seaports
func (h *InfrastructureHandler) CreateSeaport(c *gin.Context) {
	var seaport models.Seaport
	if err := c.ShouldBindJSON(&seaport); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.infrastructureService.CreateSeaport(&seaport); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, seaport)
}

// CreateInlandDepot handles creation of new inland container depots
func (h *InfrastructureHandler) CreateInlandDepot(c *gin.Context) {
	var depot models.InlandDepot
	if err := c.ShouldBindJSON(&depot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.infrastructureService.CreateInlandDepot(&depot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, depot)
}

// GetCountryInfrastructure handles retrieving all infrastructure for a country
func (h *InfrastructureHandler) GetCountryInfrastructure(c *gin.Context) {
	countryCode := c.Param("country_code")
	
	infrastructure, err := h.infrastructureService.GetCountryInfrastructure(countryCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, infrastructure)
}

// GetInfrastructureServices handles retrieving services at a location
func (h *InfrastructureHandler) GetInfrastructureServices(c *gin.Context) {
	infraID := c.Param("id")
	mode := models.TransportMode(c.Query("mode"))

	services, err := h.infrastructureService.GetInfrastructureServices(infraID, mode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// UpdateInfrastructureCapacity handles updating capacity status
func (h *InfrastructureHandler) UpdateInfrastructureCapacity(c *gin.Context) {
	infraID := c.Param("id")
	var capacity models.InfrastructureCapacity
	if err := c.ShouldBindJSON(&capacity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	capacity.InfrastructureID = infraID
	if err := h.infrastructureService.UpdateInfrastructureCapacity(infraID, &capacity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, capacity)
}

// UpdateInfrastructureSchedule handles updating operating schedules
func (h *InfrastructureHandler) UpdateInfrastructureSchedule(c *gin.Context) {
	infraID := c.Param("id")
	var schedule models.InfrastructureSchedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule.InfrastructureID = infraID
	if err := h.infrastructureService.UpdateInfrastructureSchedule(infraID, &schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// GetServiceLocations handles retrieving locations where a service is offered
func (h *InfrastructureHandler) GetServiceLocations(c *gin.Context) {
	serviceID := c.Param("service_id")
	
	locations, err := h.infrastructureService.GetServiceLocations(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, locations)
}

// GetNearbyInfrastructure handles finding infrastructure within a radius
func (h *InfrastructureHandler) GetNearbyInfrastructure(c *gin.Context) {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid latitude"})
		return
	}

	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid longitude"})
		return
	}

	radius, err := strconv.ParseFloat(c.Query("radius"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid radius"})
		return
	}

	infraType := models.InfrastructureType(c.Query("type"))

	infrastructure, err := h.infrastructureService.GetNearbyInfrastructure(lat, lng, radius, infraType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, infrastructure)
}

// GetInfrastructureCapacity handles retrieving current capacity status
func (h *InfrastructureHandler) GetInfrastructureCapacity(c *gin.Context) {
	infraID := c.Param("id")
	
	capacity, err := h.infrastructureService.GetInfrastructureCapacity(infraID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, capacity)
}

// GetOperatingSchedule handles retrieving current operating schedule
func (h *InfrastructureHandler) GetOperatingSchedule(c *gin.Context) {
	infraID := c.Param("id")
	
	schedule, err := h.infrastructureService.GetOperatingSchedule(infraID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}
