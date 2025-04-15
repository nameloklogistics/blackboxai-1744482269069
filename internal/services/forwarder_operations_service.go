package services

import (
	"fmt"
	"time"

	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/stellar"
)

// ForwarderOperationsService handles core freight forwarding operations
type ForwarderOperationsService struct {
	txManager    *stellar.TransactionManager
	tokenManager *stellar.TokenManager
}

// NewForwarderOperationsService creates a new ForwarderOperationsService instance
func NewForwarderOperationsService(txManager *stellar.TransactionManager, tokenManager *stellar.TokenManager) *ForwarderOperationsService {
	return &ForwarderOperationsService{
		txManager:    txManager,
		tokenManager: tokenManager,
	}
}

// CreateQuoteRequest creates a new freight quote request
func (s *ForwarderOperationsService) CreateQuoteRequest(request *models.QuoteRequest) error {
	if err := s.validateQuoteRequest(request); err != nil {
		return fmt.Errorf("invalid quote request: %w", err)
	}

	request.ID = fmt.Sprintf("QR-%d", time.Now().UnixNano())
	request.CreatedAt = time.Now()
	request.Status = "NEW"

	return nil
}

// GenerateFreightQuote generates a detailed freight quote
func (s *ForwarderOperationsService) GenerateFreightQuote(request *models.QuoteRequest) (*models.FreightQuote, error) {
	quote := &models.FreightQuote{
		ID:         fmt.Sprintf("FQ-%d", time.Now().UnixNano()),
		RequestID:  request.ID,
		ForwarderID: request.ForwarderID,
		Mode:       request.Mode,
		ServiceType: request.ServiceType,
		CreatedAt:  time.Now(),
		Status:     "DRAFT",
	}

	// Calculate route segments
	segments, err := s.calculateRouting(request)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate routing: %w", err)
	}
	quote.Routing = segments

	// Calculate transit time
	quote.TransitTime = s.calculateTransitTime(segments)

	// Calculate charges
	if err := s.calculateCharges(request, quote); err != nil {
		return nil, fmt.Errorf("failed to calculate charges: %w", err)
	}

	// Set validity period
	quote.ValidUntil = time.Now().AddDate(0, 0, 7) // Valid for 7 days

	return quote, nil
}

// ConfirmRate confirms a freight quote
func (s *ForwarderOperationsService) ConfirmRate(quote *models.FreightQuote) (*models.RateConfirmation, error) {
	if err := s.validateQuote(quote); err != nil {
		return nil, fmt.Errorf("invalid quote: %w", err)
	}

	confirmation := &models.RateConfirmation{
		ID:          fmt.Sprintf("RC-%d", time.Now().UnixNano()),
		QuoteID:     quote.ID,
		ForwarderID: quote.ForwarderID,
		Quote:       *quote,
		ConfirmedAt: time.Now(),
		ValidUntil:  quote.ValidUntil,
		Status:      "ACTIVE",
	}

	return confirmation, nil
}

// CreateBooking creates a new shipment booking
func (s *ForwarderOperationsService) CreateBooking(confirmation *models.RateConfirmation, bookingDetails *models.ShipmentBooking) (*models.ShipmentBooking, error) {
	if err := s.validateBookingDetails(bookingDetails); err != nil {
		return nil, fmt.Errorf("invalid booking details: %w", err)
	}

	booking := &models.ShipmentBooking{
		ID:             fmt.Sprintf("BK-%d", time.Now().UnixNano()),
		ConfirmationID: confirmation.ID,
		ForwarderID:    confirmation.ForwarderID,
		CustomerID:     confirmation.CustomerID,
		BookingNumber:  s.generateBookingNumber(),
		BookingDate:    time.Now(),
		Mode:           confirmation.Quote.Mode,
		ServiceType:    confirmation.Quote.ServiceType,
		Routing:        confirmation.Quote.Routing,
		Amount:         confirmation.Quote.TotalAmount,
		Currency:       confirmation.Quote.Currency,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Status:         "PENDING",
	}

	// Copy relevant details from booking request
	s.copyBookingDetails(bookingDetails, booking)

	return booking, nil
}

// ConfirmBooking confirms a shipment booking
func (s *ForwarderOperationsService) ConfirmBooking(booking *models.ShipmentBooking) (*models.BookingConfirmation, error) {
	if err := s.validateBooking(booking); err != nil {
		return nil, fmt.Errorf("invalid booking: %w", err)
	}

	// Allocate equipment
	equipment, err := s.allocateEquipment(booking)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate equipment: %w", err)
	}

	// Generate schedule
	schedule, err := s.generateSchedule(booking)
	if err != nil {
		return nil, fmt.Errorf("failed to generate schedule: %w", err)
	}

	confirmation := &models.BookingConfirmation{
		ID:                 fmt.Sprintf("BC-%d", time.Now().UnixNano()),
		BookingID:          booking.ID,
		ConfirmationNumber: s.generateConfirmationNumber(),
		ConfirmedSchedule:  schedule,
		AllocatedEquipment: equipment,
		Carriers:           make(map[string]string),
		Agents:            make(map[string]string),
		ConfirmedAt:       time.Now(),
		Status:            "CONFIRMED",
	}

	// Assign carriers and agents
	s.assignServiceProviders(booking, confirmation)

	return confirmation, nil
}

// Helper functions

