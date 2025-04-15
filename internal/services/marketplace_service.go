package services

import (
	"fmt"
	"time"

	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/stellar"
)

// MarketplaceService handles business logic for the marketplace
type MarketplaceService struct {
	txManager    *stellar.TransactionManager
	tokenManager *stellar.TokenManager
}

// NewMarketplaceService creates a new MarketplaceService instance
func NewMarketplaceService(txManager *stellar.TransactionManager, tokenManager *stellar.TokenManager) *MarketplaceService {
	return &MarketplaceService{
		txManager:    txManager,
		tokenManager: tokenManager,
	}
}

// GetServicesByCategory retrieves services by main category
func (s *MarketplaceService) GetServicesByCategory(category string) ([]models.LogisticsService, error) {
	// Validate category
	if !s.isValidCategory(category) {
		return nil, fmt.Errorf("invalid category: %s", category)
	}

	// In a real implementation, this would query a database
	// This is a placeholder that would be replaced with actual storage logic
	return []models.LogisticsService{}, nil
}

// GetServiceItems retrieves service items for a subcategory
func (s *MarketplaceService) GetServiceItems(subcategoryID string) ([]models.ServiceItem, error) {
	// In a real implementation, this would query a database
	// This is a placeholder that would be replaced with actual storage logic
	return []models.ServiceItem{}, nil
}

// CreateServiceListing creates a new service listing
func (s *MarketplaceService) CreateServiceListing(service *models.LogisticsService) error {
	if err := s.validateServiceListing(service); err != nil {
		return fmt.Errorf("invalid service listing: %w", err)
	}

	// Generate unique ID
	service.ID = fmt.Sprintf("SVC-%d", time.Now().UnixNano())
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()

	// Submit to blockchain
	result, err := s.txManager.CreateServiceListing(
		service.Provider.ID,
		uint8(s.getCategoryIndex(string(service.Category))),
		0, // shipment mode would come from service details
		service.Origin,
		service.Destination,
		uint64(service.Item.BasePrice * 100), // Convert to smallest unit
		service.Item.Description,
	)
	if err != nil {
		return fmt.Errorf("failed to create service listing: %w", err)
	}

	// Update with blockchain transaction details
	service.ID = result.TxID

	return nil
}

// GetQuotation calculates rate for given service and cargo details
func (s *MarketplaceService) GetQuotation(serviceID string, cargo *models.CargoDetails) (*models.Rate, error) {
	// Calculate base rate and surcharges
	rate := &models.Rate{
		ID:           fmt.Sprintf("RATE-%d", time.Now().UnixNano()),
		ServiceID:    serviceID,
		CargoDetails: *cargo,
		CreatedAt:    time.Now(),
		ValidUntil:   time.Now().Add(24 * time.Hour),
	}

	// Calculate surcharges
	surcharges := s.calculateSurcharges(cargo)
	rate.Surcharges = surcharges

	// Calculate total amount
	var total float64
	for _, surcharge := range surcharges {
		total += surcharge.Amount
	}
	rate.TotalAmount = rate.BaseRate + total

	return rate, nil
}

