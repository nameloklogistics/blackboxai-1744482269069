package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"logistics-marketplace/internal/models"
)

// Mock services
type MockMarketplaceService struct {
	mock.Mock
}

type MockCustomsService struct {
	mock.Mock
}

func (m *MockMarketplaceService) GetServicesByCategory(category string) ([]models.LogisticsService, error) {
	args := m.Called(category)
	return args.Get(0).([]models.LogisticsService), args.Error(1)
}

func (m *MockMarketplaceService) GetServiceItems(subcategoryID string) ([]models.ServiceItem, error) {
	args := m.Called(subcategoryID)
	return args.Get(0).([]models.ServiceItem), args.Error(1)
}

func (m *MockMarketplaceService) CreateServiceListing(service *models.LogisticsService) error {
	args := m.Called(service)
	return args.Error(0)
}

func (m *MockMarketplaceService) GetQuotation(serviceID string, cargoDetails *models.CargoDetails) (*models.Rate, error) {
	args := m.Called(serviceID, cargoDetails)
	return args.Get(0).(*models.Rate), args.Error(1)
}

func (m *MockMarketplaceService) CreateBooking(booking *models.Booking) error {
	args := m.Called(booking)
	return args.Error(0)
}

func (m *MockMarketplaceService) ProcessPayment(bookingID string, customerID string, amount float64) error {
	args := m.Called(bookingID, customerID, amount)
	return args.Error(0)
}

func (m *MockMarketplaceService) UpdateShipmentStatus(event *models.TrackingEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockCustomsService) InitiateCustomsClearance(clearance *models.CustomsClearance) error {
	args := m.Called(clearance)
	return args.Error(0)
}

func (m *MockCustomsService) SubmitDeclaration(clearanceID string, declarationType string, documents []string) error {
	args := m.Called(clearanceID, declarationType, documents)
	return args.Error(0)
}

func (m *MockCustomsService) GetRequiredDocuments(declarationType string) []string {
	args := m.Called(declarationType)
	return args.Get(0).([]string)
}

func setupTest() (*gin.Engine, *MockMarketplaceService, *MockCustomsService) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	mockMarketplace := new(MockMarketplaceService)
	mockCustoms := new(MockCustomsService)
	handler := NewMarketplaceHandler(mockMarketplace, mockCustoms)

	// Setup routes
	r.GET("/services/:category", handler.ListServicesByCategory)
	r.GET("/services/:category/subcategories", handler.ListSubCategories)
	r.GET("/services/items/:subcategory", handler.ListServiceItems)
	r.POST("/services", handler.CreateServiceListing)
	r.POST("/quotations", handler.GetQuotation)
	r.POST("/bookings", handler.CreateBooking)
	r.POST("/payments", handler.ProcessPayment)
	r.POST("/shipments/status", handler.UpdateShipmentStatus)
	r.POST("/customs/clearance", handler.InitiateCustomsClearance)
	r.POST("/customs/declaration", handler.SubmitCustomsDeclaration)
	r.GET("/customs/documents", handler.GetRequiredDocuments)

	return r, mockMarketplace, mockCustoms
}

