# Переменные
BINARY_NAME=go-fiber
MAIN_PACKAGE=.

# Команды Go
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GORUN=$(GOCMD) run

# Линтеры и инструменты
GOLINT=golangci-lint
GOIMPORTS=goimports
WIRE=$(HOME)/go/bin/wire
SWAG=$(HOME)/go/bin/swag

.PHONY: all build clean test coverage run lint fmt mod-tidy wire swagger help

all: mod-tidy wire swagger lint test build

build: wire swagger
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PACKAGE)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

test:
	$(GOTEST) ./...

coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

run: wire swagger
	$(GORUN) $(MAIN_PACKAGE)

lint:
	$(GOLINT) run

fmt:
	$(GOIMPORTS) -w .

mod-tidy:
	$(GOMOD) tidy

wire:
	$(WIRE)

swagger:
	$(SWAG) init

help:
	@echo "Доступные команды:"
	@echo "  make build      - Собрать приложение (включая wire и swagger)"
	@echo "  make clean      - Очистить артефакты сборки"
	@echo "  make test       - Запустить тесты"
	@echo "  make coverage   - Запустить тесты с покрытием"
	@echo "  make run        - Запустить приложение (включая wire и swagger)"
	@echo "  make lint       - Запустить линтер"
	@echo "  make fmt        - Отформатировать код"
	@echo "  make mod-tidy   - Обновить зависимости"
	@echo "  make wire       - Сгенерировать код с помощью Wire"
	@echo "  make swagger    - Сгенерировать документацию Swagger"
	@echo "  make all        - Выполнить mod-tidy, wire, swagger, lint, test и build"
	@echo "  make help       - Показать эту справку"