package models

import (
    "time"
)

// BookingStatus represents the status of a booking
const (
    BookingStatusPending    = "PENDING"
    BookingStatusConfirmed  = "CONFIRMED"
    BookingStatusInProgress = "IN_PROGRESS"
    BookingStatusCompleted  = "COMPLETED"
    BookingStatusCancelled  = "CANCELLED"
)

// QuoteStatus represents the status of a quote
const (
    QuoteStatusPending   = "PENDING"
    QuoteStatusAccepted  = "ACCEPTED"
    QuoteStatusRejected  = "REJECTED"
    QuoteStatusExpired   = "EXPIRED"
)

// Quote represents a service quote request
type Quote struct {
    BaseModel
    RequestID       string    `json:"request_id"`
    ServiceID       string    `json:"service_id"`
    ProviderID      string    `json:"provider_id"`
    CustomerID      string    `json:"customer_id"`
    Origin          Location  `json:"origin"`
    Destination     Location  `json:"destination"`
    CargoDetails    Cargo     `json:"cargo_details"`
    BasePrice       Currency  `json:"base_price"`
    AdditionalFees  []Fee     `json:"additional_fees"`
    TotalPrice      Currency  `json:"total_price"`
    ValidUntil      time.Time `json:"valid_until"`
    Status          string    `json:"status"`
    Terms           []string  `json:"terms"`
    Notes           string    `json:"notes,omitempty"`
}

// Booking represents a confirmed service booking
type Booking struct {
    BaseModel
    QuoteID         string     `json:"quote_id"`
    ServiceID       string     `json:"service_id"`
    CustomerID      string     `json:"customer_id"`
    ProviderID      string     `json:"provider_id"`
    Status          string     `json:"status"`
    CargoDetails    Cargo      `json:"cargo_details"`
    Schedule        TimeWindow `json:"schedule"`
    PickupAddress   Address    `json:"pickup_address"`
    DeliveryAddress Address    `json:"delivery_address"`
    Documents       []Document `json:"documents"`
    Payment         Payment    `json:"payment"`
    TrackingNumber  string     `json:"tracking_number"`
    Notes           string     `json:"notes,omitempty"`
}

// BookingEvent represents a tracking event in the booking lifecycle
type BookingEvent struct {
    BaseModel
    BookingID     string    `json:"booking_id"`
    Type          string    `json:"type"`
    Location      Location  `json:"location"`
    Timestamp     time.Time `json:"timestamp"`
    Description   string    `json:"description"`
    Status        string    `json:"status"`
    UpdatedBy     string    `json:"updated_by"`
}

// BookingDocument represents documents associated with a booking
type BookingDocument struct {
    BaseModel
    BookingID    string    `json:"booking_id"`
    Type         string    `json:"type"` // Invoice, Bill of Lading, etc.
    URL          string    `json:"url"`
    IssuedAt     time.Time `json:"issued_at"`
    IssuedBy     string    `json:"issued_by"`
    ValidUntil   time.Time `json:"valid_until,omitempty"`
    Status       string    `json:"status"`
    VerifiedAt   time.Time `json:"verified_at,omitempty"`
    VerifiedBy   string    `json:"verified_by,omitempty"`
}

// BookingPayment represents payment information for a booking
type BookingPayment struct {
    BaseModel
    BookingID      string    `json:"booking_id"`
    Amount         Currency  `json:"amount"`
    Status         string    `json:"status"`
    Method         string    `json:"method"`
    TransactionID  string    `json:"transaction_id"`
    PaidAt         time.Time `json:"paid_at,omitempty"`
    RefundedAt     time.Time `json:"refunded_at,omitempty"`
    EscrowID       string    `json:"escrow_id,omitempty"`
}

// BookingDispute represents a dispute raised for a booking
type BookingDispute struct {
    BaseModel
    BookingID     string    `json:"booking_id"`
    RaisedBy      string    `json:"raised_by"`
    Type          string    `json:"type"`
    Description   string    `json:"description"`
    Status        string    `json:"status"`
    Resolution    string    `json:"resolution,omitempty"`
    ResolvedAt    time.Time `json:"resolved_at,omitempty"`
    ResolvedBy    string    `json:"resolved_by,omitempty"`
    Documents     []Document `json:"documents,omitempty"`
}

// BookingReview represents a review for a completed booking
type BookingReview struct {
    BaseModel
    BookingID     string    `json:"booking_id"`
    ReviewerID    string    `json:"reviewer_id"`
    Rating        float32   `json:"rating"`
    Comment       string    `json:"comment"`
    Response      string    `json:"response,omitempty"`
    ResponseAt    time.Time `json:"response_at,omitempty"`
    Categories    map[string]float32 `json:"categories"` // Different aspects of the service
}
