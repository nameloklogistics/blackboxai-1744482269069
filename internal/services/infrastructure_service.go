package services

import (
	"fmt"
	"time"

	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/stellar"
)

// InfrastructureService handles infrastructure management
type InfrastructureService struct {
	txManager    *stellar.TransactionManager
	tokenManager *stellar.TokenManager
}

// NewInfrastructureService creates a new InfrastructureService instance
func NewInfrastructureService(txManager *stellar.TransactionManager, tokenManager *stellar.TokenManager) *InfrastructureService {
	return &InfrastructureService{
		txManager:    txManager,
		tokenManager: tokenManager,
	}
}

// CreateAirport creates a new international airport
func (s *InfrastructureService) CreateAirport(airport *models.Airport) error {
	if err := s.validateAirport(airport); err != nil {
		return fmt.Errorf("invalid airport: %w", err)
	}

	airport.Type = models.InternationalAirport
	airport.ID = fmt.Sprintf("APT-%d", time.Now().UnixNano())
	airport.CreatedAt = time.Now()
	airport.UpdatedAt = time.Now()

	return nil
}

// CreateSeaport creates a new seaport
func (s *InfrastructureService) CreateSeaport(seaport *models.Seaport) error {
	if err := s.validateSeaport(seaport); err != nil {
		return fmt.Errorf("invalid seaport: %w", err)
	}

	seaport.Type = models.Seaport
	seaport.ID = fmt.Sprintf("SPT-%d", time.Now().UnixNano())
	seaport.CreatedAt = time.Now()
	seaport.UpdatedAt = time.Now()

	return nil
}

// CreateInlandDepot creates a new inland container depot
func (s *InfrastructureService) CreateInlandDepot(depot *models.InlandDepot) error {
	if err := s.validateInlandDepot(depot); err != nil {
		return fmt.Errorf("invalid inland depot: %w", err)
	}

	depot.Type = models.InlandDepot
	depot.ID = fmt.Sprintf("ICD-%d", time.Now().UnixNano())
	depot.CreatedAt = time.Now()
	depot.UpdatedAt = time.Now()

	return nil
}

// GetCountryInfrastructure retrieves all infrastructure for a country
func (s *InfrastructureService) GetCountryInfrastructure(countryCode string) (*models.CountryInfrastructure, error) {
	// In a real implementation, this would query a database
	// This is a placeholder that would be replaced with actual storage logic
	return &models.CountryInfrastructure{
		Country:     "Sample Country",
		CountryCode: countryCode,
		UpdatedAt:   time.Now(),
	}, nil
}

// GetInfrastructureServices retrieves services available at a location
func (s *InfrastructureService) GetInfrastructureServices(
	infraID string,
	mode models.TransportMode,
) ([]models.TransportService, error) {
	// In a real implementation, this would query available services
	return []models.TransportService{}, nil
}

// UpdateInfrastructureCapacity updates current capacity status
func (s *InfrastructureService) UpdateInfrastructureCapacity(
	infraID string,
	capacity *models.InfrastructureCapacity,
) error {
	if capacity.StorageUtilization < 0 || capacity.StorageUtilization > 100 {
		return fmt.Errorf("invalid storage utilization percentage")
	}

	capacity.LastUpdated = time.Now()
	return nil
}

// UpdateInfrastructureSchedule updates operating schedules
func (s *InfrastructureService) UpdateInfrastructureSchedule(
	infraID string,
	schedule *models.InfrastructureSchedule,
) error {
	if err := s.validateSchedule(schedule); err != nil {
		return fmt.Errorf("invalid schedule: %w", err)
	}

	return nil
}

// GetServiceLocations retrieves all locations where a service is offered
func (s *InfrastructureService) GetServiceLocations(
	serviceID string,
) ([]models.ServiceLocation, error) {
	// In a real implementation, this would query service locations
	return []models.ServiceLocation{}, nil
}

// Helper functions

func (s *InfrastructureService) validateAirport(airport *models.Airport) error {
	if airport.IATACode == "" {
		return fmt.Errorf("IATA code is required")
	}
	if airport.Name == "" {
		return fmt.Errorf("airport name is required")
	}
	if airport.Country == "" {
		return fmt.Errorf("country is required")
	}
	if airport.CargoCapacityTons <= 0 {
		return fmt.Errorf("cargo capacity must be greater than 0")
	}
	return nil
}

func (s *InfrastructureService) validateSeaport(seaport *models.Seaport) error {
	if seaport.UNLOCode == "" {
		return fmt.Errorf("UN/LOCODE is required")
	}
	if seaport.Name == "" {
		return fmt.Errorf("seaport name is required")
	}
	if seaport.Country == "" {
		return fmt.Errorf("country is required")
	}
	if seaport.ContainerCapacity <= 0 {
		return fmt.Errorf("container capacity must be greater than 0")
	}
	return nil
}

func (s *InfrastructureService) validateInlandDepot(depot *models.InlandDepot) error {
	if depot.DepotCode == "" {
		return fmt.Errorf("depot code is required")
	}
	if depot.Name == "" {
		return fmt.Errorf("depot name is required")
	}
	if depot.Country == "" {
		return fmt.Errorf("country is required")
	}
	if depot.StorageCapacity <= 0 {
		return fmt.Errorf("storage capacity must be greater than 0")
	}
	return nil
}

func (s *InfrastructureService) validateSchedule(schedule *models.InfrastructureSchedule) error {
	if len(schedule.OperatingDays) == 0 {
		return fmt.Errorf("operating days are required")
	}
	if schedule.OperatingHours == "" {
		return fmt.Errorf("operating hours are required")
	}
	if schedule.ValidUntil.Before(schedule.ValidFrom) {
		return fmt.Errorf("invalid schedule validity period")
	}
	return nil
}

// GetNearbyInfrastructure finds infrastructure points within a radius
func (s *InfrastructureService) GetNearbyInfrastructure(
	latitude float64,
	longitude float64,
	radiusKm float64,
	infraType models.InfrastructureType,
) ([]models.Infrastructure, error) {
	// In a real implementation, this would use geospatial queries
	return []models.Infrastructure{}, nil
}

// GetInfrastructureCapacity gets current capacity status
func (s *InfrastructureService) GetInfrastructureCapacity(
	infraID string,
) (*models.InfrastructureCapacity, error) {
	// In a real implementation, this would query current capacity
	return &models.InfrastructureCapacity{
		InfrastructureID:   infraID,
		StorageUtilization: 75.5,
		EquipmentStatus:    map[string]string{"crane1": "operational"},
		LastUpdated:        time.Now(),
	}, nil
}

// GetOperatingSchedule gets current operating schedule
func (s *InfrastructureService) GetOperatingSchedule(
	infraID string,
) (*models.InfrastructureSchedule, error) {
	// In a real implementation, this would query current schedule
	return &models.InfrastructureSchedule{
		InfrastructureID: infraID,
		OperatingDays:    []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
		OperatingHours:   "08:00-18:00",
		CustomsHours:     "09:00-17:00",
		ValidFrom:        time.Now(),
		ValidUntil:       time.Now().AddDate(0, 1, 0),
	}, nil
}
