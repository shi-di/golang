package kafka

import (
	"context"
	"encoding/json"
	"log"
	"order-service/internal/cache"
	"order-service/internal/database"
	"order-service/internal/models"
	"time"

	"github.com/segmentio/kafka-go"
)

// Consumer представляет Kafka consumer
type Consumer struct {
	reader *kafka.Reader
	db     *database.DB
	cache  *cache.Cache
}

// NewConsumer создает новый Kafka consumer
func NewConsumer(broker, topic string, db *database.DB, cache *cache.Cache) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{broker},
		Topic:          topic,
		GroupID:        "order-service-group",
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	})

	return &Consumer{
		reader: r,
		db:     db,
		cache:  cache,
	}
}

// Start запускает потребление сообщений из Kafka
func (c *Consumer) Start(ctx context.Context) {
	log.Println("component=kafka_consumer event=start msg=\"starting consumer\"")

	for {
		select {
		case <-ctx.Done():
			log.Println("component=kafka_consumer event=stop msg=\"stopping consumer\"")
			return
		default:
			// Чтение сообщения из Kafka
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("level=error component=kafka_consumer event=fetch_error err=%v", err)
				continue
			}

			// Обработка сообщения
			if err := c.processMessage(msg); err != nil {
				log.Printf("level=error component=kafka_consumer event=process_error partition=%d offset=%d err=%v", msg.Partition, msg.Offset, err)
				// В случае ошибки не коммитим сообщение
				continue
			}

			// Подтверждение успешной обработки
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("level=error component=kafka_consumer event=commit_error partition=%d offset=%d err=%v", msg.Partition, msg.Offset, err)
			}
		}
	}
}

// processMessage обрабатывает одно сообщение
func (c *Consumer) processMessage(msg kafka.Message) error {
	log.Printf("component=kafka_consumer event=process_start partition=%d offset=%d key=%q", msg.Partition, msg.Offset, string(msg.Key))
	log.Printf("component=kafka_consumer event=message_content partition=%d offset=%d body=%q", msg.Partition, msg.Offset, string(msg.Value))

	// Парсинг JSON
	var order models.Order
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		log.Printf("level=warn component=kafka_consumer event=invalid_json partition=%d offset=%d err=%v", msg.Partition, msg.Offset, err)
		return nil // Игнорируем невалидные сообщения
	}

	// Валидация заказа
	if err := order.Validate(); err != nil {
		log.Printf("level=warn component=kafka_consumer event=invalid_order partition=%d offset=%d order_uid=%q err=%v", msg.Partition, msg.Offset, order.OrderUID, err)
		return nil // Игнорируем невалидные заказы
	}

	// Сохранение в базу данных
	if err := c.db.SaveOrder(&order); err != nil {
		log.Printf("level=error component=kafka_consumer event=db_save_failed partition=%d offset=%d order_uid=%q err=%v", msg.Partition, msg.Offset, order.OrderUID, err)
		return err // Возвращаем ошибку для повторной обработки
	}

	// Сохранение в кеш
	c.cache.Set(order.OrderUID, &order)

	log.Printf("level=info component=kafka_consumer event=processed partition=%d offset=%d order_uid=%q", msg.Partition, msg.Offset, order.OrderUID)
	return nil
}

// Close закрывает consumer
func (c *Consumer) Close() error {
	return c.reader.Close()
}
