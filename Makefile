export PORT ?= 3100
export DATABASE_URL ?= postgresql://postgres:postgres@localhost:5433/ecomm?sslmode=disable

## db up + migrate + run
dev:
	docker compose up -d --wait
	atlas migrate apply --env gorm
	go run .

## after changing a model: make migration name=add_foo
migration:
	atlas migrate diff $(name) --env gorm

.PHONY: dev migration