// CreateBooking creates a new booking
func (s *MarketplaceService) CreateBooking(booking *models.Booking) error {
	if err := s.validateBooking(booking); err != nil {
		return fmt.Errorf("invalid booking: %w", err)
	}

	// Submit to blockchain
	result, err := s.txManager.CreateBooking(
		booking.CustomerID,
		booking.ServiceID,
		map[string]interface{}{
			"weight":      booking.CargoDetails.Weight,
			"volume":      booking.CargoDetails.Volume,
			"type":        booking.CargoDetails.CargoType,
			"description": booking.CargoDetails.Description,
			"hazardous":   booking.CargoDetails.IsHazardous,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	// Update booking with blockchain transaction details
	booking.ID = result.TxID
	booking.Status = "PENDING"
	booking.PaymentStatus = "UNPAID"
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()

	return nil
}

// ProcessPayment handles payment for a booking
func (s *MarketplaceService) ProcessPayment(bookingID string, customerID string, amount float64) error {
	// Convert amount to token amount
	tokenAmount := fmt.Sprintf("%.0f", amount*100) // Convert to smallest unit

	// Process payment on blockchain
	result, err := s.txManager.ProcessPayment(customerID, bookingID, tokenAmount)
	if err != nil {
		return fmt.Errorf("failed to process payment: %w", err)
	}

	return nil
}

// UpdateShipmentStatus updates the status of a shipment
func (s *MarketplaceService) UpdateShipmentStatus(event *models.TrackingEvent) error {
	// Submit status update to blockchain
	result, err := s.txManager.UpdateShipmentStatus(
		event.BookingID,
		event.Location,
		event.Status,
		event.Description,
	)
	if err != nil {
		return fmt.Errorf("failed to update shipment status: %w", err)
	}

	event.ID = result.TxID
	event.Timestamp = time.Now()

	return nil
}

// Helper functions

func (s *MarketplaceService) validateServiceListing(service *models.LogisticsService) error {
	if service.Provider.ID == "" {
		return fmt.Errorf("provider ID is required")
	}
	if service.Category == "" {
		return fmt.Errorf("service category is required")
	}
	if service.Origin == "" {
		return fmt.Errorf("origin is required")
	}
	if service.Destination == "" {
		return fmt.Errorf("destination is required")
	}
	if service.Item.BasePrice <= 0 {
		return fmt.Errorf("base price must be greater than 0")
	}
	return nil
}

func (s *MarketplaceService) validateBooking(booking *models.Booking) error {
	if booking.CustomerID == "" {
		return fmt.Errorf("customer ID is required")
	}
	if booking.ServiceID == "" {
		return fmt.Errorf("service ID is required")
	}
	if booking.CargoDetails.Weight <= 0 {
		return fmt.Errorf("cargo weight must be greater than 0")
	}
	if booking.CargoDetails.Volume <= 0 {
		return fmt.Errorf("cargo volume must be greater than 0")
	}
	return nil
}

func (s *MarketplaceService) isValidCategory(category string) bool {
	switch models.ServiceCategory(category) {
	case models.ImportService, models.ExportService, models.TransitService, models.TransshipService:
		return true
	default:
		return false
	}
}

func (s *MarketplaceService) getCategoryIndex(category string) int {
	switch models.ServiceCategory(category) {
	case models.ImportService:
		return 0
	case models.ExportService:
		return 1
	case models.TransitService:
		return 2
	case models.TransshipService:
		return 3
	default:
		return 0
	}
}

func (s *MarketplaceService) calculateSurcharges(cargo *models.CargoDetails) []models.Charge {
	surcharges := make([]models.Charge, 0)

	// Add fuel surcharge
	surcharges = append(surcharges, models.Charge{
		Type:        "FUEL",
		Description: "Fuel Surcharge",
		Amount:      cargo.Weight * 0.5, // Example calculation
		Currency:    "USD",
	})

	// Add security surcharge for hazardous cargo
	if cargo.IsHazardous {
		surcharges = append(surcharges, models.Charge{
			Type:        "SECURITY",
			Description: "Hazardous Cargo Security Fee",
			Amount:      100.0, // Fixed fee
			Currency:    "USD",
		})
	}

	// Add peak season surcharge if applicable
	if s.isPeakSeason() {
		surcharges = append(surcharges, models.Charge{
			Type:        "PEAK",
			Description: "Peak Season Surcharge",
			Amount:      cargo.Weight * 0.3, // Example calculation
			Currency:    "USD",
		})
	}

	return surcharges
}

func (s *MarketplaceService) isPeakSeason() bool {
	// Example peak season logic
	currentMonth := time.Now().Month()
	return currentMonth >= time.October && currentMonth <= time.December
}
