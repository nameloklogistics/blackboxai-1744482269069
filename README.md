# Logistics Marketplace on Stellar Blockchain

A decentralized two-sided marketplace for logistics and supply chain services built on the Stellar blockchain network. This platform facilitates interactions between freight forwarders, customs brokers, consignees, and shippers.

## Features

### User Profiles

#### Service Providers (Freight Forwarders)
- Detailed company information
- Contact person details
- Infrastructure management:
  - Major sea ports
  - International airports
  - Inland container terminals
- Business certifications (IATA, FIATA)
- Customs broker licenses

#### Service Buyers (Shippers)
- Company information
- Contact details
- Business type (Importer/Exporter)
- Trade documentation (Import/Export codes)

### Service Categories

1. Import Services
   - FCL Import
   - LCL Import
   - Air Import
   - Customs Clearance
   - Door Delivery

2. Export Services
   - FCL Export
   - LCL Export
   - Air Export
   - Export Documentation
   - Pickup Service

3. Transit Services
   - Land Transit
   - Rail Transit
   - Multimodal Transit
   - Transit Documentation

4. Transshipment Services
   - Sea-Sea Transshipment
   - Air-Air Transshipment
   - Sea-Air Transshipment
   - Air-Sea Transshipment

## Technical Architecture

### Smart Contracts
- Token Contract: Manages the logistics marketplace token (LMT)
- Marketplace Contract: Handles service listings, bookings, and payments

### Core Components
- Stellar Network Integration
- Business Logic Services
- RESTful API Layer

## API Endpoints

### Profile Management
```
POST   /api/v1/profiles/providers          # Create service provider profile
PUT    /api/v1/profiles/providers/:id      # Update provider profile
POST   /api/v1/profiles/providers/:id/ports      # Add port
POST   /api/v1/profiles/providers/:id/airports   # Add airport
POST   /api/v1/profiles/providers/:id/terminals  # Add terminal

POST   /api/v1/profiles/buyers             # Create service buyer profile
PUT    /api/v1/profiles/buyers/:id         # Update buyer profile
```

### Service Management
```
# Import Services
POST   /api/v1/services/import            # Create import service
GET    /api/v1/services/import            # List import services
GET    /api/v1/services/import/subcategories    # List subcategories
GET    /api/v1/services/import/items/:subcategory  # List service items

# Export Services
POST   /api/v1/services/export            # Create export service
GET    /api/v1/services/export            # List export services
GET    /api/v1/services/export/subcategories    # List subcategories
GET    /api/v1/services/export/items/:subcategory  # List service items

# Transit Services
POST   /api/v1/services/transit           # Create transit service
GET    /api/v1/services/transit           # List transit services
GET    /api/v1/services/transit/subcategories   # List subcategories
GET    /api/v1/services/transit/items/:subcategory # List service items

# Transshipment Services
POST   /api/v1/services/transshipment     # Create transshipment service
GET    /api/v1/services/transshipment     # List transshipment services
GET    /api/v1/services/transshipment/subcategories  # List subcategories
GET    /api/v1/services/transshipment/items/:subcategory  # List service items
```

### Booking and Tracking
```
POST   /api/v1/bookings                   # Create booking
POST   /api/v1/bookings/:id/payment       # Process payment
POST   /api/v1/bookings/:id/status        # Update status

POST   /api/v1/tracking/events            # Add tracking event
GET    /api/v1/tracking/:booking_id       # Get tracking history
POST   /api/v1/tracking/route/optimal     # Get optimal route
```

## Prerequisites

- Go 1.21 or higher
- Stellar Network account
- Soroban CLI
- Environment variables setup

## Environment Variables

```bash
export PORT=8080
export ISSUER_KEY=<your-stellar-issuer-key>
export CONTRACT_ID=<deployed-contract-id>
```

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd logistics-marketplace
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
go build -o logistics-marketplace cmd/api/main.go
```

## Running the Application

1. Start the API server:
```bash
./logistics-marketplace
```

The server will start on port 8080 (or the port specified in environment variables).

## Token Economics

- Token Name: Logistics Marketplace Token (LMT)
- Total Supply: 100,000,000,000 tokens
- Use Cases:
  - Service payments
  - Booking deposits
  - Customs duty payments
  - Platform fees

## Security

- JWT-based authentication
- Blockchain-based transaction security
- Smart contract security measures
- Rate limiting and CORS protection

## Development

### Testing
```bash
go test ./...
```

### Local Development
1. Use Stellar testnet for development
2. Deploy contracts using Soroban CLI
3. Set up environment variables
4. Run the application

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License

## Support

For support, please open an issue in the repository or contact the development team.
