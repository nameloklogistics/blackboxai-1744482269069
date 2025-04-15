package models

import (
	"time"
)

// ServiceCategory represents different categories of services
type ServiceCategory string

const (
	// Import Categories
	ImportDirect      ServiceCategory = "IMPORT_DIRECT"
	ImportTransit     ServiceCategory = "IMPORT_TRANSIT"
	ImportTransshipment ServiceCategory = "IMPORT_TRANSSHIPMENT"

	// Export Categories
	ExportDirect      ServiceCategory = "EXPORT_DIRECT"
	ExportTransit     ServiceCategory = "EXPORT_TRANSIT"

	// Transit Category
	Transit           ServiceCategory = "TRANSIT"

	// Transshipment Category
	Transshipment     ServiceCategory = "TRANSSHIPMENT"
)

// TransportMode represents different modes of transport
type TransportMode string

const (
	SeaTransport  TransportMode = "SEA"
	AirTransport  TransportMode = "AIR"
	RailTransport TransportMode = "RAIL"
	LandTransport TransportMode = "LAND"
)

// PackagingMode represents different packaging modes
type PackagingMode string

const (
	// Container Packaging
	FCL PackagingMode = "FCL" // Full Container Load
	LCL PackagingMode = "LCL" // Less than Container Load
	
	// Air Cargo Packaging
	BulkAir    PackagingMode = "BULK_AIR"    // Bulk Air Cargo
	ULD        PackagingMode = "ULD"         // Unit Load Device
	Palletized PackagingMode = "PALLETIZED"  // Palletized Air Cargo
	
	// Rail Cargo Packaging
	ContainerRail PackagingMode = "CONTAINER_RAIL" // Container on Rail
	BulkRail     PackagingMode = "BULK_RAIL"      // Bulk Rail Cargo
	CarLoad      PackagingMode = "CAR_LOAD"       // Full Car Load
	WagonLoad    PackagingMode = "WAGON_LOAD"     // Full Wagon Load
	
	// Land/Road Packaging
	FTL         PackagingMode = "FTL"          // Full Truck Load
	LTL         PackagingMode = "LTL"          // Less than Truck Load
	Parcel      PackagingMode = "PARCEL"       // Parcel Delivery
	Consolidated PackagingMode = "CONSOLIDATED" // Consolidated Load
)

// PackagingConfiguration represents packaging-specific configuration
type PackagingConfiguration struct {
	Mode           PackagingMode `json:"mode"`
	Capacity       float64       `json:"capacity"`       // Capacity in relevant unit (TEU, KG, CBM)
	MaxWeight      float64       `json:"max_weight"`     // Maximum weight in KG
	MaxVolume      float64       `json:"max_volume"`     // Maximum volume in CBM
	MaxDimensions  Dimensions    `json:"max_dimensions"` // Maximum dimensions
	SpecialHandling []string     `json:"special_handling,omitempty"`
}

// Dimensions represents physical dimensions
type Dimensions struct {
	Length float64 `json:"length"` // Length in meters
	Width  float64 `json:"width"`  // Width in meters
	Height float64 `json:"height"` // Height in meters
}

