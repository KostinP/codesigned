# =====================================================================
#                          GLOBAL CONFIG
# =====================================================================

ENV_FILE?=.env
DEV_ENV_FILE?=.env.dev
STAGE_ENV_FILE?=.env.stage
PROD_ENV_FILE?=.env.prod

# Проверяем существование .env файла
ifneq (,$(wildcard $(ENV_FILE)))
    include $(ENV_FILE)
    export
else
    $(warning File $(ENV_FILE) not found. Run 'make init-env' first)
endif

# Используем strip для удаления пробелов в переменных
DB_HOST_STRIP := $(strip $(DB_HOST))
DB_USER_STRIP := $(strip $(DB_USER))
DB_PASSWORD_STRIP := $(strip $(DB_PASSWORD))
DB_NAME_STRIP := $(strip $(DB_NAME))
DB_PORT_STRIP := $(strip $(DB_PORT))

DB_URL?=postgres://$(DB_USER_STRIP):$(DB_PASSWORD_STRIP)@$(DB_HOST_STRIP):$(DB_PORT_STRIP)/$(DB_NAME_STRIP)?sslmode=disable
SWAGGER_DIRS?=./backend/cmd,./backend/internal/user/transport/http,./backend/internal/user/entity
MIGRATIONS_DIR?=./backend/migrations/postgres

GREEN := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
BLUE := $(shell tput -Txterm setaf 4)
RED := $(shell tput -Txterm setaf 1)
NC := $(shell tput -Txterm sgr0)

ANALYTICS_PROFILE?=--profile analytics

.DEFAULT_GOAL := help

.PHONY: help run test build tidy swagger wire migrate-up migrate-down migrate-create migrate-status docker-build deploy setup-vps upload-certs setup-webhook vps-logs vps-info

# =====================================================================
#                   UNIVERSAL ENVIRONMENT LOADER
# =====================================================================

define load_env
	@if [ -z "$(1)" ]; then \
		echo "$(RED)Usage: make <command> env=<dev|stage|prod>$(NC)"; \
		exit 1; \
	fi
	@if [ ! -f .env.$(1) ]; then \
		echo "$(RED)File .env.$(1) does not exist$(NC)"; \
		exit 1; \
	fi
	@cp -f .env.$(1) .env || true
	@echo "$(BLUE)Using environment: $(1)$(NC)"
	@APP_ENV=$(1) $(2)
endef

# =====================================================================
#                      ENVIRONMENT INITIALIZATION
# =====================================================================

init-env: ## Initialize environment files
	@echo "$(BLUE)Initializing environment files...$(NC)"
	@mkdir -p backend/migrations/postgres backend/migrations/clickhouse
	@if [ ! -f .env.dev ]; then \
		cp .env.example .env.dev 2>/dev/null || touch .env.dev; \
		echo "$(GREEN)Created .env.dev$(NC)"; \
	else \
		echo "$(YELLOW).env.dev already exists$(NC)"; \
	fi
	@if [ ! -f .env.stage ]; then \
		cp .env.example .env.stage 2>/dev/null || touch .env.stage; \
		echo "$(GREEN)Created .env.stage$(NC)"; \
	else \
		echo "$(YELLOW).env.stage already exists$(NC)"; \
	fi
	@if [ ! -f .env.prod ]; then \
		cp .env.example .env.prod 2>/dev/null || touch .env.prod; \
		echo "$(GREEN)Created .env.prod$(NC)"; \
	else \
		echo "$(YELLOW).env.prod already exists$(NC)"; \
	fi
	@if [ ! -f .env ]; then \
		cp -f .env.dev .env; \
		echo "$(GREEN)Created .env from .env.dev$(NC)"; \
	else \
		echo "$(YELLOW).env already exists$(NC)"; \
	fi
	@echo "$(GREEN)Environment files initialized!$(NC)"
	@echo "$(YELLOW)Please configure the .env.* files before proceeding.$(NC)"

env-switch: ## Switch environment (usage: make env-switch env=dev)
	$(call load_env,$(env),true)