func (s *ForwarderOperationsService) validateQuoteRequest(request *models.QuoteRequest) error {
	if request.ForwarderID == "" {
		return fmt.Errorf("forwarder ID is required")
	}
	if request.CustomerID == "" {
		return fmt.Errorf("customer ID is required")
	}
	if request.Origin.Code == "" {
		return fmt.Errorf("origin location is required")
	}
	if request.Destination.Code == "" {
		return fmt.Errorf("destination location is required")
	}
	if len(request.Cargo) == 0 {
		return fmt.Errorf("cargo details are required")
	}
	return nil
}

func (s *ForwarderOperationsService) calculateRouting(request *models.QuoteRequest) ([]models.RouteSegment, error) {
	// In a real implementation, this would:
	// 1. Calculate optimal routing based on mode and service type
	// 2. Consider via points if specified
	// 3. Include schedules if available
	return []models.RouteSegment{
		{
			SegmentType: "MAIN_CARRIAGE",
			Mode:        request.Mode,
			From:        request.Origin,
			To:          request.Destination,
			TransitTime: "3 days",
		},
	}, nil
}

func (s *ForwarderOperationsService) calculateTransitTime(segments []models.RouteSegment) string {
	// In a real implementation, this would:
	// 1. Sum up transit times across segments
	// 2. Add buffer for connections
	// 3. Consider schedules and cutoff times
	return "3 days"
}

func (s *ForwarderOperationsService) calculateCharges(request *models.QuoteRequest, quote *models.FreightQuote) error {
	// Base charges
	quote.BaseCharges = models.RateComponent{
		ChargeType:  "FREIGHT",
		Description: "Main freight charges",
		Basis:      "PER_CONTAINER",
		UnitRate:   1000.0,
		Units:      float64(len(request.Cargo)),
		Amount:     1000.0 * float64(len(request.Cargo)),
		Currency:   "USD",
		Mandatory:  true,
	}

	// Local charges
	quote.LocalCharges = models.RateComponent{
		ChargeType:  "LOCAL",
		Description: "Local handling charges",
		Basis:      "PER_SHIPMENT",
		UnitRate:   200.0,
		Units:      1.0,
		Amount:     200.0,
		Currency:   "USD",
		Mandatory:  true,
	}

	// Calculate totals
	quote.SubTotal = quote.BaseCharges.Amount + quote.LocalCharges.Amount
	quote.Tax = quote.SubTotal * 0.1 // 10% tax
	quote.TotalAmount = quote.SubTotal + quote.Tax
	quote.Currency = "USD"

	return nil
}

func (s *ForwarderOperationsService) validateQuote(quote *models.FreightQuote) error {
	if quote.ID == "" {
		return fmt.Errorf("quote ID is required")
	}
	if quote.TotalAmount <= 0 {
		return fmt.Errorf("invalid total amount")
	}
	return nil
}

func (s *ForwarderOperationsService) validateBookingDetails(booking *models.ShipmentBooking) error {
	if booking.CustomerID == "" {
		return fmt.Errorf("customer ID is required")
	}
	if len(booking.Cargo) == 0 {
		return fmt.Errorf("cargo details are required")
	}
	return nil
}

func (s *ForwarderOperationsService) validateBooking(booking *models.ShipmentBooking) error {
	if booking.ID == "" {
		return fmt.Errorf("booking ID is required")
	}
	if booking.BookingNumber == "" {
		return fmt.Errorf("booking number is required")
	}
	return nil
}

func (s *ForwarderOperationsService) generateBookingNumber() string {
	return fmt.Sprintf("BK%d", time.Now().UnixNano())
}

func (s *ForwarderOperationsService) generateConfirmationNumber() string {
	return fmt.Sprintf("BC%d", time.Now().UnixNano())
}

func (s *ForwarderOperationsService) copyBookingDetails(from, to *models.ShipmentBooking) {
	to.Cargo = from.Cargo
	to.TotalWeight = from.TotalWeight
	to.TotalVolume = from.TotalVolume
	to.PickupDate = from.PickupDate
	to.DeliveryDate = from.DeliveryDate
	to.PickupInstructions = from.PickupInstructions
	to.DeliveryInstructions = from.DeliveryInstructions
	to.CargoInstructions = from.CargoInstructions
	to.CustomsInstructions = from.CustomsInstructions
}

func (s *ForwarderOperationsService) allocateEquipment(booking *models.ShipmentBooking) ([]models.Equipment, error) {
	// In a real implementation, this would:
	// 1. Check equipment availability
	// 2. Reserve equipment
	// 3. Assign specific units
	return []models.Equipment{}, nil
}

func (s *ForwarderOperationsService) generateSchedule(booking *models.ShipmentBooking) (models.Schedule, error) {
	// In a real implementation, this would:
	// 1. Check service schedules
	// 2. Consider transit times
	// 3. Add handling times
	return models.Schedule{
		Pickup: models.ScheduleEvent{
			Location: booking.Routing[0].From,
			Planned:  booking.PickupDate,
		},
		Delivery: models.ScheduleEvent{
			Location: booking.Routing[len(booking.Routing)-1].To,
			Planned:  booking.DeliveryDate,
		},
	}, nil
}

func (s *ForwarderOperationsService) assignServiceProviders(booking *models.ShipmentBooking, confirmation *models.BookingConfirmation) {
	// In a real implementation, this would:
	// 1. Assign carriers for each segment
	// 2. Assign agents at each location
	// 3. Update service provider details
}
