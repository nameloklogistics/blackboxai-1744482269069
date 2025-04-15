package services

import (
	"fmt"
	"time"

	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/stellar"
)

// UserOperationsService handles operations for all user types
type UserOperationsService struct {
	txManager    *stellar.TransactionManager
	tokenManager *stellar.TokenManager
}

// NewUserOperationsService creates a new UserOperationsService instance
func NewUserOperationsService(txManager *stellar.TransactionManager, tokenManager *stellar.TokenManager) *UserOperationsService {
	return &UserOperationsService{
		txManager:    txManager,
		tokenManager: tokenManager,
	}
}

// CreateQuoteRequest creates a new quote request for any user type
func (s *UserOperationsService) CreateQuoteRequest(request *models.UserQuoteRequest) error {
	if err := s.validateQuoteRequest(request); err != nil {
		return fmt.Errorf("invalid quote request: %w", err)
	}

	request.ID = fmt.Sprintf("QR-%d", time.Now().UnixNano())
	request.CreatedAt = time.Now()
	request.Status = "NEW"

	return nil
}

// GenerateQuoteResponse generates a quote response
func (s *UserOperationsService) GenerateQuoteResponse(request *models.UserQuoteRequest) (*models.UserQuoteResponse, error) {
	response := &models.UserQuoteResponse{
		ID:          fmt.Sprintf("QT-%d", time.Now().UnixNano()),
		RequestID:   request.ID,
		RespondedBy: models.FreightForwarder,
		Parties:     request.Parties,
		Mode:        request.Mode,
		ServiceType: request.ServiceType,
		CreatedAt:   time.Now(),
		Status:      "DRAFT",
	}

	// Calculate route segments
	segments, err := s.calculateRouting(request)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate routing: %w", err)
	}
	response.Routing = segments

	// Calculate transit time
	response.TransitTime = s.calculateTransitTime(segments)

	// Calculate charges
	if err := s.calculateCharges(request, response); err != nil {
		return nil, fmt.Errorf("failed to calculate charges: %w", err)
	}

	// Set validity period
	response.ValidUntil = time.Now().AddDate(0, 0, 7) // Valid for 7 days

	return response, nil
}

// ConfirmQuote confirms a quote by any user type
func (s *UserOperationsService) ConfirmQuote(quote *models.UserQuoteResponse, confirmedBy models.UserType) (*models.UserQuoteConfirmation, error) {
	if err := s.validateQuoteResponse(quote); err != nil {
		return nil, fmt.Errorf("invalid quote: %w", err)
	}

	confirmation := &models.UserQuoteConfirmation{
		ID:          fmt.Sprintf("QC-%d", time.Now().UnixNano()),
		QuoteID:     quote.ID,
		ConfirmedBy: confirmedBy,
		Parties:     quote.Parties,
		Quote:       *quote,
		ConfirmedAt: time.Now(),
		ValidUntil:  quote.ValidUntil,
		Status:      "ACTIVE",
	}

	return confirmation, nil
}

// CreateBooking creates a new booking request
func (s *UserOperationsService) CreateBooking(confirmation *models.UserQuoteConfirmation, request *models.UserBookingRequest) (*models.UserBookingRequest, error) {
	if err := s.validateBookingRequest(request); err != nil {
		return nil, fmt.Errorf("invalid booking request: %w", err)
	}

	request.ID = fmt.Sprintf("BK-%d", time.Now().UnixNano())
	request.ConfirmationID = confirmation.ID
	request.CreatedAt = time.Now()
	request.Status = "PENDING"

	return request, nil
}

