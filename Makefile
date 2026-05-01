.PHONY: order-migrate-up order-migrate-down payment-migrate-up payment-migrate-down

include .env
export

PROJECT_ROOT := $(CURDIR)

run-all:
	go run $(PROJECT_ROOT)/order-service/cmd/main.go & \
	go run $(PROJECT_ROOT)/payment-service/cmd/main.go & \
	go run $(PROJECT_ROOT)/notification-service/cmd/main.go

env-up:
	@docker compose up -d order-db
	@docker compose up -d payment-db
	@docker compose up -d rabbitmq

order-migrate-create:
ifndef name
	@echo Error: name is undefined.
	@echo Usage: make migrate-create name=init
	@exit 1
endif
	docker compose run --rm order-db-migrate create -ext sql -dir /migrations -seq $(name)

payment-migrate-create:
ifndef name
	@echo Error: name is undefined.
	@echo Usage: make migrate-create name=init
	@exit 1
endif
	docker compose run --rm payment-db-migrate create -ext sql -dir /migrations -seq $(name)

order-migrate-up:
	@make order-migrate-action action=up

order-migrate-down:
	@make order-migrate-action action=down

order-migrate-action:
	@docker compose run --rm order-db-migrate \
		-path=/migrations/ \
		-database "postgres://${ORDER_POSTGRES_USER}:${ORDER_POSTGRES_PASSWORD}@order-db:5432/${ORDER_POSTGRES_DB}?sslmode=disable" \
		$(action)

payment-migrate-up:
	@make payment-migrate-action action=up

payment-migrate-down:
	@make payment-migrate-action action=down

payment-migrate-action:
	@docker compose run --rm payment-db-migrate \
		-path=/migrations/ \
		-database "postgres://${PAYMENT_POSTGRES_USER}:${PAYMENT_POSTGRES_PASSWORD}@payment-db:5432/${PAYMENT_POSTGRES_DB}?sslmode=disable" \
		$(action)