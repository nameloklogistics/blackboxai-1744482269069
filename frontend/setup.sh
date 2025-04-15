#!/bin/bash

# Make script exit on first error
set -e

echo "🚀 Setting up frontend development environment..."

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js first."
    exit 1
fi

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install npm first."
    exit 1
fi

echo "📦 Installing dependencies..."
npm install

# Create necessary directories if they don't exist
echo "📁 Creating necessary directories..."
mkdir -p src/assets
mkdir -p src/hooks
mkdir -p src/pages/Auth
mkdir -p src/store/middleware

# Set up environment variables
echo "🔧 Setting up environment variables..."
if [ ! -f .env ]; then
    cat > .env << EOL
VITE_API_URL=http://localhost:8080/api
VITE_STELLAR_NETWORK=TESTNET
VITE_TOKEN_ISSUER=your_token_issuer_key_here
EOL
    echo "✅ Created .env file with default values"
fi

# Set up git hooks
echo "🔨 Setting up git hooks..."
if [ -d .git ]; then
    # Create pre-commit hook
    cat > .git/hooks/pre-commit << EOL
#!/bin/sh
npm run lint
npm run test
EOL
    chmod +x .git/hooks/pre-commit
    echo "✅ Created pre-commit hook"
fi

# Run type checking
echo "🔍 Running type checking..."
npm run tsc --noEmit

# Run linting
echo "🧹 Running linting..."
npm run lint

# Run tests
echo "🧪 Running tests..."
npm run test

echo "✨ Setup complete! You can now start the development server with:"
echo "npm run dev"

# Optional: Start the development server
read -p "Would you like to start the development server now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🚀 Starting development server..."
    npm run dev
fi