func TestListServicesByCategory(t *testing.T) {
	r, mock, _ := setupTest()

	t.Run("valid category returns services", func(t *testing.T) {
		services := []models.LogisticsService{
			{
				ID:          "1",
				Name:        "Import Service 1",
				Category:    models.ImportService,
				Provider:    models.Provider{ID: "provider1"},
				IsActive:    true,
				BaseRate:    100.0,
			},
		}

		mock.On("GetServicesByCategory", "IMPORT").Return(services, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/services/IMPORT", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response []models.LogisticsService
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, services, response)
	})

	t.Run("invalid category returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/services/INVALID", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid service category", response["error"])
	})

	t.Run("service error returns internal server error", func(t *testing.T) {
		mock.On("GetServicesByCategory", "IMPORT").Return([]models.LogisticsService{}, errors.New("database error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/services/IMPORT", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "database error", response["error"])
	})

	t.Run("empty service list returns success", func(t *testing.T) {
		mock.On("GetServicesByCategory", "IMPORT").Return([]models.LogisticsService{}, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/services/IMPORT", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []models.LogisticsService
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Empty(t, response)
	})
}

func TestCreateServiceListing(t *testing.T) {
	r, mock, _ := setupTest()

	t.Run("valid service", func(t *testing.T) {
		service := models.LogisticsService{
			Name:     "Test Service",
			Category: models.ImportService,
			Provider: models.Provider{ID: "provider1"},
			IsActive: true,
			BaseRate: 100.0,
		}

		mock.On("CreateServiceListing", &service).Return(nil).Once()

		body, _ := json.Marshal(service)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		// Set user context
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("user_id", "provider1")
		req = req.WithContext(ctx)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid service data", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/services", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetQuotation(t *testing.T) {
	r, mock, _ := setupTest()

	validRequest := struct {
		ServiceID    string            `json:"service_id"`
		CargoDetails models.CargoDetails `json:"cargo_details"`
	}{
		ServiceID: "service1",
		CargoDetails: models.CargoDetails{
			Weight:      1000,
			Volume:      10,
			Type:        "General",
			Description: "Test cargo",
		},
	}

	expectedRate := &models.Rate{
		BaseRate:     100.0,
		Surcharges:   20.0,
		TotalRate:    120.0,
		CurrencyCode: "USD",
	}

	t.Run("valid quotation request returns rate", func(t *testing.T) {
		mock.On("GetQuotation", validRequest.ServiceID, &validRequest.CargoDetails).Return(expectedRate, nil).Once()

		body, _ := json.Marshal(validRequest)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/quotations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Rate
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedRate, &response)
	})

	t.Run("invalid JSON returns bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/quotations", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "invalid")
	})

	t.Run("missing required fields returns bad request", func(t *testing.T) {
		invalidRequest := struct {
			ServiceID string `json:"service_id"`
		}{
			ServiceID: "service1",
		}

		body, _ := json.Marshal(invalidRequest)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/quotations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error returns internal server error", func(t *testing.T) {
		mock.On("GetQuotation", validRequest.ServiceID, &validRequest.CargoDetails).Return(nil, errors.New("rate calculation error")).Once()

		body, _ := json.Marshal(validRequest)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/quotations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "rate calculation error", response["error"])
	})
}

func TestCreateBooking(t *testing.T) {
	r, mock, _ := setupTest()

	validBooking := models.Booking{
		ServiceID: "service1",
		CargoDetails: models.CargoDetails{
			Weight:      1000,
			Volume:      10,
			Type:        "General",
			Description: "Test cargo",
		},
	}

	t.Run("valid booking creates successfully", func(t *testing.T) {
		bookingWithCustomer := validBooking
		bookingWithCustomer.CustomerID = "customer1"
		mock.On("CreateBooking", &bookingWithCustomer).Return(nil).Once()

		body, _ := json.Marshal(validBooking)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("user_id", "customer1")
		req = req.WithContext(ctx)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.Booking
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "customer1", response.CustomerID)
	})

	t.Run("invalid JSON returns bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/bookings", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "invalid")
	})

	t.Run("missing user ID returns error", func(t *testing.T) {
		body, _ := json.Marshal(validBooking)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "user_id")
	})

	t.Run("booking creation error returns internal server error", func(t *testing.T) {
		bookingWithCustomer := validBooking
		bookingWithCustomer.CustomerID = "customer1"
		mock.On("CreateBooking", &bookingWithCustomer).Return(errors.New("database error")).Once()

		body, _ := json.Marshal(validBooking)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("user_id", "customer1")
		req = req.WithContext(ctx)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "database error", response["error"])
	})

	t.Run("invalid booking fields return validation error", func(t *testing.T) {
		invalidBooking := models.Booking{
			ServiceID: "", // Empty service ID should fail validation
			CargoDetails: models.CargoDetails{
				Weight:      -1000, // Negative weight should fail validation
				Volume:      -10,   // Negative volume should fail validation
				Type:        "",    // Empty type should fail validation
				Description: "Test cargo",
			},
		}

		body, _ := json.Marshal(invalidBooking)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("user_id", "customer1")
		req = req.WithContext(ctx)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "validation")
	})
}

