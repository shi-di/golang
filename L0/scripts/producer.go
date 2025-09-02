package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Структуры для тестовых данных (упрощенные версии из models)
type TestOrder struct {
	OrderUID          string       `json:"order_uid"`
	TrackNumber       string       `json:"track_number"`
	Entry             string       `json:"entry"`
	Delivery          TestDelivery `json:"delivery"`
	Payment           TestPayment  `json:"payment"`
	Items             []TestItem   `json:"items"`
	Locale            string       `json:"locale"`
	InternalSignature string       `json:"internal_signature"`
	CustomerID        string       `json:"customer_id"`
	DeliveryService   string       `json:"delivery_service"`
	Shardkey          string       `json:"shardkey"`
	SmID              int          `json:"sm_id"`
	DateCreated       string       `json:"date_created"`
	OofShard          string       `json:"oof_shard"`
}

type TestDelivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type TestPayment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type TestItem struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {
	// Настройка Kafka writer
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	log.Println("Starting Kafka producer...")

	// Отправка тестовых заказов
	testOrders := generateTestOrders()

	for i, order := range testOrders {
		// Конвертация в JSON
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("Error marshaling order %d: %v", i+1, err)
			continue
		}

		// Отправка сообщения
		message := kafka.Message{
			Key:   []byte(order.OrderUID),
			Value: orderJSON,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err = writer.WriteMessages(ctx, message)
		cancel()

		if err != nil {
			log.Printf("Error sending order %s: %v", order.OrderUID, err)
		} else {
			log.Printf("Successfully sent order: %s", order.OrderUID)
		}

		// Пауза между отправками
		time.Sleep(2 * time.Second)
	}

	// Отправка невалидного сообщения для тестирования обработки ошибок
	log.Println("Sending invalid message for error handling test...")
	invalidMessage := kafka.Message{
		Key:   []byte("invalid"),
		Value: []byte(`{"order_uid": "invalid_test_uid", "track_number": "", "items": [{"chrt_id": 999, "track_number": "TEST", "price": 100, "rid": "test_rid", "name": "Test Item", "sale": 0, "size": "M", "total_price": 100, "nm_id": 999999, "brand": "Test Brand", "status": 200}]}`),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err := writer.WriteMessages(ctx, invalidMessage)
	cancel()

	if err != nil {
		log.Printf("Error sending invalid message: %v", err)
	} else {
		log.Println("Successfully sent invalid message")
	}

	log.Println("Producer finished")
}

func generateTestOrders() []TestOrder {
	now := time.Now().Format(time.RFC3339)

	return []TestOrder{
		{
			OrderUID:    "b563feb7b2b84b6test",
			TrackNumber: "WBILMTESTTRACK",
			Entry:       "WBIL",
			Delivery: TestDelivery{
				Name:    "Test Testov",
				Phone:   "+9720000000",
				Zip:     "2639809",
				City:    "Kiryat Mozkin",
				Address: "Ploshad Mira 15",
				Region:  "Kraiot",
				Email:   "test@gmail.com",
			},
			Payment: TestPayment{
				Transaction:  "b563feb7b2b84b6test",
				RequestID:    "",
				Currency:     "USD",
				Provider:     "wbpay",
				Amount:       1817,
				PaymentDt:    1637907727,
				Bank:         "alpha",
				DeliveryCost: 1500,
				GoodsTotal:   317,
				CustomFee:    0,
			},
			Items: []TestItem{
				{
					ChrtID:      9934930,
					TrackNumber: "WBILMTESTTRACK",
					Price:       453,
					Rid:         "ab4219087a764ae0btest",
					Name:        "Mascaras",
					Sale:        30,
					Size:        "0",
					TotalPrice:  317,
					NmID:        2389212,
					Brand:       "Vivienne Sabo",
					Status:      202,
				},
			},
			Locale:            "en",
			InternalSignature: "",
			CustomerID:        "test",
			DeliveryService:   "meest",
			Shardkey:          "9",
			SmID:              99,
			DateCreated:       now,
			OofShard:          "1",
		},
		{
			OrderUID:    "a123def4c5d67e8test",
			TrackNumber: "WBILMTESTTRACK2",
			Entry:       "WBIL",
			Delivery: TestDelivery{
				Name:    "Ivan Petrov",
				Phone:   "+79161234567",
				Zip:     "123456",
				City:    "Moscow",
				Address: "Tverskaya 1",
				Region:  "Moscow",
				Email:   "ivan@example.com",
			},
			Payment: TestPayment{
				Transaction:  "a123def4c5d67e8test",
				RequestID:    "req_001",
				Currency:     "RUB",
				Provider:     "sberbank",
				Amount:       2500,
				PaymentDt:    time.Now().Unix(),
				Bank:         "sberbank",
				DeliveryCost: 300,
				GoodsTotal:   2200,
				CustomFee:    0,
			},
			Items: []TestItem{
				{
					ChrtID:      1234567,
					TrackNumber: "WBILMTESTTRACK2",
					Price:       1200,
					Rid:         "rid_test_001",
					Name:        "T-Shirt",
					Sale:        10,
					Size:        "M",
					TotalPrice:  1080,
					NmID:        1111111,
					Brand:       "Nike",
					Status:      200,
				},
				{
					ChrtID:      2345678,
					TrackNumber: "WBILMTESTTRACK2",
					Price:       1300,
					Rid:         "rid_test_002",
					Name:        "Jeans",
					Sale:        15,
					Size:        "32",
					TotalPrice:  1105,
					NmID:        2222222,
					Brand:       "Levi's",
					Status:      200,
				},
			},
			Locale:            "ru",
			InternalSignature: "sig_001",
			CustomerID:        "customer_123",
			DeliveryService:   "cdek",
			Shardkey:          "5",
			SmID:              55,
			DateCreated:       now,
			OofShard:          "2",
		},
		{
			OrderUID:    "c789ghi0j1k23l4test",
			TrackNumber: "WBILMTESTTRACK3",
			Entry:       "WBIL",
			Delivery: TestDelivery{
				Name:    "Elena Sidorova",
				Phone:   "+79267654321",
				Zip:     "190000",
				City:    "Saint Petersburg",
				Address: "Nevsky Prospekt 20",
				Region:  "Leningrad Oblast",
				Email:   "elena@example.com",
			},
			Payment: TestPayment{
				Transaction:  "c789ghi0j1k23l4test",
				RequestID:    "req_002",
				Currency:     "RUB",
				Provider:     "yandex_money",
				Amount:       3200,
				PaymentDt:    time.Now().Unix(),
				Bank:         "vtb",
				DeliveryCost: 400,
				GoodsTotal:   2800,
				CustomFee:    50,
			},
			Items: []TestItem{
				{
					ChrtID:      3456789,
					TrackNumber: "WBILMTESTTRACK3",
					Price:       2800,
					Rid:         "rid_test_003",
					Name:        "Winter Jacket",
					Sale:        20,
					Size:        "L",
					TotalPrice:  2240,
					NmID:        3333333,
					Brand:       "The North Face",
					Status:      201,
				},
			},
			Locale:            "ru",
			InternalSignature: "sig_002",
			CustomerID:        "customer_456",
			DeliveryService:   "boxberry",
			Shardkey:          "7",
			SmID:              77,
			DateCreated:       now,
			OofShard:          "3",
		},
	}
}
