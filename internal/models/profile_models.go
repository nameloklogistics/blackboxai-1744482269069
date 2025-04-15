package models

import (
	"time"
)

// ServiceProvider represents a logistics service provider (forwarder) profile
type ServiceProvider struct {
	ID                    string    `json:"id"`
	CompanyName           string    `json:"company_name"`
	ContactPerson         string    `json:"contact_person"`
	Telephone            string    `json:"telephone"`
	Email                string    `json:"email"`
	Country              string    `json:"country"`
	Address              string    `json:"address"`
	BusinessLicense      string    `json:"business_license"`
	TaxID                string    `json:"tax_id"`
	
	// Logistics Infrastructure
	MajorSeaPorts        []Port    `json:"major_sea_ports"`
	InternationalAirports []Airport `json:"international_airports"`
	ContainerTerminals   []Terminal `json:"container_terminals"`
	
	// Certifications and Memberships
	IATA                 bool      `json:"iata_member"`
	FIATA                bool      `json:"fiata_member"`
	CustomsBrokerLicense string    `json:"customs_broker_license"`
	
	// Blockchain Details
	WalletAddress        string    `json:"wallet_address"`
	
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ServiceBuyer represents a shipper/customer profile
type ServiceBuyer struct {
	ID                string    `json:"id"`
	CompanyName       string    `json:"company_name"`
	ContactPerson     string    `json:"contact_person"`
	Telephone        string    `json:"telephone"`
	Email            string    `json:"email"`
	Country          string    `json:"country"`
	Address          string    `json:"address"`
	BusinessType     string    `json:"business_type"` // Importer/Exporter/Both
	TaxID            string    `json:"tax_id"`
	
	// Trade Information
	ImporterCode     string    `json:"importer_code,omitempty"`
	ExporterCode     string    `json:"exporter_code,omitempty"`
	
	// Blockchain Details
	WalletAddress    string    `json:"wallet_address"`
	
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Port represents a sea port
type Port struct {
	Name        string `json:"name"`
	Code        string `json:"code"` // UN/LOCODE
	Country     string `json:"country"`
	Coordinates string `json:"coordinates"`
}

// Airport represents an international airport
type Airport struct {
	Name        string `json:"name"`
	Code        string `json:"code"` // IATA code
	Country     string `json:"country"`
	Coordinates string `json:"coordinates"`
}

// Terminal represents an inland container terminal
type Terminal struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Location    string `json:"location"`
	Country     string `json:"country"`
	Coordinates string `json:"coordinates"`
	Facilities  []string `json:"facilities"` // e.g., "Container Storage", "Customs Clearance"
}
