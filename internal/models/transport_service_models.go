package models

import (
	"time"
)

// TransportMode represents different modes of transport
type TransportMode string

const (
	SeaTransport  TransportMode = "SEA"
	AirTransport  TransportMode = "AIR"
	RailTransport TransportMode = "RAIL"
	RoadTransport TransportMode = "ROAD"
)

// TransportService represents a service offering for a specific transport mode
type TransportService struct {
	ID            string         `json:"id"`
	Mode          TransportMode  `json:"mode"`
	Category      ServiceCategory `json:"category"` // Import, Export, Transit, Transshipment
	Provider      ServiceProvider `json:"provider"`
	
	// Route Information
	Origin        string         `json:"origin"`
	Destination   string         `json:"destination"`
	TransitPoints []string       `json:"transit_points,omitempty"`
	TransitTime   string         `json:"transit_time"`
	Frequency     string         `json:"frequency"` // Daily, Weekly, etc.
	
	// Equipment and Capacity
	Equipment     []Equipment    `json:"equipment,omitempty"`
	Capacity      Capacity       `json:"capacity"`
	
	// Pricing
	BaseRate      float64        `json:"base_rate"`
	Currency      string         `json:"currency"`
	Surcharges    []Surcharge    `json:"surcharges"`
	
	// Service Schedule
	ValidFrom     time.Time      `json:"valid_from"`
	ValidUntil    time.Time      `json:"valid_until"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// Equipment represents transport equipment details
type Equipment struct {
	Type        string `json:"type"` // Container, Trailer, Wagon, Aircraft
	Size        string `json:"size"` // 20ft, 40ft, etc.
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

// Capacity represents service capacity
type Capacity struct {
	WeightLimit float64 `json:"weight_limit"`
	VolumeLimit float64 `json:"volume_limit"`
	UnitType    string  `json:"unit_type"` // TEU, CBM, KG, etc.
}

// Sea Transport Specific Models
type SeaService struct {
	TransportService
	VesselType     string   `json:"vessel_type"`
	ShippingLine   string   `json:"shipping_line"`
	ContainerTypes []string `json:"container_types"` // Dry, Reefer, Open Top, etc.
	PortServices   []string `json:"port_services"`   // Loading, Unloading, Storage
}

// Air Transport Specific Models
type AirService struct {
	TransportService
	AircraftType   string   `json:"aircraft_type"`
	Carrier        string   `json:"carrier"`
	CargoTypes     []string `json:"cargo_types"` // General, Perishable, Dangerous
	AirportServices []string `json:"airport_services"`
}

// Rail Transport Specific Models
type RailService struct {
	TransportService
	TrainType      string   `json:"train_type"`
	RailOperator   string   `json:"rail_operator"`
	WagonTypes     []string `json:"wagon_types"`
	TerminalServices []string `json:"terminal_services"`
}

// Road Transport Specific Models
type RoadService struct {
	TransportService
	VehicleType    string   `json:"vehicle_type"`
	Carrier        string   `json:"carrier"`
	TruckTypes     []string `json:"truck_types"`
	DeliveryServices []string `json:"delivery_services"`
}

// Service Templates by Transport Mode
var (
	// Sea Transport Services
	SeaServiceTemplates = map[ServiceCategory][]string{
		ImportService: {
			"FCL Import",
			"LCL Import",
			"Break Bulk Import",
			"RoRo Import",
		},
		ExportService: {
			"FCL Export",
			"LCL Export",
			"Break Bulk Export",
			"RoRo Export",
		},
		TransitService: {
			"Sea-Sea Transit",
			"Sea-Rail Transit",
			"Sea-Road Transit",
		},
		TransshipService: {
			"Port Transshipment",
			"Feeder Service",
			"Hub Connection",
		},
	}

	// Air Transport Services
	AirServiceTemplates = map[ServiceCategory][]string{
		ImportService: {
			"General Cargo Import",
			"Express Freight Import",
			"Charter Service Import",
		},
		ExportService: {
			"General Cargo Export",
			"Express Freight Export",
			"Charter Service Export",
		},
		TransitService: {
			"Air-Air Transit",
			"Air-Road Transit",
			"Air-Rail Transit",
		},
		TransshipService: {
			"Airport Transshipment",
			"Express Hub Connection",
			"Interline Transfer",
		},
	}

	// Rail Transport Services
	RailServiceTemplates = map[ServiceCategory][]string{
		ImportService: {
			"Container Rail Import",
			"Bulk Cargo Rail Import",
			"Intermodal Import",
		},
		ExportService: {
			"Container Rail Export",
			"Bulk Cargo Rail Export",
			"Intermodal Export",
		},
		TransitService: {
			"Rail-Rail Transit",
			"Rail-Road Transit",
			"Cross-Border Rail",
		},
		TransshipService: {
			"Terminal Transshipment",
			"Rail Hub Connection",
			"Intermodal Transfer",
		},
	}

	// Road Transport Services
	RoadServiceTemplates = map[ServiceCategory][]string{
		ImportService: {
			"FTL Import Delivery",
			"LTL Import Delivery",
			"Special Equipment Import",
		},
		ExportService: {
			"FTL Export Pickup",
			"LTL Export Pickup",
			"Special Equipment Export",
		},
		TransitService: {
			"Cross-Border Road",
			"Interstate Transit",
			"City-to-City",
		},
		TransshipService: {
			"Truck Terminal Transfer",
			"Cross-Dock Service",
			"Distribution Hub",
		},
	}
)

// GetServiceTemplatesByMode returns service templates for a specific transport mode
func GetServiceTemplatesByMode(mode TransportMode) map[ServiceCategory][]string {
	switch mode {
	case SeaTransport:
		return SeaServiceTemplates
	case AirTransport:
		return AirServiceTemplates
	case RailTransport:
		return RailServiceTemplates
	case RoadTransport:
		return RoadServiceTemplates
	default:
		return nil
	}
}

// GetEquipmentByMode returns available equipment types for a specific transport mode
func GetEquipmentByMode(mode TransportMode) []Equipment {
	switch mode {
	case SeaTransport:
		return []Equipment{
			{Type: "Container", Size: "20ft", Description: "Standard 20ft Container"},
			{Type: "Container", Size: "40ft", Description: "Standard 40ft Container"},
			{Type: "Container", Size: "40ft HC", Description: "High Cube Container"},
			{Type: "Reefer", Size: "20ft", Description: "Refrigerated Container"},
			{Type: "Reefer", Size: "40ft", Description: "Refrigerated Container"},
		}
	case AirTransport:
		return []Equipment{
			{Type: "ULD", Size: "PMC", Description: "Pallet Wide Body"},
			{Type: "ULD", Size: "AKE", Description: "Container Narrow Body"},
			{Type: "Bulk", Size: "Various", Description: "Bulk Loading"},
		}
	case RailTransport:
		return []Equipment{
			{Type: "Wagon", Size: "20ft", Description: "Container Wagon"},
			{Type: "Wagon", Size: "40ft", Description: "Container Wagon"},
			{Type: "Wagon", Size: "Bulk", Description: "Bulk Cargo Wagon"},
		}
	case RoadTransport:
		return []Equipment{
			{Type: "Truck", Size: "20ft", Description: "Container Truck"},
			{Type: "Truck", Size: "40ft", Description: "Container Truck"},
			{Type: "Trailer", Size: "53ft", Description: "Box Trailer"},
		}
	default:
		return nil
	}
}
