package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/services"
)

type ImportServiceHandler struct {
	importService *services.ImportService
}

func NewImportServiceHandler(importService *services.ImportService) *ImportServiceHandler {
	return &ImportServiceHandler{
		importService: importService,
	}
}

// CreateSeaImportService handles creation of sea import services
func (h *ImportServiceHandler) CreateSeaImportService(c *gin.Context) {
	var service models.SeaImportService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.importService.CreateSeaImportService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// CreateAirImportService handles creation of air import services
func (h *ImportServiceHandler) CreateAirImportService(c *gin.Context) {
	var service models.AirImportService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.importService.CreateAirImportService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// CreateRailImportService handles creation of rail import services
func (h *ImportServiceHandler) CreateRailImportService(c *gin.Context) {
	var service models.RailImportService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.importService.CreateRailImportService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// CreateLandImportService handles creation of land import services
func (h *ImportServiceHandler) CreateLandImportService(c *gin.Context) {
	var service models.LandImportService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.importService.CreateLandImportService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// GetImportServiceSchedule handles retrieving service schedules
func (h *ImportServiceHandler) GetImportServiceSchedule(c *gin.Context) {
	serviceID := c.Param("id")

	schedule, err := h.importService.GetImportServiceSchedule(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// GetImportServiceAvailability handles checking service availability
func (h *ImportServiceHandler) GetImportServiceAvailability(c *gin.Context) {
	serviceID := c.Param("id")

	availability, err := h.importService.GetImportServiceAvailability(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, availability)
}

// GetImportServiceRate handles retrieving service rates
func (h *ImportServiceHandler) GetImportServiceRate(c *gin.Context) {
	serviceID := c.Param("id")

	rate, err := h.importService.GetImportServiceRate(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rate)
}

// GetImportServiceRequirements handles retrieving service requirements
func (h *ImportServiceHandler) GetImportServiceRequirements(c *gin.Context) {
	serviceType := models.ImportServiceType(c.Param("type"))

	requirements, err := h.importService.GetImportServiceRequirements(serviceType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requirements)
}

// UpdateImportServiceSchedule handles updating service schedules
func (h *ImportServiceHandler) UpdateImportServiceSchedule(c *gin.Context) {
	var schedule models.ImportServiceSchedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.importService.UpdateImportServiceSchedule(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// UpdateImportServiceRate handles updating service rates
func (h *ImportServiceHandler) UpdateImportServiceRate(c *gin.Context) {
	var rate models.ImportServiceRate
	if err := c.ShouldBindJSON(&rate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.importService.UpdateImportServiceRate(&rate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rate)
}

// ListImportServices handles listing import services by mode
func (h *ImportServiceHandler) ListImportServices(c *gin.Context) {
	mode := models.TransportMode(c.Param("mode"))

	services, err := h.importService.ListImportServices(mode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// GetImportService handles retrieving a specific import service
func (h *ImportServiceHandler) GetImportService(c *gin.Context) {
	serviceID := c.Param("id")

	service, err := h.importService.GetImportService(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

// ListSeaImportServices handles listing sea import services
func (h *ImportServiceHandler) ListSeaImportServices(c *gin.Context) {
	services, err := h.importService.ListImportServices(models.SeaTransport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// ListAirImportServices handles listing air import services
func (h *ImportServiceHandler) ListAirImportServices(c *gin.Context) {
	services, err := h.importService.ListImportServices(models.AirTransport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// ListRailImportServices handles listing rail import services
func (h *ImportServiceHandler) ListRailImportServices(c *gin.Context) {
	services, err := h.importService.ListImportServices(models.RailTransport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// ListLandImportServices handles listing land import services
func (h *ImportServiceHandler) ListLandImportServices(c *gin.Context) {
	services, err := h.importService.ListImportServices(models.RoadTransport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}
