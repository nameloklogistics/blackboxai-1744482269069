package models

import (
    "time"
)

// Transport Modes
const (
    ModeSea  = "SEA"
    ModeAir  = "AIR"
    ModeRail = "RAIL"
    ModeRoad = "ROAD"
)

// TransportService represents a transport service offering
type TransportService struct {
    BaseModel
    ProviderID      string     `json:"provider_id"`
    Mode            string     `json:"mode"`            // Sea, Air, Rail, Road
    Origin          Location   `json:"origin"`
    Destination     Location   `json:"destination"`
    TransitPoints   []Location `json:"transit_points,omitempty"`
    Schedule        Schedule   `json:"schedule"`
    Rates           []Rate     `json:"rates"`
    Requirements    []string   `json:"requirements"`    // Required documents, certifications
    Restrictions    []string   `json:"restrictions"`    // Weight limits, cargo restrictions
    Status          string     `json:"status"`
    AvailableSpace  float64    `json:"available_space"` // Available capacity
    Rating          float32    `json:"rating"`
    ReviewCount     int        `json:"review_count"`
}

// Schedule represents service timing information
type Schedule struct {
    Frequency       string      `json:"frequency"`        // Daily, Weekly, Monthly
    DepartureDays   []string    `json:"departure_days"`   // Days of week
    DepartureTime   string      `json:"departure_time"`   // HH:MM format
    TransitTime     string      `json:"transit_time"`     // Duration
    ValidFrom       time.Time   `json:"valid_from"`
    ValidUntil      time.Time   `json:"valid_until"`
    Availability    []TimeSlot  `json:"availability"`
}

// TimeSlot represents a specific time window
type TimeSlot struct {
    TimeWindow
    Capacity        float64    `json:"capacity"`
    Available       bool       `json:"available"`
}

// Rate represents pricing information
type Rate struct {
    BaseModel
    ServiceID       string     `json:"service_id"`
    ContainerType   string     `json:"container_type,omitempty"` // For sea freight
    WeightRange     Range      `json:"weight_range"`
    VolumeRange     Range      `json:"volume_range"`
    BasePrice       Currency   `json:"base_price"`
    AdditionalFees  []Fee      `json:"additional_fees"`
    ValidFrom       time.Time  `json:"valid_from"`
    ValidUntil      time.Time  `json:"valid_until"`
    Conditions      []string   `json:"conditions"`
}

// Range represents a numeric range
type Range struct {
    Min     float64 `json:"min"`
    Max     float64 `json:"max"`
    Unit    string  `json:"unit"`
}

// Fee represents additional charges
type Fee struct {
    Name        string   `json:"name"`
    Amount      Currency `json:"amount"`
    Type        string   `json:"type"`  // Fixed, Percentage
    Mandatory   bool     `json:"mandatory"`
    Description string   `json:"description"`
}

// Cargo represents shipment details
type Cargo struct {
    Type            string    `json:"type"`
    Description     string    `json:"description"`
    Weight          float64   `json:"weight"`
    WeightUnit      string    `json:"weight_unit"`
    Volume          float64   `json:"volume"`
    VolumeUnit      string    `json:"volume_unit"`
    Pieces          int       `json:"pieces"`
    ContainerType   string    `json:"container_type,omitempty"`
    DangerousGoods  bool      `json:"dangerous_goods"`
    SpecialHandling []string  `json:"special_handling,omitempty"`
}

// Payment represents payment information
type Payment struct {
    Status          string    `json:"status"`
    Amount          Currency  `json:"amount"`
    Method          string    `json:"method"`
    TransactionID   string    `json:"transaction_id"`
    PaidAt          time.Time `json:"paid_at"`
    EscrowID        string    `json:"escrow_id,omitempty"`
}
