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


build-backend:
	docker build -t muqeeth26832/mess-registration-backend:latest ./backend
build-frontend:
	docker build -t muqeeth26832/mess-registration-frontend:latest ./frontend

push:
	docker push muqeeth26832/mess-registration-backend:latest
	docker push muqeeth26832/mess-registration-frontend:latest

# Optional combo command
build:
	make build-backend && make build-frontend && make push

create-context:
	docker context create mess-registration --docker "host=ssh://lambda@cabsharing.iith.dev"

use-context:
	docker context use mess-registration

deploy:
	docker stack deploy -c docker-compose.prod.yml mess-registration

.PHONY: up-prod down-prod logs-prod restart-prod

up-prod:
	docker compose -f docker-compose.prod.yml up --build -d

down-prod:
	docker compose -f docker-compose.prod.yml down

logs-prod:
	docker compose -f docker-compose.prod.yml logs -f

restart-prod:
	docker compose -f docker-compose.prod.yml down
	docker compose -f docker-compose.prod.yml up --build -d