func TestProcessPayment(t *testing.T) {
	r, mock, _ := setupTest()

	validRequest := struct {
		BookingID string  `json:"booking_id"`
		Amount    float64 `json:"amount"`
	}{
		BookingID: "booking1",
		Amount:    100.0,
	}

	setupRequest := func(t *testing.T, payload interface{}, userID string) (*httptest.ResponseRecorder, *http.Request) {
		body, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/payments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		if userID != "" {
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set("user_id", userID)
			req = req.WithContext(ctx)
		}
		
		return w, req
	}

	t.Run("valid payment processes successfully", func(t *testing.T) {
		mock.On("ProcessPayment", validRequest.BookingID, "customer1", validRequest.Amount).Return(nil).Once()

		w, req := setupRequest(t, validRequest, "customer1")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "payment processed", response["status"])
	})

	t.Run("invalid JSON returns bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/payments", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "invalid")
	})

	t.Run("missing user ID returns error", func(t *testing.T) {
		w, req := setupRequest(t, validRequest, "")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "user_id")
	})

	t.Run("payment processing error returns internal server error", func(t *testing.T) {
		mock.On("ProcessPayment", validRequest.BookingID, "customer1", validRequest.Amount).Return(errors.New("payment failed")).Once()

		w, req := setupRequest(t, validRequest, "customer1")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "payment failed", response["error"])
	})

	t.Run("invalid payment fields return validation error", func(t *testing.T) {
		invalidRequests := []struct {
			name    string
			request interface{}
		}{
			{
				name: "negative amount",
				request: struct {
					BookingID string  `json:"booking_id"`
					Amount    float64 `json:"amount"`
				}{
					BookingID: "booking1",
					Amount:    -100.0,
				},
			},
			{
				name: "empty booking ID",
				request: struct {
					BookingID string  `json:"booking_id"`
					Amount    float64 `json:"amount"`
				}{
					BookingID: "",
					Amount:    100.0,
				},
			},
			{
				name: "zero amount",
				request: struct {
					BookingID string  `json:"booking_id"`
					Amount    float64 `json:"amount"`
				}{
					BookingID: "booking1",
					Amount:    0.0,
				},
			},
		}

		for _, tc := range invalidRequests {
			t.Run(tc.name, func(t *testing.T) {
				w, req := setupRequest(t, tc.request, "customer1")
				r.ServeHTTP(w, req)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				var response gin.H
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], "validation")
			})
		}
	})
}

func TestInitiateCustomsClearance(t *testing.T) {
	r, _, mockCustoms := setupTest()

	validClearance := models.CustomsClearance{
		BookingID:       "booking1",
		DeclarationType: "IMPORT",
		Status:         "PENDING",
		Documents:      []string{"invoice.pdf", "packing_list.pdf"},
		CargoDetails: models.CargoDetails{
			Weight:      1000,
			Volume:      10,
			Type:        "General",
			Description: "Test cargo",
		},
	}

	setupRequest := func(t *testing.T, payload interface{}) (*httptest.ResponseRecorder, *http.Request) {
		body, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/customs/clearance", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		return w, req
	}

	t.Run("valid clearance initiates successfully", func(t *testing.T) {
		mockCustoms.On("InitiateCustomsClearance", &validClearance).Return(nil).Once()

		w, req := setupRequest(t, validClearance)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.CustomsClearance
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, validClearance.BookingID, response.BookingID)
		assert.Equal(t, validClearance.DeclarationType, response.DeclarationType)
		assert.Equal(t, validClearance.Status, response.Status)
		assert.Equal(t, len(validClearance.Documents), len(response.Documents))
	})

	t.Run("invalid JSON returns bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/customs/clearance", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "invalid")
	})

	t.Run("service error returns internal server error", func(t *testing.T) {
		mockCustoms.On("InitiateCustomsClearance", &validClearance).Return(errors.New("service error")).Once()

		w, req := setupRequest(t, validClearance)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "service error", response["error"])
	})

	t.Run("invalid clearance fields return validation error", func(t *testing.T) {
		invalidCases := []struct {
			name      string
			clearance models.CustomsClearance
			errorMsg  string
		}{
			{
				name: "empty booking ID",
				clearance: models.CustomsClearance{
					BookingID:       "",
					DeclarationType: "IMPORT",
					Status:         "PENDING",
					Documents:      []string{"invoice.pdf"},
				},
				errorMsg: "booking_id",
			},
			{
				name: "invalid declaration type",
				clearance: models.CustomsClearance{
					BookingID:       "booking1",
					DeclarationType: "INVALID",
					Status:         "PENDING",
					Documents:      []string{"invoice.pdf"},
				},
				errorMsg: "declaration_type",
			},
			{
				name: "invalid status",
				clearance: models.CustomsClearance{
					BookingID:       "booking1",
					DeclarationType: "IMPORT",
					Status:         "INVALID",
					Documents:      []string{"invoice.pdf"},
				},
				errorMsg: "status",
			},
			{
				name: "missing required documents",
				clearance: models.CustomsClearance{
					BookingID:       "booking1",
					DeclarationType: "IMPORT",
					Status:         "PENDING",
					Documents:      []string{},
				},
				errorMsg: "documents",
			},
			{
				name: "invalid cargo details",
				clearance: models.CustomsClearance{
					BookingID:       "booking1",
					DeclarationType: "IMPORT",
					Status:         "PENDING",
					Documents:      []string{"invoice.pdf"},
					CargoDetails: models.CargoDetails{
						Weight:      -1000, // Negative weight
						Volume:      -10,   // Negative volume
						Type:        "",    // Empty type
						Description: "",    // Empty description
					},
				},
				errorMsg: "cargo_details",
			},
		}

		for _, tc := range invalidCases {
			t.Run(tc.name, func(t *testing.T) {
				w, req := setupRequest(t, tc.clearance)
				r.ServeHTTP(w, req)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				var response gin.H
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tc.errorMsg)
			})
		}
	})

	t.Run("duplicate clearance returns conflict error", func(t *testing.T) {
		mockCustoms.On("InitiateCustomsClearance", &validClearance).Return(errors.New("clearance already exists")).Once()

		w, req := setupRequest(t, validClearance)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "already exists")
	})
}

