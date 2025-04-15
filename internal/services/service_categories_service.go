package services

import (
	"fmt"
	"time"

	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/stellar"
)

// ServiceCategoriesService handles operations for all service categories
type ServiceCategoriesService struct {
	txManager    *stellar.TransactionManager
	tokenManager *stellar.TokenManager
}

// NewServiceCategoriesService creates a new ServiceCategoriesService instance
func NewServiceCategoriesService(txManager *stellar.TransactionManager, tokenManager *stellar.TokenManager) *ServiceCategoriesService {
	return &ServiceCategoriesService{
		txManager:    txManager,
		tokenManager: tokenManager,
	}
}

// CreateSeaService creates a new sea service for any category
func (s *ServiceCategoriesService) CreateSeaService(service *models.SeaService) error {
	if err := s.validateSeaService(service); err != nil {
		return fmt.Errorf("invalid sea service: %w", err)
	}

	service.ID = fmt.Sprintf("SS-%d", time.Now().UnixNano())
	service.Mode = models.SeaTransport
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	service.Status = "ACTIVE"

	return nil
}

// CreateAirService creates a new air service for any category
func (s *ServiceCategoriesService) CreateAirService(service *models.AirService) error {
	if err := s.validateAirService(service); err != nil {
		return fmt.Errorf("invalid air service: %w", err)
	}

	service.ID = fmt.Sprintf("AS-%d", time.Now().UnixNano())
	service.Mode = models.AirTransport
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	service.Status = "ACTIVE"

	return nil
}

// CreateRailService creates a new rail service for any category
func (s *ServiceCategoriesService) CreateRailService(service *models.RailService) error {
	if err := s.validateRailService(service); err != nil {
		return fmt.Errorf("invalid rail service: %w", err)
	}

	service.ID = fmt.Sprintf("RS-%d", time.Now().UnixNano())
	service.Mode = models.RailTransport
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	service.Status = "ACTIVE"

	return nil
}

// CreateLandService creates a new land service for any category
func (s *ServiceCategoriesService) CreateLandService(service *models.LandService) error {
	if err := s.validateLandService(service); err != nil {
		return fmt.Errorf("invalid land service: %w", err)
	}

	service.ID = fmt.Sprintf("LS-%d", time.Now().UnixNano())
	service.Mode = models.LandTransport
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	service.Status = "ACTIVE"

	return nil
}

// AddTransitDetails adds transit details to a service
func (s *ServiceCategoriesService) AddTransitDetails(serviceID string, details *models.TransitDetails) error {
	if err := s.validateTransitDetails(details); err != nil {
		return fmt.Errorf("invalid transit details: %w", err)
	}
	return nil
}

// AddTransshipmentDetails adds transshipment details to a service
func (s *ServiceCategoriesService) AddTransshipmentDetails(serviceID string, details *models.TransshipmentDetails) error {
	if err := s.validateTransshipmentDetails(details); err != nil {
		return fmt.Errorf("invalid transshipment details: %w", err)
	}
	return nil
}

// GetServiceSchedule retrieves schedule for a service
func (s *ServiceCategoriesService) GetServiceSchedule(serviceID string) (*models.ServiceSchedule, error) {
	// In a real implementation, this would retrieve from storage
	return &models.ServiceSchedule{
		ServiceID:   serviceID,
		Frequency:   "WEEKLY",
		ValidFrom:   time.Now(),
		ValidUntil:  time.Now().AddDate(0, 1, 0),
	}, nil
}

// GetServiceAvailability checks availability for a service
func (s *ServiceCategoriesService) GetServiceAvailability(serviceID string) (*models.ServiceAvailability, error) {
	// In a real implementation, this would check real-time availability
	return &models.ServiceAvailability{
		ServiceID:      serviceID,
		AvailableSpace: 1000.0,
		NextAvailable:  time.Now().AddDate(0, 0, 1),
	}, nil
}

// GetServiceRate retrieves rates for a service
func (s *ServiceCategoriesService) GetServiceRate(serviceID string) (*models.ServiceRate, error) {
	// In a real implementation, this would retrieve current rates
	return &models.ServiceRate{
		ServiceID:  serviceID,
		BaseRate:   1000.0,
		Currency:   "USD",
		ValidFrom:  time.Now(),
		ValidUntil: time.Now().AddDate(0, 1, 0),
	}, nil
}

