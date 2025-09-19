# Project variables
COMPOSE = docker compose
BACKEND = mess_registration_backend
DB = mess_registration_db

.PHONY: up down restart logs backend-logs db-logs ps build

# Start all services (backend + db)
up:
	$(COMPOSE) up -d --build

# Stop and remove containers
down:
	$(COMPOSE) down

# Restart services
restart: down up

# Build images without running
build:
	$(COMPOSE) build

# Logs
logs:
	$(COMPOSE) logs -f

backend-logs:
	$(COMPOSE) logs -f $(BACKEND)

db-logs:
	$(COMPOSE) logs -f $(DB)

# Run only database
db:
	$(COMPOSE) up -d db

# Run only backend
backend:
	$(COMPOSE) up -d backend

# Open a psql shell into DB
psql:
	docker exec -it $(DB) psql -U postgres -d mess_registration