func TestGetRequiredDocuments(t *testing.T) {
	r, _, mockCustoms := setupTest()

	setupRequest := func(t *testing.T, declarationType string) (*httptest.ResponseRecorder, *http.Request) {
		w := httptest.NewRecorder()
		url := "/customs/documents"
		if declarationType != "" {
			url += "?declaration_type=" + declarationType
		}
		req, _ := http.NewRequest("GET", url, nil)
		return w, req
	}

	t.Run("valid declaration types return correct documents", func(t *testing.T) {
		testCases := []struct {
			name            string
			declarationType string
			expectedDocs    []string
		}{
			{
				name:            "import declaration",
				declarationType: "IMPORT",
				expectedDocs: []string{
					"Commercial Invoice",
					"Packing List",
					"Bill of Lading",
					"Import License",
					"Customs Declaration Form",
				},
			},
			{
				name:            "export declaration",
				declarationType: "EXPORT",
				expectedDocs: []string{
					"Commercial Invoice",
					"Packing List",
					"Export License",
					"Certificate of Origin",
					"Shipping Bill",
				},
			},
			{
				name:            "transit declaration",
				declarationType: "TRANSIT",
				expectedDocs: []string{
					"Transit Permit",
					"Cargo Manifest",
					"Transport Document",
					"Insurance Certificate",
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				mockCustoms.On("GetRequiredDocuments", tc.declarationType).Return(tc.expectedDocs).Once()

				w, req := setupRequest(t, tc.declarationType)
				r.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)
				var response struct {
					Documents []string `json:"documents"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDocs, response.Documents)
				assert.True(t, len(response.Documents) > 0, "Document list should not be empty")
			})
		}
	})

	t.Run("missing declaration type returns bad request", func(t *testing.T) {
		w, req := setupRequest(t, "")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "declaration_type is required", response["error"])
	})

	t.Run("invalid declaration type returns bad request", func(t *testing.T) {
		w, req := setupRequest(t, "INVALID")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "invalid declaration type")
	})

	t.Run("empty document list returns success with empty array", func(t *testing.T) {
		mockCustoms.On("GetRequiredDocuments", "IMPORT").Return([]string{}).Once()

		w, req := setupRequest(t, "IMPORT")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response struct {
			Documents []string `json:"documents"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Empty(t, response.Documents)
	})

	t.Run("special declaration types return additional documents", func(t *testing.T) {
		specialCases := []struct {
			name            string
			declarationType string
			expectedDocs    []string
		}{
			{
				name:            "dangerous goods",
				declarationType: "DANGEROUS_GOODS",
				expectedDocs: []string{
					"Dangerous Goods Declaration",
					"Safety Data Sheet",
					"Hazmat Certificate",
					"Emergency Contact Details",
					"Transport Emergency Card",
					"Container Packing Certificate",
				},
			},
			{
				name:            "perishable goods",
				declarationType: "PERISHABLE",
				expectedDocs: []string{
					"Temperature Requirements",
					"Health Certificate",
					"Inspection Certificate",
					"Storage Instructions",
					"Time Sensitivity Declaration",
					"Quality Control Report",
				},
			},
			{
				name:            "high value goods",
				declarationType: "HIGH_VALUE",
				expectedDocs: []string{
					"Insurance Certificate",
					"Valuation Certificate",
					"Security Protocol",
					"Bank Guarantee",
					"Special Handling Instructions",
				},
			},
		}

		for _, tc := range specialCases {
			t.Run(tc.name, func(t *testing.T) {
				mockCustoms.On("GetRequiredDocuments", tc.declarationType).Return(tc.expectedDocs).Once()

				w, req := setupRequest(t, tc.declarationType)
				r.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)
				var response struct {
					Documents []string `json:"documents"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDocs, response.Documents)
				assert.True(t, len(response.Documents) >= 5, "Special cases should require at least 5 documents")
			})
		}
	})
}

func TestSubmitCustomsDeclaration(t *testing.T) {
	r, _, mockCustoms := setupTest()

	validRequest := struct {
		ClearanceID     string   `json:"clearance_id"`
		DeclarationType string   `json:"declaration_type"`
		Documents       []string `json:"documents"`
	}{
		ClearanceID:     "clearance1",
		DeclarationType: "IMPORT",
		Documents:       []string{"invoice.pdf", "packing_list.pdf", "bill_of_lading.pdf"},
	}

	setupRequest := func(t *testing.T, payload interface{}) (*httptest.ResponseRecorder, *http.Request) {
		body, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/customs/declaration", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		return w, req
	}

	t.Run("valid declaration submits successfully", func(t *testing.T) {
		mockCustoms.On("SubmitDeclaration",
			validRequest.ClearanceID,
			validRequest.DeclarationType,
			validRequest.Documents,
		).Return(nil).Once()

		w, req := setupRequest(t, validRequest)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "declaration submitted", response["status"])
	})

	t.Run("invalid JSON returns bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/customs/declaration", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "invalid")
	})

	t.Run("missing required fields return validation error", func(t *testing.T) {
		invalidCases := []struct {
			name    string
			request interface{}
			errMsg  string
		}{
			{
				name: "missing clearance ID",
				request: struct {
					DeclarationType string   `json:"declaration_type"`
					Documents       []string `json:"documents"`
				}{
					DeclarationType: "IMPORT",
					Documents:       []string{"invoice.pdf"},
				},
				errMsg: "clearance_id is required",
			},
			{
				name: "missing declaration type",
				request: struct {
					ClearanceID string   `json:"clearance_id"`
					Documents   []string `json:"documents"`
				}{
					ClearanceID: "clearance1",
					Documents:   []string{"invoice.pdf"},
				},
				errMsg: "declaration_type is required",
			},
			{
				name: "empty documents array",
				request: struct {
					ClearanceID     string   `json:"clearance_id"`
					DeclarationType string   `json:"declaration_type"`
					Documents       []string `json:"documents"`
				}{
					ClearanceID:     "clearance1",
					DeclarationType: "IMPORT",
					Documents:       []string{},
				},
				errMsg: "at least one document is required",
			},
		}

		for _, tc := range invalidCases {
			t.Run(tc.name, func(t *testing.T) {
				w, req := setupRequest(t, tc.request)
				r.ServeHTTP(w, req)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				var response gin.H
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tc.errMsg)
			})
		}
	})

	t.Run("service error returns internal server error", func(t *testing.T) {
		mockCustoms.On("SubmitDeclaration",
			validRequest.ClearanceID,
			validRequest.DeclarationType,
			validRequest.Documents,
		).Return(errors.New("service error")).Once()

		w, req := setupRequest(t, validRequest)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "service error", response["error"])
	})

	t.Run("invalid document formats return validation error", func(t *testing.T) {
		invalidFormats := []struct {
			name     string
			document string
		}{
			{name: "executable file", document: "document.exe"},
			{name: "shell script", document: "script.sh"},
			{name: "system file", document: "system.dll"},
			{name: "no extension", document: "document"},
			{name: "invalid extension", document: "doc.invalid"},
		}

		for _, tc := range invalidFormats {
			t.Run(tc.name, func(t *testing.T) {
				invalidRequest := struct {
					ClearanceID     string   `json:"clearance_id"`
					DeclarationType string   `json:"declaration_type"`
					Documents       []string `json:"documents"`
				}{
					ClearanceID:     "clearance1",
					DeclarationType: "IMPORT",
					Documents:       []string{tc.document},
				}

				w, req := setupRequest(t, invalidRequest)
				r.ServeHTTP(w, req)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				var response gin.H
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], "invalid document format")
			})
		}
	})

	t.Run("clearance not found returns not found error", func(t *testing.T) {
		mockCustoms.On("SubmitDeclaration",
			validRequest.ClearanceID,
			validRequest.DeclarationType,
			validRequest.Documents,
		).Return(errors.New("clearance not found")).Once()

		w, req := setupRequest(t, validRequest)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "not found")
	})
}
