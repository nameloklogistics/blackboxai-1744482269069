package models

import (
	"time"
)

// ServiceCategory represents the main service categories
type ServiceCategory string

const (
	ImportService     ServiceCategory = "IMPORT"
	ExportService     ServiceCategory = "EXPORT"
	TransitService    ServiceCategory = "TRANSIT"
	TransshipService  ServiceCategory = "TRANSSHIPMENT"
)

// ServiceSubCategory represents service subcategories
type ServiceSubCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    ServiceCategory `json:"category"`
}

// ServiceItem represents a specific service item with pricing
type ServiceItem struct {
	ID              string    `json:"id"`
	SubCategoryID   string    `json:"subcategory_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	BasePrice       float64   `json:"base_price"`
	Currency        string    `json:"currency"`
	PriceUnit       string    `json:"price_unit"` // per container, per kg, per CBM, etc.
	MinQuantity     float64   `json:"min_quantity"`
	MaxQuantity     float64   `json:"max_quantity"`
	ValidFrom       time.Time `json:"valid_from"`
	ValidUntil      time.Time `json:"valid_until"`
	ProviderID      string    `json:"provider_id"`
	Terms           string    `json:"terms"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// LogisticsService represents a complete service offering
type LogisticsService struct {
	ID              string          `json:"id"`
	Category        ServiceCategory `json:"category"`
	SubCategory     ServiceSubCategory `json:"subcategory"`
	Item            ServiceItem     `json:"item"`
	Provider        ServiceProvider `json:"provider"`
	
	// Service-specific details
	Origin          string    `json:"origin"`
	Destination     string    `json:"destination"`
	TransitTime     string    `json:"transit_time"`
	Frequency       string    `json:"frequency"` // Daily, Weekly, etc.
	RouteType       string    `json:"route_type"` // Direct, Via Hub
	
	// Additional charges
	Surcharges      []Surcharge `json:"surcharges"`
	
	// Service availability
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Surcharge represents additional charges
type Surcharge struct {
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
	IsRequired  bool    `json:"is_required"`
}

// Common service subcategories
var (
	// Import Service Subcategories
	ImportSubCategories = []ServiceSubCategory{
		{Name: "FCL Import", Description: "Full Container Load Import Services", Category: ImportService},
		{Name: "LCL Import", Description: "Less than Container Load Import Services", Category: ImportService},
		{Name: "Air Import", Description: "Air Freight Import Services", Category: ImportService},
		{Name: "Customs Clearance", Description: "Import Customs Clearance Services", Category: ImportService},
		{Name: "Door Delivery", Description: "Last Mile Delivery Services", Category: ImportService},
	}

	// Export Service Subcategories
	ExportSubCategories = []ServiceSubCategory{
		{Name: "FCL Export", Description: "Full Container Load Export Services", Category: ExportService},
		{Name: "LCL Export", Description: "Less than Container Load Export Services", Category: ExportService},
		{Name: "Air Export", Description: "Air Freight Export Services", Category: ExportService},
		{Name: "Export Documentation", Description: "Export Documentation Services", Category: ExportService},
		{Name: "Pickup Service", Description: "Cargo Pickup Services", Category: ExportService},
	}

	// Transit Service Subcategories
	TransitSubCategories = []ServiceSubCategory{
		{Name: "Land Transit", Description: "Cross-border Land Transportation", Category: TransitService},
		{Name: "Rail Transit", Description: "Rail Freight Services", Category: TransitService},
		{Name: "Multimodal Transit", Description: "Combined Transportation Services", Category: TransitService},
		{Name: "Transit Documentation", Description: "Transit Documentation Services", Category: TransitService},
	}

	// Transshipment Service Subcategories
	TransshipSubCategories = []ServiceSubCategory{
		{Name: "Sea-Sea Transshipment", Description: "Ocean Freight Transshipment", Category: TransshipService},
		{Name: "Air-Air Transshipment", Description: "Air Freight Transshipment", Category: TransshipService},
		{Name: "Sea-Air Transshipment", Description: "Sea to Air Transshipment", Category: TransshipService},
		{Name: "Air-Sea Transshipment", Description: "Air to Sea Transshipment", Category: TransshipService},
	}
)

// ServiceItemTemplate represents common service items within subcategories
type ServiceItemTemplate struct {
	Name        string
	Description string
	PriceUnit   string
	Category    ServiceCategory
	SubCategory string
}

// Common service items
var (
	// FCL Import Items
	FCLImportItems = []ServiceItemTemplate{
		{
			Name: "20' Container Import",
			Description: "Import service for 20 foot container",
			PriceUnit: "per container",
			Category: ImportService,
			SubCategory: "FCL Import",
		},
		{
			Name: "40' Container Import",
			Description: "Import service for 40 foot container",
			PriceUnit: "per container",
			Category: ImportService,
			SubCategory: "FCL Import",
		},
		{
			Name: "40' HC Container Import",
			Description: "Import service for 40 foot high cube container",
			PriceUnit: "per container",
			Category: ImportService,
			SubCategory: "FCL Import",
		},
	}

	// LCL Import Items
	LCLImportItems = []ServiceItemTemplate{
		{
			Name: "LCL Import Basic",
			Description: "Import service for loose cargo",
			PriceUnit: "per CBM",
			Category: ImportService,
			SubCategory: "LCL Import",
		},
		{
			Name: "LCL Import Premium",
			Description: "Premium import service for loose cargo with priority handling",
			PriceUnit: "per CBM",
			Category: ImportService,
			SubCategory: "LCL Import",
		},
	}

	// Air Import Items
	AirImportItems = []ServiceItemTemplate{
		{
			Name: "Air Import Standard",
			Description: "Standard air freight import service",
			PriceUnit: "per kg",
			Category: ImportService,
			SubCategory: "Air Import",
		},
		{
			Name: "Air Import Express",
			Description: "Express air freight import service",
			PriceUnit: "per kg",
			Category: ImportService,
			SubCategory: "Air Import",
		},
	}
)
