#!/bin/bash

# Deployment script for logistics marketplace

set -e

echo "Starting deployment..."

# Build backend
echo "Building backend..."
go build -o logistics-api cmd/api/main.go

# Build frontend
echo "Building frontend..."
cd frontend
npm install
npm run build
cd ..

# Set environment variables (example, replace with your actual values)
export JWT_SECRET="your_jwt_secret_here"
export ISSUER_KEY="your_issuer_key_here"
export CONTRACT_ID="your_contract_id_here"
export GOVERNANCE_CONTRACT_ID="your_governance_contract_id_here"
export PORT=8080

# Run backend in background
echo "Starting backend server..."
nohup ./logistics-api > backend.log 2>&1 &

# Serve frontend (using simple HTTP server for demo, replace with your production server)
echo "Serving frontend on port 3000..."
cd frontend/dist
nohup python3 -m http.server 3000 > ../frontend.log 2>&1 &

echo "Deployment complete."
echo "Backend logs: backend.log"
echo "Frontend logs: frontend/frontend.log"
echo "Backend running on port $PORT"
echo "Frontend running on port 3000"