env: ## Show current environment
	@if [ -f .env ]; then \
		echo "$(GREEN)Current environment: $$(basename .env)$(NC)"; \
	else \
		echo "$(RED)No .env file found$(NC)"; \
	fi

check-env: ## Check required env variables
	$(call load_env,$(env),true)
	@if [ -z "$(DB_HOST_STRIP)" ] || [ -z "$(DB_PORT_STRIP)" ] || [ -z "$(DB_NAME_STRIP)" ] || [ -z "$(DB_USER_STRIP)" ] || [ -z "$(DB_PASSWORD_STRIP)" ] || [ -z "$(TELEGRAM_BOT_TOKEN)" ] || [ -z "$(DOMAIN)" ] || [ -z "$(REPOLINK)" ] || [ -z "$(SUPERSET_SECRET_KEY)" ]; then \
		echo "$(RED)Missing required env variables! Check .env.$(env)$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Env check passed for $(env)$(NC)"

# =====================================================================
#                             COMMANDS
# =====================================================================

run: ## Run server with Air (usage: make run env=dev)
	$(call load_env,$(env),cd backend && air)

test: ## Run tests (usage: make test env=dev)
	$(call load_env,$(env),cd backend && go test ./... -v)

build: ## Build app (usage: make build env=dev)
	$(call load_env,$(env),cd backend && go build -o bin/edu-platform ./cmd)

tidy: ## Tidy Go modules
	cd backend && go mod tidy

swagger: ## Generate Swagger docs
	cd backend && swag init --dir $(SWAGGER_DIRS) --output ./docs

wire: ## Generate wire dependencies
	cd backend && wire ./cmd

migrate-up: ## Apply DB migrations (usage: make migrate-up env=dev)
	$(call load_env,$(env),migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up)

migrate-down: ## Rollback DB migrations
	$(call load_env,$(env),migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down)

migrate-create: ## Create new migration file (name=required)
	$(call load_env,$(env),migrate create -ext sql -dir $(MIGRATIONS_DIR) $(name))

migrate-status: ## Show migration status
	$(call load_env,$(env),migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version)

docker-build: ## Build Docker image (usage: make docker-build env=dev)
	$(call load_env,$(env),docker build -t edu-platform:$(env) .)

deploy: ## Deploy to VPS (usage: make deploy env=stage)
	$(call load_env,$(env),)
	@echo "$(BLUE)Loading environment from .env.$(env)$(NC)"
	@./scripts/deploy_vps.sh

setup-vps: ## Setup VPS
	$(call load_env,$(env),./scripts/setup_vps.sh)

upload-certs: ## Upload certificates to VPS
	$(call load_env,$(env),./scripts/upload_certificates.sh)

setup-webhook: ## Configure Telegram webhook
	$(call load_env,$(env),./scripts/configure_telegram_webhook.sh)

vps-logs: ## Show VPS logs
	$(call load_env,$(env),ssh $(VPS_USER)@$(VPS_IP) "docker compose logs -f")

vps-info: ## VPS system info
	$(call load_env,$(env),ssh $(VPS_USER)@$(VPS_IP) "uname -a && df -h && free -h")

# =====================================================================
#                          DOCKER COMMANDS
# =====================================================================

docker-up: ## Start with docker-compose (usage: make docker-up env=dev)
	$(call load_env,$(env),\
		if [ "${ENABLE_ANALYTICS}" = "true" ]; then \
			docker-compose $(ANALYTICS_PROFILE) up --build; \
		else \
			docker-compose up --build; \
		fi)

docker-up-with-analytics: ## Start with analytics services (usage: make docker-up-with-analytics env=dev)
	$(call load_env,$(env),docker-compose --profile analytics up --build)

docker-up-core: ## Start only core services without analytics
	$(call load_env,$(env),docker-compose up --build postgres app)

docker-down: ## Stop docker-compose
	docker-compose down

docker-clean: ## Stop and remove containers, networks, and volumes
	docker-compose down -v --remove-orphans

docker-logs: ## Show docker logs
	docker-compose logs -f

