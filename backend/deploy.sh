#!/bin/bash

# Start script for mess registration backend
set -e

echo "🚀 Starting Mess Registration Backend..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Create .env if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "⚠️  Please update .env file with your configuration"
fi

# Build and start containers
echo "🏗️  Building and starting containers..."
docker-compose down -v
docker-compose up --build -d

# Wait for database to be ready
echo "⏳ Waiting for database to be ready..."
sleep 10

# Check if services are running
if docker-compose ps | grep -q "Up"; then
    echo "✅ Services started successfully!"
    echo ""
    echo "🌐 Application is running at: http://localhost:8080"
    echo "📊 Health check: http://localhost:8080/health"
    echo ""
    echo "📝 API Documentation: http://localhost:8080/docs"
    echo "🧪 Test the API using test_api.http file"
    echo ""
    echo "To view logs: docker-compose logs -f"
    echo "To stop: docker-compose down"
else
    echo "❌ Failed to start services"
    docker-compose logs
    exit 1
fi
