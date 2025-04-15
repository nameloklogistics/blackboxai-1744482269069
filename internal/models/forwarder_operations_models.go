package models

import (
	"time"
)

// QuoteRequest represents a request for freight rate quotation
type QuoteRequest struct {
	ID          string         `json:"id"`
	ForwarderID string         `json:"forwarder_id"`
	CustomerID  string         `json:"customer_id"`
	
	// Transport Details
	Mode        TransportMode  `json:"mode"`
	ServiceType ServiceCategory `json:"service_type"`
	
	// Route
	Origin      Location       `json:"origin"`
	Destination Location       `json:"destination"`
	ViaPoints   []Location    `json:"via_points,omitempty"`
	
	// Schedule
	ReadyDate   time.Time     `json:"ready_date"`
	DeadlineDate time.Time    `json:"deadline_date"`
	
	// Cargo Details
	CargoType   string        `json:"cargo_type"`
	Incoterms   string        `json:"incoterms"`
	Cargo       []CargoUnit   `json:"cargo"`
	
	// Service Requirements
	RequiredServices []string  `json:"required_services"` // customs, insurance, etc.
	SpecialHandling []string   `json:"special_handling"`  // dangerous goods, temperature controlled, etc.
	
	CreatedAt   time.Time     `json:"created_at"`
	Status      string        `json:"status"` // NEW, PROCESSING, QUOTED, EXPIRED
}

// Location represents a service point in the transport chain
type Location struct {
	Type        string        `json:"type"` // AIRPORT, SEAPORT, DEPOT, ADDRESS
	Code        string        `json:"code"` // IATA/UN LOCODE/Depot Code
	Name        string        `json:"name"`
	Country     string        `json:"country"`
	Address     string        `json:"address,omitempty"`
}

// CargoUnit represents individual cargo units
type CargoUnit struct {
	Type        string        `json:"type"` // CONTAINER, PALLET, CARTON, etc.
	Quantity    int           `json:"quantity"`
	Weight      float64       `json:"weight"`
	Volume      float64       `json:"volume"`
	Length      float64       `json:"length,omitempty"`
	Width       float64       `json:"width,omitempty"`
	Height      float64       `json:"height,omitempty"`
	StackableFlag bool        `json:"stackable_flag"`
}

// FreightQuote represents a detailed freight quotation
type FreightQuote struct {
	ID          string        `json:"id"`
	RequestID   string        `json:"request_id"`
	ForwarderID string        `json:"forwarder_id"`
	
	// Service Details
	Mode        TransportMode  `json:"mode"`
	ServiceType ServiceCategory `json:"service_type"`
	TransitTime string         `json:"transit_time"`
	Routing     []RouteSegment `json:"routing"`
	
	// Rate Components
	BaseCharges RateComponent  `json:"base_charges"`
	LocalCharges RateComponent `json:"local_charges"`
	AdditionalCharges []RateComponent `json:"additional_charges"`
	
	// Totals
	SubTotal    float64        `json:"sub_total"`
	Tax         float64        `json:"tax"`
	TotalAmount float64        `json:"total_amount"`
	Currency    string         `json:"currency"`
	
	// Terms
	ValidUntil  time.Time      `json:"valid_until"`
	PaymentTerms string        `json:"payment_terms"`
	Remarks     string         `json:"remarks"`
	
	CreatedAt   time.Time      `json:"created_at"`
	Status      string         `json:"status"` // DRAFT, SENT, ACCEPTED, REJECTED
}

// RouteSegment represents a segment in the transport route
type RouteSegment struct {
	SegmentType string         `json:"segment_type"` // MAIN_CARRIAGE, PRE_CARRIAGE, ON_CARRIAGE
	Mode        TransportMode  `json:"mode"`
	From        Location       `json:"from"`
	To          Location       `json:"to"`
	Carrier     string         `json:"carrier,omitempty"`
	Service     string         `json:"service,omitempty"`
	TransitTime string         `json:"transit_time"`
	Schedule    SegmentSchedule `json:"schedule,omitempty"`
}

// SegmentSchedule represents timing for a route segment
type SegmentSchedule struct {
	Departure   time.Time      `json:"departure"`
	Arrival     time.Time      `json:"arrival"`
	Frequency   string         `json:"frequency,omitempty"`
}

