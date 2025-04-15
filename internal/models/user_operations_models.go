package models

import (
	"time"
)

// UserType represents different types of users
type UserType string

const (
	Consignee       UserType = "CONSIGNEE"
	Shipper         UserType = "SHIPPER"
	FreightForwarder UserType = "FREIGHT_FORWARDER"
)

// ShipmentParty represents a party involved in the shipment
type ShipmentParty struct {
	ID          string    `json:"id"`
	UserType    UserType  `json:"user_type"`
	Name        string    `json:"name"`
	Company     string    `json:"company"`
	Address     string    `json:"address"`
	Country     string    `json:"country"`
	Contact     string    `json:"contact"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
}

// QuoteRequestParties represents parties involved in a quote request
type QuoteRequestParties struct {
	Consignee   ShipmentParty `json:"consignee"`
	Shipper     ShipmentParty `json:"shipper"`
	Forwarder   ShipmentParty `json:"forwarder,omitempty"`
}

// UserQuoteRequest represents a quote request from any user type
type UserQuoteRequest struct {
	ID          string             `json:"id"`
	RequestedBy UserType           `json:"requested_by"`
	Parties     QuoteRequestParties `json:"parties"`
	
	// Transport Details
	Mode        TransportMode      `json:"mode"`
	ServiceType ServiceCategory    `json:"service_type"`
	
	// Route
	Origin      Location          `json:"origin"`
	Destination Location          `json:"destination"`
	ViaPoints   []Location        `json:"via_points,omitempty"`
	
	// Schedule
	ReadyDate   time.Time         `json:"ready_date"`
	DeadlineDate time.Time        `json:"deadline_date"`
	
	// Cargo Details
	CargoType   string            `json:"cargo_type"`
	Incoterms   string            `json:"incoterms"`
	Cargo       []CargoUnit       `json:"cargo"`
	
	// Service Requirements
	RequiredServices []string      `json:"required_services"`
	SpecialHandling []string       `json:"special_handling"`
	
	CreatedAt   time.Time         `json:"created_at"`
	Status      string            `json:"status"` // NEW, PROCESSING, QUOTED, EXPIRED
}

// UserQuoteResponse represents a quote response to any user type
type UserQuoteResponse struct {
	ID          string            `json:"id"`
	RequestID   string            `json:"request_id"`
	RespondedBy UserType          `json:"responded_by"`
	Parties     QuoteRequestParties `json:"parties"`
	
	// Service Details
	Mode        TransportMode     `json:"mode"`
	ServiceType ServiceCategory   `json:"service_type"`
	TransitTime string            `json:"transit_time"`
	Routing     []RouteSegment    `json:"routing"`
	
	// Rate Components
	BaseCharges RateComponent     `json:"base_charges"`
	LocalCharges RateComponent    `json:"local_charges"`
	AdditionalCharges []RateComponent `json:"additional_charges"`
	
	// Totals
	SubTotal    float64           `json:"sub_total"`
	Tax         float64           `json:"tax"`
	TotalAmount float64           `json:"total_amount"`
	Currency    string            `json:"currency"`
	
	// Terms
	ValidUntil  time.Time         `json:"valid_until"`
	PaymentTerms string           `json:"payment_terms"`
	Remarks     string            `json:"remarks"`
	
	CreatedAt   time.Time         `json:"created_at"`
	Status      string            `json:"status"` // DRAFT, SENT, ACCEPTED, REJECTED
}

// UserQuoteConfirmation represents a quote confirmation from any user type
type UserQuoteConfirmation struct {
	ID          string            `json:"id"`
	QuoteID     string            `json:"quote_id"`
	ConfirmedBy UserType          `json:"confirmed_by"`
	Parties     QuoteRequestParties `json:"parties"`
	
	// Confirmed Details
	Quote       UserQuoteResponse `json:"quote"`
	
	// Additional Terms
	SpecialTerms []string         `json:"special_terms"`
	Restrictions []string         `json:"restrictions"`
	
	ConfirmedAt time.Time         `json:"confirmed_at"`
	ValidUntil  time.Time         `json:"valid_until"`
	Status      string            `json:"status"` // ACTIVE, EXPIRED
}

// UserBookingRequest represents a booking request from any user type
type UserBookingRequest struct {
	ID              string         `json:"id"`
	ConfirmationID  string         `json:"confirmation_id"`
	RequestedBy     UserType       `json:"requested_by"`
	Parties         QuoteRequestParties `json:"parties"`
	
	// Booking Details
	Mode            TransportMode  `json:"mode"`
	ServiceType     ServiceCategory `json:"service_type"`
	Routing         []RouteSegment `json:"routing"`
	
	// Cargo Details
	Cargo           []CargoUnit    `json:"cargo"`
	TotalWeight     float64        `json:"total_weight"`
	TotalVolume     float64        `json:"total_volume"`
	
	// Schedule
	PickupDate      time.Time      `json:"pickup_date"`
	DeliveryDate    time.Time      `json:"delivery_date"`
	
	// Instructions
	PickupInstructions    string   `json:"pickup_instructions"`
	DeliveryInstructions  string   `json:"delivery_instructions"`
	CargoInstructions     string   `json:"cargo_instructions"`
	CustomsInstructions   string   `json:"customs_instructions"`
	
	CreatedAt       time.Time      `json:"created_at"`
	Status          string         `json:"status"` // PENDING, CONFIRMED
}

// UserBookingConfirmation represents a booking confirmation for any user type
type UserBookingConfirmation struct {
	ID              string         `json:"id"`
	BookingID       string         `json:"booking_id"`
	ConfirmedBy     UserType       `json:"confirmed_by"`
	Parties         QuoteRequestParties `json:"parties"`
	
	// Confirmation Details
	BookingNumber   string         `json:"booking_number"`
	ConfirmedSchedule Schedule     `json:"confirmed_schedule"`
	
	// Equipment Details
	AllocatedEquipment []Equipment `json:"allocated_equipment"`
	
	// Service Providers
	Carriers          map[string]string `json:"carriers"` // segment -> carrier
	Agents            map[string]string `json:"agents"`   // location -> agent
	
	// Financial
	Amount            float64       `json:"amount"`
	Currency          string        `json:"currency"`
	PaymentTerms      string        `json:"payment_terms"`
	
	// Instructions
	OperationalNotes   string       `json:"operational_notes"`
	SpecialInstructions string      `json:"special_instructions"`
	
	ConfirmedAt       time.Time     `json:"confirmed_at"`
	Status            string        `json:"status"` // CONFIRMED, CANCELLED
}
