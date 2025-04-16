package models

import (
    "time"
)

// SearchType represents different types of searchable entities
const (
    SearchTypeService    = "SERVICE"
    SearchTypeBooking    = "BOOKING"
    SearchTypeProvider   = "PROVIDER"
    SearchTypeLocation   = "LOCATION"
    SearchTypeDocument   = "DOCUMENT"
)

// SearchFilter represents search filtering options
type SearchFilter struct {
    Field           string      `json:"field"`
    Operator        string      `json:"operator"` // eq, gt, lt, contains, etc.
    Value           interface{} `json:"value"`
}

// SearchSort represents search sorting options
type SearchSort struct {
    Field      string `json:"field"`
    Direction  string `json:"direction"` // asc, desc
}

// SearchQuery represents a search request
type SearchQuery struct {
    SearchType      string         `json:"search_type"`
    Keywords        string         `json:"keywords"`
    Filters         []SearchFilter `json:"filters"`
    Sort           []SearchSort   `json:"sort"`
    Page           int            `json:"page"`
    PageSize       int            `json:"page_size"`
    IncludeDeleted bool           `json:"include_deleted"`
}

// SearchResult represents a generic search result
type SearchResult struct {
    Type            string      `json:"type"`
    ID              string      `json:"id"`
    Title           string      `json:"title"`
    Description     string      `json:"description"`
    Highlights      []string    `json:"highlights"`
    Score           float64     `json:"score"`
    Data            interface{} `json:"data"`
    LastUpdated     time.Time   `json:"last_updated"`
}

// SearchResponse represents the response to a search query
type SearchResponse struct {
    Query           SearchQuery    `json:"query"`
    Results         []SearchResult `json:"results"`
    TotalResults    int           `json:"total_results"`
    Page            int           `json:"page"`
    PageSize        int           `json:"page_size"`
    TotalPages      int           `json:"total_pages"`
    ExecutionTime   float64       `json:"execution_time"`
}

// ServiceSearchParams represents service-specific search parameters
type ServiceSearchParams struct {
    Origin          string    `json:"origin"`
    Destination     string    `json:"destination"`
    ServiceType     string    `json:"service_type"`
    TransportMode   string    `json:"transport_mode"`
    DepartureDate   time.Time `json:"departure_date"`
    CargoType       string    `json:"cargo_type"`
    Weight          float64   `json:"weight"`
    Volume          float64   `json:"volume"`
    PriceRange      struct {
        Min float64 `json:"min"`
        Max float64 `json:"max"`
    } `json:"price_range"`
}

// LocationSearchParams represents location-specific search parameters
type LocationSearchParams struct {
    Type            string    `json:"type"` // Port, Airport, Warehouse
    Country         string    `json:"country"`
    Region          string    `json:"region"`
    Coordinates     struct {
        Latitude    float64   `json:"latitude"`
        Longitude   float64   `json:"longitude"`
        Radius      float64   `json:"radius"` // in kilometers
    } `json:"coordinates"`
    Services        []string  `json:"services"`
}

// BookingSearchParams represents booking-specific search parameters
type BookingSearchParams struct {
    BookingID       string    `json:"booking_id"`
    Status          string    `json:"status"`
    DateRange       struct {
        Start       time.Time `json:"start"`
        End         time.Time `json:"end"`
    } `json:"date_range"`
    CustomerID      string    `json:"customer_id"`
    ProviderID      string    `json:"provider_id"`
}

// ProviderSearchParams represents provider-specific search parameters
type ProviderSearchParams struct {
    ServiceTypes    []string  `json:"service_types"`
    Locations       []string  `json:"locations"`
    Rating          float64   `json:"rating"`
    Certifications  []string  `json:"certifications"`
    AvailableNow    bool      `json:"available_now"`
}

// SearchSuggestion represents an auto-complete suggestion
type SearchSuggestion struct {
    Type            string    `json:"type"`
    Value           string    `json:"value"`
    Label           string    `json:"label"`
    Description     string    `json:"description,omitempty"`
    Icon            string    `json:"icon,omitempty"`
    Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// SearchHistory represents a user's search history
type SearchHistory struct {
    BaseModel
    UserID          string      `json:"user_id"`
    Query           SearchQuery `json:"query"`
    ResultCount     int         `json:"result_count"`
    ExecutionTime   float64     `json:"execution_time"`
    Successful      bool        `json:"successful"`
}

// SavedSearch represents a user's saved search
type SavedSearch struct {
    BaseModel
    UserID          string      `json:"user_id"`
    Name            string      `json:"name"`
    Query           SearchQuery `json:"query"`
    NotifyResults   bool        `json:"notify_results"`
    NotifyFrequency string      `json:"notify_frequency"` // DAILY, WEEKLY
    LastNotified    time.Time   `json:"last_notified"`
    IsActive        bool        `json:"is_active"`
}

// SearchAnalytics represents analytics data for searches
type SearchAnalytics struct {
    BaseModel
    SearchType      string    `json:"search_type"`
    QueryPattern    string    `json:"query_pattern"`
    ResultCount     int       `json:"result_count"`
    AverageScore    float64   `json:"average_score"`
    ExecutionTime   float64   `json:"execution_time"`
    UserCount       int       `json:"user_count"`
    FailureRate     float64   `json:"failure_rate"`
    PopularFilters  []string  `json:"popular_filters"`
}
