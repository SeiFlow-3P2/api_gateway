env:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env created from .env.example"; \
	else \
		echo ".env already exists"; \
	fi

api-gateway-build:
	docker compose --env-file .env build
	@echo "Build completed"

api-gateway-up:
	docker compose --env-file .env up -d
	@echo "Service started"

api-gateway-down:
	docker compose --env-file .env down
	@echo "Service stopped"