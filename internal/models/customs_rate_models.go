package models

import (
	"fmt"
	"time"
)

// Address represents a physical address
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`     // ISO country code
}

// CustomsBrokerBranch represents a branch office of a customs broker
type CustomsBrokerBranch struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`           // Branch code
	Type          string    `json:"type"`           // HQ, BRANCH, SATELLITE
	Address       Address   `json:"address"`
	ContactPerson string    `json:"contact_person"`
	ContactEmail  string    `json:"contact_email"`
	ContactPhone  string    `json:"contact_phone"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// LicenseVerificationStatus represents the verification status of a license
type LicenseVerificationStatus string

const (
	LicenseVerified    LicenseVerificationStatus = "VERIFIED"
	LicenseUnverified  LicenseVerificationStatus = "UNVERIFIED"
	LicensePending     LicenseVerificationStatus = "PENDING"
	LicenseInvalid     LicenseVerificationStatus = "INVALID"
)

// LicenseVerification represents the verification details of a customs broker license
type LicenseVerification struct {
	VerificationID     string                   `json:"verification_id"`
	AuthorityID        string                   `json:"authority_id"`     // ID of the licensing authority
	AuthorityName      string                   `json:"authority_name"`   // Name of the licensing authority
	VerificationStatus LicenseVerificationStatus `json:"status"`
	VerifiedAt         time.Time                `json:"verified_at"`
	ValidUntil         time.Time                `json:"valid_until"`
	LastCheckedAt      time.Time                `json:"last_checked_at"`
	NextCheckDue       time.Time                `json:"next_check_due"`
}

// LicenseAuthority represents a customs licensing authority
type LicenseAuthority struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CountryCode string    `json:"country_code"`
	APIEndpoint string    `json:"api_endpoint"`
	PublicKey   string    `json:"public_key"`     // For verifying authority signatures
}

// CustomsBroker represents a registered customs broker
type CustomsBroker struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	RegistrationNo  string              `json:"registration_no"`
	CountryCode     string              `json:"country_code"`     // ISO country code where broker is registered
	LicenseNumber   string              `json:"license_number"`   // Customs broker license number
	LicenseExpiry   time.Time           `json:"license_expiry"`
	Status          string              `json:"status"`           // ACTIVE, SUSPENDED, EXPIRED
	HeadOffice      Address             `json:"head_office"`
	BranchOffices   []CustomsBrokerBranch `json:"branch_offices"`
	
	// License Verification
	LicenseAuthority   string               `json:"license_authority"`   // ID of the issuing authority
	LicenseClass       string               `json:"license_class"`       // Class/Type of license
	LicenseScope       []string             `json:"license_scope"`       // Authorized activities
	LicenseVerification LicenseVerification  `json:"license_verification"`
	
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

// LicenseVerificationError represents an error during license verification
type LicenseVerificationError struct {
	BrokerID      string
	LicenseNumber string
	AuthorityID   string
	Reason        string
}

func (e *LicenseVerificationError) Error() string {
	return fmt.Sprintf("license verification failed for broker %s (license: %s) with authority %s: %s",
		e.BrokerID, e.LicenseNumber, e.AuthorityID, e.Reason)
}

// CustomsRateType represents different types of customs clearance rates
type CustomsRateType string

const (
	OriginCustoms       CustomsRateType = "ORIGIN"
	TransshipmentCustoms CustomsRateType = "TRANSSHIPMENT"
	TransitCustoms      CustomsRateType = "TRANSIT"
	DestinationCustoms  CustomsRateType = "DESTINATION"
)

// CustomsRate represents the rate structure for customs clearance
type CustomsRate struct {
	ID          string         `json:"id"`
	Type        CustomsRateType `json:"type"`
	Country     string         `json:"country"`
	Port        string         `json:"port,omitempty"`
	Airport     string         `json:"airport,omitempty"`
	
	// Transport Mode Specific Rates
	TransportModeRates map[TransportMode]TransportModeRate `json:"transport_mode_rates"`
	
	// Base Charges
	BasicHandling float64 `json:"basic_handling"`
	Documentation float64 `json:"documentation"`
	
	// Government Fees
	CustomsDuty   float64 `json:"customs_duty"` // Percentage
	VAT           float64 `json:"vat"`          // Percentage
	OtherTaxes    float64 `json:"other_taxes"`  // Percentage
	
	// Additional Services
	Inspection    float64 `json:"inspection"`
	Storage       float64 `json:"storage_per_day"`
	ExamFees      float64 `json:"examination_fees"`
	
	// Special Handling
	DangerousGoods float64 `json:"dangerous_goods_fee"`
	Refrigerated   float64 `json:"refrigerated_cargo_fee"`
	OversizeCargo  float64 `json:"oversize_cargo_fee"`
	
	// Time-based Charges
	PeakSeasonSurcharge float64 `json:"peak_season_surcharge"`
	ExpressHandling    float64 `json:"express_handling"`
	
	Currency     string    `json:"currency"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidUntil   time.Time `json:"valid_until"`
	ProviderID   string    `json:"provider_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TransportModeRate contains rates specific to a transport mode
type TransportModeRate struct {
	// Sea Transport Rates
	FCLRates map[string]float64 `json:"fcl_rates,omitempty"` // Key: container size (20GP, 40GP, etc.)
	LCLRate  float64           `json:"lcl_rate,omitempty"`   // Per CBM
	
	// Air Transport Rates
	ULDRates map[string]float64 `json:"uld_rates,omitempty"` // Key: ULD type (AKE, PMC, etc.)
	BulkAirRate float64        `json:"bulk_air_rate,omitempty"` // Per kg
	PalletRate  float64        `json:"pallet_rate,omitempty"`   // Per pallet
	
	// Rail Transport Rates
	ContainerRailRates map[string]float64 `json:"container_rail_rates,omitempty"` // Key: container type
	BulkRailRate float64                 `json:"bulk_rail_rate,omitempty"`        // Per ton
	WagonRate    float64                 `json:"wagon_rate,omitempty"`            // Per wagon
	
	// Land Transport Rates
	FTLRate     float64 `json:"ftl_rate,omitempty"`     // Per truck
	LTLRate     float64 `json:"ltl_rate,omitempty"`     // Per kg/CBM
	ParcelRate  float64 `json:"parcel_rate,omitempty"`  // Per parcel
	
	// Mode-specific Handling Fees
	TerminalHandling float64 `json:"terminal_handling"`
	SecurityFee      float64 `json:"security_fee"`
	SealFee         float64 `json:"seal_fee,omitempty"`
	
	// Special Equipment Charges
	SpecialEquipmentFees map[string]float64 `json:"special_equipment_fees,omitempty"`
}

// CustomsQuotation represents a quotation for customs clearance
type CustomsQuotation struct {
	BrokerID       string         `json:"broker_id"`
	BranchOfficeID string         `json:"branch_office_id"` // Branch office providing the quote
	ID            string         `json:"id"`
	RateID        string         `json:"rate_id"`
	Type          CustomsRateType `json:"type"`
	
	// Transport and Packaging Information
	TransportMode TransportMode   `json:"transport_mode"`
	PackagingMode PackagingMode   `json:"packaging_mode"`
	
	// Cargo Information
	CargoValue    float64        `json:"cargo_value"`
	CargoType     string         `json:"cargo_type"`
	HSCode        string         `json:"hs_code"`
	Weight        float64        `json:"weight"`
	Volume        float64        `json:"volume"`
	
	// Packaging Details
	PackagingDetails interface{} `json:"packaging_details"` // FCLDetails, ULDDetails, etc.
	
	// Calculated Charges
	BasicCharges  float64        `json:"basic_charges"`
	CustomsDuty   float64        `json:"customs_duty"`
	VAT           float64        `json:"vat"`
	OtherTaxes    float64        `json:"other_taxes"`
	
	// Mode-specific Charges
	TransportModeCharges float64  `json:"transport_mode_charges"`
	PackagingCharges    float64   `json:"packaging_charges"`
	HandlingFees        float64   `json:"handling_fees"`
	
	AdditionalFees float64       `json:"additional_fees"`
	TotalAmount   float64        `json:"total_amount"`
	
	// Breakdown of Additional Fees
	FeeBreakdown  []CustomsFee   `json:"fee_breakdown"`
	
	Currency      string         `json:"currency"`
	ValidUntil    time.Time      `json:"valid_until"`
	CreatedAt     time.Time      `json:"created_at"`
}

// Packaging Details Structures
type FCLDetails struct {
	ContainerSize string `json:"container_size"` // 20GP, 40GP, etc.
	ContainerType string `json:"container_type"` // Standard, Reefer, etc.
	Quantity      int    `json:"quantity"`
}

type LCLDetails struct {
	CBM           float64 `json:"cbm"`
	Pieces        int     `json:"pieces"`
	Palletized    bool    `json:"palletized"`
}

type ULDDetails struct {
	ULDType     string  `json:"uld_type"`     // AKE, PMC, etc.
	Quantity    int     `json:"quantity"`
	BuildUpRequired bool `json:"build_up_required"`
}

type BulkAirDetails struct {
	Pieces    int     `json:"pieces"`
	Weight    float64 `json:"weight"`
	Palletized bool   `json:"palletized"`
}

type RailDetails struct {
	WagonType string  `json:"wagon_type"`
	Quantity  int     `json:"quantity"`
	LoadType  string  `json:"load_type"` // Container, Bulk, etc.
}

type TruckDetails struct {
	TruckType string  `json:"truck_type"`
	Quantity  int     `json:"quantity"`
	LoadType  string  `json:"load_type"` // FTL, LTL, etc.
}

type ParcelDetails struct {
	Pieces    int       `json:"pieces"`
	Weight    float64   `json:"weight"`
	Express   bool      `json:"express"`
}

// CustomsFee represents a single fee component in customs clearance
type CustomsFee struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	IsRequired  bool    `json:"is_required"`
}

// CustomsRateCalculator calculates customs charges based on cargo details
type CustomsRateCalculator struct {
	Rate           *CustomsRate
	CargoDetails   *CargoDetails
	TransportMode  TransportMode
	PackagingMode  PackagingMode
	PackagingDetails interface{}
}

func (c *CustomsRateCalculator) CalculateBasicCharges() float64 {
	baseCharges := c.Rate.BasicHandling + c.Rate.Documentation
	
	// Add transport mode specific charges
	if modeRate, exists := c.Rate.TransportModeRates[c.TransportMode]; exists {
		baseCharges += modeRate.TerminalHandling + modeRate.SecurityFee
		if modeRate.SealFee > 0 {
			baseCharges += modeRate.SealFee
		}
	}
	
	return baseCharges
}

func (c *CustomsRateCalculator) CalculateCustomsDuty(cargoValue float64) float64 {
	return cargoValue * (c.Rate.CustomsDuty / 100)
}

func (c *CustomsRateCalculator) CalculateVAT(cargoValue float64) float64 {
	return cargoValue * (c.Rate.VAT / 100)
}

func (c *CustomsRateCalculator) CalculateOtherTaxes(cargoValue float64) float64 {
	return cargoValue * (c.Rate.OtherTaxes / 100)
}

func (c *CustomsRateCalculator) CalculateAdditionalFees() []CustomsFee {
	fees := make([]CustomsFee, 0)

	// Add inspection fee if applicable
	if c.Rate.Inspection > 0 {
		fees = append(fees, CustomsFee{
			Description: "Inspection Fee",
			Amount:      c.Rate.Inspection,
			Currency:    c.Rate.Currency,
			IsRequired:  true,
		})
	}

	// Add storage fee if applicable
	if c.Rate.Storage > 0 {
		fees = append(fees, CustomsFee{
			Description: "Storage Fee (per day)",
			Amount:      c.Rate.Storage,
			Currency:    c.Rate.Currency,
			IsRequired:  false,
		})
	}

	// Add dangerous goods fee if applicable
	if c.CargoDetails.IsHazardous && c.Rate.DangerousGoods > 0 {
		fees = append(fees, CustomsFee{
			Description: "Dangerous Goods Handling",
			Amount:      c.Rate.DangerousGoods,
			Currency:    c.Rate.Currency,
			IsRequired:  true,
		})
	}

	// Add transport mode specific fees
	if modeRate, exists := c.Rate.TransportModeRates[c.TransportMode]; exists {
		// Add special equipment fees if applicable
		for equipment, fee := range modeRate.SpecialEquipmentFees {
			fees = append(fees, CustomsFee{
				Description: fmt.Sprintf("%s Equipment Fee", equipment),
				Amount:      fee,
				Currency:    c.Rate.Currency,
				IsRequired:  true,
			})
		}
	}

	// Add packaging mode specific fees
	fees = append(fees, c.calculatePackagingFees()...)

	// Add peak season surcharge if applicable
	if isPeakSeason() && c.Rate.PeakSeasonSurcharge > 0 {
		fees = append(fees, CustomsFee{
			Description: "Peak Season Surcharge",
			Amount:      c.Rate.PeakSeasonSurcharge,
			Currency:    c.Rate.Currency,
			IsRequired:  true,
		})
	}

	return fees
}

func (c *CustomsRateCalculator) calculatePackagingFees() []CustomsFee {
	fees := make([]CustomsFee, 0)
	modeRate, exists := c.Rate.TransportModeRates[c.TransportMode]
	if !exists {
		return fees
	}

	switch c.PackagingMode {
	case FCL:
		if details, ok := c.PackagingDetails.(FCLDetails); ok {
			if rate, exists := modeRate.FCLRates[details.ContainerSize]; exists {
				fees = append(fees, CustomsFee{
					Description: fmt.Sprintf("FCL Handling - %s", details.ContainerSize),
					Amount:      rate * float64(details.Quantity),
					Currency:    c.Rate.Currency,
					IsRequired:  true,
				})
			}
		}
	
	case LCL:
		if details, ok := c.PackagingDetails.(LCLDetails); ok {
			fees = append(fees, CustomsFee{
				Description: "LCL Handling",
				Amount:      modeRate.LCLRate * details.CBM,
				Currency:    c.Rate.Currency,
				IsRequired:  true,
			})
		}
	
	case ULD:
		if details, ok := c.PackagingDetails.(ULDDetails); ok {
			if rate, exists := modeRate.ULDRates[details.ULDType]; exists {
				fees = append(fees, CustomsFee{
					Description: fmt.Sprintf("ULD Handling - %s", details.ULDType),
					Amount:      rate * float64(details.Quantity),
					Currency:    c.Rate.Currency,
					IsRequired:  true,
				})
			}
		}
	
	case BulkAir:
		if details, ok := c.PackagingDetails.(BulkAirDetails); ok {
			fees = append(fees, CustomsFee{
				Description: "Bulk Air Cargo Handling",
				Amount:      modeRate.BulkAirRate * details.Weight,
				Currency:    c.Rate.Currency,
				IsRequired:  true,
			})
		}
	
	case FTL:
		if details, ok := c.PackagingDetails.(TruckDetails); ok {
			fees = append(fees, CustomsFee{
				Description: "FTL Handling",
				Amount:      modeRate.FTLRate * float64(details.Quantity),
				Currency:    c.Rate.Currency,
				IsRequired:  true,
			})
		}
	
	case LTL:
		if details, ok := c.PackagingDetails.(TruckDetails); ok {
			fees = append(fees, CustomsFee{
				Description: "LTL Handling",
				Amount:      modeRate.LTLRate * c.CargoDetails.Weight,
				Currency:    c.Rate.Currency,
				IsRequired:  true,
			})
		}
	
	case ContainerRail:
		if details, ok := c.PackagingDetails.(RailDetails); ok {
			if rate, exists := modeRate.ContainerRailRates[details.WagonType]; exists {
				fees = append(fees, CustomsFee{
					Description: fmt.Sprintf("Rail Container Handling - %s", details.WagonType),
					Amount:      rate * float64(details.Quantity),
					Currency:    c.Rate.Currency,
					IsRequired:  true,
				})
			}
		}
	
	case BulkRail:
		if details, ok := c.PackagingDetails.(RailDetails); ok {
			fees = append(fees, CustomsFee{
				Description: "Bulk Rail Cargo Handling",
				Amount:      modeRate.BulkRailRate * c.CargoDetails.Weight,
				Currency:    c.Rate.Currency,
				IsRequired:  true,
			})
		}
	
	case WagonLoad:
		if details, ok := c.PackagingDetails.(RailDetails); ok {
			fees = append(fees, CustomsFee{
				Description: fmt.Sprintf("Wagon Load Handling - %s", details.WagonType),
				Amount:      modeRate.WagonRate * float64(details.Quantity),
				Currency:    c.Rate.Currency,
				IsRequired:  true,
			})
		}
	}

	return fees
}

func (c *CustomsRateCalculator) CalculateTotalAmount(cargoValue float64) float64 {
	// Base charges
	total := c.CalculateBasicCharges()
	
	// Government charges
	total += c.CalculateCustomsDuty(cargoValue)
	total += c.CalculateVAT(cargoValue)
	total += c.CalculateOtherTaxes(cargoValue)

	// Transport mode specific charges
	if modeRate, exists := c.Rate.TransportModeRates[c.TransportMode]; exists {
		total += modeRate.TerminalHandling
		total += modeRate.SecurityFee
		if modeRate.SealFee > 0 {
			total += modeRate.SealFee
		}
	}

	// Additional fees including packaging-specific charges
	fees := c.CalculateAdditionalFees()
	for _, fee := range fees {
		if fee.IsRequired {
			total += fee.Amount
		}
	}

	return total
}

// Helper function to determine peak season
func isPeakSeason() bool {
	currentMonth := time.Now().Month()
	return currentMonth >= time.October && currentMonth <= time.December
}
