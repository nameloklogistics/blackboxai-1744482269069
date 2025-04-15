package models

import (
    "time"
)

// Base model with common fields
type BaseModel struct {
    ID        string    `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// User roles
const (
    RoleAdmin           = "ADMIN"
    RoleFreightForwarder = "FREIGHT_FORWARDER"
    RoleCustomsBroker    = "CUSTOMS_BROKER"
    RoleShipper         = "SHIPPER"
    RoleConsignee       = "CONSIGNEE"
)

// User represents a system user
type User struct {
    BaseModel
    Email     string `json:"email"`
    Name      string `json:"name"`
    Role      string `json:"role"`
    CompanyID string `json:"company_id"`
    Active    bool   `json:"active"`
}

// Company represents a business entity
type Company struct {
    BaseModel
    Name            string `json:"name"`
    RegistrationNo  string `json:"registration_no"`
    Type            string `json:"type"` // e.g., Freight Forwarder, Shipper, etc.
    Address         string `json:"address"`
    Country        string `json:"country"`
    ContactPerson   string `json:"contact_person"`
    ContactEmail    string `json:"contact_email"`
    ContactPhone    string `json:"contact_phone"`
    VerificationStatus string `json:"verification_status"`
}

// Wallet represents a user's blockchain wallet
type Wallet struct {
    BaseModel
    UserID      string `json:"user_id"`
    Address     string `json:"address"`
    PublicKey   string `json:"public_key"`
    Balance     string `json:"balance"`
    TokenBalance string `json:"token_balance"`
}

// Transaction represents a blockchain transaction
type Transaction struct {
    BaseModel
    FromAddress string `json:"from_address"`
    ToAddress   string `json:"to_address"`
    Amount      string `json:"amount"`
    TokenCode   string `json:"token_code"`
    Status      string `json:"status"`
    TxHash      string `json:"tx_hash"`
    Type        string `json:"type"` // e.g., Transfer, Escrow, Release
}

// Location represents a geographical location
type Location struct {
    BaseModel
    Name        string  `json:"name"`
    Type        string  `json:"type"` // e.g., Port, Airport, Warehouse
    Country     string  `json:"country"`
    City        string  `json:"city"`
    Latitude    float64 `json:"latitude"`
    Longitude   float64 `json:"longitude"`
    Code        string  `json:"code"` // e.g., Port code, Airport code
    Status      string  `json:"status"`
}

// Address represents a detailed address
type Address struct {
    Street      string `json:"street"`
    City        string `json:"city"`
    State       string `json:"state"`
    Country     string `json:"country"`
    PostalCode  string `json:"postal_code"`
    Coordinates struct {
        Latitude  float64 `json:"latitude"`
        Longitude float64 `json:"longitude"`
    } `json:"coordinates"`
}

// Contact represents contact information
type Contact struct {
    Name        string `json:"name"`
    Email       string `json:"email"`
    Phone       string `json:"phone"`
    Role        string `json:"role"`
    Department  string `json:"department"`
}

// Document represents any type of document in the system
type Document struct {
    BaseModel
    Type        string `json:"type"` // e.g., Invoice, Bill of Lading, etc.
    Number      string `json:"number"`
    IssueDate   time.Time `json:"issue_date"`
    ExpiryDate  time.Time `json:"expiry_date"`
    Status      string `json:"status"`
    URL         string `json:"url"` // Document storage URL
    IssuedBy    string `json:"issued_by"`
    IssuedTo    string `json:"issued_to"`
    VerifiedAt  *time.Time `json:"verified_at,omitempty"`
    VerifiedBy  string `json:"verified_by,omitempty"`
}

// Currency represents monetary values
type Currency struct {
    Amount   float64 `json:"amount"`
    Code     string  `json:"code"` // e.g., USD, EUR
}

// TimeWindow represents a time period
type TimeWindow struct {
    Start    time.Time `json:"start"`
    End      time.Time `json:"end"`
    Duration string    `json:"duration"` // e.g., 2h30m
}

// Status represents the current state of an entity
type Status struct {
    Code        string    `json:"code"`
    Description string    `json:"description"`
    UpdatedAt   time.Time `json:"updated_at"`
    UpdatedBy   string    `json:"updated_by"`
}
