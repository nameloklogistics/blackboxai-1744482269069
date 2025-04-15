package models

import "time"

// SearchFilter represents the search criteria for logistics services
type SearchFilter struct {
	// Basic filters
	Query       string          `json:"query"`        // General search query
	Category    ServiceCategory `json:"category"`     // Service category
	SubCategory string          `json:"subcategory"`  // Service subcategory
	
	// Location filters
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	
	// Price range
	MinPrice    *float64 `json:"min_price"`
	MaxPrice    *float64 `json:"max_price"`
	
	// Service attributes
	RouteType   string `json:"route_type"`   // Direct, Via Hub
	TransitTime string `json:"transit_time"` // Express, Standard, etc.
	
	// Provider filters
	ProviderID  string `json:"provider_id"`
	
	// Availability
	IsActive    *bool     `json:"is_active"`
	ValidFrom   time.Time `json:"valid_from"`
	ValidUntil  time.Time `json:"valid_until"`
	
	// Pagination
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	
	// Sorting
	SortBy    string `json:"sort_by"`    // price, transit_time, rating
	SortOrder string `json:"sort_order"` // asc, desc
}

// SearchResult represents the search results with pagination info
type SearchResult struct {
	Services    []LogisticsService `json:"services"`
	TotalCount  int               `json:"total_count"`
	CurrentPage int               `json:"current_page"`
	PageSize    int               `json:"page_size"`
	TotalPages  int               `json:"total_pages"`
}

// SearchSuggestion represents search suggestions for autocomplete
type SearchSuggestion struct {
	Type  string `json:"type"`  // location, service, provider
	Value string `json:"value"` // The suggested value
	Label string `json:"label"` // Display label
}

// SearchStats represents search analytics data
type SearchStats struct {
	PopularSearches []struct {
		Query     string    `json:"query"`
		Count     int       `json:"count"`
		LastUsed  time.Time `json:"last_used"`
	} `json:"popular_searches"`
	
	PopularFilters []struct {
		Filter string `json:"filter"`
		Value  string `json:"value"`
		Count  int    `json:"count"`
	} `json:"popular_filters"`
	
	PopularRoutes []struct {
		Origin      string `json:"origin"`
		Destination string `json:"destination"`
		Count       int    `json:"count"`
	} `json:"popular_routes"`
}
