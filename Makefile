.PHONY: help dev dev-local up down install setup env db db-down backend frontend seed \
	fuzz fuzz-deps railway-login railway-deploy railway-api railway-web railway-postgres

COMPOSE := docker compose
BACKEND_DIR := backend
FRONTEND_DIR := frontend
FUZZ_DIR := fuzzing
FUZZ_VENV := $(FUZZ_DIR)/.venv
FUZZ_PY := $(FUZZ_VENV)/bin/python
FUZZ_PIP := $(FUZZ_VENV)/bin/pip
RAILWAY := npx @railway/cli
RAILWAY_API_URL ?= https://api-production-3cf0.up.railway.app/api/v1
RAILWAY_WEB_URL ?= https://web-production-196fd.up.railway.app

help:
	@echo "Targets:"
	@echo "  make up         - full stack in Docker (db + api + web)"
	@echo "  make dev        - alias for make up"
	@echo "  make dev-local  - postgres in Docker, API + Vite on host"
	@echo "  make install    - go modules + npm dependencies"
	@echo "  make setup      - install + .env files"
	@echo "  make db         - only PostgreSQL in Docker"
	@echo "  make seed       - demo organizer + published hackathon"
	@echo "  make down       - stop all compose services"
	@echo "  make fuzz       - API fuzz (Schemathesis + business scenarios)"
	@echo "  make fuzz-deps  - install Python deps for fuzz (fuzzing/.venv)"
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
	$(COMPOSE) up -d db --wait

db-down:
	$(COMPOSE) stop db

down:
	$(COMPOSE) down

up: env
	$(COMPOSE) up --build -d --wait
	@echo ""
	@echo "  Web:  http://127.0.0.1:$${WEB_PORT:-5173}"
	@echo "  API:  http://127.0.0.1:$${API_PORT:-8080}/api/v1/health"
	@echo "  Demo: admin@admin.ru / admin, user@user.ru / user"
	@echo ""

seed: setup db
	@set -a; [ -f .env ] && . ./.env; set +a; \
	cd $(BACKEND_DIR) && go run ./cmd/seed

dev: up

fuzz-deps:
	@test -d $(FUZZ_VENV) || python3 -m venv $(FUZZ_VENV)
	@$(FUZZ_PIP) install -q -r $(FUZZ_DIR)/requirements.txt

fuzz: env fuzz-deps
	@$(FUZZ_PY) $(FUZZ_DIR)/fuzz.py

dev-local: setup db
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
