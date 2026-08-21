.PHONY: test test-backend test-frontend build smoke smoke-docker e2e-docker integration-smoke-docker docker-up docker-down docker-prod docker-integrations backup-db worker check-ops

test: test-backend test-frontend

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && npm run build

build: test-backend test-frontend

smoke:
	cd frontend && npm run test:smoke

smoke-docker: docker-up
	./scripts/wait-for-api.sh
	cd frontend && npm run test:smoke

check-ops: docker-up
	./scripts/wait-for-api.sh
	./scripts/check-ops.sh

e2e-docker: docker-up
	./scripts/wait-for-api.sh
	./scripts/wait-for-frontend.sh
	cd frontend && PLAYWRIGHT_SKIP_WEBSERVER=1 npm run test:e2e

integration-smoke-docker: docker-integrations
	./scripts/wait-for-api.sh
	cd frontend && npm run test:integration

docker-up:
	docker compose up --build -d

docker-integrations:
	docker compose -f docker-compose.yml -f docker-compose.integrations.yml up --build -d

docker-prod:
	docker compose --env-file .env.production -f docker-compose.yml -f docker-compose.prod.yml up -d --build

docker-down:
	docker compose down

backup-db:
	chmod +x scripts/backup-postgres.sh
	./scripts/backup-postgres.sh

worker:
	cd backend && go run ./cmd/worker
