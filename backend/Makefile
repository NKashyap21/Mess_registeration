# Project variables
APP_NAME = mess_registration
MAIN_FILE = main.go
DOCKER_COMPOSE = docker-compose.yml

# Run app locally without Docker
run:
	go run $(MAIN_FILE)

# Build Docker image for API
build:
	docker build -t $(APP_NAME)-api .

# Start services (Postgres + API)
up:
	docker compose -f $(DOCKER_COMPOSE) up -d

# Stop services
down:
	docker compose -f $(DOCKER_COMPOSE) down

# Restart API service only
restart-api:
	docker compose -f $(DOCKER_COMPOSE) up -d --build api

# Show logs
logs:
	docker compose -f $(DOCKER_COMPOSE) logs -f

# Remove volumes + containers
clean:
	docker compose -f $(DOCKER_COMPOSE) down -v --remove-orphans
