.PHONY: build run docker-up docker-down producer clean help

# Переменные
BINARY_NAME=order-service
PRODUCER_BINARY=producer

# Сборка основного приложения
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) cmd/main.go

# Сборка producer
build-producer:
	@echo "Building $(PRODUCER_BINARY)..."
	@go build -o bin/$(PRODUCER_BINARY) scripts/producer.go

# Запуск основного приложения
run: build
	@echo "Running $(BINARY_NAME)..."
	@./bin/$(BINARY_NAME)

# Запуск producer
run-producer: build-producer
	@echo "Running $(PRODUCER_BINARY)..."
	@./bin/$(PRODUCER_BINARY)

# Запуск Docker Compose
docker-up:
	@echo "Starting Docker services..."
	@docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 10

# Остановка Docker Compose
docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down

# Полная остановка с удалением volumes
docker-clean:
	@echo "Stopping Docker services and removing volumes..."
	@docker-compose down -v
	@docker system prune -f

# Инициализация проекта
init: docker-up
	@echo "Initializing project..."
	@sleep 15
	@echo "Services should be ready!"

# Загрузка зависимостей
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Тестирование
test:
	@echo "Running tests..."
	@go test -v ./...

# Очистка бинарных файлов
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

# Создание директории bin
bin:
	@mkdir -p bin

# Полная сборка
build-all: bin build build-producer

# Демонстрация работы сервиса
demo: docker-up
	@echo "Starting demo..."
	@sleep 15
	@echo "Building and starting service..."
	@make build
	@./bin/$(BINARY_NAME) &
	@SERVICE_PID=$$!
	@sleep 5
	@echo "Starting producer to send test orders..."
	@make run-producer
	@echo ""
	@echo "Demo completed!"
	@echo "Open http://localhost:8081 in your browser"
	@echo "Try searching for order: b563feb7b2b84b6test"

# Справка
help:
	@echo "Available commands:"
	@echo "  build         - Build the main application"
	@echo "  build-producer - Build the producer application"
	@echo "  build-all     - Build both applications"
	@echo "  run           - Build and run the main application"
	@echo "  run-producer  - Build and run the producer"
	@echo "  docker-up     - Start Docker services"
	@echo "  docker-down   - Stop Docker services"
	@echo "  docker-clean  - Stop services and remove volumes"
	@echo "  init          - Initialize project (start Docker)"
	@echo "  deps          - Download Go dependencies"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean binary files"
	@echo "  demo          - Full demo (Docker + Service + Producer)"
	@echo "  help          - Show this help"
