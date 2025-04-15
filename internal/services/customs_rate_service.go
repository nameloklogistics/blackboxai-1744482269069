package services

import (
	"context"
	"fmt"
	"time"

	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/stellar"
)

// CustomsRateService handles customs rate operations
type CustomsRateService struct {
	txManager              *stellar.TransactionManager
	tokenManager           *stellar.TokenManager
	licenseVerificationSvc LicenseVerificationService
	membershipSvc         *MembershipService
}

// NewCustomsRateService creates a new CustomsRateService instance
func NewCustomsRateService(
	txManager *stellar.TransactionManager,
	tokenManager *stellar.TokenManager,
	licenseVerificationSvc LicenseVerificationService,
	membershipSvc *MembershipService,
) *CustomsRateService {
	return &CustomsRateService{
		txManager:              txManager,
		tokenManager:           tokenManager,
		licenseVerificationSvc: licenseVerificationSvc,
		membershipSvc:         membershipSvc,
	}
}

// RegisterCustomsBroker registers a new customs broker
func (s *CustomsRateService) RegisterCustomsBroker(ctx context.Context, broker *models.CustomsBroker) error {
	if err := s.validateCustomsBroker(broker); err != nil {
		return fmt.Errorf("invalid customs broker: %w", err)
	}

	broker.ID = fmt.Sprintf("CB-%d", time.Now().UnixNano())
	broker.Status = "ACTIVE"
	broker.CreatedAt = time.Now()
	broker.UpdatedAt = time.Now()

	// Create HQ branch office
	hqOffice := models.CustomsBrokerBranch{
		ID:            fmt.Sprintf("CBB-%d", time.Now().UnixNano()),
		Name:          broker.Name + " - Head Office",
		Code:          "HQ-001",
		Type:          "HQ",
		Address:       broker.HeadOffice,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	broker.BranchOffices = []models.CustomsBrokerBranch{hqOffice}

	// Create trial membership
	membership, err := s.membershipSvc.CreateMembership(models.CustomsBrokerMembership, broker.ID)
	if err != nil {
		return fmt.Errorf("failed to create membership: %w", err)
	}

	// Schedule membership reminders
	if err := s.membershipSvc.ScheduleMembershipReminders(ctx, membership.ID); err != nil {
		return fmt.Errorf("failed to schedule reminders: %w", err)
	}

	return nil
}

// AddBrokerBranch adds a new branch office to a customs broker
func (s *CustomsRateService) AddBrokerBranch(ctx context.Context, brokerID string, branch *models.CustomsBrokerBranch) error {
	// Validate membership status
	if err := s.membershipSvc.ValidateMembership(ctx, models.CustomsBrokerMembership, brokerID); err != nil {
		return fmt.Errorf("invalid membership: %w", err)
	}
	// Get broker details
	broker, err := s.getCustomsBroker(brokerID)
	if err != nil {
		return fmt.Errorf("failed to get customs broker: %w", err)
	}

	// Validate branch office is in the same country
	if branch.Address.Country != broker.CountryCode {
		return fmt.Errorf("branch office must be in the same country as the customs broker (%s)", broker.CountryCode)
	}

	// Validate branch code uniqueness
	for _, office := range broker.BranchOffices {
		if office.Code == branch.Code {
			return fmt.Errorf("branch code %s already exists", branch.Code)
		}
	}

	branch.ID = fmt.Sprintf("CBB-%d", time.Now().UnixNano())
	branch.IsActive = true
	branch.CreatedAt = time.Now()
	branch.UpdatedAt = time.Now()

	// Add to broker's branch offices
	broker.BranchOffices = append(broker.BranchOffices, *branch)
	broker.UpdatedAt = time.Now()

	return nil
}

// GetBrokerBranch retrieves a branch office by ID
func (s *CustomsRateService) GetBrokerBranch(ctx context.Context, brokerID, branchID string) (*models.CustomsBrokerBranch, error) {
	// Validate membership status
	if err := s.membershipSvc.ValidateMembership(ctx, models.CustomsBrokerMembership, brokerID); err != nil {
		return nil, fmt.Errorf("invalid membership: %w", err)
	}
	broker, err := s.getCustomsBroker(brokerID)
	if err != nil {
		return nil, err
	}

	for _, office := range broker.BranchOffices {
		if office.ID == branchID {
			return &office, nil
		}
	}

	return nil, fmt.Errorf("branch office not found")
}

// getCustomsBroker retrieves a customs broker by ID
func (s *CustomsRateService) getCustomsBroker(id string) (*models.CustomsBroker, error) {
	// In a real implementation, this would query a database
	return &models.CustomsBroker{
		ID:             id,
		Name:           "Sample Broker",
		RegistrationNo: "CB123456",
		CountryCode:    "SG",
		LicenseNumber:  "LIC789012",
		LicenseExpiry:  time.Now().AddDate(1, 0, 0),
		Status:         "ACTIVE",
		HeadOffice: models.Address{
			Street:     "123 Main Street",
			City:       "Singapore",
			State:      "Singapore",
			PostalCode: "123456",
			Country:    "SG",
		},
		LicenseAuthority: "SG-CUSTOMS-AUTH",
		LicenseClass:     "CLASS-A",
		LicenseScope:     []string{"CUSTOMS_CLEARANCE", "DOCUMENTATION"},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}, nil
}

// CreateCustomsRate creates a new customs rate
func (s *CustomsRateService) CreateCustomsRate(ctx context.Context, rate *models.CustomsRate) error {
	// Validate membership status
	if err := s.membershipSvc.ValidateMembership(ctx, models.CustomsBrokerMembership, rate.BrokerID); err != nil {
		return fmt.Errorf("invalid membership: %w", err)
	}
	if err := s.validateCustomsRate(rate); err != nil {
		return fmt.Errorf("invalid customs rate: %w", err)
	}

	rate.ID = fmt.Sprintf("RATE-%d", time.Now().UnixNano())
	rate.CreatedAt = time.Now()
	rate.UpdatedAt = time.Now()

	return nil
}

// GetQuotation calculates customs clearance quotation
func (s *CustomsRateService) GetQuotation(
	ctx context.Context,
	brokerID string,
	branchID string,
	rateID string,
	cargoDetails *models.CargoDetails,
	cargoValue float64,
	transportMode models.TransportMode,
	packagingMode models.PackagingMode,
	packagingDetails interface{},
) (*models.CustomsQuotation, error) {
	// Validate membership status
	if err := s.membershipSvc.ValidateMembership(ctx, models.CustomsBrokerMembership, brokerID); err != nil {
		return nil, fmt.Errorf("invalid membership: %w", err)
	}
	// Get broker details
	broker, err := s.getCustomsBroker(brokerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customs broker: %w", err)
	}

	// Verify current license status
	status, err := s.licenseVerificationSvc.GetVerificationStatus(ctx, broker)
	if err != nil {
		return nil, fmt.Errorf("failed to get license status: %w", err)
	}

	if status != models.LicenseVerified {
		return nil, &models.LicenseVerificationError{
			BrokerID:      broker.ID,
			LicenseNumber: broker.LicenseNumber,
			AuthorityID:   broker.LicenseAuthority,
			Reason:        fmt.Sprintf("current license status is %s", status),
		}
	}

	// Validate branch office
	branch, err := s.GetBrokerBranch(brokerID, branchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch office: %w", err)
	}

	if !branch.IsActive {
		return nil, fmt.Errorf("branch office %s is not active", branch.ID)
	}

	if branch.Address.Country != broker.CountryCode {
		return nil, fmt.Errorf("branch office must be in the same country as the customs broker")
	}

	// Get rate details
	rate, err := s.getCustomsRate(rateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customs rate: %w", err)
	}

	// Validate broker can operate in the rate's country
	if rate.Country != broker.CountryCode {
		return &models.CustomsBrokerRestrictionError{
			BrokerID:        broker.ID,
			BrokerCountry:   broker.CountryCode,
			RequestedCountry: rate.Country,
		}
	}

	// Validate transport and packaging modes
	if err := s.validateModes(transportMode, packagingMode); err != nil {
		return nil, fmt.Errorf("invalid modes: %w", err)
	}

	// Create calculator
	calculator := &models.CustomsRateCalculator{
		Rate:            rate,
		CargoDetails:    cargoDetails,
		TransportMode:   transportMode,
		PackagingMode:   packagingMode,
		PackagingDetails: packagingDetails,
	}

	// Calculate all components
	basicCharges := calculator.CalculateBasicCharges()
	customsDuty := calculator.CalculateCustomsDuty(cargoValue)
	vat := calculator.CalculateVAT(cargoValue)
	otherTaxes := calculator.CalculateOtherTaxes(cargoValue)
	additionalFees := calculator.CalculateAdditionalFees()

	// Calculate total additional fees
	var totalAdditionalFees float64
	for _, fee := range additionalFees {
		if fee.IsRequired {
			totalAdditionalFees += fee.Amount
		}
	}

	// Create quotation
	quotation := &models.CustomsQuotation{
		ID:                  fmt.Sprintf("QUOTE-%d", time.Now().UnixNano()),
		BrokerID:           brokerID,
		BranchOfficeID:     branchID,
		RateID:             rateID,
		Type:               rate.Type,
		TransportMode:      transportMode,
		PackagingMode:      packagingMode,
		CargoValue:         cargoValue,
		CargoType:          cargoDetails.CargoType,
		HSCode:             cargoDetails.HsCode,
		Weight:             cargoDetails.Weight,
		Volume:             cargoDetails.Volume,
		PackagingDetails:   packagingDetails,
		BasicCharges:       basicCharges,
		CustomsDuty:        customsDuty,
		VAT:                vat,
		OtherTaxes:         otherTaxes,
		TransportModeCharges: s.calculateTransportModeCharges(rate, transportMode),
		PackagingCharges:    s.calculatePackagingCharges(rate, transportMode, packagingMode, packagingDetails),
		HandlingFees:        s.calculateHandlingFees(rate, transportMode),
		AdditionalFees:      totalAdditionalFees,
		FeeBreakdown:        additionalFees,
		TotalAmount:         calculator.CalculateTotalAmount(cargoValue),
		Currency:            rate.Currency,
		ValidUntil:          time.Now().Add(24 * time.Hour),
		CreatedAt:           time.Now(),
	}

	return quotation, nil
}

// Helper functions

func (s *CustomsRateService) validateCustomsBroker(broker *models.CustomsBroker) error {
	if broker.Name == "" {
		return fmt.Errorf("name is required")
	}
	if broker.RegistrationNo == "" {
		return fmt.Errorf("registration number is required")
	}
	if broker.CountryCode == "" {
		return fmt.Errorf("country code is required")
	}
	if broker.LicenseNumber == "" {
		return fmt.Errorf("license number is required")
	}
	if broker.LicenseExpiry.Before(time.Now()) {
		return fmt.Errorf("license has expired")
	}
	if broker.HeadOffice.Country != broker.CountryCode {
		return fmt.Errorf("head office must be in the same country as registration")
	}
	if broker.LicenseAuthority == "" {
		return fmt.Errorf("license authority is required")
	}

	// Verify license with authority
	ctx := context.Background()
	verification, err := s.licenseVerificationSvc.VerifyLicense(ctx, broker)
	if err != nil {
		return fmt.Errorf("license verification failed: %w", err)
	}

	// Check verification status
	if verification.VerificationStatus != models.LicenseVerified {
		return &models.LicenseVerificationError{
			BrokerID:      broker.ID,
			LicenseNumber: broker.LicenseNumber,
			AuthorityID:   broker.LicenseAuthority,
			Reason:        fmt.Sprintf("license status is %s", verification.VerificationStatus),
		}
	}

	// Validate license scope for customs operations
	requiredScope := []string{"CUSTOMS_CLEARANCE", "DOCUMENTATION"}
	if err := s.licenseVerificationSvc.ValidateLicenseScope(ctx, broker, requiredScope); err != nil {
		return fmt.Errorf("invalid license scope: %w", err)
	}

	// Schedule next verification
	if err := s.licenseVerificationSvc.ScheduleVerification(ctx, broker, 30*24*time.Hour); err != nil {
		return fmt.Errorf("failed to schedule verification: %w", err)
	}

	// Update broker's verification details
	broker.LicenseVerification = *verification
	broker.UpdatedAt = time.Now()

	return nil
}

func (s *CustomsRateService) validateCustomsRate(rate *models.CustomsRate) error {
	if rate.Type == "" {
		return fmt.Errorf("rate type is required")
	}
	if !s.isValidRateType(rate.Type) {
		return fmt.Errorf("invalid rate type: %s", rate.Type)
	}
	if rate.Country == "" {
		return fmt.Errorf("country is required")
	}
	if rate.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if rate.BasicHandling <= 0 {
		return fmt.Errorf("basic handling charge must be greater than 0")
	}
	if rate.Documentation <= 0 {
		return fmt.Errorf("documentation charge must be greater than 0")
	}

	// Validate broker and branch if provided
	if rate.BrokerID != "" {
		broker, err := s.getCustomsBroker(rate.BrokerID)
		if err != nil {
			return fmt.Errorf("failed to get customs broker: %w", err)
		}

		if rate.Country != broker.CountryCode {
			return &models.CustomsBrokerRestrictionError{
				BrokerID:        broker.ID,
				BrokerCountry:   broker.CountryCode,
				RequestedCountry: rate.Country,
			}
		}

		if rate.BranchOfficeID != "" {
			branch, err := s.GetBrokerBranch(rate.BrokerID, rate.BranchOfficeID)
			if err != nil {
				return fmt.Errorf("failed to get branch office: %w", err)
			}

			if !branch.IsActive {
				return fmt.Errorf("branch office %s is not active", branch.ID)
			}

			if branch.Address.Country != broker.CountryCode {
				return fmt.Errorf("branch office must be in the same country as the customs broker")
			}
		}
	}

	return nil
}

func (s *CustomsRateService) isValidRateType(rateType models.CustomsRateType) bool {
	switch rateType {
	case models.OriginCustoms,
		models.TransshipmentCustoms,
		models.TransitCustoms,
		models.DestinationCustoms:
		return true
	default:
		return false
	}
}

func (s *CustomsRateService) validateModes(transportMode models.TransportMode, packagingMode models.PackagingMode) error {
	// Validate transport mode
	switch transportMode {
	case models.SeaTransport, models.AirTransport, models.RailTransport, models.LandTransport:
	default:
		return fmt.Errorf("invalid transport mode: %s", transportMode)
	}

	// Validate packaging mode based on transport mode
	validPackaging := false
	switch transportMode {
	case models.SeaTransport:
		switch packagingMode {
		case models.FCL, models.LCL:
			validPackaging = true
		}
	case models.AirTransport:
		switch packagingMode {
		case models.ULD, models.BulkAir, models.Palletized:
			validPackaging = true
		}
	case models.RailTransport:
		switch packagingMode {
		case models.ContainerRail, models.BulkRail, models.WagonLoad:
			validPackaging = true
		}
	case models.LandTransport:
		switch packagingMode {
		case models.FTL, models.LTL, models.Parcel:
			validPackaging = true
		}
	}

	if !validPackaging {
		return fmt.Errorf("invalid packaging mode %s for transport mode %s", packagingMode, transportMode)
	}

	return nil
}

func (s *CustomsRateService) calculateTransportModeCharges(rate *models.CustomsRate, mode models.TransportMode) float64 {
	if modeRate, exists := rate.TransportModeRates[mode]; exists {
		return modeRate.TerminalHandling + modeRate.SecurityFee
	}
	return 0
}

func (s *CustomsRateService) calculateHandlingFees(rate *models.CustomsRate, mode models.TransportMode) float64 {
	if modeRate, exists := rate.TransportModeRates[mode]; exists {
		total := modeRate.TerminalHandling + modeRate.SecurityFee
		if modeRate.SealFee > 0 {
			total += modeRate.SealFee
		}
		return total
	}
	return 0
}

func (s *CustomsRateService) calculatePackagingCharges(
	rate *models.CustomsRate,
	transportMode models.TransportMode,
	packagingMode models.PackagingMode,
	packagingDetails interface{},
) float64 {
	modeRate, exists := rate.TransportModeRates[transportMode]
	if !exists {
		return 0
	}

	switch packagingMode {
	case models.FCL:
		if details, ok := packagingDetails.(models.FCLDetails); ok {
			if rate, exists := modeRate.FCLRates[details.ContainerSize]; exists {
				return rate * float64(details.Quantity)
			}
		}
	case models.LCL:
		if details, ok := packagingDetails.(models.LCLDetails); ok {
			return modeRate.LCLRate * details.CBM
		}
	case models.ULD:
		if details, ok := packagingDetails.(models.ULDDetails); ok {
			if rate, exists := modeRate.ULDRates[details.ULDType]; exists {
				return rate * float64(details.Quantity)
			}
		}
	case models.BulkAir:
		if details, ok := packagingDetails.(models.BulkAirDetails); ok {
			return modeRate.BulkAirRate * details.Weight
		}
	case models.ContainerRail:
		if details, ok := packagingDetails.(models.RailDetails); ok {
			if rate, exists := modeRate.ContainerRailRates[details.WagonType]; exists {
				return rate * float64(details.Quantity)
			}
		}
	case models.BulkRail:
		return modeRate.BulkRailRate
	case models.WagonLoad:
		if details, ok := packagingDetails.(models.RailDetails); ok {
			return modeRate.WagonRate * float64(details.Quantity)
		}
	case models.FTL:
		if details, ok := packagingDetails.(models.TruckDetails); ok {
			return modeRate.FTLRate * float64(details.Quantity)
		}
	case models.LTL:
		return modeRate.LTLRate
	case models.Parcel:
		return modeRate.ParcelRate
	}

	return 0
}

func (s *CustomsRateService) getCustomsRate(rateID string) (*models.CustomsRate, error) {
	// In a real implementation, this would query a database
	// This is a placeholder that returns a sample rate
	return &models.CustomsRate{
		ID:              rateID,
		Type:            models.OriginCustoms,
		Country:         "Sample Country",
		BasicHandling:   100.0,
		Documentation:   50.0,
		CustomsDuty:     5.0,  // 5%
		VAT:            10.0,  // 10%
		OtherTaxes:      2.0,  // 2%
		Inspection:      75.0,
		Storage:         25.0,
		DangerousGoods:  150.0,
		Currency:        "USD",
		ValidFrom:       time.Now(),
		ValidUntil:      time.Now().AddDate(0, 1, 0), // Valid for 1 month
		TransportModeRates: map[models.TransportMode]models.TransportModeRate{
			models.SeaTransport: {
				FCLRates: map[string]float64{
					"20GP": 200.0,
					"40GP": 350.0,
					"40HC": 400.0,
				},
				LCLRate: 25.0, // per CBM
				TerminalHandling: 150.0,
				SecurityFee:     50.0,
				SealFee:        15.0,
			},
			models.AirTransport: {
				ULDRates: map[string]float64{
					"AKE": 300.0,
					"PMC": 500.0,
				},
				BulkAirRate: 2.5, // per kg
				PalletRate:  150.0,
				TerminalHandling: 200.0,
				SecurityFee:     75.0,
			},
			models.RailTransport: {
				ContainerRailRates: map[string]float64{
					"20ft": 180.0,
					"40ft": 320.0,
				},
				BulkRailRate: 15.0, // per ton
				WagonRate:    400.0,
				TerminalHandling: 120.0,
				SecurityFee:     40.0,
			},
			models.LandTransport: {
				FTLRate:    300.0,
				LTLRate:    1.5, // per kg
				ParcelRate: 20.0,
				TerminalHandling: 100.0,
				SecurityFee:     30.0,
			},
		},
	}, nil
}