// RateComponent represents a charge component in the quote
type RateComponent struct {
	ChargeType  string        `json:"charge_type"`
	Description string        `json:"description"`
	Basis       string        `json:"basis"` // per container, per kg, per cbm
	UnitRate    float64       `json:"unit_rate"`
	Units       float64       `json:"units"`
	Amount      float64       `json:"amount"`
	Currency    string        `json:"currency"`
	Mandatory   bool          `json:"mandatory"`
}

// RateConfirmation represents a confirmed freight rate
type RateConfirmation struct {
	ID          string        `json:"id"`
	QuoteID     string        `json:"quote_id"`
	ForwarderID string        `json:"forwarder_id"`
	CustomerID  string        `json:"customer_id"`
	
	// Confirmed Details
	Quote       FreightQuote  `json:"quote"`
	
	// Additional Terms
	SpecialTerms []string     `json:"special_terms"`
	Restrictions []string     `json:"restrictions"`
	
	ConfirmedAt time.Time     `json:"confirmed_at"`
	ValidUntil  time.Time     `json:"valid_until"`
	Status      string        `json:"status"` // ACTIVE, EXPIRED
}

// ShipmentBooking represents a confirmed transport booking
type ShipmentBooking struct {
	ID              string    `json:"id"`
	ConfirmationID  string    `json:"confirmation_id"`
	ForwarderID     string    `json:"forwarder_id"`
	CustomerID      string    `json:"customer_id"`
	
	// Booking Details
	BookingNumber   string    `json:"booking_number"`
	BookingDate     time.Time `json:"booking_date"`
	
	// Service Details
	Mode            TransportMode `json:"mode"`
	ServiceType     ServiceCategory `json:"service_type"`
	Routing         []RouteSegment `json:"routing"`
	
	// Cargo Details
	Cargo           []CargoUnit  `json:"cargo"`
	TotalWeight     float64      `json:"total_weight"`
	TotalVolume     float64      `json:"total_volume"`
	
	// Equipment Details
	Equipment       []Equipment  `json:"equipment"`
	
	// Schedule
	PickupDate      time.Time   `json:"pickup_date"`
	DeliveryDate    time.Time   `json:"delivery_date"`
	
	// Instructions
	PickupInstructions    string `json:"pickup_instructions"`
	DeliveryInstructions  string `json:"delivery_instructions"`
	CargoInstructions     string `json:"cargo_instructions"`
	CustomsInstructions   string `json:"customs_instructions"`
	
	// Documents
	RequiredDocuments     []string `json:"required_documents"`
	DocumentStatus        map[string]string `json:"document_status"`
	
	// Financial
	Amount          float64      `json:"amount"`
	Currency        string       `json:"currency"`
	PaymentStatus   string       `json:"payment_status"`
	PaymentTerms    string       `json:"payment_terms"`
	
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	Status          string       `json:"status"` // PENDING, CONFIRMED, IN_PROGRESS, COMPLETED
}

// BookingConfirmation represents a confirmed transport booking
type BookingConfirmation struct {
	ID              string    `json:"id"`
	BookingID       string    `json:"booking_id"`
	
	// Confirmation Details
	ConfirmationNumber string   `json:"confirmation_number"`
	ConfirmedSchedule  Schedule `json:"confirmed_schedule"`
	
	// Equipment Details
	AllocatedEquipment []Equipment `json:"allocated_equipment"`
	
	// Service Providers
	Carriers          map[string]string `json:"carriers"` // segment -> carrier
	Agents            map[string]string `json:"agents"`   // location -> agent
	
	// Instructions
	OperationalNotes   string    `json:"operational_notes"`
	SpecialInstructions string   `json:"special_instructions"`
	
	ConfirmedAt       time.Time  `json:"confirmed_at"`
	Status            string     `json:"status"`
}

// Schedule represents detailed transport schedule
type Schedule struct {
	Pickup     ScheduleEvent `json:"pickup"`
	Departure  ScheduleEvent `json:"departure"`
	Arrival    ScheduleEvent `json:"arrival"`
	Delivery   ScheduleEvent `json:"delivery"`
}

// ScheduleEvent represents a scheduled event
type ScheduleEvent struct {
	Location   Location   `json:"location"`
	Planned    time.Time  `json:"planned"`
	Estimated  time.Time  `json:"estimated"`
	Actual     time.Time  `json:"actual,omitempty"`
}
