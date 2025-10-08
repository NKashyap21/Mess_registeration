.PHONY: up down logs

# Up dev environment
up:
	docker compose -f docker-compose.dev.yml up --build

# Down dev environment
down:
	docker compose -f docker-compose.dev.yml down

# View logs for all services
logs:
	docker compose -f docker-compose.dev.yml logs -f


.PHONY: up-prod down-prod logs-prod

up-prod:
	docker compose -f docker-compose.prod.yml up --build -d

down-prod:
	docker compose -f docker-compose.prod.yml down

logs-prod:
	docker compose -f docker-compose.prod.yml logs -f
