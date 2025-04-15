package services

import (
	"context"
	"fmt"
	"time"

	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/stellar"
)

// CustomsRateService handles customs rate operations
type CustomsRateService struct {
	txManager              *stellar.TransactionManager
	tokenManager           *stellar.TokenManager
	licenseVerificationSvc LicenseVerificationService
}

// NewCustomsRateService creates a new CustomsRateService instance
func NewCustomsRateService(
	txManager *stellar.TransactionManager,
	tokenManager *stellar.TokenManager,
	licenseVerificationSvc LicenseVerificationService,
) *CustomsRateService {
	return &CustomsRateService{
		txManager:              txManager,
		tokenManager:           tokenManager,
		licenseVerificationSvc: licenseVerificationSvc,
	}
}
