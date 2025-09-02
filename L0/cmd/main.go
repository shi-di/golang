package main

import (
	"context"
	"log"
	"net/http"
	"order-service/internal/cache"
	"order-service/internal/database"
	"order-service/internal/handlers"
	"order-service/internal/kafka"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load("config.env"); err != nil {
		log.Printf("level=warn component=bootstrap event=env_load_warning err=%v", err)
	}

	// Получение переменных окружения
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "orders_db")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	kafkaBroker := getEnv("KAFKA_BROKER", "localhost:9092")
	kafkaTopic := getEnv("KAFKA_TOPIC", "orders")

	httpPort := getEnv("HTTP_PORT", "8081")

	log.Printf("level=info component=bootstrap event=config db_host=%q db_port=%q db_name=%q kafka_broker=%q topic=%q http_port=%q", dbHost, dbPort, dbName, kafkaBroker, kafkaTopic, httpPort)

	// Подключение к базе данных
	db, err := database.New(dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)
	if err != nil {
		log.Fatalf("level=fatal component=bootstrap event=db_connect_failed err=%v", err)
	}
	defer db.Close()

	// Создание кеша
	orderCache := cache.New()

	// Восстановление кеша из базы данных
	log.Println("level=info component=bootstrap event=cache_load msg=\"loading cache from database\"")
	orders, err := db.GetAllOrders()
	if err != nil {
		log.Printf("level=warn component=bootstrap event=cache_load_failed err=%v", err)
	} else {
		orderCache.LoadFromDB(orders)
	}

	// Создание HTTP handlers
	orderHandler := handlers.NewOrderHandler(db, orderCache)

	// Настройка роутера
	router := mux.NewRouter()

	// Health endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// API endpoints
	router.HandleFunc("/order/{order_uid}", orderHandler.GetOrder).Methods("GET")
	router.HandleFunc("/cache/stats", orderHandler.GetCacheStats).Methods("GET")

	// Статические файлы
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/static/")))

	// Создание HTTP сервера
	server := &http.Server{
		Addr:         ":" + httpPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Создание Kafka consumer
	consumer := kafka.NewConsumer(kafkaBroker, kafkaTopic, db, orderCache)

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запуск Kafka consumer в отдельной горутине
	go consumer.Start(ctx)

	// Запуск HTTP сервера в отдельной горутине
	go func() {
		log.Printf("level=info component=http event=start port=%s", httpPort)
		log.Printf("level=info component=http event=endpoints web=\"http://localhost:%s\" api_example=\"http://localhost:%s/order/<order_uid>\"", httpPort, httpPort)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("level=fatal component=http event=server_failed err=%v", err)
		}
	}()

	// Обработка сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Ожидание сигнала завершения
	<-sigChan
	log.Println("level=info component=bootstrap event=shutdown msg=\"shutting down gracefully\"")

	// Остановка контекста для Kafka consumer
	cancel()

	// Остановка HTTP сервера
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("level=error component=http event=shutdown_error err=%v", err)
	}

	// Закрытие Kafka consumer
	if err := consumer.Close(); err != nil {
		log.Printf("level=error component=kafka_consumer event=close_error err=%v", err)
	}

	log.Println("level=info component=bootstrap event=stopped msg=\"server stopped\"")
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
