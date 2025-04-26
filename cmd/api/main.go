package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"logistics-marketplace/cmd/api/handlers"
	"logistics-marketplace/internal/services"
	"logistics-marketplace/internal/stellar"
)

var jwtSecret []byte

func main() {
	jwtSecretEnv := os.Getenv("JWT_SECRET")
	if jwtSecretEnv == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	jwtSecret = []byte(jwtSecretEnv)

	// Initialize Stellar components
	accountManager := stellar.NewAccountManager(true) // Use testnet for development
	tokenManager := stellar.NewTokenManager(
		accountManager,
		"LMT",                    // Token code
		os.Getenv("ISSUER_KEY"), // Get from environment
	)
	txManager := stellar.NewTransactionManager(
		accountManager,
		tokenManager,
		os.Getenv("CONTRACT_ID"), // Get from environment
	)

	// Initialize services
	governanceService := services.NewGovernanceService(
		accountManager,
		tokenManager,
		os.Getenv("GOVERNANCE_CONTRACT_ID"),
	)
	marketplaceService := services.NewMarketplaceService(txManager, tokenManager)
	customsService := services.NewCustomsService(txManager, tokenManager)
	trackingService := services.NewTrackingService(txManager, tokenManager)
	profileService := services.NewProfileService(txManager, tokenManager)
	customsRateService := services.NewCustomsRateService(txManager, tokenManager)
	infrastructureService := services.NewInfrastructureService(txManager, tokenManager)
	userOperationsService := services.NewUserOperationsService(txManager, tokenManager)
	serviceCategoriesService := services.NewServiceCategoriesService(txManager, tokenManager)

	// Initialize handlers
	governanceHandler := handlers.NewGovernanceHandler(governanceService)
	marketplaceHandler := handlers.NewMarketplaceHandler(marketplaceService, customsService)
	trackingHandler := handlers.NewTrackingHandler(trackingService)
	profileHandler := handlers.NewProfileHandler(profileService)
	customsRateHandler := handlers.NewCustomsRateHandler(customsRateService)
	transportHandler := handlers.NewTransportHandler(marketplaceService)
	infrastructureHandler := handlers.NewInfrastructureHandler(infrastructureService)
	userOperationsHandler := handlers.NewUserOperationsHandler(userOperationsService)
	serviceCategoriesHandler := handlers.NewServiceCategoriesHandler(serviceCategoriesService)

	// Initialize Gin router
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(loggingMiddleware())
	router.Use(securityHeadersMiddleware())
	router.Use(corsMiddleware())
	router.Use(authMiddleware())
	router.Use(rateLimitMiddleware(100, time.Minute)) // 100 requests per minute

	// API Routes
	api := router.Group("/api/v1")
	{
		// Governance Routes
		governance := api.Group("/governance")
		{
			// Proposals
			governance.POST("/proposals", governanceHandler.CreateProposal)
			governance.GET("/proposals", governanceHandler.ListProposals)
			governance.GET("/proposals/:id", governanceHandler.GetProposal)
			governance.POST("/proposals/:id/vote", governanceHandler.CastVote)
			governance.POST("/proposals/:id/execute", governanceHandler.ExecuteProposal)
			governance.GET("/proposals/:id/votes/:voter", governanceHandler.GetVote)

			// Parameters
			governance.GET("/parameters/:name", governanceHandler.GetParameter)
		}

		// Service Categories
		services := api.Group("/services")
		{
			// Import Services (Direct imports only)
			imports := services.Group("/imports")
			{
				// Sea Import
				seaImport := imports.Group("/sea")
				{
					seaImport.POST("/", serviceCategoriesHandler.CreateImportSeaService)
					seaImport.GET("/", func(c *gin.Context) {
						c.Param("category", "IMPORT_DIRECT")
						c.Param("mode", "SEA")
						serviceCategoriesHandler.ListServices(c)
					})
				}

				// Air Import
				airImport := imports.Group("/air")
				{
					airImport.POST("/", serviceCategoriesHandler.CreateImportAirService)
					airImport.GET("/", func(c *gin.Context) {
						c.Param("category", "IMPORT_DIRECT")
						c.Param("mode", "AIR")
						serviceCategoriesHandler.ListServices(c)
					})
				}

				// Rail Import
				railImport := imports.Group("/rail")
				{
					railImport.POST("/", serviceCategoriesHandler.CreateImportRailService)
					railImport.GET("/", func(c *gin.Context) {
						c.Param("category", "IMPORT_DIRECT")
						c.Param("mode", "RAIL")
						serviceCategoriesHandler.ListServices(c)
					})
				}

				// Land Import
				landImport := imports.Group("/land")
				{
					landImport.POST("/", serviceCategoriesHandler.CreateImportLandService)
					landImport.GET("/", func(c *gin.Context) {
						c.Param("category", "IMPORT_DIRECT")
						c.Param("mode", "LAND")
						serviceCategoriesHandler.ListServices(c)
					})
				}
			}

			// Export Services (Direct exports only)
			exports := services.Group("/exports")
			{
				// Sea Export
				seaExport := exports.Group("/sea")
				{
					seaExport.POST("/", serviceCategoriesHandler.CreateExportSeaService)
					seaExport.GET("/", func(c *gin.Context) {
						c.Param("category", "EXPORT_DIRECT")
						c.Param("mode", "SEA")
						serviceCategoriesHandler.ListServices(c)
					})
				}

				// Air Export
				airExport := exports.Group("/air")
				{
					airExport.POST("/", serviceCategoriesHandler.CreateExportAirService)
					airExport.GET("/", func(c *gin.Context) {
						c.Param("category", "EXPORT_DIRECT")
						c.Param("mode", "AIR")
						serviceCategoriesHandler.ListServices(c)
					})
				}

				// Rail Export
				railExport := exports.Group("/rail")
				{
					railExport.POST("/", serviceCategoriesHandler.CreateExportRailService)
					railExport.GET("/", func(c *gin.Context) {
						c.Param("category", "EXPORT_DIRECT")
						c.Param("mode", "RAIL")
						serviceCategoriesHandler.ListServices(c)
					})
				}

				// Land Export
				landExport := exports.Group("/land")
				{
					landExport.POST("/", serviceCategoriesHandler.CreateExportLandService)
					landExport.GET("/", func(c *gin.Context) {
						c.Param("category", "EXPORT_DIRECT")
						c.Param("mode", "LAND")
						serviceCategoriesHandler.ListServices(c)
					})
				}
			}

			// Transit Services
			transit := services.Group("/transit")
			{
				// Sea Transit
				seaTransit := transit.Group("/sea")
				{
					seaTransit.POST("/", func(c *gin.Context) {
						c.Param("category", "TRANSIT")
						c.Param("mode", "SEA")
						serviceCategoriesHandler.CreateSeaService(c)
					})
					seaTransit.GET("/", func(c *gin.Context) {
						c.Param("category", "TRANSIT")
						c.Param("mode", "SEA")
						serviceCategoriesHandler.ListServices(c)
					})
				}

				// Air Transit
				airTransit := transit.Group("/air")
				{
					airTransit.POST("/", func(c *gin.Context) {
						c.Param("category", "TRANSIT")
						c.Param("mode", "AIR")
						serviceCategoriesHandler.CreateAirService(c)
					})
					airTransit.GET("/", func(c *gin.Context) {
						c.Param("category", "TRANSIT")
						c.Param("mode", "AIR")
						serviceCategoriesHandler.ListServices(c)
					})
				}

				// Rail Transit
				railTransit := transit.Group("/rail")
				{
					railTransit.POST("/", func(c *gin.Context) {
						c.Param("category", "TRANSIT")
						c.Param("mode", "RAIL")
						serviceCategoriesHandler.CreateRailService(c)
					})
					railTransit.GET("/", func(c *gin.Context) {
						c.Param("category", "TRANSIT")
						c.Param("mode", "RAIL")
						serviceCategoriesHandler.ListServices(c)
					})
				}

				// Land Transit
				landTransit := transit.Group("/land")
				{
					landTransit.POST("/", func(c *gin.Context) {
						c.Param("category", "TRANSIT")
						c.Param("mode", "LAND")
						serviceCategoriesHandler.CreateLandService(c)
					})
					landTransit.GET("/", func(c *gin.Context) {
						c.Param("category", "TRANSIT")
						c.Param("mode", "LAND")
						serviceCategoriesHandler.ListServices(c)
					})
				}
			}

			// Transshipment Services
			transshipment := services.Group("/transshipment")
			{
				// Sea Transshipment
				seaTransshipment := transshipment.Group("/sea")
				{
					seaTransshipment.POST("/", func(c *gin.Context) {
						c.Param("category", "TRANSSHIPMENT")
						c.Param("mode", "SEA")
						serviceCategoriesHandler.CreateSeaService(c)
					})
					seaTransshipment.GET("/", func(c *gin.Context) {
						c.Param("category", "TRANSSHIPMENT")
						c.Param("mode", "SEA")
						serviceCategoriesHandler.ListServices(c)
					})
				}

				// Air Transshipment
				airTransshipment := transshipment.Group("/air")
				{
					airTransshipment.POST("/", func(c *gin.Context) {
						c.Param("category", "TRANSSHIPMENT")
						c.Param("mode", "AIR")
						serviceCategoriesHandler.CreateAirService(c)
					})
					airTransshipment.GET("/", func(c *gin.Context) {
						c.Param("category", "TRANSSHIPMENT")
						c.Param("mode", "AIR")
						serviceCategoriesHandler.ListServices(c)
					})
				}
			}

			// Common Service Operations
			services.GET("/:id", serviceCategoriesHandler.GetService)
			services.GET("/:id/schedule", serviceCategoriesHandler.GetServiceSchedule)
			services.GET("/:id/availability", serviceCategoriesHandler.GetServiceAvailability)
			services.GET("/:id/rates", serviceCategoriesHandler.GetServiceRate)
			services.PUT("/:id/schedule", serviceCategoriesHandler.UpdateServiceSchedule)
			services.PUT("/:id/rates", serviceCategoriesHandler.UpdateServiceRate)
			services.GET("/requirements/:category/:mode", serviceCategoriesHandler.GetServiceRequirements)
		}

		// User Operations
		users := api.Group("/users")
		{
			// Consignee Operations
			consignee := users.Group("/consignee")
			{
				// Quote Management
				consignee.POST("/quotes/request", func(c *gin.Context) {
					c.Param("user_type", "CONSIGNEE")
					userOperationsHandler.CreateQuoteRequest(c)
				})
				consignee.GET("/quotes/request/:id", func(c *gin.Context) {
					c.Param("user_type", "CONSIGNEE")
					userOperationsHandler.GetQuoteRequest(c)
				})
				consignee.GET("/quotes/requests", func(c *gin.Context) {
					c.Param("user_type", "CONSIGNEE")
					userOperationsHandler.ListQuoteRequests(c)
				})
				consignee.POST("/quotes/confirm", func(c *gin.Context) {
					c.Param("user_type", "CONSIGNEE")
					userOperationsHandler.ConfirmQuote(c)
				})

				// Booking Management
				consignee.POST("/bookings", func(c *gin.Context) {
					c.Param("user_type", "CONSIGNEE")
					userOperationsHandler.CreateBooking(c)
				})
				consignee.GET("/bookings/:id", func(c *gin.Context) {
					c.Param("user_type", "CONSIGNEE")
					userOperationsHandler.GetBooking(c)
				})
				consignee.GET("/bookings", func(c *gin.Context) {
					c.Param("user_type", "CONSIGNEE")
					userOperationsHandler.ListBookings(c)
				})
			}

			// Shipper Operations
			shipper := users.Group("/shipper")
			{
				// Quote Management
				shipper.POST("/quotes/request", func(c *gin.Context) {
					c.Param("user_type", "SHIPPER")
					userOperationsHandler.CreateQuoteRequest(c)
				})
				shipper.GET("/quotes/request/:id", func(c *gin.Context) {
					c.Param("user_type", "SHIPPER")
					userOperationsHandler.GetQuoteRequest(c)
				})
				shipper.GET("/quotes/requests", func(c *gin.Context) {
					c.Param("user_type", "SHIPPER")
					userOperationsHandler.ListQuoteRequests(c)
				})
				shipper.POST("/quotes/confirm", func(c *gin.Context) {
					c.Param("user_type", "SHIPPER")
					userOperationsHandler.ConfirmQuote(c)
				})

				// Booking Management
				shipper.POST("/bookings", func(c *gin.Context) {
					c.Param("user_type", "SHIPPER")
					userOperationsHandler.CreateBooking(c)
				})
				shipper.GET("/bookings/:id", func(c *gin.Context) {
					c.Param("user_type", "SHIPPER")
					userOperationsHandler.GetBooking(c)
				})
				shipper.GET("/bookings", func(c *gin.Context) {
					c.Param("user_type", "SHIPPER")
					userOperationsHandler.ListBookings(c)
				})
			}

			// Freight Forwarder Operations
			forwarder := users.Group("/forwarder")
			{
				// Quote Management
				forwarder.GET("/quotes/requests", func(c *gin.Context) {
					c.Param("user_type", "FREIGHT_FORWARDER")
					userOperationsHandler.ListQuoteRequests(c)
				})
				forwarder.POST("/quotes/response", userOperationsHandler.GenerateQuoteResponse)
				forwarder.GET("/quotes/responses", func(c *gin.Context) {
					c.Param("user_type", "FREIGHT_FORWARDER")
					userOperationsHandler.ListQuoteResponses(c)
				})

				// Booking Management
				forwarder.POST("/bookings/confirm", userOperationsHandler.ConfirmBooking)
				forwarder.GET("/bookings", func(c *gin.Context) {
					c.Param("user_type", "FREIGHT_FORWARDER")
					userOperationsHandler.ListBookings(c)
				})
			}
		}

		// Infrastructure Management
		infrastructure := api.Group("/infrastructure")
		{
			infrastructure.GET("/country/:country_code", infrastructureHandler.GetCountryInfrastructure)

			// Airports
			airports := infrastructure.Group("/airports")
			{
				airports.POST("/", infrastructureHandler.CreateAirport)
				airports.GET("/:id/services", infrastructureHandler.GetInfrastructureServices)
				airports.PUT("/:id/capacity", infrastructureHandler.UpdateInfrastructureCapacity)
				airports.PUT("/:id/schedule", infrastructureHandler.UpdateInfrastructureSchedule)
				airports.GET("/:id/capacity", infrastructureHandler.GetInfrastructureCapacity)
				airports.GET("/:id/schedule", infrastructureHandler.GetOperatingSchedule)
			}

			// Seaports
			seaports := infrastructure.Group("/seaports")
			{
				seaports.POST("/", infrastructureHandler.CreateSeaport)
				seaports.GET("/:id/services", infrastructureHandler.GetInfrastructureServices)
				seaports.PUT("/:id/capacity", infrastructureHandler.UpdateInfrastructureCapacity)
				seaports.PUT("/:id/schedule", infrastructureHandler.UpdateInfrastructureSchedule)
				seaports.GET("/:id/capacity", infrastructureHandler.GetInfrastructureCapacity)
				seaports.GET("/:id/schedule", infrastructureHandler.GetOperatingSchedule)
			}

			// Inland Depots
			depots := infrastructure.Group("/depots")
			{
				depots.POST("/", infrastructureHandler.CreateInlandDepot)
				depots.GET("/:id/services", infrastructureHandler.GetInfrastructureServices)
				depots.PUT("/:id/capacity", infrastructureHandler.UpdateInfrastructureCapacity)
				depots.PUT("/:id/schedule", infrastructureHandler.UpdateInfrastructureSchedule)
				depots.GET("/:id/capacity", infrastructureHandler.GetInfrastructureCapacity)
				depots.GET("/:id/schedule", infrastructureHandler.GetOperatingSchedule)
			}

			infrastructure.GET("/nearby", infrastructureHandler.GetNearbyInfrastructure)
			infrastructure.GET("/services/:service_id/locations", infrastructureHandler.GetServiceLocations)
		}

		// Transport Services
		transport := api.Group("/transport")
		{
			// Sea Transport
			sea := transport.Group("/sea")
			{
				sea.POST("/services", transportHandler.CreateSeaService)
				sea.GET("/services", func(c *gin.Context) {
					c.Param("mode", "SEA")
					transportHandler.ListServicesByMode(c)
				})
				sea.GET("/templates", func(c *gin.Context) {
					c.Param("mode", "SEA")
					transportHandler.GetServiceTemplates(c)
				})
				sea.GET("/equipment", func(c *gin.Context) {
					c.Param("mode", "SEA")
					transportHandler.GetEquipmentTypes(c)
				})
			}

			// Air Transport
			air := transport.Group("/air")
			{
				air.POST("/services", transportHandler.CreateAirService)
				air.GET("/services", func(c *gin.Context) {
					c.Param("mode", "AIR")
					transportHandler.ListServicesByMode(c)
				})
				air.GET("/templates", func(c *gin.Context) {
					c.Param("mode", "AIR")
					transportHandler.GetServiceTemplates(c)
				})
				air.GET("/equipment", func(c *gin.Context) {
					c.Param("mode", "AIR")
					transportHandler.GetEquipmentTypes(c)
				})
			}

			// Rail Transport
			rail := transport.Group("/rail")
			{
				rail.POST("/services", transportHandler.CreateRailService)
				rail.GET("/services", func(c *gin.Context) {
					c.Param("mode", "RAIL")
					transportHandler.ListServicesByMode(c)
				})
				rail.GET("/templates", func(c *gin.Context) {
					c.Param("mode", "RAIL")
					transportHandler.GetServiceTemplates(c)
				})
				rail.GET("/equipment", func(c *gin.Context) {
					c.Param("mode", "RAIL")
					transportHandler.GetEquipmentTypes(c)
				})
			}

			// Road Transport
			road := transport.Group("/road")
			{
				road.POST("/services", transportHandler.CreateRoadService)
				road.GET("/services", func(c *gin.Context) {
					c.Param("mode", "ROAD")
					transportHandler.ListServicesByMode(c)
				})
				road.GET("/templates", func(c *gin.Context) {
					c.Param("mode", "ROAD")
					transportHandler.GetServiceTemplates(c)
				})
				road.GET("/equipment", func(c *gin.Context) {
					c.Param("mode", "ROAD")
					transportHandler.GetEquipmentTypes(c)
				})
			}
		}

		// Tracking and Routing
		tracking := api.Group("/tracking")
		{
			tracking.POST("/events", trackingHandler.AddTrackingEvent)
			tracking.GET("/:booking_id", trackingHandler.GetShipmentTracking)
			tracking.POST("/transshipment", trackingHandler.AddTransshipmentPoint)
			tracking.PUT("/routing/:booking_id", trackingHandler.UpdateRouting)
			tracking.POST("/route/optimal", trackingHandler.GetOptimalRoute)
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		log.Printf("%s %s %d %s", c.Request.Method, c.Request.URL.Path, status, latency)
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// Check token expiration
		if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			c.Abort()
			return
		}

		userAddress, ok := claims["user_address"].(string)
		if !ok || userAddress == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_address claim missing"})
			c.Abort()
			return
		}

		c.Set("user_address", userAddress)
		c.Next()
	}
}

func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Referrer-Policy", "no-referrer")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		c.Next()
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := []string{
			"https://yourdomain.com",
			"https://www.yourdomain.com",
		}
		origin := c.GetHeader("Origin")
		allowed := false
		for _, o := range allowedOrigins {
			if o == origin {
				allowed = true
				break
			}
		}
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func rateLimitMiddleware(rps int, burstDuration time.Duration) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(burstDuration/time.Duration(rps)), rps)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}
