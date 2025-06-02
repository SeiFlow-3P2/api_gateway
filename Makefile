env:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Создан файл .env из .env.example"; \
	else \
		echo "Файл .env уже существует"; \
	fi

api-gateway-build:
	docker compose --env-file .env build
	@echo "Сборка завершена"

api-gateway-up:
	docker compose --env-file .env up -d
	@echo "Сервис запущен"

api-gateway-down:
	docker compose --env-file .env down
	@echo "Сервис остановлен"