.PHONY: help up down restart logs ps build rebuild pull clean test

DC = docker compose

help:
	@echo "Available targets:"
	@echo "  make up       - Start all services in background"
	@echo "  make down     - Stop and remove containers"
	@echo "  make restart  - Restart the full stack"
	@echo "  make logs     - Follow logs for all services"
	@echo "  make ps       - Show running containers"
	@echo "  make build    - Build all images"
	@echo "  make rebuild  - Rebuild images and start stack"
	@echo "  make pull     - Pull latest base images"
	@echo "  make clean    - Stop stack and remove volumes"
	@echo "  make test     - Run go tests in all services"

up:
	$(DC) up -d --build

down:
	$(DC) down

restart: down up

logs:
	$(DC) logs -f

ps:
	$(DC) ps

build:
	$(DC) build

rebuild:
	$(DC) up -d --build --force-recreate

pull:
	$(DC) pull

clean:
	$(DC) down -v --remove-orphans

test:
	cd order-service && go test ./...
	cd payment-service && go test ./...
	cd notification-service && go test ./...
