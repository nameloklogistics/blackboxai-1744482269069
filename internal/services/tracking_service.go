package services

import (
	"fmt"
	"time"

	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/stellar"
)

// TrackingService handles shipment tracking and routing operations
type TrackingService struct {
	txManager    *stellar.TransactionManager
	tokenManager *stellar.TokenManager
}

// NewTrackingService creates a new TrackingService instance
func NewTrackingService(txManager *stellar.TransactionManager, tokenManager *stellar.TokenManager) *TrackingService {
	return &TrackingService{
		txManager:    txManager,
		tokenManager: tokenManager,
	}
}

// AddTrackingEvent adds a new tracking event for a shipment
func (s *TrackingService) AddTrackingEvent(event *models.TrackingEvent) error {
	if err := s.validateTrackingEvent(event); err != nil {
		return fmt.Errorf("invalid tracking event: %w", err)
	}

	// Update shipment status on blockchain
	result, err := s.txManager.UpdateShipmentStatus(
		event.BookingID,
		event.Location,
		event.Status,
		event.Description,
	)
	if err != nil {
		return fmt.Errorf("failed to update tracking status: %w", err)
	}

	event.ID = result.TxID
	event.Timestamp = time.Now()
	return nil
}

// AddTransshipmentPoint adds a new transshipment point to the route
func (s *TrackingService) AddTransshipmentPoint(point *models.TransshipmentPoint) error {
	if err := s.validateTransshipmentPoint(point); err != nil {
		return fmt.Errorf("invalid transshipment point: %w", err)
	}

	// Create tracking event for transshipment
	event := &models.TrackingEvent{
		BookingID:   point.BookingID,
		Location:    point.Location,
		Status:      fmt.Sprintf("TRANSSHIPMENT_%s", point.Status),
		Description: fmt.Sprintf("Transshipment at %s via %s", point.Location, point.Carrier),
		Timestamp:   time.Now(),
	}

	if err := s.AddTrackingEvent(event); err != nil {
		return fmt.Errorf("failed to add transshipment event: %w", err)
	}

	return nil
}

// UpdateRouting updates the routing information for a shipment
func (s *TrackingService) UpdateRouting(
	bookingID string,
	transshipmentPoints []models.TransshipmentPoint,
) error {
	for _, point := range transshipmentPoints {
		point.BookingID = bookingID
		if err := s.AddTransshipmentPoint(&point); err != nil {
			return fmt.Errorf("failed to update routing: %w", err)
		}
	}

	return nil
}

// GetShipmentTracking retrieves tracking history for a shipment
func (s *TrackingService) GetShipmentTracking(bookingID string) ([]models.TrackingEvent, error) {
	// In a real implementation, this would query a database
	// This is a placeholder that would be replaced with actual storage logic
	return []models.TrackingEvent{}, nil
}

// CalculateETA calculates estimated time of arrival based on route and conditions
func (s *TrackingService) CalculateETA(
	transshipmentPoints []models.TransshipmentPoint,
	shipmentMode string,
) time.Time {
	var totalDuration time.Duration

	// Calculate transit time between points
	for i := 0; i < len(transshipmentPoints)-1; i++ {
		current := transshipmentPoints[i]
		next := transshipmentPoints[i+1]

		// Add transit time between points
		transitTime := s.calculateTransitTime(
			current.Location,
			next.Location,
			shipmentMode,
		)
		totalDuration += transitTime

		// Add processing time at transshipment point
		totalDuration += s.calculateProcessingTime(next.Location, shipmentMode)
	}

	return time.Now().Add(totalDuration)
}

// Helper functions

func (s *TrackingService) validateTrackingEvent(event *models.TrackingEvent) error {
	if event.BookingID == "" {
		return fmt.Errorf("booking ID is required")
	}
	if event.Location == "" {
		return fmt.Errorf("location is required")
	}
	if event.Status == "" {
		return fmt.Errorf("status is required")
	}
	return nil
}

func (s *TrackingService) validateTransshipmentPoint(point *models.TransshipmentPoint) error {
	if point.BookingID == "" {
		return fmt.Errorf("booking ID is required")
	}
	if point.Location == "" {
		return fmt.Errorf("location is required")
	}
	if point.Carrier == "" {
		return fmt.Errorf("carrier is required")
	}
	return nil
}

func (s *TrackingService) calculateTransitTime(
	origin string,
	destination string,
	mode string,
) time.Duration {
	// In a real implementation, this would use a routing algorithm and historical data
	// This is a simplified example
	baseTime := 24 * time.Hour // 1 day base transit time

	switch mode {
	case "Air":
		return baseTime
	case "Sea":
		return baseTime * 7 // 7 days
	case "Road":
		return baseTime * 2 // 2 days
	case "Rail":
		return baseTime * 3 // 3 days
	default:
		return baseTime
	}
}

func (s *TrackingService) calculateProcessingTime(
	location string,
	mode string,
) time.Duration {
	// In a real implementation, this would use historical data and port/terminal efficiency metrics
	// This is a simplified example
	baseTime := 12 * time.Hour // 12 hours base processing time

	switch mode {
	case "Air":
		return baseTime
	case "Sea":
		return baseTime * 2 // 24 hours
	case "Road":
		return baseTime / 2 // 6 hours
	case "Rail":
		return baseTime // 12 hours
	default:
		return baseTime
	}
}

// GetOptimalRoute calculates the optimal route based on various factors
func (s *TrackingService) GetOptimalRoute(
	origin string,
	destination string,
	shipmentMode string,
	cargoDetails models.CargoDetails,
) ([]models.TransshipmentPoint, error) {
	// In a real implementation, this would use a routing algorithm considering:
	// - Distance
	// - Cost
	// - Transit time
	// - Port/terminal efficiency
	// - Historical performance
	// - Current conditions
	// This is a simplified example

	route := []models.TransshipmentPoint{
		{
			Location:    origin,
			Status:     "ORIGIN",
			ArrivalDate: time.Now(),
			DepartDate: time.Now().Add(s.calculateProcessingTime(origin, shipmentMode)),
		},
	}

	// Add intermediate points based on mode
	switch shipmentMode {
	case "Air":
		// Direct route for air shipments
		route = append(route, models.TransshipmentPoint{
			Location:    destination,
			Status:     "DESTINATION",
			ArrivalDate: time.Now().Add(s.calculateTransitTime(origin, destination, shipmentMode)),
		})
	case "Sea":
		// Add major port transshipment
		transitPort := s.getNearestPort(destination)
		route = append(route, models.TransshipmentPoint{
			Location:    transitPort,
			Status:     "TRANSIT",
			Carrier:    "Ocean Carrier",
			ArrivalDate: time.Now().Add(s.calculateTransitTime(origin, transitPort, shipmentMode)),
			DepartDate: time.Now().Add(s.calculateTransitTime(origin, transitPort, shipmentMode)).
				Add(s.calculateProcessingTime(transitPort, shipmentMode)),
		})
		route = append(route, models.TransshipmentPoint{
			Location:    destination,
			Status:     "DESTINATION",
			ArrivalDate: time.Now().Add(s.calculateTransitTime(transitPort, destination, shipmentMode)),
		})
	}

	return route, nil
}

func (s *TrackingService) getNearestPort(location string) string {
	// In a real implementation, this would query a database of ports and calculate distances
	// This is a simplified example
	return "Major Transit Port"
}
