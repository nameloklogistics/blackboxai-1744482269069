package services

import (
	"fmt"
	"time"

	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/stellar"
)

// ImportService handles import operations across all transport modes
type ImportService struct {
	txManager    *stellar.TransactionManager
	tokenManager *stellar.TokenManager
}

// NewImportService creates a new ImportService instance
func NewImportService(txManager *stellar.TransactionManager, tokenManager *stellar.TokenManager) *ImportService {
	return &ImportService{
		txManager:    txManager,
		tokenManager: tokenManager,
	}
}

// CreateSeaImportService creates a new sea import service
func (s *ImportService) CreateSeaImportService(service *models.SeaImportService) error {
	if err := s.validateSeaImportService(service); err != nil {
		return fmt.Errorf("invalid sea import service: %w", err)
	}

	service.ID = fmt.Sprintf("SIS-%d", time.Now().UnixNano())
	service.Mode = models.SeaTransport
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	service.Status = "ACTIVE"

	return nil
}

// CreateAirImportService creates a new air import service
func (s *ImportService) CreateAirImportService(service *models.AirImportService) error {
	if err := s.validateAirImportService(service); err != nil {
		return fmt.Errorf("invalid air import service: %w", err)
	}

	service.ID = fmt.Sprintf("AIS-%d", time.Now().UnixNano())
	service.Mode = models.AirTransport
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	service.Status = "ACTIVE"

	return nil
}

// CreateRailImportService creates a new rail import service
func (s *ImportService) CreateRailImportService(service *models.RailImportService) error {
	if err := s.validateRailImportService(service); err != nil {
		return fmt.Errorf("invalid rail import service: %w", err)
	}

	service.ID = fmt.Sprintf("RIS-%d", time.Now().UnixNano())
	service.Mode = models.RailTransport
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	service.Status = "ACTIVE"

	return nil
}

// CreateLandImportService creates a new land import service
func (s *ImportService) CreateLandImportService(service *models.LandImportService) error {
	if err := s.validateLandImportService(service); err != nil {
		return fmt.Errorf("invalid land import service: %w", err)
	}

	service.ID = fmt.Sprintf("LIS-%d", time.Now().UnixNano())
	service.Mode = models.RoadTransport
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	service.Status = "ACTIVE"

	return nil
}

// GetImportServiceSchedule retrieves schedule for an import service
func (s *ImportService) GetImportServiceSchedule(serviceID string) (*models.ImportServiceSchedule, error) {
	// In a real implementation, this would retrieve from storage
	return &models.ImportServiceSchedule{
		ServiceID:   serviceID,
		Frequency:   "WEEKLY",
		ValidFrom:   time.Now(),
		ValidUntil:  time.Now().AddDate(0, 1, 0),
	}, nil
}

// GetImportServiceAvailability checks availability for an import service
func (s *ImportService) GetImportServiceAvailability(serviceID string) (*models.ImportServiceAvailability, error) {
	// In a real implementation, this would check real-time availability
	return &models.ImportServiceAvailability{
		ServiceID:      serviceID,
		AvailableSpace: 1000.0,
		NextAvailable:  time.Now().AddDate(0, 0, 1),
	}, nil
}

// GetImportServiceRate retrieves rates for an import service
func (s *ImportService) GetImportServiceRate(serviceID string) (*models.ImportServiceRate, error) {
	// In a real implementation, this would retrieve current rates
	return &models.ImportServiceRate{
		ServiceID:  serviceID,
		BaseRate:   1000.0,
		Currency:   "USD",
		ValidFrom:  time.Now(),
		ValidUntil: time.Now().AddDate(0, 1, 0),
	}, nil
}

// GetImportServiceRequirements retrieves requirements for an import service
func (s *ImportService) GetImportServiceRequirements(serviceType models.ImportServiceType) (*models.ImportServiceRequirement, error) {
	// In a real implementation, this would retrieve from configuration
	return &models.ImportServiceRequirement{
		ServiceType: serviceType,
		Documents:   []string{"Commercial Invoice", "Packing List"},
	}, nil
}

// Helper functions

func (s *ImportService) validateSeaImportService(service *models.SeaImportService) error {
	if service.Origin.Code == "" {
		return fmt.Errorf("origin port is required")
	}
	if service.Destination.Code == "" {
		return fmt.Errorf("destination port is required")
	}
	if service.ShippingLine == "" {
		return fmt.Errorf("shipping line is required")
	}
	if service.Type == models.SeaFCLImport && len(service.ContainerTypes) == 0 {
		return fmt.Errorf("container types are required for FCL service")
	}
	return nil
}

func (s *ImportService) validateAirImportService(service *models.AirImportService) error {
	if service.Origin.Code == "" {
		return fmt.Errorf("origin airport is required")
	}
	if service.Destination.Code == "" {
		return fmt.Errorf("destination airport is required")
	}
	if service.Carrier == "" {
		return fmt.Errorf("carrier is required")
	}
	if len(service.CargoTypes) == 0 {
		return fmt.Errorf("cargo types are required")
	}
	return nil
}

func (s *ImportService) validateRailImportService(service *models.RailImportService) error {
	if service.Origin.Code == "" {
		return fmt.Errorf("origin terminal is required")
	}
	if service.Destination.Code == "" {
		return fmt.Errorf("destination terminal is required")
	}
	if service.RailOperator == "" {
		return fmt.Errorf("rail operator is required")
	}
	if len(service.WagonTypes) == 0 {
		return fmt.Errorf("wagon types are required")
	}
	return nil
}

func (s *ImportService) validateLandImportService(service *models.LandImportService) error {
	if service.Origin.Code == "" {
		return fmt.Errorf("origin location is required")
	}
	if service.Destination.Code == "" {
		return fmt.Errorf("destination location is required")
	}
	if service.TransportOperator == "" {
		return fmt.Errorf("transport operator is required")
	}
	if len(service.VehicleTypes) == 0 {
		return fmt.Errorf("vehicle types are required")
	}
	return nil
}

// UpdateImportServiceSchedule updates schedule for an import service
func (s *ImportService) UpdateImportServiceSchedule(schedule *models.ImportServiceSchedule) error {
	if schedule.ServiceID == "" {
		return fmt.Errorf("service ID is required")
	}
	if schedule.ValidUntil.Before(schedule.ValidFrom) {
		return fmt.Errorf("invalid validity period")
	}
	return nil
}

// UpdateImportServiceRate updates rates for an import service
func (s *ImportService) UpdateImportServiceRate(rate *models.ImportServiceRate) error {
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

// ListImportServices lists all import services by mode
func (s *ImportService) ListImportServices(mode models.TransportMode) ([]models.ImportService, error) {
	// In a real implementation, this would retrieve from storage
	return []models.ImportService{
		{
			ID:   "SIS-1",
			Mode: models.SeaTransport,
			Type: models.SeaFCLImport,
		},
	}, nil
}

// GetImportService retrieves an import service by ID
func (s *ImportService) GetImportService(serviceID string) (*models.ImportService, error) {
	// In a real implementation, this would retrieve from storage
	return &models.ImportService{
		ID:     serviceID,
		Status: "ACTIVE",
	}, nil
}
