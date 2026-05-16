.PHONY: help dev up down install setup env db db-down backend frontend seed \
	railway-login railway-deploy railway-api railway-web railway-postgres

COMPOSE := docker compose
BACKEND_DIR := backend
FRONTEND_DIR := frontend
RAILWAY := npx @railway/cli
RAILWAY_API_URL ?= https://api-production-3cf0.up.railway.app/api/v1
RAILWAY_WEB_URL ?= https://web-production-196fd.up.railway.app

help:
	@echo "Targets:"
	@echo "  make dev     - postgres + API (auto-seed) + frontend"
	@echo "  make install - go modules + npm dependencies"
	@echo "  make setup   - install + .env files"
	@echo "  make db      - only PostgreSQL in Docker"
	@echo "  make seed    - demo organizer + published hackathon"
	@echo "  make down    - stop Docker postgres"
	@echo "  make railway-login    - browser login to Railway"
	@echo "  make railway-deploy   - Postgres + API + web (production)"
	@echo "  make railway-api      - deploy API only"
	@echo "  make railway-web      - deploy frontend only"

setup: install env

install:
	cd $(BACKEND_DIR) && go mod download
	cd $(FRONTEND_DIR) && npm install

env:
	@test -f .env || cp .env.example .env
	@test -f $(FRONTEND_DIR)/.env || cp $(FRONTEND_DIR)/.env.example $(FRONTEND_DIR)/.env

db:
	$(COMPOSE) up -d --wait

db-down:
	$(COMPOSE) down

down: db-down

seed: setup db
	@set -a; [ -f .env ] && . ./.env; set +a; \
	cd $(BACKEND_DIR) && go run ./cmd/seed

dev: setup db
	@echo ""
	@echo "  API:  http://localhost:8080/api/v1/health"
	@echo "  Web:  http://localhost:5173"
	@echo ""
	@echo "Press Ctrl+C to stop API and frontend (postgres keeps running; use 'make down' to stop it)"
	@set -a; [ -f .env ] && . ./.env; set +a; \
	trap 'kill 0' EXIT INT TERM; \
	(cd $(BACKEND_DIR) && go run ./cmd/app) & \
	(cd $(FRONTEND_DIR) && npm run dev) & \
	wait

up: dev

backend: setup db
	@set -a; [ -f .env ] && . ./.env; set +a; \
	cd $(BACKEND_DIR) && go run ./cmd/app

frontend: setup
	@set -a; [ -f $(FRONTEND_DIR)/.env ] && . ./$(FRONTEND_DIR)/.env; set +a; \
	cd $(FRONTEND_DIR) && npm run dev

railway-login:
	$(RAILWAY) login

railway-postgres:
	$(RAILWAY) redeploy --service Postgres -y

railway-api:
	cd $(BACKEND_DIR) && $(RAILWAY) service link api
	cd $(BACKEND_DIR) && $(RAILWAY) up --detach -m "deploy api"

railway-web:
	$(RAILWAY) variable set VITE_API_URL="$(RAILWAY_API_URL)" --service web
	cd $(FRONTEND_DIR) && $(RAILWAY) service link web
	cd $(FRONTEND_DIR) && $(RAILWAY) up --detach -m "deploy web"

railway-deploy: railway-postgres
	$(RAILWAY) variable set CORS_ORIGIN="$(RAILWAY_WEB_URL)" --service api
	$(RAILWAY) variable set VITE_API_URL="$(RAILWAY_API_URL)" --service web
	@$(MAKE) railway-api
	@$(MAKE) railway-web
	@echo ""
	@echo "  API:  $(RAILWAY_API_URL)"
	@echo "  Web:  $(RAILWAY_WEB_URL)"
	@echo "  Wait ~2 min, then: curl $(RAILWAY_API_URL)/health"
