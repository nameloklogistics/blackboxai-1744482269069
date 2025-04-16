package models

import (
    "time"
)

// ProfileType represents different types of profiles
const (
    ProfileTypePersonal   = "PERSONAL"
    ProfileTypeCompany    = "COMPANY"
    ProfileTypeEmployee   = "EMPLOYEE"
)

// VerificationStatus represents profile verification status
const (
    VerificationStatusPending   = "PENDING"
    VerificationStatusVerified  = "VERIFIED"
    VerificationStatusRejected  = "REJECTED"
    VerificationStatusExpired   = "EXPIRED"
)

// UserProfile represents a user's profile information
type UserProfile struct {
    BaseModel
    UserID          string    `json:"user_id"`
    ProfileType     string    `json:"profile_type"`
    FirstName       string    `json:"first_name"`
    LastName        string    `json:"last_name"`
    Email           string    `json:"email"`
    Phone           string    `json:"phone"`
    Title           string    `json:"title"`
    Department      string    `json:"department"`
    CompanyID       string    `json:"company_id,omitempty"`
    Address         Address   `json:"address"`
    Timezone        string    `json:"timezone"`
    Language        string    `json:"language"`
    ProfilePicture  string    `json:"profile_picture"`
    IsVerified      bool      `json:"is_verified"`
    LastLogin       time.Time `json:"last_login"`
}

// CompanyProfile represents a company's profile information
type CompanyProfile struct {
    BaseModel
    Name                string    `json:"name"`
    LegalName           string    `json:"legal_name"`
    RegistrationNumber  string    `json:"registration_number"`
    TaxID               string    `json:"tax_id"`
    Industry           string    `json:"industry"`
    CompanySize        string    `json:"company_size"`
    YearEstablished    int       `json:"year_established"`
    Website            string    `json:"website"`
    Description        string    `json:"description"`
    Logo               string    `json:"logo"`
    HeadOffice         Address   `json:"head_office"`
    BillingAddress     Address   `json:"billing_address"`
    PrimaryContact     Contact   `json:"primary_contact"`
    VerificationStatus string    `json:"verification_status"`
    VerifiedAt        time.Time  `json:"verified_at,omitempty"`
    Documents         []Document `json:"documents"`
    ServiceTypes      []string  `json:"service_types"`
    OperatingRegions []string  `json:"operating_regions"`
    Certifications   []Certification `json:"certifications"`
    Licenses         []License  `json:"licenses"`
}

// Certification represents a business certification
type Certification struct {
    BaseModel
    Type            string    `json:"type"`
    Number          string    `json:"number"`
    IssuedBy        string    `json:"issued_by"`
    IssuedDate      time.Time `json:"issued_date"`
    ExpiryDate      time.Time `json:"expiry_date"`
    Status          string    `json:"status"`
    DocumentURL     string    `json:"document_url"`
}

// CompanyBranch represents a company branch office
type CompanyBranch struct {
    BaseModel
    CompanyID       string    `json:"company_id"`
    Name            string    `json:"name"`
    Type            string    `json:"type"` // HQ, BRANCH, SATELLITE
    Address         Address   `json:"address"`
    Contact         Contact   `json:"contact"`
    OperatingHours  Schedule  `json:"operating_hours"`
    Services        []string  `json:"services"`
    IsActive        bool      `json:"is_active"`
}

// ProfileSettings represents user profile settings
type ProfileSettings struct {
    BaseModel
    UserID          string    `json:"user_id"`
    Notifications   struct {
        Email       bool      `json:"email"`
        SMS         bool      `json:"sms"`
        Push        bool      `json:"push"`
    } `json:"notifications"`
    Privacy         struct {
        ShowEmail   bool      `json:"show_email"`
        ShowPhone   bool      `json:"show_phone"`
        ShowProfile bool      `json:"show_profile"`
    } `json:"privacy"`
    Communication   struct {
        Language    string    `json:"language"`
        TimeZone    string    `json:"timezone"`
        Currency    string    `json:"currency"`
    } `json:"communication"`
}

// ProfileVerification represents verification details
type ProfileVerification struct {
    BaseModel
    ProfileID       string    `json:"profile_id"`
    Type            string    `json:"type"`
    Status          string    `json:"status"`
    VerifiedBy      string    `json:"verified_by"`
    VerifiedAt      time.Time `json:"verified_at"`
    ExpiresAt       time.Time `json:"expires_at"`
    Documents       []Document `json:"documents"`
    Notes           string    `json:"notes"`
}

// ProfileActivity represents profile activity logs
type ProfileActivity struct {
    BaseModel
    ProfileID       string    `json:"profile_id"`
    Type            string    `json:"type"`
    Description     string    `json:"description"`
    IPAddress       string    `json:"ip_address"`
    UserAgent       string    `json:"user_agent"`
    Location        string    `json:"location"`
}

// CompanyReview represents a review for a company
type CompanyReview struct {
    BaseModel
    CompanyID       string    `json:"company_id"`
    ReviewerID      string    `json:"reviewer_id"`
    Rating          float32   `json:"rating"`
    Title           string    `json:"title"`
    Content         string    `json:"content"`
    Categories      map[string]float32 `json:"categories"`
    Response        string    `json:"response,omitempty"`
    ResponseDate    time.Time `json:"response_date,omitempty"`
    IsVerified      bool      `json:"is_verified"`
}
