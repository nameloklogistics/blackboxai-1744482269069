package services

import (
	"context"
	"fmt"
	"time"

	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/stellar"
)

// MembershipService handles membership operations
type MembershipService struct {
	txManager    *stellar.TransactionManager
	tokenManager *stellar.TokenManager
	config       models.MembershipConfig
}

// NewMembershipService creates a new MembershipService instance
func NewMembershipService(
	txManager *stellar.TransactionManager,
	tokenManager *stellar.TokenManager,
) *MembershipService {
	return &MembershipService{
		txManager:    txManager,
		tokenManager: tokenManager,
		config:       models.DefaultMembershipConfig(),
	}
}

// CreateMembership creates a new membership with trial period
func (s *MembershipService) CreateMembership(memberType models.MembershipType, memberID string) (*models.Membership, error) {
	now := time.Now()
	trialEnd := now.AddDate(0, 0, s.config.TrialPeriodDays)
	membershipEnd := trialEnd.AddDate(1, 0, 0) // 1 year after trial ends

	membership := &models.Membership{
		ID:              fmt.Sprintf("MEM-%d", time.Now().UnixNano()),
		MemberType:      memberType,
		MemberID:        memberID,
		Status:          models.MembershipTrial,
		StartDate:       now,
		EndDate:         membershipEnd,
		TrialEndDate:    trialEnd,
		LastRenewalDate: now,
		NextRenewalDate: trialEnd,
		AnnualFeeUSD:    s.config.AnnualFeeUSD,
		IsAutoRenew:     false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	return membership, nil
}

// ActivateMembership activates a membership after trial or renewal
func (s *MembershipService) ActivateMembership(ctx context.Context, membershipID string, payment *models.Payment) error {
	membership, err := s.GetMembership(membershipID)
	if err != nil {
		return err
	}

	if membership.Status != models.MembershipTrial && membership.Status != models.MembershipExpired {
		return &models.MembershipError{
			MemberID:   membership.MemberID,
			MemberType: membership.MemberType,
			Reason:     fmt.Sprintf("cannot activate membership in %s status", membership.Status),
		}
	}

	if payment.Amount < membership.AnnualFeeUSD {
		return fmt.Errorf("payment amount %.2f USD is less than required fee %.2f USD", 
			payment.Amount, membership.AnnualFeeUSD)
	}

	now := time.Now()
	membership.Status = models.MembershipActive
	membership.LastRenewalDate = now
	membership.NextRenewalDate = now.AddDate(1, 0, 0)
	membership.EndDate = membership.NextRenewalDate
	membership.UpdatedAt = now

	return nil
}

// RenewMembership renews an active membership
func (s *MembershipService) RenewMembership(ctx context.Context, membershipID string, payment *models.Payment) error {
	membership, err := s.GetMembership(membershipID)
	if err != nil {
		return err
	}

	if membership.Status != models.MembershipActive {
		return &models.MembershipError{
			MemberID:   membership.MemberID,
			MemberType: membership.MemberType,
			Reason:     fmt.Sprintf("cannot renew membership in %s status", membership.Status),
		}
	}

	if payment.Amount < membership.AnnualFeeUSD {
		return fmt.Errorf("payment amount %.2f USD is less than required fee %.2f USD", 
			payment.Amount, membership.AnnualFeeUSD)
	}

	now := time.Now()
	membership.LastRenewalDate = now
	membership.NextRenewalDate = now.AddDate(1, 0, 0)
	membership.EndDate = membership.NextRenewalDate
	membership.UpdatedAt = now

	return nil
}

// CancelMembership cancels a membership
func (s *MembershipService) CancelMembership(ctx context.Context, membershipID string) error {
	membership, err := s.GetMembership(membershipID)
	if err != nil {
		return err
	}

	if membership.Status != models.MembershipActive && membership.Status != models.MembershipTrial {
		return &models.MembershipError{
			MemberID:   membership.MemberID,
			MemberType: membership.MemberType,
			Reason:     fmt.Sprintf("cannot cancel membership in %s status", membership.Status),
		}
	}

	membership.Status = models.MembershipCanceled
	membership.UpdatedAt = time.Now()
	membership.IsAutoRenew = false

	return nil
}

// GetMembership retrieves a membership by ID
func (s *MembershipService) GetMembership(membershipID string) (*models.Membership, error) {
	// In a real implementation, this would query a database
	return nil, fmt.Errorf("membership not found: %s", membershipID)
}

// GetMembershipByMember retrieves a membership by member ID and type
func (s *MembershipService) GetMembershipByMember(memberType models.MembershipType, memberID string) (*models.Membership, error) {
	// In a real implementation, this would query a database
	return nil, fmt.Errorf("membership not found for %s %s", memberType, memberID)
}

// ValidateMembership checks if a membership is valid and active
func (s *MembershipService) ValidateMembership(ctx context.Context, memberType models.MembershipType, memberID string) error {
	membership, err := s.GetMembershipByMember(memberType, memberID)
	if err != nil {
		return err
	}

	now := time.Now()

	switch membership.Status {
	case models.MembershipTrial:
		if now.After(membership.TrialEndDate) {
			membership.Status = models.MembershipExpired
			membership.UpdatedAt = now
			return &models.MembershipError{
				MemberID:   memberID,
				MemberType: memberType,
				Reason:     "trial period has expired",
			}
		}
	case models.MembershipActive:
		if now.After(membership.EndDate) {
			membership.Status = models.MembershipExpired
			membership.UpdatedAt = now
			return &models.MembershipError{
				MemberID:   memberID,
				MemberType: memberType,
				Reason:     "membership has expired",
			}
		}
	case models.MembershipExpired:
		return &models.MembershipError{
			MemberID:   memberID,
			MemberType: memberType,
			Reason:     "membership is expired",
		}
	case models.MembershipCanceled:
		return &models.MembershipError{
			MemberID:   memberID,
			MemberType: memberType,
			Reason:     "membership is canceled",
		}
	}

	return nil
}

// ProcessPayment processes a membership payment
func (s *MembershipService) ProcessPayment(ctx context.Context, membershipID string, amount float64, method string) (*models.Payment, error) {
	membership, err := s.GetMembership(membershipID)
	if err != nil {
		return nil, err
	}

	payment := &models.Payment{
		ID:            fmt.Sprintf("PAY-%d", time.Now().UnixNano()),
		MembershipID:  membershipID,
		Amount:        amount,
		Currency:      "USD",
		PaymentMethod: method,
		Status:        "COMPLETED",
		PaidAt:        time.Now(),
		CreatedAt:     time.Now(),
	}

	// In a real implementation, this would:
	// 1. Process payment through payment gateway
	// 2. Store payment record
	// 3. Update membership status

	return payment, nil
}

// ScheduleMembershipReminders schedules renewal reminders for a membership
func (s *MembershipService) ScheduleMembershipReminders(ctx context.Context, membershipID string) error {
	membership, err := s.GetMembership(membershipID)
	if err != nil {
		return err
	}

	// In a real implementation, this would schedule reminders for each day in config.RenewalReminders
	// before membership.NextRenewalDate

	return nil
}
