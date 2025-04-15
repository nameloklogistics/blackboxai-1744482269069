package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/services"
)

type CustomsRateHandler struct {
	customsRateService *services.CustomsRateService
	membershipService  *services.MembershipService
}

func NewCustomsRateHandler(
	customsRateService *services.CustomsRateService,
	membershipService *services.MembershipService,
) *CustomsRateHandler {
	return &CustomsRateHandler{
		customsRateService: customsRateService,
		membershipService:  membershipService,
	}
}

// CreateCustomsRate handles creation of new customs rates
func (h *CustomsRateHandler) CreateCustomsRate(c *gin.Context) {
	var rate models.CustomsRate
	if err := c.ShouldBindJSON(&rate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get broker ID from authenticated user
	brokerID := c.GetString("user_id")
	rate.BrokerID = brokerID

	// Create rate with context
	if err := h.customsRateService.CreateCustomsRate(c, &rate); err != nil {
		switch err := err.(type) {
		case *models.MembershipError:
			c.JSON(http.StatusPaymentRequired, gin.H{
				"error": err.Error(),
				"code":  "MEMBERSHIP_REQUIRED",
			})
		case *models.LicenseVerificationError:
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
				"code":  "LICENSE_INVALID",
			})
		case *models.CustomsBrokerRestrictionError:
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
				"code":  "COUNTRY_RESTRICTED",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, rate)
}

// GetCustomsQuotation handles customs clearance quotation requests
func (h *CustomsRateHandler) GetCustomsQuotation(c *gin.Context) {
	var request struct {
		BrokerID        string              `json:"broker_id"`
		BranchOfficeID  string              `json:"branch_office_id"`
		RateID          string              `json:"rate_id"`
		CargoDetails    models.CargoDetails `json:"cargo_details"`
		CargoValue      float64             `json:"cargo_value"`
		TransportMode   models.TransportMode `json:"transport_mode"`
		PackagingMode   models.PackagingMode `json:"packaging_mode"`
		PackagingDetails interface{}         `json:"packaging_details"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quotation, err := h.customsRateService.GetQuotation(
		c,
		request.BrokerID,
		request.BranchOfficeID,
		request.RateID,
		&request.CargoDetails,
		request.CargoValue,
		request.TransportMode,
		request.PackagingMode,
		request.PackagingDetails,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quotation)
}

// GetRatesByType handles retrieving customs rates by type
func (h *CustomsRateHandler) GetRatesByType(c *gin.Context) {
	rateType := models.CustomsRateType(c.Param("type"))
	
	// Get broker ID from authenticated user
	brokerID := c.GetString("user_id")

	// Validate membership
	if err := h.membershipService.ValidateMembership(c, models.CustomsBrokerMembership, brokerID); err != nil {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"error": err.Error(),
			"code":  "MEMBERSHIP_REQUIRED",
		})
		return
	}

	rates, err := h.customsRateService.GetCustomsRatesByType(rateType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rates)
}

// UpdateCustomsRate handles updates to customs rates
func (h *CustomsRateHandler) UpdateCustomsRate(c *gin.Context) {
	rateID := c.Param("id")
	var rate models.CustomsRate
	if err := c.ShouldBindJSON(&rate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify broker ID matches authenticated user
	brokerID := c.GetString("user_id")
	if rate.BrokerID != brokerID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	rate.ID = rateID
	if err := h.customsRateService.UpdateCustomsRate(c, &rate); err != nil {
		switch err := err.(type) {
		case *models.MembershipError:
			c.JSON(http.StatusPaymentRequired, gin.H{
				"error": err.Error(),
				"code":  "MEMBERSHIP_REQUIRED",
			})
		case *models.LicenseVerificationError:
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
				"code":  "LICENSE_INVALID",
			})
		case *models.CustomsBrokerRestrictionError:
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
				"code":  "COUNTRY_RESTRICTED",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, rate)
}

// GetRateComparison handles comparing rates across providers
func (h *CustomsRateHandler) GetRateComparison(c *gin.Context) {
	var request struct {
		RateType         models.CustomsRateType `json:"rate_type"`
		Country          string                 `json:"country"`
		CargoDetails     models.CargoDetails    `json:"cargo_details"`
		CargoValue       float64                `json:"cargo_value"`
		TransportMode    models.TransportMode   `json:"transport_mode"`
		PackagingMode    models.PackagingMode   `json:"packaging_mode"`
		PackagingDetails interface{}            `json:"packaging_details"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get broker ID from authenticated user
	brokerID := c.GetString("user_id")

	// Validate membership
	if err := h.membershipService.ValidateMembership(c, models.CustomsBrokerMembership, brokerID); err != nil {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"error": err.Error(),
			"code":  "MEMBERSHIP_REQUIRED",
		})
		return
	}

	comparisons, err := h.customsRateService.GetRateComparison(
		request.RateType,
		request.Country,
		&request.CargoDetails,
		request.CargoValue,
		request.TransportMode,
		request.PackagingMode,
		request.PackagingDetails,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comparisons)
}

// GetHistoricalRates handles retrieving historical rate data
func (h *CustomsRateHandler) GetHistoricalRates(c *gin.Context) {
	rateType := models.CustomsRateType(c.Query("type"))
	country := c.Query("country")
	
	// Get broker ID from authenticated user
	brokerID := c.GetString("user_id")

	// Validate membership
	if err := h.membershipService.ValidateMembership(c, models.CustomsBrokerMembership, brokerID); err != nil {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"error": err.Error(),
			"code":  "MEMBERSHIP_REQUIRED",
		})
		return
	}

	var request struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse(time.RFC3339, request.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, request.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format"})
		return
	}

	rates, err := h.customsRateService.GetHistoricalRates(
		rateType,
		country,
		startDate,
		endDate,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rates)
}

// ValidateHSCode handles HS code validation and duty rate retrieval
func (h *CustomsRateHandler) ValidateHSCode(c *gin.Context) {
	hsCode := c.Query("hs_code")
	country := c.Query("country")

	if hsCode == "" || country == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "HS code and country are required"})
		return
	}

	// Get broker ID from authenticated user
	brokerID := c.GetString("user_id")

	// Validate membership
	if err := h.membershipService.ValidateMembership(c, models.CustomsBrokerMembership, brokerID); err != nil {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"error": err.Error(),
			"code":  "MEMBERSHIP_REQUIRED",
		})
		return
	}

	dutyRate, err := h.customsRateService.ValidateHSCode(hsCode, country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hs_code":    hsCode,
		"country":    country,
		"duty_rate":  dutyRate,
		"validated":  true,
	})
}
