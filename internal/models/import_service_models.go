package models

import (
	"time"
)

// ImportServiceType represents different types of import services
type ImportServiceType string

const (
	// Sea Import Services
	SeaFCLImport    ImportServiceType = "SEA_FCL_IMPORT"
	SeaLCLImport    ImportServiceType = "SEA_LCL_IMPORT"
	SeaBulkImport   ImportServiceType = "SEA_BULK_IMPORT"
	SeaRoRoImport   ImportServiceType = "SEA_RORO_IMPORT"

	// Air Import Services
	AirGeneralImport    ImportServiceType = "AIR_GENERAL_IMPORT"
	AirExpressImport    ImportServiceType = "AIR_EXPRESS_IMPORT"
	AirCharterImport    ImportServiceType = "AIR_CHARTER_IMPORT"
	AirConsolidatedImport ImportServiceType = "AIR_CONSOLIDATED_IMPORT"

	// Rail Import Services
	RailContainerImport ImportServiceType = "RAIL_CONTAINER_IMPORT"
	RailBulkImport     ImportServiceType = "RAIL_BULK_IMPORT"
	RailWagonImport    ImportServiceType = "RAIL_WAGON_IMPORT"

	// Land Import Services
	LandFTLImport      ImportServiceType = "LAND_FTL_IMPORT"
	LandLTLImport      ImportServiceType = "LAND_LTL_IMPORT"
	LandParcelImport   ImportServiceType = "LAND_PARCEL_IMPORT"
)

// ImportService represents a base import service
type ImportService struct {
	ID            string           `json:"id"`
	Type          ImportServiceType `json:"type"`
	Mode          TransportMode    `json:"mode"`
	
	// Route Information
	Origin        Location         `json:"origin"`
	Destination   Location         `json:"destination"`
	ViaPoints     []Location       `json:"via_points,omitempty"`
	
	// Service Details
	TransitTime   string           `json:"transit_time"`
	Frequency     string           `json:"frequency"`
	Cutoff        string           `json:"cutoff"`
	
	// Rates and Charges
	BaseRate      float64          `json:"base_rate"`
	LocalCharges  float64          `json:"local_charges"`
	CustomsCharges float64         `json:"customs_charges"`
	OtherCharges  []Charge         `json:"other_charges"`
	Currency      string           `json:"currency"`
	
	// Validity
	ValidFrom     time.Time        `json:"valid_from"`
	ValidUntil    time.Time        `json:"valid_until"`
	
	// Status
	Status        string           `json:"status"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// SeaImportService represents sea-specific import service
type SeaImportService struct {
	ImportService
	VesselType     string    `json:"vessel_type"`
	ShippingLine   string    `json:"shipping_line"`
	
	// Container Details (for FCL/LCL)
	ContainerTypes []string  `json:"container_types,omitempty"`
	ContainerRates map[string]float64 `json:"container_rates,omitempty"`
	
	// Bulk/RoRo Details
	CargoTypes     []string  `json:"cargo_types,omitempty"`
	VolumeRates    bool      `json:"volume_rates"`
	WeightRates    bool      `json:"weight_rates"`
	
	// Port Services
	PortServices   []string  `json:"port_services"`
	TerminalHandling bool    `json:"terminal_handling"`
}

// AirImportService represents air-specific import service
type AirImportService struct {
	ImportService
	AircraftType   string    `json:"aircraft_type"`
	Carrier        string    `json:"carrier"`
	
	// Cargo Specifications
	CargoTypes     []string  `json:"cargo_types"`
	ULDTypes       []string  `json:"uld_types,omitempty"`
	
	// Rate Structure
	MinimumCharge  float64   `json:"minimum_charge"`
	WeightBreaks   []float64 `json:"weight_breaks"`
	RatePerKg      []float64 `json:"rate_per_kg"`
	
	// Airport Services
	AirportServices []string `json:"airport_services"`
	SecurityScreening bool   `json:"security_screening"`
}

// RailImportService represents rail-specific import service
type RailImportService struct {
	ImportService
	RailOperator   string    `json:"rail_operator"`
	TrainService   string    `json:"train_service"`
	
	// Equipment
	WagonTypes     []string  `json:"wagon_types"`
	WagonCapacity  map[string]float64 `json:"wagon_capacity"`
	
	// Terminal Operations
	TerminalServices []string `json:"terminal_services"`
	Intermodal      bool     `json:"intermodal"`
	LastMile        bool     `json:"last_mile"`
}

// LandImportService represents land-specific import service
type LandImportService struct {
	ImportService
	TransportOperator string   `json:"transport_operator"`
	VehicleTypes     []string `json:"vehicle_types"`
	
	// Service Options
	DoorDelivery    bool     `json:"door_delivery"`
	CrossBorder     bool     `json:"cross_border"`
	CustomsClearance bool    `json:"customs_clearance"`
	
	// Rate Structure
	DistanceRate    float64  `json:"distance_rate"`
	WeightRate      float64  `json:"weight_rate"`
	MinimumCharge   float64  `json:"minimum_charge"`
}

// ImportServiceRequirement represents requirements for import service
type ImportServiceRequirement struct {
	ServiceType    ImportServiceType `json:"service_type"`
	Documents      []string         `json:"documents"`
	Licenses       []string         `json:"licenses"`
	Certifications []string         `json:"certifications"`
	CustomsRegimes []string         `json:"customs_regimes"`
}

// ImportServiceSchedule represents service schedule details
type ImportServiceSchedule struct {
	ServiceID      string           `json:"service_id"`
	ServiceType    ImportServiceType `json:"service_type"`
	Departure      time.Time        `json:"departure"`
	Arrival        time.Time        `json:"arrival"`
	TransitPoints  []TransitPoint   `json:"transit_points"`
	Frequency      string           `json:"frequency"`
	ValidFrom      time.Time        `json:"valid_from"`
	ValidUntil     time.Time        `json:"valid_until"`
}

// TransitPoint represents a point in the service route
type TransitPoint struct {
	Location       Location         `json:"location"`
	ArrivalTime    time.Time        `json:"arrival_time"`
	DepartureTime  time.Time        `json:"departure_time"`
	Services       []string         `json:"services"`
	Charges        []Charge         `json:"charges"`
}

// ImportServiceAvailability represents service availability
type ImportServiceAvailability struct {
	ServiceID      string           `json:"service_id"`
	ServiceType    ImportServiceType `json:"service_type"`
	AvailableSpace float64          `json:"available_space"`
	AvailableWeight float64         `json:"available_weight"`
	Equipment      []Equipment      `json:"equipment"`
	NextAvailable  time.Time        `json:"next_available"`
	Restrictions   []string         `json:"restrictions"`
}

// ImportServiceRate represents detailed rate structure
type ImportServiceRate struct {
	ServiceID      string           `json:"service_id"`
	ServiceType    ImportServiceType `json:"service_type"`
	BaseRate       float64          `json:"base_rate"`
	LocalCharges   []Charge         `json:"local_charges"`
	Surcharges     []Charge         `json:"surcharges"`
	Discounts      []Charge         `json:"discounts"`
	Currency       string           `json:"currency"`
	ValidFrom      time.Time        `json:"valid_from"`
	ValidUntil     time.Time        `json:"valid_until"`
}