// ConfirmBooking confirms a booking
func (s *UserOperationsService) ConfirmBooking(request *models.UserBookingRequest) (*models.UserBookingConfirmation, error) {
	if err := s.validateBookingRequest(request); err != nil {
		return nil, fmt.Errorf("invalid booking request: %w", err)
	}

	// Allocate equipment
	equipment, err := s.allocateEquipment(request)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate equipment: %w", err)
	}

	// Generate schedule
	schedule, err := s.generateSchedule(request)
	if err != nil {
		return nil, fmt.Errorf("failed to generate schedule: %w", err)
	}

	confirmation := &models.UserBookingConfirmation{
		ID:                fmt.Sprintf("BC-%d", time.Now().UnixNano()),
		BookingID:         request.ID,
		ConfirmedBy:       models.FreightForwarder,
		Parties:           request.Parties,
		BookingNumber:     s.generateBookingNumber(),
		ConfirmedSchedule: schedule,
		AllocatedEquipment: equipment,
		Carriers:          make(map[string]string),
		Agents:            make(map[string]string),
		Amount:            0, // Will be set based on quote
		Currency:          "USD",
		ConfirmedAt:       time.Now(),
		Status:            "CONFIRMED",
	}

	// Assign carriers and agents
	s.assignServiceProviders(request, confirmation)

	return confirmation, nil
}

// Helper functions

func (s *UserOperationsService) validateQuoteRequest(request *models.UserQuoteRequest) error {
	if request.RequestedBy == "" {
		return fmt.Errorf("requester type is required")
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

func (s *UserOperationsService) validateQuoteResponse(quote *models.UserQuoteResponse) error {
	if quote.ID == "" {
		return fmt.Errorf("quote ID is required")
	}
	if quote.TotalAmount <= 0 {
		return fmt.Errorf("invalid total amount")
	}
	return nil
}

func (s *UserOperationsService) validateBookingRequest(request *models.UserBookingRequest) error {
	if request.ConfirmationID == "" {
		return fmt.Errorf("quote confirmation ID is required")
	}
	if request.RequestedBy == "" {
		return fmt.Errorf("requester type is required")
	}
	if len(request.Cargo) == 0 {
		return fmt.Errorf("cargo details are required")
	}
	return nil
}

func (s *UserOperationsService) calculateRouting(request *models.UserQuoteRequest) ([]models.RouteSegment, error) {
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

func (s *UserOperationsService) calculateTransitTime(segments []models.RouteSegment) string {
	// In a real implementation, this would:
	// 1. Sum up transit times across segments
	// 2. Add buffer for connections
	// 3. Consider schedules and cutoff times
	return "3 days"
}

func (s *UserOperationsService) calculateCharges(request *models.UserQuoteRequest, response *models.UserQuoteResponse) error {
	// Base charges
	response.BaseCharges = models.RateComponent{
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
	response.LocalCharges = models.RateComponent{
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
	response.SubTotal = response.BaseCharges.Amount + response.LocalCharges.Amount
	response.Tax = response.SubTotal * 0.1 // 10% tax
	response.TotalAmount = response.SubTotal + response.Tax
	response.Currency = "USD"

	return nil
}

func (s *UserOperationsService) generateBookingNumber() string {
	return fmt.Sprintf("BK%d", time.Now().UnixNano())
}

func (s *UserOperationsService) allocateEquipment(request *models.UserBookingRequest) ([]models.Equipment, error) {
	// In a real implementation, this would:
	// 1. Check equipment availability
	// 2. Reserve equipment
	// 3. Assign specific units
	return []models.Equipment{}, nil
}

func (s *UserOperationsService) generateSchedule(request *models.UserBookingRequest) (models.Schedule, error) {
	// In a real implementation, this would:
	// 1. Check service schedules
	// 2. Consider transit times
	// 3. Add handling times
	return models.Schedule{
		Pickup: models.ScheduleEvent{
			Location: request.Routing[0].From,
			Planned:  request.PickupDate,
		},
		Delivery: models.ScheduleEvent{
			Location: request.Routing[len(request.Routing)-1].To,
			Planned:  request.DeliveryDate,
		},
	}, nil
}

func (s *UserOperationsService) assignServiceProviders(request *models.UserBookingRequest, confirmation *models.UserBookingConfirmation) {
	// In a real implementation, this would:
	// 1. Assign carriers for each segment
	// 2. Assign agents at each location
	// 3. Update service provider details
}
