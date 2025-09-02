package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"order-service/internal/cache"
	"order-service/internal/database"
	"order-service/internal/models"

	"github.com/gorilla/mux"
)

// OrderHandler обрабатывает HTTP запросы для заказов
type OrderHandler struct {
	db    *database.DB
	cache *cache.Cache
}

// NewOrderHandler создает новый handler для заказов
func NewOrderHandler(db *database.DB, cache *cache.Cache) *OrderHandler {
	return &OrderHandler{
		db:    db,
		cache: cache,
	}
}

// GetOrder возвращает заказ по UID
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["order_uid"]

	if orderUID == "" {
		http.Error(w, "Order UID is required", http.StatusBadRequest)
		return
	}

	// Сначала ищем в кеше
	order, found := h.cache.Get(orderUID)
	if found {
		log.Printf("level=info component=http_handler route=get_order source=cache event=hit order_uid=%q", orderUID)
		h.writeJSONResponse(w, order)
		return
	}

	// Если в кеше нет, ищем в базе данных
	log.Printf("level=info component=http_handler route=get_order source=cache event=miss order_uid=%q", orderUID)
	order, err := h.db.GetOrder(orderUID)
	if err != nil {
		if err == models.ErrOrderNotFound {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		log.Printf("level=error component=http_handler route=get_order event=db_error order_uid=%q err=%v", orderUID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Добавляем в кеш для следующих запросов
	h.cache.Set(orderUID, order)
	log.Printf("level=info component=http_handler route=get_order source=db event=cached order_uid=%q", orderUID)

	h.writeJSONResponse(w, order)
}

// GetCacheStats возвращает статистику кеша (для отладки)
func (h *OrderHandler) GetCacheStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"cache_size": h.cache.Size(),
		"orders":     []string{},
	}

	// Получаем список всех заказов в кеше
	allOrders := h.cache.GetAll()
	orderUIDs := make([]string, 0, len(allOrders))
	for uid := range allOrders {
		orderUIDs = append(orderUIDs, uid)
	}
	stats["orders"] = orderUIDs

	h.writeJSONResponse(w, stats)
}

// writeJSONResponse записывает JSON ответ
func (h *OrderHandler) writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("level=error component=http_handler route=get_order event=json_encode_error err=%v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
