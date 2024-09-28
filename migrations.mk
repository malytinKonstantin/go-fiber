include .env
export $(shell sed 's/=.*//' .env)

MIGRATION_DIR=./db/migration
SCHEMA_DIR=./db/schema

# Генерация и выполнение миграции
migrate-auto:
	@echo "Генерация миграции..."
	@SCHEMA_FILE=$(SCHEMA_FILE) python scripts/generate_migration.py
	@echo "Выполнение миграций..."
	@migrate -database "$(DB_URL)" -path $(MIGRATION_DIR) up

# Остальные команды остаются без изменений...

migrate-up:
	@echo "Выполнение миграций..."
	@migrate -database "$(DB_URL)" -path $(MIGRATION_DIR) up

migrate-down:
	@echo "Откат последней миграции..."
	@migrate -database "$(DB_URL)" -path $(MIGRATION_DIR) down 1

migrate-create:
	@read -p "Введите название миграции: " name; \
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $${name}

migrate-status:
	@echo "Статус миграций:"
	@migrate -database "$(DB_URL)" -path $(MIGRATION_DIR) version

migrate-force:
	@read -p "Введите версию для принудительной установки: " version; \
	migrate -database "$(DB_URL)" -path $(MIGRATION_DIR) force $$version

help:
	@echo "Доступные команды:"
	@echo "  make migrate-auto   - Автоматически сгенерировать и выполнить миграцию"
	@echo "  make migrate-up     - Выполнить все доступные миграции"
	@echo "  make migrate-down   - Откатить последнюю миграцию"
	@echo "  make migrate-create - Создать новую миграцию"
	@echo "  make migrate-status - Проверить статус миграций"
	@echo "  make migrate-force  - Принудительно установить версию миграции"