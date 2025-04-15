package marketplace

import (
	"github.com/stellar/soroban-sdk/go/soroban"
)

// ServiceType represents different types of logistics services
type ServiceType uint8

const (
	FreightForwarding ServiceType = iota
	CustomsBrokerage
	Shipping
	Transshipment
)

// ShipmentMode represents different modes of transport
type ShipmentMode uint8

const (
	Air ShipmentMode = iota
	Sea
	Road
	Rail
)

// ServiceListing represents a service offered in the marketplace
type ServiceListing struct {
	ID            string
	Provider      string
	ServiceType   ServiceType
	ShipmentMode  ShipmentMode
	Origin        string
	Destination   string
	Rate          uint64
	Description   string
	IsActive      bool
}

// Booking represents a cargo booking
type Booking struct {
	ID            string
	ServiceID     string
	Customer      string
	Provider      string
	Status        string
	CargoDetails  CargoDetails
	PaymentStatus string
	TrackingInfo  []TrackingEvent
}

// CargoDetails contains information about the cargo
type CargoDetails struct {
	Weight      float64
	Volume      float64
	Type        string
	Description string
	Hazardous   bool
}

// TrackingEvent represents a shipment tracking event
type TrackingEvent struct {
	Timestamp   uint64
	Location    string
	Status      string
	Description string
}

type MarketplaceContract struct {
	soroban.Contract
	tokenContract *TokenContract
	listings      map[string]ServiceListing
	bookings     map[string]Booking
}

func (c *MarketplaceContract) Initialize(env soroban.Env, tokenContractAddress string) {
	if c.tokenContract != nil {
		panic("Contract already initialized")
	}

	// Initialize storage
	c.listings = make(map[string]ServiceListing)
	c.bookings = make(map[string]Booking)
}

// CreateServiceListing creates a new service listing
func (c *MarketplaceContract) CreateServiceListing(env soroban.Env, listing ServiceListing) string {
	// Verify provider is authorized
	provider := env.Current().Auth().Address()
	listing.Provider = provider.String()
	listing.IsActive = true

	// Generate unique ID
	listing.ID = env.GenerateUUID()
	
	// Store listing
	c.listings[listing.ID] = listing

	// Emit event
	env.Events().Publish("service_listed", map[string]interface{}{
		"id":       listing.ID,
		"provider": listing.Provider,
	})

	return listing.ID
}

// GetQuotation calculates shipping rate for given parameters
func (c *MarketplaceContract) GetQuotation(
	serviceID string,
	cargoDetails CargoDetails,
) uint64 {
	listing, exists := c.listings[serviceID]
	if !exists {
		return 0
	}

	// Basic rate calculation (would be more complex in production)
	baseRate := listing.Rate
	volumeMultiplier := uint64(cargoDetails.Volume * 100)
	weightMultiplier := uint64(cargoDetails.Weight * 100)

	return baseRate * (volumeMultiplier + weightMultiplier) / 100
}

// CreateBooking creates a new cargo booking
func (c *MarketplaceContract) CreateBooking(
	env soroban.Env,
	serviceID string,
	cargoDetails CargoDetails,
) string {
	listing, exists := c.listings[serviceID]
	if !exists || !listing.IsActive {
		panic("Invalid or inactive service")
	}

	customer := env.Current().Auth().Address()
	
	booking := Booking{
		ID:            env.GenerateUUID(),
		ServiceID:     serviceID,
		Customer:      customer.String(),
		Provider:      listing.Provider,
		Status:        "PENDING",
		CargoDetails:  cargoDetails,
		PaymentStatus: "UNPAID",
		TrackingInfo:  make([]TrackingEvent, 0),
	}

	c.bookings[booking.ID] = booking

	// Emit event
	env.Events().Publish("booking_created", map[string]interface{}{
		"id":       booking.ID,
		"customer": booking.Customer,
		"provider": booking.Provider,
	})

	return booking.ID
}

// ProcessPayment handles payment for a booking
func (c *MarketplaceContract) ProcessPayment(env soroban.Env, bookingID string) bool {
	booking, exists := c.bookings[bookingID]
	if !exists {
		return false
	}

	customer := env.Current().Auth().Address()
	if customer.String() != booking.Customer {
		return false
	}

	// Calculate payment amount
	amount := c.GetQuotation(booking.ServiceID, booking.CargoDetails)

	// Process payment using token contract
	success := c.tokenContract.Transfer(
		env,
		booking.Customer,
		booking.Provider,
		amount,
	)

	if success {
		booking.PaymentStatus = "PAID"
		booking.Status = "CONFIRMED"
		c.bookings[bookingID] = booking

		// Emit payment event
		env.Events().Publish("payment_processed", map[string]interface{}{
			"booking_id": bookingID,
			"amount":     amount,
		})
	}

	return success
}

// UpdateShipmentStatus updates tracking information for a booking
func (c *MarketplaceContract) UpdateShipmentStatus(
	env soroban.Env,
	bookingID string,
	location string,
	status string,
	description string,
) bool {
	booking, exists := c.bookings[bookingID]
	if !exists {
		return false
	}

	// Verify provider authorization
	provider := env.Current().Auth().Address()
	if provider.String() != booking.Provider {
		return false
	}

	event := TrackingEvent{
		Timestamp:   uint64(env.Current().Ledger().Timestamp()),
		Location:    location,
		Status:      status,
		Description: description,
	}

	booking.TrackingInfo = append(booking.TrackingInfo, event)
	booking.Status = status
	c.bookings[bookingID] = booking

	// Emit tracking update event
	env.Events().Publish("tracking_updated", map[string]interface{}{
		"booking_id":  bookingID,
		"status":      status,
		"location":    location,
		"timestamp":   event.Timestamp,
	})

	return true
}

// GetBooking retrieves booking information
func (c *MarketplaceContract) GetBooking(bookingID string) *Booking {
	booking, exists := c.bookings[bookingID]
	if !exists {
		return nil
	}
	return &booking
}

// GetServiceListing retrieves service listing information
func (c *MarketplaceContract) GetServiceListing(listingID string) *ServiceListing {
	listing, exists := c.listings[listingID]
	if !exists {
		return nil
	}
	return &listing
}