docker-restart: ## Restart services (usage: make docker-restart env=dev)
	$(call load_env,$(env),docker-compose restart)

docker-ps: ## Show running containers
	docker-compose ps

# Обновите существующую команду docker-build чтобы использовала env
docker-build: ## Build Docker image (usage: make docker-build env=dev)
	$(call load_env,$(env),docker-compose build)

# Команда для полной очистки и перезапуска с инициализацией Superset
docker-fresh: ## Fresh start with clean volumes and Superset initialization (usage: make docker-fresh env=dev)
	@echo "$(YELLOW)Stopping and removing all containers and volumes...$(NC)"
	@docker-compose down -v --remove-orphans
	@echo "$(BLUE)Using environment: $(env)$(NC)"
	@cp -f .env.$(env) .env
	@APP_ENV=$(env) docker-compose --profile analytics up --build -d
	@echo "$(BLUE)Waiting for services to start...$(NC)"
	@sleep 30
	@if [ "$(env)" = "dev" ] || [ "$(env)" = "stage" ]; then \
		echo "$(BLUE)Running Superset initialization for $(env)...$(NC)"; \
		$(MAKE) superset-init env=$(env); \
	fi
	@echo "$(GREEN)All services are ready!$(NC)"

# =====================================================================
#                         SUPERSET MANAGEMENT
# =====================================================================

superset-init: ## Initialize Superset (usage: make superset-init env=dev)
	$(call load_env,$(env),)
	@echo "$(BLUE)Initializing Superset...$(NC)"
	chmod +x scripts/init_superset.sh
	./scripts/init_superset.sh $(env)

superset-init-with-examples: ## Initialize Superset with examples
	$(call load_env,$(env),)
	@echo "$(BLUE)Initializing Superset with examples...$(NC)"
	chmod +x scripts/init_superset.sh
	./scripts/init_superset.sh $(env) admin admin123 admin@edu-platform.com true

superset-create-admin: ## Create Superset admin user
	$(call load_env,$(env),)
	@echo "$(BLUE)Creating Superset admin user...$(NC)"
	docker-compose exec superset superset fab create-admin \
		--username admin \
		--firstname Admin \
		--lastname User \
		--email admin@edu-platform.com \
		--password admin123

superset-db-upgrade: ## Upgrade Superset database
	$(call load_env,$(env),)
	@echo "$(BLUE)Upgrading Superset database...$(NC)"
	docker-compose exec superset superset db upgrade

superset-setup-clickhouse: ## Setup ClickHouse connection
	$(call load_env,$(env),)
	@echo "$(BLUE)Setting up ClickHouse connection...$(NC)"
	chmod +x scripts/init_superset.sh
	./scripts/init_superset.sh $(env)

superset-full-setup: superset-init superset-setup-clickhouse ## Full Superset setup
	@echo "$(GREEN)Superset full setup completed!$(NC)"

superset-logs: ## Show Superset logs
	docker-compose logs superset -f

superset-bash: ## Access Superset container
	docker-compose exec superset bash

superset-restart: ## Restart Superset
	docker-compose restart superset

superset-clean: ## Clean Superset data
	@echo "$(YELLOW)Cleaning Superset data...$(NC)"
	docker-compose down -v superset
	docker volume rm $$(docker volume ls -q | grep superset) 2>/dev/null || true


# =====================================================================
#                      QUICK START PRODUCTION
# =====================================================================

prod-full-deploy: ## Complete production deployment with Superset
	@echo "$(BLUE)Starting complete production deployment...$(NC)"
	@$(MAKE) setup-vps env=prod
	@$(MAKE) upload-certs env=prod
	@$(MAKE) deploy env=prod
	@$(MAKE) setup-webhook env=prod
	@$(MAKE) superset-init env=prod
	@echo "$(GREEN)🎉 Production deployment completed!$(NC)"

# =====================================================================
#                           HELP GENERATOR
# =====================================================================

help:
	@echo ""
	@echo "$(GREEN)==================== AVAILABLE COMMANDS ====================$(NC)"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z0-9_-]+:.*?## / {printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}' Makefile
	@echo ""