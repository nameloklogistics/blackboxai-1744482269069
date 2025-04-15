package services

import (
	"fmt"
	"time"

	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/stellar"
)

// FreightForwarderService handles freight forwarding operations
type FreightForwarderService struct {
	txManager    *stellar.TransactionManager
	tokenManager *stellar.TokenManager
}

// RegisterFreightForwarder registers a new freight forwarder
func (s *FreightForwarderService) RegisterFreightForwarder(forwarder *models.FreightForwarder) error {
	if err := s.validateFreightForwarder(forwarder); err != nil {
		return fmt.Errorf("invalid freight forwarder: %w", err)
	}

	forwarder.ID = fmt.Sprintf("FF-%d", time.Now().UnixNano())
	forwarder.Status = "ACTIVE"
	forwarder.CreatedAt = time.Now()
	forwarder.UpdatedAt = time.Now()

	// Create HQ branch office
	hqOffice := models.BranchOffice{
		ID:            fmt.Sprintf("BR-%d", time.Now().UnixNano()),
		Name:          forwarder.Name + " - Head Office",
		Code:          "HQ-001",
		Type:          "HQ",
		Address:       forwarder.HeadOffice,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	forwarder.BranchOffices = []models.BranchOffice{hqOffice}

	return nil
}

// AddBranchOffice adds a new branch office to a freight forwarder
func (s *FreightForwarderService) AddBranchOffice(request *models.BranchOfficeRequest) error {
	// Get forwarder details
	forwarder, err := s.GetFreightForwarder(request.ForwarderID)
	if err != nil {
		return fmt.Errorf("failed to get forwarder details: %w", err)
	}

	// Validate branch office is in the same country
	if request.Address.Country != forwarder.CountryCode {
		return fmt.Errorf("branch office must be in the same country as the freight forwarder (%s)", forwarder.CountryCode)
	}

	// Validate branch code uniqueness
	for _, office := range forwarder.BranchOffices {
		if office.Code == request.Code {
			return fmt.Errorf("branch code %s already exists", request.Code)
		}
	}

	// Create new branch office
	office := models.BranchOffice{
		ID:            fmt.Sprintf("BR-%d", time.Now().UnixNano()),
		Name:          request.Name,
		Code:          request.Code,
		Type:          request.Type,
		Address:       request.Address,
		ContactPerson: request.ContactPerson,
		ContactEmail:  request.ContactEmail,
		ContactPhone:  request.ContactPhone,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Add to forwarder's branch offices
	forwarder.BranchOffices = append(forwarder.BranchOffices, office)
	forwarder.UpdatedAt = time.Now()

	return nil
}

// GetBranchOffice retrieves a branch office by ID
func (s *FreightForwarderService) GetBranchOffice(forwarderID, branchID string) (*models.BranchOffice, error) {
	forwarder, err := s.GetFreightForwarder(forwarderID)
	if err != nil {
		return nil, err
	}

	for _, office := range forwarder.BranchOffices {
		if office.ID == branchID {
			return &office, nil
		}
	}

	return nil, fmt.Errorf("branch office not found")
}

// GetFreightForwarder retrieves a freight forwarder by ID
func (s *FreightForwarderService) GetFreightForwarder(ctx context.Context, id string) (*models.FreightForwarder, error) {
	// In a real implementation, this would query a database
	// This is a placeholder that returns a sample forwarder
	return &models.FreightForwarder{
		ID:             id,
		Name:           "Sample Forwarder",
		RegistrationNo: "FF123456",
		CountryCode:    "SG",
		LicenseNumber:  "CB789012",
		LicenseExpiry:  time.Now().AddDate(1, 0, 0),
		Status:         "ACTIVE",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

// NewFreightForwarderService creates a new FreightForwarderService instance
func NewFreightForwarderService(txManager *stellar.TransactionManager, tokenManager *stellar.TokenManager) *FreightForwarderService {
	return &FreightForwarderService{
		txManager:    txManager,
		tokenManager: tokenManager,
	}
}

// CreateRateRequest creates a new rate request
func (s *FreightForwarderService) CreateRateRequest(ctx context.Context, request *models.RateRequest) error {
	// Validate membership status
	if err := s.membershipSvc.ValidateMembership(ctx, models.FreightForwarderMembership, request.ForwarderID); err != nil {
		return fmt.Errorf("invalid membership: %w", err)
	}
	if err := s.validateRateRequest(request); err != nil {
		return fmt.Errorf("invalid rate request: %w", err)
	}

	request.ID = fmt.Sprintf("RR-%d", time.Now().UnixNano())
	request.CreatedAt = time.Now()

	return nil
}

// CreateRateQuotation creates a new rate quotation
func (s *FreightForwarderService) CreateRateQuotation(quotation *models.RateQuotation) error {
	if err := s.validateRateQuotation(quotation); err != nil {
		return fmt.Errorf("invalid rate quotation: %w", err)
	}

	// Calculate total amount
	var total float64
	total += quotation.FreightCharges
	total += quotation.LocalCharges
	total += quotation.CustomsCharges
	
	for _, charge := range quotation.OtherCharges {
		total += charge.Amount
	}
	
	quotation.TotalAmount = total
	quotation.ID = fmt.Sprintf("QT-%d", time.Now().UnixNano())
	quotation.CreatedAt = time.Now()
	quotation.Status = "DRAFT"

	return nil
}

// ConfirmRate confirms a rate quotation
func (s *FreightForwarderService) ConfirmRate(confirmation *models.RateConfirmation) error {
	if err := s.validateRateConfirmation(confirmation); err != nil {
		return fmt.Errorf("invalid rate confirmation: %w", err)
	}

	confirmation.ID = fmt.Sprintf("RC-%d", time.Now().UnixNano())
	confirmation.ConfirmedAt = time.Now()
	confirmation.Status = "CONFIRMED"

	return nil
}

// CreateBooking creates a new shipment booking
func (s *FreightForwarderService) CreateBooking(booking *models.ShipmentBooking) error {
	if err := s.validateBooking(booking); err != nil {
		return fmt.Errorf("invalid booking: %w", err)
	}

	booking.ID = fmt.Sprintf("BK-%d", time.Now().UnixNano())
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()
	booking.Status = "PENDING"

	return nil
}

// ConfirmBooking confirms a shipment booking
func (s *FreightForwarderService) ConfirmBooking(confirmation *models.BookingConfirmation) error {
	if err := s.validateBookingConfirmation(confirmation); err != nil {
		return fmt.Errorf("invalid booking confirmation: %w", err)
	}

	confirmation.ID = fmt.Sprintf("BC-%d", time.Now().UnixNano())
	confirmation.ConfirmedAt = time.Now()
	confirmation.Status = "CONFIRMED"

	return nil
}

// CalculateRates calculates rates for a given request
func (s *FreightForwarderService) CalculateRates(request *models.RateRequest) (*models.RateCalculation, error) {
	calculation := &models.RateCalculation{}

	// Calculate base rates
	baseRates, err := s.calculateBaseRates(request)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate base rates: %w", err)
	}
	calculation.FreightRate = baseRates.FreightRate
	calculation.LocalRate = baseRates.LocalRate

	// Calculate surcharges
	surcharges, err := s.calculateSurcharges(request)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate surcharges: %w", err)
	}
	calculation.FuelSurcharge = surcharges.FuelSurcharge
	calculation.SecurityCharge = surcharges.SecurityCharge
	calculation.PeakSeasonCharge = surcharges.PeakSeasonCharge

	// Calculate additional services
	if err := s.calculateAdditionalServices(request, calculation); err != nil {
		return nil, fmt.Errorf("failed to calculate additional services: %w", err)
	}

	// Calculate totals
	calculation.SubTotal = s.calculateSubTotal(calculation)
	calculation.Tax = s.calculateTax(calculation.SubTotal)
	calculation.Total = calculation.SubTotal + calculation.Tax

	// Generate charge breakdown
	calculation.ChargeBreakdown = s.generateChargeBreakdown(calculation)

	return calculation, nil
}

// AddBookingInstruction adds instructions to a booking
func (s *FreightForwarderService) AddBookingInstruction(instruction *models.BookingInstruction) error {
	if err := s.validateBookingInstruction(instruction); err != nil {
		return fmt.Errorf("invalid booking instruction: %w", err)
	}

	instruction.CreatedAt = time.Now()
	return nil
}

// UpdateDocumentRequirement updates document requirements for a booking
func (s *FreightForwarderService) UpdateDocumentRequirement(requirement *models.DocumentRequirement) error {
	if err := s.validateDocumentRequirement(requirement); err != nil {
		return fmt.Errorf("invalid document requirement: %w", err)
	}

	return nil
}

// Helper functions

func (s *FreightForwarderService) validateRateRequest(request *models.RateRequest) error {
	if request.ForwarderID == "" {
		return fmt.Errorf("forwarder ID is required")
	}
	if request.ShipperID == "" {
		return fmt.Errorf("shipper ID is required")
	}
	if request.Origin == "" {
		return fmt.Errorf("origin is required")
	}
	if request.Destination == "" {
		return fmt.Errorf("destination is required")
	}
	if request.OriginCountry == "" {
		return fmt.Errorf("origin country is required")
	}
	if request.DestCountry == "" {
		return fmt.Errorf("destination country is required")
	}

	// Get forwarder details
	forwarder, err := s.GetFreightForwarder(request.ForwarderID)
	if err != nil {
		return fmt.Errorf("failed to get forwarder details: %w", err)
	}

	// Check if forwarder can operate in the requested countries
	switch request.ServiceType {
	case models.ImportService:
		if request.DestCountry != forwarder.CountryCode {
			return &models.CountryRestrictionError{
				ForwarderID:     forwarder.ID,
				ForwarderCountry: forwarder.CountryCode,
				RequestedCountry: request.DestCountry,
			}
		}
	case models.ExportService:
		if request.OriginCountry != forwarder.CountryCode {
			return &models.CountryRestrictionError{
				ForwarderID:     forwarder.ID,
				ForwarderCountry: forwarder.CountryCode,
				RequestedCountry: request.OriginCountry,
			}
		}
	case models.TransitService, models.TransshipmentService:
		if request.OriginCountry != forwarder.CountryCode && request.DestCountry != forwarder.CountryCode {
			return &models.CountryRestrictionError{
				ForwarderID:     forwarder.ID,
				ForwarderCountry: forwarder.CountryCode,
				RequestedCountry: fmt.Sprintf("%s/%s", request.OriginCountry, request.DestCountry),
			}
		}
	}

	return nil
}

func (s *FreightForwarderService) validateRateQuotation(quotation *models.RateQuotation) error {
	if quotation.RequestID == "" {
		return fmt.Errorf("request ID is required")
	}
	if quotation.ForwarderID == "" {
		return fmt.Errorf("forwarder ID is required")
	}
	if quotation.BranchOfficeID == "" {
		return fmt.Errorf("branch office ID is required")
	}
	if quotation.FreightCharges <= 0 {
		return fmt.Errorf("freight charges must be greater than 0")
	}

	// Get forwarder details
	forwarder, err := s.GetFreightForwarder(quotation.ForwarderID)
	if err != nil {
		return fmt.Errorf("failed to get forwarder details: %w", err)
	}

	// Validate branch office
	branchOffice, err := s.GetBranchOffice(quotation.ForwarderID, quotation.BranchOfficeID)
	if err != nil {
		return fmt.Errorf("failed to get branch office details: %w", err)
	}

	if !branchOffice.IsActive {
		return fmt.Errorf("branch office %s is not active", branchOffice.ID)
	}

	if branchOffice.Address.Country != forwarder.CountryCode {
		return fmt.Errorf("branch office must be in the same country as the freight forwarder")
	}

	// Set forwarder's registered country
	quotation.ForwarderCountry = forwarder.CountryCode

	// Check if forwarder can operate in the requested countries
	switch quotation.ServiceType {
	case models.ImportService:
		if quotation.DestCountry != forwarder.CountryCode {
			return &models.CountryRestrictionError{
				ForwarderID:     forwarder.ID,
				ForwarderCountry: forwarder.CountryCode,
				RequestedCountry: quotation.DestCountry,
			}
		}
	case models.ExportService:
		if quotation.OriginCountry != forwarder.CountryCode {
			return &models.CountryRestrictionError{
				ForwarderID:     forwarder.ID,
				ForwarderCountry: forwarder.CountryCode,
				RequestedCountry: quotation.OriginCountry,
			}
		}
	case models.TransitService, models.TransshipmentService:
		if quotation.OriginCountry != forwarder.CountryCode && quotation.DestCountry != forwarder.CountryCode {
			return &models.CountryRestrictionError{
				ForwarderID:     forwarder.ID,
				ForwarderCountry: forwarder.CountryCode,
				RequestedCountry: fmt.Sprintf("%s/%s", quotation.OriginCountry, quotation.DestCountry),
			}
		}
	}

	return nil
}

func (s *FreightForwarderService) validateFreightForwarder(forwarder *models.FreightForwarder) error {
	if forwarder.Name == "" {
		return fmt.Errorf("name is required")
	}
	if forwarder.RegistrationNo == "" {
		return fmt.Errorf("registration number is required")
	}
	if forwarder.CountryCode == "" {
		return fmt.Errorf("country code is required")
	}
	if forwarder.LicenseNumber == "" {
		return fmt.Errorf("license number is required")
	}
	if forwarder.LicenseExpiry.Before(time.Now()) {
		return fmt.Errorf("license has expired")
	}
	return nil
}

func (s *FreightForwarderService) validateRateConfirmation(confirmation *models.RateConfirmation) error {
	if confirmation.QuotationID == "" {
		return fmt.Errorf("quotation ID is required")
	}
	if confirmation.ForwarderID == "" {
		return fmt.Errorf("forwarder ID is required")
	}
	if confirmation.ShipperID == "" {
		return fmt.Errorf("shipper ID is required")
	}
	return nil
}

func (s *FreightForwarderService) validateBooking(booking *models.ShipmentBooking) error {
	if booking.ConfirmationID == "" {
		return fmt.Errorf("confirmation ID is required")
	}
	if booking.ForwarderID == "" {
		return fmt.Errorf("forwarder ID is required")
	}
	if booking.ShipperID == "" {
		return fmt.Errorf("shipper ID is required")
	}
	return nil
}

func (s *FreightForwarderService) validateBookingConfirmation(confirmation *models.BookingConfirmation) error {
	if confirmation.BookingID == "" {
		return fmt.Errorf("booking ID is required")
	}
	if confirmation.BookingNumber == "" {
		return fmt.Errorf("booking number is required")
	}
	return nil
}

func (s *FreightForwarderService) validateBookingInstruction(instruction *models.BookingInstruction) error {
	if instruction.BookingID == "" {
		return fmt.Errorf("booking ID is required")
	}
	if instruction.InstructionType == "" {
		return fmt.Errorf("instruction type is required")
	}
	if instruction.Instructions == "" {
		return fmt.Errorf("instructions are required")
	}
	return nil
}

func (s *FreightForwarderService) validateDocumentRequirement(requirement *models.DocumentRequirement) error {
	if requirement.BookingID == "" {
		return fmt.Errorf("booking ID is required")
	}
	if requirement.DocumentType == "" {
		return fmt.Errorf("document type is required")
	}
	return nil
}

func (s *FreightForwarderService) calculateBaseRates(request *models.RateRequest) (*models.RateCalculation, error) {
	// In a real implementation, this would calculate rates based on:
	// - Transport mode
	// - Route
	// - Cargo details
	// - Market rates
	return &models.RateCalculation{
		FreightRate: 1000.0,
		LocalRate:   200.0,
	}, nil
}

func (s *FreightForwarderService) calculateSurcharges(request *models.RateRequest) (*models.RateCalculation, error) {
	// In a real implementation, this would calculate surcharges based on:
	// - Fuel prices
	// - Security requirements
	// - Season
	return &models.RateCalculation{
		FuelSurcharge:    50.0,
		SecurityCharge:   30.0,
		PeakSeasonCharge: 20.0,
	}, nil
}

func (s *FreightForwarderService) calculateAdditionalServices(request *models.RateRequest, calculation *models.RateCalculation) error {
	// In a real implementation, this would calculate charges for:
	// - Customs clearance
	// - Insurance
	// - Documentation
	calculation.CustomsClearance = 150.0
	calculation.Insurance = 100.0
	calculation.Documentation = 50.0
	return nil
}

func (s *FreightForwarderService) calculateSubTotal(calculation *models.RateCalculation) float64 {
	return calculation.FreightRate +
		calculation.LocalRate +
		calculation.FuelSurcharge +
		calculation.SecurityCharge +
		calculation.PeakSeasonCharge +
		calculation.CustomsClearance +
		calculation.Insurance +
		calculation.Documentation
}

func (s *FreightForwarderService) calculateTax(subTotal float64) float64 {
	// In a real implementation, this would calculate applicable taxes
	return subTotal * 0.1 // 10% tax
}

func (s *FreightForwarderService) generateChargeBreakdown(calculation *models.RateCalculation) []models.Charge {
	return []models.Charge{
		{Type: "Freight", Amount: calculation.FreightRate, Currency: "USD"},
		{Type: "Local Charges", Amount: calculation.LocalRate, Currency: "USD"},
		{Type: "Fuel Surcharge", Amount: calculation.FuelSurcharge, Currency: "USD"},
		{Type: "Security", Amount: calculation.SecurityCharge, Currency: "USD"},
		{Type: "Peak Season", Amount: calculation.PeakSeasonCharge, Currency: "USD"},
		{Type: "Customs", Amount: calculation.CustomsClearance, Currency: "USD"},
		{Type: "Insurance", Amount: calculation.Insurance, Currency: "USD"},
		{Type: "Documentation", Amount: calculation.Documentation, Currency: "USD"},
	}
}
