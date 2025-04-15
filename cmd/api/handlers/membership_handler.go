package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"logistics-marketplace/internal/models"
	"logistics-marketplace/internal/services"
)

type MembershipHandler struct {
	membershipService *services.MembershipService
}

func NewMembershipHandler(membershipService *services.MembershipService) *MembershipHandler {
	return &MembershipHandler{
		membershipService: membershipService,
	}
}

// ActivateMembership handles activation of membership after trial
func (h *MembershipHandler) ActivateMembership(c *gin.Context) {
	membershipID := c.Param("id")
	var request struct {
		PaymentMethod string  `json:"payment_method"`
		Amount        float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process payment
	payment, err := h.membershipService.ProcessPayment(c, membershipID, request.Amount, request.PaymentMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Activate membership
	if err := h.membershipService.ActivateMembership(c, membershipID, payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "activated",
		"payment": payment,
	})
}

// RenewMembership handles membership renewal
func (h *MembershipHandler) RenewMembership(c *gin.Context) {
	membershipID := c.Param("id")
	var request struct {
		PaymentMethod string  `json:"payment_method"`
		Amount        float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process payment
	payment, err := h.membershipService.ProcessPayment(c, membershipID, request.Amount, request.PaymentMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Renew membership
	if err := h.membershipService.RenewMembership(c, membershipID, payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "renewed",
		"payment": payment,
	})
}

// CancelMembership handles membership cancellation
func (h *MembershipHandler) CancelMembership(c *gin.Context) {
	membershipID := c.Param("id")

	if err := h.membershipService.CancelMembership(c, membershipID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}

// GetMembershipStatus handles membership status retrieval
func (h *MembershipHandler) GetMembershipStatus(c *gin.Context) {
	memberType := models.MembershipType(c.Query("type"))
	memberID := c.Query("member_id")

	if memberType == "" || memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member type and ID are required"})
		return
	}

	// Validate membership
	err := h.membershipService.ValidateMembership(c, memberType, memberID)
	if err != nil {
		if _, ok := err.(*models.MembershipError); ok {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get membership details
	membership, err := h.membershipService.GetMembershipByMember(memberType, memberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, membership)
}
