package models

import (
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

// BranchOffice represents a branch office of a freight forwarder
type BranchOffice struct {
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

// FreightForwarder represents a registered freight forwarding company
type FreightForwarder struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	RegistrationNo  string         `json:"registration_no"`
	CountryCode     string         `json:"country_code"`     // ISO country code where forwarder is registered
	LicenseNumber   string         `json:"license_number"`   // Customs broker license number
	LicenseExpiry   time.Time      `json:"license_expiry"`
	Status          string         `json:"status"`           // ACTIVE, SUSPENDED, EXPIRED
	HeadOffice      Address        `json:"head_office"`
	BranchOffices   []BranchOffice `json:"branch_offices"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// BranchOfficeRequest represents a request to add a new branch office
type BranchOfficeRequest struct {
	ForwarderID    string    `json:"forwarder_id"`
	Name           string    `json:"name"`
	Code           string    `json:"code"`
	Type           string    `json:"type"`
	Address        Address   `json:"address"`
	ContactPerson  string    `json:"contact_person"`
	ContactEmail   string    `json:"contact_email"`
	ContactPhone   string    `json:"contact_phone"`
}

// RateRequest represents a request for rate quotation
type RateRequest struct {
	ID              string         `json:"id"`
	ForwarderID     string         `json:"forwarder_id"`
	BranchOfficeID  string         `json:"branch_office_id"` // Branch office handling the request
	ShipperID       string         `json:"shipper_id"`
	
	// Route Information
	Origin          string         `json:"origin"`
	OriginCountry   string         `json:"origin_country"`      // ISO country code
	Destination     string         `json:"destination"`
	DestCountry     string         `json:"destination_country"` // ISO country code
	TransportMode   TransportMode  `json:"transport_mode"`
	ServiceType     ServiceCategory `json:"service_type"` // Import, Export, Transit, Transshipment
	
	// Cargo Details
	CargoDetails    CargoDetails   `json:"cargo_details"`
	RequiredServices []string      `json:"required_services"` // e.g., customs clearance, insurance
	
	// Timing
	ReadyDate       time.Time      `json:"ready_date"`
	DeliveryDate    time.Time      `json:"delivery_date"`
	CreatedAt       time.Time      `json:"created_at"`
}

// CountryRestrictionError represents an error when a forwarder attempts to quote outside their registered country
type CountryRestrictionError struct {
	ForwarderID     string
	ForwarderCountry string
	RequestedCountry string
}

func (e *CountryRestrictionError) Error() string {
	return fmt.Sprintf("forwarder %s registered in %s cannot provide quotes for %s", 
		e.ForwarderID, e.ForwarderCountry, e.RequestedCountry)
}

// RateQuotation represents a rate quotation from a freight forwarder
type RateQuotation struct {
	ID              string         `json:"id"`
	RequestID       string         `json:"request_id"`
	ForwarderID     string         `json:"forwarder_id"`
	BranchOfficeID  string         `json:"branch_office_id"` // Branch office providing the quote
	
	// Service Details
	TransportMode   TransportMode  `json:"transport_mode"`
	ServiceType     ServiceCategory `json:"service_type"`
	Route           []string       `json:"route"` // Via points
	TransitTime     string         `json:"transit_time"`
	
	// Country Information
	OriginCountry   string         `json:"origin_country"`      // ISO country code
	DestCountry     string         `json:"destination_country"` // ISO country code
	ForwarderCountry string        `json:"forwarder_country"`   // ISO country code of registered country
	
	// Rate Components
	FreightCharges  float64        `json:"freight_charges"`
	LocalCharges    float64        `json:"local_charges"`
	CustomsCharges  float64        `json:"customs_charges"`
	OtherCharges    []Charge       `json:"other_charges"`
	TotalAmount     float64        `json:"total_amount"`
	Currency        string         `json:"currency"`
	
	// Validity
	ValidUntil      time.Time      `json:"valid_until"`
	CreatedAt       time.Time      `json:"created_at"`
	Status          string         `json:"status"` // DRAFT, SENT, ACCEPTED, REJECTED
}

// RateConfirmation represents a confirmed rate
type RateConfirmation struct {
	ID              string         `json:"id"`
	QuotationID     string         `json:"quotation_id"`
	ForwarderID     string         `json:"forwarder_id"`
	ShipperID       string         `json:"shipper_id"`
	
	// Confirmed Details
	FreightCharges  float64        `json:"freight_charges"`
	LocalCharges    float64        `json:"local_charges"`
	CustomsCharges  float64        `json:"customs_charges"`
	OtherCharges    []Charge       `json:"other_charges"`
	TotalAmount     float64        `json:"total_amount"`
	Currency        string         `json:"currency"`
	
	// Terms
	PaymentTerms    string         `json:"payment_terms"`
	ValidityPeriod  time.Time      `json:"validity_period"`
	SpecialTerms    []string       `json:"special_terms"`
	
	ConfirmedAt     time.Time      `json:"confirmed_at"`
	Status          string         `json:"status"` // CONFIRMED, EXPIRED
}

// ShipmentBooking represents a booking made with confirmed rates
type ShipmentBooking struct {
	ID                  string         `json:"id"`
	ConfirmationID      string         `json:"confirmation_id"`
	ForwarderID         string         `json:"forwarder_id"`
	ShipperID           string         `json:"shipper_id"`
	
	// Transport Details
	TransportMode       TransportMode  `json:"transport_mode"`
	ServiceType         ServiceCategory `json:"service_type"`
	Origin             string         `json:"origin"`
	Destination        string         `json:"destination"`
	Route              []string       `json:"route"`
	
	// Schedule
	RequestedPickup     time.Time      `json:"requested_pickup"`
	RequestedDelivery   time.Time      `json:"requested_delivery"`
	EstimatedDeparture  time.Time      `json:"estimated_departure"`
	EstimatedArrival    time.Time      `json:"estimated_arrival"`
	
	// Cargo Details
	CargoDetails        CargoDetails   `json:"cargo_details"`
	Equipment           []Equipment    `json:"equipment"`
	
	// Documents
	RequiredDocuments   []string       `json:"required_documents"`
	DocumentStatus      map[string]string `json:"document_status"`
	
	// Service Instructions
	PickupInstructions  string         `json:"pickup_instructions"`
	DeliveryInstructions string        `json:"delivery_instructions"`
	CustomsInstructions string         `json:"customs_instructions"`
	
	// Financial
	BookingAmount       float64        `json:"booking_amount"`
	Currency           string         `json:"currency"`
	PaymentStatus      string         `json:"payment_status"`
	
	// Status
	Status             string         `json:"status"` // PENDING, CONFIRMED, CANCELLED
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

// BookingConfirmation represents a confirmed booking
type BookingConfirmation struct {
	ID                string         `json:"id"`
	BookingID         string         `json:"booking_id"`
	
	// Confirmation Details
	BookingNumber     string         `json:"booking_number"`
	ConfirmedSchedule struct {
		Pickup    time.Time      `json:"pickup"`
		Departure time.Time      `json:"departure"`
		Arrival   time.Time      `json:"arrival"`
		Delivery  time.Time      `json:"delivery"`
	} `json:"confirmed_schedule"`
	
	// Equipment Details
	AllocatedEquipment []Equipment    `json:"allocated_equipment"`
	
	// Service Details
	ServiceProviders   map[string]string `json:"service_providers"` // service type -> provider ID
	
	// Instructions
	OperationalNotes   string         `json:"operational_notes"`
	SpecialInstructions string        `json:"special_instructions"`
	
	ConfirmedAt       time.Time      `json:"confirmed_at"`
	Status            string         `json:"status"`
}

// RateCalculation represents rate calculation components
type RateCalculation struct {
	// Base Rates
	FreightRate     float64 `json:"freight_rate"`
	LocalRate       float64 `json:"local_rate"`
	
	// Surcharges
	FuelSurcharge   float64 `json:"fuel_surcharge"`
	SecurityCharge  float64 `json:"security_charge"`
	PeakSeasonCharge float64 `json:"peak_season_charge"`
	
	// Additional Services
	CustomsClearance float64 `json:"customs_clearance"`
	Insurance       float64 `json:"insurance"`
	Documentation   float64 `json:"documentation"`
	
	// Calculations
	SubTotal       float64 `json:"sub_total"`
	Tax            float64 `json:"tax"`
	Total          float64 `json:"total"`
	
	// Breakdown
	ChargeBreakdown []Charge `json:"charge_breakdown"`
}

// BookingInstruction represents detailed booking instructions
type BookingInstruction struct {
	BookingID       string    `json:"booking_id"`
	InstructionType string    `json:"instruction_type"` // PICKUP, DELIVERY, CUSTOMS, HANDLING
	Instructions    string    `json:"instructions"`
	Attachments     []string  `json:"attachments"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedBy       string    `json:"updated_by"`
}

// DocumentRequirement represents required documents for a booking
type DocumentRequirement struct {
	BookingID       string    `json:"booking_id"`
	DocumentType    string    `json:"document_type"`
	Required        bool      `json:"required"`
	Deadline        time.Time `json:"deadline"`
	Instructions    string    `json:"instructions"`
	Status          string    `json:"status"` // PENDING, SUBMITTED, APPROVED, REJECTED
}