// BaseService represents common fields for all service types
type BaseService struct {
	ID            string         `json:"id"`
	Category      ServiceCategory `json:"category"`
	Mode          TransportMode  `json:"mode"`
	
	// Route Information
	Origin        Location       `json:"origin"`
	Destination   Location       `json:"destination"`
	ViaPoints     []Location     `json:"via_points,omitempty"`
	
	// Service Details
	TransitTime   string         `json:"transit_time"`
	Frequency     string         `json:"frequency"`
	Cutoff        string         `json:"cutoff"`
	
	// Packaging Information
	PackagingModes []PackagingConfiguration `json:"packaging_modes"`
	
	// Rates and Charges
	BaseRate      float64        `json:"base_rate"`
	LocalCharges  float64        `json:"local_charges"`
	OtherCharges  []Charge       `json:"other_charges"`
	Currency      string         `json:"currency"`
	
	// Validity
	ValidFrom     time.Time      `json:"valid_from"`
	ValidUntil    time.Time      `json:"valid_until"`
	
	// Status
	Status        string         `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// SeaService represents sea-specific service attributes
type SeaService struct {
	BaseService
	VesselType     string    `json:"vessel_type"`
	ShippingLine   string    `json:"shipping_line"`
	
	// Packaging-specific Details
	FCLCapabilities struct {
		ContainerSizes []string `json:"container_sizes"` // 20GP, 40GP, 40HC, etc.
		SpecialTypes   []string `json:"special_types"`   // Reefer, OpenTop, FlatRack, etc.
		SOC           bool     `json:"soc"`             // Shipper Owned Container
		DirectRouting bool     `json:"direct_routing"`  // Direct routing available
	} `json:"fcl_capabilities,omitempty"`
	
	LCLCapabilities struct {
		MinVolume     float64 `json:"min_volume"`      // Minimum CBM
		MinWeight     float64 `json:"min_weight"`      // Minimum KG
		Consolidation bool    `json:"consolidation"`   // Consolidation available
		Deconsolidation bool  `json:"deconsolidation"` // Deconsolidation available
	} `json:"lcl_capabilities,omitempty"`
	
	// Port Services
	PortServices   []string  `json:"port_services"`
	TerminalHandling bool    `json:"terminal_handling"`
}

// AirService represents air-specific service attributes
type AirService struct {
	BaseService
	AircraftType   string    `json:"aircraft_type"`
	Carrier        string    `json:"carrier"`
	
	// Packaging-specific Details
	ULDCapabilities struct {
		Types        []string `json:"types"`          // AKE, PMC, etc.
		Availability bool     `json:"availability"`   // ULD availability
		BuildBreak   bool     `json:"build_break"`    // Build/break available
	} `json:"uld_capabilities,omitempty"`
	
	BulkCapabilities struct {
		MaxPieces    int      `json:"max_pieces"`     // Maximum pieces per shipment
		LooseLoad    bool     `json:"loose_load"`     // Loose load accepted
	} `json:"bulk_capabilities,omitempty"`
	
	PalletCapabilities struct {
		Types        []string `json:"types"`          // Standard, Euro, etc.
		MaxHeight    float64  `json:"max_height"`     // Maximum height in meters
		Stackable    bool     `json:"stackable"`      // Stackable pallets accepted
	} `json:"pallet_capabilities,omitempty"`
	
	// Rate Structure
	MinimumCharge  float64   `json:"minimum_charge"`
	WeightBreaks   []float64 `json:"weight_breaks"`
	RatePerKg      []float64 `json:"rate_per_kg"`
	
	// Airport Services
	AirportServices []string `json:"airport_services"`
	SecurityScreening bool   `json:"security_screening"`
}

// RailService represents rail-specific service attributes
type RailService struct {
	BaseService
	RailOperator   string    `json:"rail_operator"`
	TrainService   string    `json:"train_service"`
	
	// Packaging-specific Details
	ContainerRailCapabilities struct {
		Types        []string `json:"types"`          // 20ft, 40ft, etc.
		StackingLevel int     `json:"stacking_level"` // Maximum stacking level
		Restrictions []string `json:"restrictions"`   // Container restrictions
	} `json:"container_rail_capabilities,omitempty"`
	
	BulkRailCapabilities struct {
		CargoTypes   []string `json:"cargo_types"`    // Bulk cargo types
		LoadingGear  bool     `json:"loading_gear"`   // Loading gear available
	} `json:"bulk_rail_capabilities,omitempty"`
	
	WagonCapabilities struct {
		Types        []string `json:"types"`          // Wagon types
		Capacity     float64  `json:"capacity"`       // Capacity in tons
		Specifications map[string]string `json:"specifications"` // Wagon specifications
	} `json:"wagon_capabilities,omitempty"`
	
	// Terminal Operations
	TerminalServices []string `json:"terminal_services"`
	Intermodal      bool     `json:"intermodal"`
	LastMile        bool     `json:"last_mile"`
}

// LandService represents land-specific service attributes
type LandService struct {
	BaseService
	TransportOperator string   `json:"transport_operator"`
	// Packaging-specific Details
	FTLCapabilities struct {
		VehicleTypes []string `json:"vehicle_types"`  // Truck types
		MaxPayload   float64  `json:"max_payload"`    // Maximum payload in tons
		Specifications map[string]string `json:"specifications"` // Vehicle specifications
	} `json:"ftl_capabilities,omitempty"`
	
	LTLCapabilities struct {
		MinVolume    float64  `json:"min_volume"`     // Minimum volume in CBM
		MinWeight    float64  `json:"min_weight"`     // Minimum weight in KG
		MaxPieces    int      `json:"max_pieces"`     // Maximum pieces per shipment
		Consolidation bool    `json:"consolidation"`  // Consolidation available
	} `json:"ltl_capabilities,omitempty"`
	
	ParcelCapabilities struct {
		MaxDimensions Dimensions `json:"max_dimensions"` // Maximum dimensions
		MaxWeight     float64    `json:"max_weight"`    // Maximum weight in KG
		Express       bool       `json:"express"`       // Express delivery available
	} `json:"parcel_capabilities,omitempty"`
	
	// Service Options
	DoorDelivery    bool     `json:"door_delivery"`
	CrossBorder     bool     `json:"cross_border"`
	CustomsClearance bool    `json:"customs_clearance"`
	
	// Rate Structure
	DistanceRate    float64  `json:"distance_rate"`
	WeightRate      float64  `json:"weight_rate"`
	MinimumCharge   float64  `json:"minimum_charge"`
}

// TransitDetails represents transit-specific details
type TransitDetails struct {
	TransitCountries []string  `json:"transit_countries"`
	TransitPermits   []string  `json:"transit_permits"`
	CustomsRegimes   []string  `json:"customs_regimes"`
	BondDetails     string     `json:"bond_details"`
	TransitTime     string     `json:"transit_time"`
	Restrictions    []string   `json:"restrictions"`
}

// TransshipmentDetails represents transshipment-specific details
type TransshipmentDetails struct {
	TransshipmentPort Location   `json:"transshipment_port"`
	ConnectionTime    string     `json:"connection_time"`
	StorageAllowance  string     `json:"storage_allowance"`
	StorageCharges    []Charge   `json:"storage_charges"`
	HandlingCharges   []Charge   `json:"handling_charges"`
	Restrictions      []string   `json:"restrictions"`
}

// ServiceSchedule represents service schedule details
type ServiceSchedule struct {
	ServiceID      string         `json:"service_id"`
	Category       ServiceCategory `json:"category"`
	Mode          TransportMode   `json:"mode"`
	Departure      time.Time      `json:"departure"`
	Arrival        time.Time      `json:"arrival"`
	TransitPoints  []TransitPoint `json:"transit_points"`
	Frequency      string         `json:"frequency"`
	ValidFrom      time.Time      `json:"valid_from"`
	ValidUntil     time.Time      `json:"valid_until"`
}

// ServiceAvailability represents service availability
type ServiceAvailability struct {
	ServiceID      string         `json:"service_id"`
	Category       ServiceCategory `json:"category"`
	Mode          TransportMode   `json:"mode"`
	AvailableSpace float64        `json:"available_space"`
	AvailableWeight float64       `json:"available_weight"`
	Equipment      []Equipment    `json:"equipment"`
	NextAvailable  time.Time      `json:"next_available"`
	Restrictions   []string       `json:"restrictions"`
}

// ServiceRate represents detailed rate structure
type ServiceRate struct {
	ServiceID      string         `json:"service_id"`
	Category       ServiceCategory `json:"category"`
	Mode          TransportMode   `json:"mode"`
	BaseRate       float64        `json:"base_rate"`
	LocalCharges   []Charge       `json:"local_charges"`
	TransitCharges []Charge       `json:"transit_charges,omitempty"`
	TransshipmentCharges []Charge `json:"transshipment_charges,omitempty"`
	Surcharges     []Charge       `json:"surcharges"`
	Discounts      []Charge       `json:"discounts"`
	Currency       string         `json:"currency"`
	ValidFrom      time.Time      `json:"valid_from"`
	ValidUntil     time.Time      `json:"valid_until"`
}

// ServiceRequirement represents requirements for a service
type ServiceRequirement struct {
	Category       ServiceCategory `json:"category"`
	Mode          TransportMode   `json:"mode"`
	Documents      []string       `json:"documents"`
	Licenses       []string       `json:"licenses"`
	Certifications []string       `json:"certifications"`
	CustomsRegimes []string       `json:"customs_regimes"`
	Restrictions   []string       `json:"restrictions"`
}
