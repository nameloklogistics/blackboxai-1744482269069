package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/services"
)

type ServiceCategoriesHandler struct {
	categoriesService *services.ServiceCategoriesService
}

func NewServiceCategoriesHandler(categoriesService *services.ServiceCategoriesService) *ServiceCategoriesHandler {
	return &ServiceCategoriesHandler{
		categoriesService: categoriesService,
	}
}

// Import Service Handlers

// CreateImportSeaService handles creation of sea import services
func (h *ServiceCategoriesHandler) CreateImportSeaService(c *gin.Context) {
	var service models.SeaService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.Category = models.ImportDirect
	if err := h.categoriesService.CreateSeaService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// CreateImportAirService handles creation of air import services
func (h *ServiceCategoriesHandler) CreateImportAirService(c *gin.Context) {
	var service models.AirService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.Category = models.ImportDirect
	if err := h.categoriesService.CreateAirService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// CreateImportRailService handles creation of rail import services
func (h *ServiceCategoriesHandler) CreateImportRailService(c *gin.Context) {
	var service models.RailService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.Category = models.ImportDirect
	if err := h.categoriesService.CreateRailService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// CreateImportLandService handles creation of land import services
func (h *ServiceCategoriesHandler) CreateImportLandService(c *gin.Context) {
	var service models.LandService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.Category = models.ImportDirect
	if err := h.categoriesService.CreateLandService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// Export Service Handlers

// CreateExportSeaService handles creation of sea export services
func (h *ServiceCategoriesHandler) CreateExportSeaService(c *gin.Context) {
	var service models.SeaService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.Category = models.ExportDirect
	if err := h.categoriesService.CreateSeaService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// CreateExportAirService handles creation of air export services
func (h *ServiceCategoriesHandler) CreateExportAirService(c *gin.Context) {
	var service models.AirService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.Category = models.ExportDirect
	if err := h.categoriesService.CreateAirService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// CreateExportRailService handles creation of rail export services
func (h *ServiceCategoriesHandler) CreateExportRailService(c *gin.Context) {
	var service models.RailService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.Category = models.ExportDirect
	if err := h.categoriesService.CreateRailService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// CreateExportLandService handles creation of land export services
func (h *ServiceCategoriesHandler) CreateExportLandService(c *gin.Context) {
	var service models.LandService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.Category = models.ExportDirect
	if err := h.categoriesService.CreateLandService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// Transit Service Handlers

// AddTransitDetails handles adding transit details to a service
func (h *ServiceCategoriesHandler) AddTransitDetails(c *gin.Context) {
	serviceID := c.Param("id")
	var details models.TransitDetails
	if err := c.ShouldBindJSON(&details); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.categoriesService.AddTransitDetails(serviceID, &details); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, details)
}

// Transshipment Service Handlers

// AddTransshipmentDetails handles adding transshipment details to a service
func (h *ServiceCategoriesHandler) AddTransshipmentDetails(c *gin.Context) {
	serviceID := c.Param("id")
	var details models.TransshipmentDetails
	if err := c.ShouldBindJSON(&details); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.categoriesService.AddTransshipmentDetails(serviceID, &details); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, details)
}

// Common Service Operations

// GetServiceSchedule handles retrieving service schedules
func (h *ServiceCategoriesHandler) GetServiceSchedule(c *gin.Context) {
	serviceID := c.Param("id")

	schedule, err := h.categoriesService.GetServiceSchedule(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// GetServiceAvailability handles checking service availability
func (h *ServiceCategoriesHandler) GetServiceAvailability(c *gin.Context) {
	serviceID := c.Param("id")

	availability, err := h.categoriesService.GetServiceAvailability(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, availability)
}

// GetServiceRate handles retrieving service rates
func (h *ServiceCategoriesHandler) GetServiceRate(c *gin.Context) {
	serviceID := c.Param("id")

	rate, err := h.categoriesService.GetServiceRate(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rate)
}

// GetServiceRequirements handles retrieving service requirements
func (h *ServiceCategoriesHandler) GetServiceRequirements(c *gin.Context) {
	category := models.ServiceCategory(c.Param("category"))
	mode := models.TransportMode(c.Param("mode"))

	requirements, err := h.categoriesService.GetServiceRequirements(category, mode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requirements)
}

// ListServices handles listing services by category and mode
func (h *ServiceCategoriesHandler) ListServices(c *gin.Context) {
	category := models.ServiceCategory(c.Param("category"))
	mode := models.TransportMode(c.Param("mode"))

	services, err := h.categoriesService.ListServices(category, mode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// GetService handles retrieving a specific service
func (h *ServiceCategoriesHandler) GetService(c *gin.Context) {
	serviceID := c.Param("id")

	service, err := h.categoriesService.GetService(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

// UpdateServiceSchedule handles updating service schedules
func (h *ServiceCategoriesHandler) UpdateServiceSchedule(c *gin.Context) {
	var schedule models.ServiceSchedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.categoriesService.UpdateServiceSchedule(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// UpdateServiceRate handles updating service rates
func (h *ServiceCategoriesHandler) UpdateServiceRate(c *gin.Context) {
	var rate models.ServiceRate
	if err := c.ShouldBindJSON(&rate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.categoriesService.UpdateServiceRate(&rate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rate)
}