// GetServiceRequirements retrieves requirements for a service
func (s *ServiceCategoriesService) GetServiceRequirements(category models.ServiceCategory, mode models.TransportMode) (*models.ServiceRequirement, error) {
	// In a real implementation, this would retrieve from configuration
	return &models.ServiceRequirement{
		Category:  category,
		Mode:     mode,
		Documents: []string{"Commercial Invoice", "Packing List"},
	}, nil
}

// ListServices lists all services by category and mode
func (s *ServiceCategoriesService) ListServices(category models.ServiceCategory, mode models.TransportMode) ([]models.BaseService, error) {
	// In a real implementation, this would retrieve from storage
	return []models.BaseService{
		{
			ID:       "SS-1",
			Category: category,
			Mode:     mode,
			Status:   "ACTIVE",
		},
	}, nil
}

// Helper functions

func (s *ServiceCategoriesService) validateSeaService(service *models.SeaService) error {
	if service.Origin.Code == "" {
		return fmt.Errorf("origin port is required")
	}
	if service.Destination.Code == "" {
		return fmt.Errorf("destination port is required")
	}
	if service.ShippingLine == "" {
		return fmt.Errorf("shipping line is required")
	}
	return nil
}

func (s *ServiceCategoriesService) validateAirService(service *models.AirService) error {
	if service.Origin.Code == "" {
		return fmt.Errorf("origin airport is required")
	}
	if service.Destination.Code == "" {
		return fmt.Errorf("destination airport is required")
	}
	if service.Carrier == "" {
		return fmt.Errorf("carrier is required")
	}
	return nil
}

func (s *ServiceCategoriesService) validateRailService(service *models.RailService) error {
	if service.Origin.Code == "" {
		return fmt.Errorf("origin terminal is required")
	}
	if service.Destination.Code == "" {
		return fmt.Errorf("destination terminal is required")
	}
	if service.RailOperator == "" {
		return fmt.Errorf("rail operator is required")
	}
	return nil
}

func (s *ServiceCategoriesService) validateLandService(service *models.LandService) error {
	if service.Origin.Code == "" {
		return fmt.Errorf("origin location is required")
	}
	if service.Destination.Code == "" {
		return fmt.Errorf("destination location is required")
	}
	if service.TransportOperator == "" {
		return fmt.Errorf("transport operator is required")
	}
	return nil
}

func (s *ServiceCategoriesService) validateTransitDetails(details *models.TransitDetails) error {
	if len(details.TransitCountries) == 0 {
		return fmt.Errorf("transit countries are required")
	}
	if len(details.CustomsRegimes) == 0 {
		return fmt.Errorf("customs regimes are required")
	}
	return nil
}

func (s *ServiceCategoriesService) validateTransshipmentDetails(details *models.TransshipmentDetails) error {
	if details.TransshipmentPort.Code == "" {
		return fmt.Errorf("transshipment port is required")
	}
	if details.ConnectionTime == "" {
		return fmt.Errorf("connection time is required")
	}
	return nil
}

// UpdateServiceSchedule updates schedule for a service
func (s *ServiceCategoriesService) UpdateServiceSchedule(schedule *models.ServiceSchedule) error {
	if schedule.ServiceID == "" {
		return fmt.Errorf("service ID is required")
	}
	if schedule.ValidUntil.Before(schedule.ValidFrom) {
		return fmt.Errorf("invalid validity period")
	}
	return nil
}

// UpdateServiceRate updates rates for a service
func (s *ServiceCategoriesService) UpdateServiceRate(rate *models.ServiceRate) error {
	if rate.ServiceID == "" {
		return fmt.Errorf("service ID is required")
	}
	if rate.BaseRate <= 0 {
		return fmt.Errorf("base rate must be greater than 0")
	}
	if rate.ValidUntil.Before(rate.ValidFrom) {
		return fmt.Errorf("invalid validity period")
	}
	return nil
}

// GetService retrieves a service by ID
func (s *ServiceCategoriesService) GetService(serviceID string) (*models.BaseService, error) {
	// In a real implementation, this would retrieve from storage
	return &models.BaseService{
		ID:     serviceID,
		Status: "ACTIVE",
	}, nil
}
