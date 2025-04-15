package models

import (
	"time"
)

// InfrastructureType represents different types of logistics infrastructure
type InfrastructureType string

const (
	InternationalAirport InfrastructureType = "AIRPORT"
	Seaport             InfrastructureType = "SEAPORT"
	InlandDepot         InfrastructureType = "INLAND_DEPOT"
)

// Infrastructure represents a logistics infrastructure point
type Infrastructure struct {
	ID          string             `json:"id"`
	Type        InfrastructureType `json:"type"`
	Name        string             `json:"name"`
	Code        string             `json:"code"` // IATA/UN LOCODE
	Country     string             `json:"country"`
	Region      string             `json:"region"`
	City        string             `json:"city"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	Status      string    `json:"status"` // ACTIVE, INACTIVE, MAINTENANCE
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Airport represents an international airport
type Airport struct {
	Infrastructure
	IATACode     string   `json:"iata_code"`
	ICAOCode     string   `json:"icao_code"`
	Terminals    []string `json:"terminals"`
	CargoZones   []string `json:"cargo_zones"`
	
	// Capabilities
	MaxAircraftSize    string   `json:"max_aircraft_size"`
	CargoCapacityTons  float64  `json:"cargo_capacity_tons"`
	HasCustoms         bool     `json:"has_customs"`
	HasColdStorage     bool     `json:"has_cold_storage"`
	HasDGRFacility     bool     `json:"has_dgr_facility"`
	
	// Operating Hours
	OperatingHours     string   `json:"operating_hours"`
	CargoOperations    string   `json:"cargo_operations"`
	CustomsOperations  string   `json:"customs_operations"`
	
	// Handling Equipment
	HandlingEquipment  []string `json:"handling_equipment"`
	
	// Connected Routes
	DirectDestinations []string `json:"direct_destinations"`
}

// Seaport represents a seaport
type Seaport struct {
	Infrastructure
	UNLOCode     string   `json:"un_locode"`
	Terminals    []string `json:"terminals"`
	
	// Capabilities
	MaxVesselSize      string   `json:"max_vessel_size"`
	MaxDraft           float64  `json:"max_draft"`
	ContainerCapacity  int      `json:"container_capacity"`
	HasCustoms         bool     `json:"has_customs"`
	HasColdStorage     bool     `json:"has_cold_storage"`
	HasDGRFacility     bool     `json:"has_dgr_facility"`
	
	// Equipment
	Cranes            []string `json:"cranes"`
	YardEquipment     []string `json:"yard_equipment"`
	
	// Storage
	ContainerYard     bool     `json:"container_yard"`
	BulkStorage       bool     `json:"bulk_storage"`
	LiquidStorage     bool     `json:"liquid_storage"`
	
	// Operating Hours
	OperatingHours    string   `json:"operating_hours"`
	CustomsHours      string   `json:"customs_hours"`
	
	// Connected Routes
	ShippingRoutes    []string `json:"shipping_routes"`
	FeederServices    []string `json:"feeder_services"`
}

// InlandDepot represents an inland container depot
type InlandDepot struct {
	Infrastructure
	DepotCode    string   `json:"depot_code"`
	
	// Capabilities
	StorageCapacity    int      `json:"storage_capacity"`
	HasCustoms         bool     `json:"has_customs"`
	HasColdStorage     bool     `json:"has_cold_storage"`
	HasDGRFacility     bool     `json:"has_dgr_facility"`
	
	// Connections
	NearestSeaport     string   `json:"nearest_seaport"`
	NearestAirport     string   `json:"nearest_airport"`
	RailConnection     bool     `json:"rail_connection"`
	
	// Equipment
	HandlingEquipment  []string `json:"handling_equipment"`
	
	// Storage Types
	ContainerStorage   bool     `json:"container_storage"`
	BulkStorage       bool     `json:"bulk_storage"`
	CustomsBonded     bool     `json:"customs_bonded"`
	
	// Operating Hours
	OperatingHours    string   `json:"operating_hours"`
	CustomsHours      string   `json:"customs_hours"`
	
	// Transport Connectivity
	RoadAccess        bool     `json:"road_access"`
	RailAccess        bool     `json:"rail_access"`
	LastMileDelivery  bool     `json:"last_mile_delivery"`
}

// ServiceLocation represents a location where a transport service is offered
type ServiceLocation struct {
	InfrastructureID   string             `json:"infrastructure_id"`
	Type               InfrastructureType  `json:"type"`
	Name               string             `json:"name"`
	Code               string             `json:"code"`
	Country            string             `json:"country"`
	Services           []string           `json:"services"` // List of service IDs
	OperatingHours     string             `json:"operating_hours"`
	CustomsAvailable   bool               `json:"customs_available"`
	StorageAvailable   bool               `json:"storage_available"`
	HandlingEquipment  []string           `json:"handling_equipment"`
}

// CountryInfrastructure represents all logistics infrastructure in a country
type CountryInfrastructure struct {
	Country           string        `json:"country"`
	CountryCode       string        `json:"country_code"`
	Airports          []Airport     `json:"airports"`
	Seaports          []Seaport     `json:"seaports"`
	InlandDepots      []InlandDepot `json:"inland_depots"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

// InfrastructureCapacity represents the current capacity status of infrastructure
type InfrastructureCapacity struct {
	InfrastructureID   string    `json:"infrastructure_id"`
	StorageUtilization float64   `json:"storage_utilization"` // Percentage
	EquipmentStatus    map[string]string `json:"equipment_status"`
	LastUpdated        time.Time `json:"last_updated"`
}

// InfrastructureSchedule represents operating schedules
type InfrastructureSchedule struct {
	InfrastructureID   string    `json:"infrastructure_id"`
	OperatingDays      []string  `json:"operating_days"`
	OperatingHours     string    `json:"operating_hours"`
	CustomsHours       string    `json:"customs_hours"`
	SpecialClosures    []string  `json:"special_closures"`
	ValidFrom          time.Time `json:"valid_from"`
	ValidUntil         time.Time `json:"valid_until"`
}

// GetInfrastructureType returns the type based on code format
func GetInfrastructureType(code string) InfrastructureType {
	// IATA airport codes are 3 letters
	if len(code) == 3 && isAllLetters(code) {
		return InternationalAirport
	}
	// UN/LOCODE are 5 characters (2 country + 3 location)
	if len(code) == 5 && isUNLOCODE(code) {
		return Seaport
	}
	// Inland depot codes typically start with ICD
	if len(code) >= 3 && code[:3] == "ICD" {
		return InlandDepot
	}
	return ""
}

// Helper functions

func isAllLetters(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

func isUNLOCODE(s string) bool {
	if len(s) != 5 {
		return false
	}
	// First 2 characters should be letters (country code)
	if !isAllLetters(s[:2]) {
		return false
	}
	// Last 3 characters can be letters or numbers
	for _, r := range s[2:] {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}
