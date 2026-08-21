GOOSE ?= goose
DATABASE_URL ?=

.PHONY: migrate-up migrate-down migrate-status seed-dev

migrate-up:
	$(GOOSE) -dir ./migrations postgres "$(DATABASE_URL)" up

migrate-down:
	$(GOOSE) -dir ./migrations postgres "$(DATABASE_URL)" down

migrate-status:
	$(GOOSE) -dir ./migrations postgres "$(DATABASE_URL)" status

seed-dev:
	psql "$(DATABASE_URL)" -f ./seeds/dev.sql
