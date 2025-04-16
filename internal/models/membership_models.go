package models

import (
    "time"
)

// MembershipTier represents different membership levels
const (
    TierBasic    = "BASIC"
    TierPremium  = "PREMIUM"
    TierBusiness = "BUSINESS"
    TierEnterprise = "ENTERPRISE"
)

// MembershipStatus represents the status of a membership
const (
    StatusActive    = "ACTIVE"
    StatusInactive  = "INACTIVE"
    StatusSuspended = "SUSPENDED"
    StatusExpired   = "EXPIRED"
)

// Membership represents a user's membership details
type Membership struct {
    BaseModel
    UserID          string    `json:"user_id"`
    Tier            string    `json:"tier"`
    Status          string    `json:"status"`
    StartDate       time.Time `json:"start_date"`
    EndDate         time.Time `json:"end_date"`
    AutoRenew       bool      `json:"auto_renew"`
    PaymentMethod   string    `json:"payment_method"`
    LastBillingDate time.Time `json:"last_billing_date"`
    NextBillingDate time.Time `json:"next_billing_date"`
}

// MembershipBenefit represents benefits available to different tiers
type MembershipBenefit struct {
    BaseModel
    Tier            string    `json:"tier"`
    Name            string    `json:"name"`
    Description     string    `json:"description"`
    Type            string    `json:"type"` // DISCOUNT, FEATURE, SERVICE
    Value           string    `json:"value"`
    IsActive        bool      `json:"is_active"`
}

// MembershipPlan represents available membership plans
type MembershipPlan struct {
    BaseModel
    Tier            string    `json:"tier"`
    Name            string    `json:"name"`
    Description     string    `json:"description"`
    Price           Currency  `json:"price"`
    BillingCycle    string    `json:"billing_cycle"` // MONTHLY, QUARTERLY, YEARLY
    Benefits        []MembershipBenefit `json:"benefits"`
    IsActive        bool      `json:"is_active"`
}

// MembershipTransaction represents a membership-related transaction
type MembershipTransaction struct {
    BaseModel
    UserID          string    `json:"user_id"`
    MembershipID    string    `json:"membership_id"`
    Type            string    `json:"type"` // PURCHASE, RENEWAL, UPGRADE, CANCELLATION
    Amount          Currency  `json:"amount"`
    PaymentMethod   string    `json:"payment_method"`
    Status          string    `json:"status"`
    TransactionID   string    `json:"transaction_id"`
}

// MembershipDiscount represents discounts available to members
type MembershipDiscount struct {
    BaseModel
    Tier            string    `json:"tier"`
    ServiceType     string    `json:"service_type"`
    DiscountType    string    `json:"discount_type"` // PERCENTAGE, FIXED
    DiscountValue   float64   `json:"discount_value"`
    MinSpend        float64   `json:"min_spend"`
    MaxDiscount     float64   `json:"max_discount"`
    ValidFrom       time.Time `json:"valid_from"`
    ValidUntil      time.Time `json:"valid_until"`
    IsActive        bool      `json:"is_active"`
}

// MembershipUsage represents usage of membership benefits
type MembershipUsage struct {
    BaseModel
    UserID          string    `json:"user_id"`
    MembershipID    string    `json:"membership_id"`
    BenefitID       string    `json:"benefit_id"`
    UsageDate       time.Time `json:"usage_date"`
    UsageAmount     float64   `json:"usage_amount"`
    RemainingQuota  float64   `json:"remaining_quota"`
}

// MembershipInvitation represents invitations for enterprise memberships
type MembershipInvitation struct {
    BaseModel
    InviterID       string    `json:"inviter_id"`
    Email           string    `json:"email"`
    CompanyID       string    `json:"company_id"`
    Role            string    `json:"role"`
    Status          string    `json:"status"`
    ExpiresAt       time.Time `json:"expires_at"`
    AcceptedAt      time.Time `json:"accepted_at,omitempty"`
}

// MembershipNotification represents notifications for membership events
type MembershipNotification struct {
    BaseModel
    UserID          string    `json:"user_id"`
    Type            string    `json:"type"` // EXPIRY, RENEWAL, UPGRADE
    Title           string    `json:"title"`
    Message         string    `json:"message"`
    SentAt          time.Time `json:"sent_at"`
    ReadAt          time.Time `json:"read_at,omitempty"`
}

// MembershipAudit represents audit logs for membership changes
type MembershipAudit struct {
    BaseModel
    UserID          string    `json:"user_id"`
    MembershipID    string    `json:"membership_id"`
    Action          string    `json:"action"`
    PreviousState   string    `json:"previous_state"`
    NewState        string    `json:"new_state"`
    ChangedBy       string    `json:"changed_by"`
    ChangeReason    string    `json:"change_reason"`
}

// MembershipSettings represents user-specific membership settings
type MembershipSettings struct {
    BaseModel
    UserID          string    `json:"user_id"`
    NotifyExpiry    bool      `json:"notify_expiry"`
    NotifyRenewal   bool      `json:"notify_renewal"`
    NotifyUpgrades  bool      `json:"notify_upgrades"`
    AutoRenew       bool      `json:"auto_renew"`
    PreferredPaymentMethod string `json:"preferred_payment_method"`
}
