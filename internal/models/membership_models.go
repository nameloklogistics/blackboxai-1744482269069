package models

import (
	"fmt"
	"time"
)

// MembershipStatus represents the status of a membership
type MembershipStatus string

const (
	MembershipTrial    MembershipStatus = "TRIAL"
	MembershipActive   MembershipStatus = "ACTIVE"
	MembershipExpired  MembershipStatus = "EXPIRED"
	MembershipCanceled MembershipStatus = "CANCELED"
)

// MembershipType represents the type of membership
type MembershipType string

const (
	FreightForwarderMembership MembershipType = "FREIGHT_FORWARDER"
	CustomsBrokerMembership    MembershipType = "CUSTOMS_BROKER"
)

// Membership represents a subscription membership
type Membership struct {
	ID              string           `json:"id"`
	MemberType      MembershipType   `json:"member_type"`
	MemberID        string           `json:"member_id"`      // ID of the freight forwarder or customs broker
	Status          MembershipStatus `json:"status"`
	StartDate       time.Time        `json:"start_date"`
	EndDate         time.Time        `json:"end_date"`
	TrialEndDate    time.Time        `json:"trial_end_date"`
	LastRenewalDate time.Time        `json:"last_renewal_date"`
	NextRenewalDate time.Time        `json:"next_renewal_date"`
	AnnualFeeUSD    float64         `json:"annual_fee_usd"`
	IsAutoRenew     bool            `json:"is_auto_renew"`
	PaymentMethod   string          `json:"payment_method,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// MembershipConfig contains configuration for memberships
type MembershipConfig struct {
	AnnualFeeUSD     float64
	TrialPeriodDays  int
	GracePeriodDays  int
	RenewalReminders []int // Days before expiry to send reminders
}

// DefaultMembershipConfig returns the default membership configuration
func DefaultMembershipConfig() MembershipConfig {
	return MembershipConfig{
		AnnualFeeUSD:    150.0,
		TrialPeriodDays: 30,
		GracePeriodDays: 15,
		RenewalReminders: []int{30, 15, 7, 3, 1},
	}
}

// MembershipError represents an error related to membership operations
type MembershipError struct {
	MemberID   string
	MemberType MembershipType
	Reason     string
}

func (e *MembershipError) Error() string {
	return fmt.Sprintf("membership error for %s %s: %s", e.MemberType, e.MemberID, e.Reason)
}

// Payment represents a membership payment
type Payment struct {
	ID            string      `json:"id"`
	MembershipID  string      `json:"membership_id"`
	Amount        float64     `json:"amount"`
	Currency      string      `json:"currency"`
	PaymentMethod string      `json:"payment_method"`
	Status        string      `json:"status"`
	PaidAt        time.Time   `json:"paid_at"`
	CreatedAt     time.Time   `json:"created_at"`
}
