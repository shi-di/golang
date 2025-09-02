# Order Service - Демонстрационный сервис с Kafka и PostgreSQL

Микросервис на Go для обработки заказов с использованием Kafka, PostgreSQL и in-memory кеширования.

## Особенности

- **Kafka Consumer** - получение заказов из очереди сообщений
- **PostgreSQL** - надежное хранение данных с транзакциями
- **In-Memory Cache** - быстрый доступ к данным
- **REST API** - получение заказов по ID
- **Web Interface** - простой интерфейс для поиска заказов
- **Graceful Shutdown** - корректное завершение работы
- **Error Handling** - обработка невалидных сообщений

## Быстрый старт

### 1. Запуск инфраструктуры

```bash
# Запуск PostgreSQL и Kafka
make docker-up

# Или
docker-compose up -d
```

### 2. Установка зависимостей

```bash
make deps
```

### 3. Запуск сервиса

```bash
# Сборка и запуск
make run

# Или вручную
go build -o bin/order-service cmd/main.go
./bin/order-service
```

### 4. Отправка тестовых данных

```bash
# В отдельном терминале
make run-producer
```

### 5. Тестирование

Откройте в браузере: http://localhost:8081

Попробуйте найти заказ: `b563feb7b2b84b6test`

## API Endpoints

- `GET /order/{order_uid}` - получить заказ по UID
- `GET /cache/stats` - статистика кеша
- `GET /` - веб-интерфейс

### Пример запроса

```bash
curl http://localhost:8081/order/b563feb7b2b84b6test
```

## Структура проекта

```
├── cmd/
│   └── main.go              # Точка входа
├── internal/
│   ├── models/              # Модели данных
│   ├── database/            # Работа с PostgreSQL
│   ├── kafka/               # Kafka consumer
│   ├── cache/               # In-memory кеш
│   └── handlers/            # HTTP handlers
├── web/static/              # Веб-интерфейс
├── scripts/                 # Утилиты
├── migrations/              # SQL миграции
├── docker-compose.yml       # Docker окружение
└── Makefile                # Команды сборки
```

## Конфигурация

Настройки в файле `config.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=orders_db

KAFKA_BROKER=localhost:9092
KAFKA_TOPIC=orders

HTTP_PORT=8081
```

## Команды Make

```bash
make help                 # Справка по командам
make build               # Сборка приложения
make run                 # Запуск сервиса
make docker-up           # Запуск Docker
make docker-down         # Остановка Docker
make run-producer        # Отправка тестовых сообщений
make demo                # Полная демонстрация
```

## Схема базы данных

### Таблица `orders`
- `order_uid` (PK) - уникальный идентификатор заказа
- `track_number` - номер отслеживания
- `entry` - точка входа
- `locale`, `customer_id`, `delivery_service` и др.

### Таблица `deliveries`
- `order_uid` (FK) - связь с заказом
- `name`, `phone`, `email` - данные получателя
- `address`, `city`, `zip` - адрес доставки

### Таблица `payments`
- `order_uid` (FK) - связь с заказом
- `transaction`, `amount`, `currency` - данные платежа
- `provider`, `bank` - платежная информация

### Таблица `items`
- `order_uid` (FK) - связь с заказом
- `name`, `brand`, `price` - информация о товаре
- `sale`, `total_price` - цены и скидки

## Обработка ошибок

- Валидация входящих JSON сообщений
- Игнорирование невалидных заказов
- Транзакции для целостности данных
- Подтверждение сообщений Kafka
- Graceful shutdown при ошибках

## Производительность

- **Кеш** ускоряет повторные запросы
- **Пул соединений** к базе данных
- **Асинхронная** обработка Kafka
- **Индексы** в PostgreSQL

## Мониторинг

- Логирование всех операций
- Статистика кеша через API
- Health checks для Docker

## Остановка сервисов

```bash
# Остановка приложения: Ctrl+C

# Остановка Docker
make docker-down

# Полная очистка
make docker-clean
```
