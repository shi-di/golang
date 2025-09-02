package cache

import (
	"log"
	"order-service/internal/models"
	"sync"
)

// Cache представляет in-memory кеш для заказов
type Cache struct {
	mu     sync.RWMutex
	orders map[string]*models.Order
}

// New создает новый экземпляр кеша
func New() *Cache {
	return &Cache{
		orders: make(map[string]*models.Order),
	}
}

// Set сохраняет заказ в кеше
func (c *Cache) Set(orderUID string, order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[orderUID] = order
}

// Get получает заказ из кеша
func (c *Cache) Get(orderUID string) (*models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.orders[orderUID]
	return order, exists
}

// LoadFromDB восстанавливает кеш из базы данных
func (c *Cache) LoadFromDB(orders []*models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.orders = make(map[string]*models.Order, len(orders))

	for _, order := range orders {
		c.orders[order.OrderUID] = order
	}

	log.Printf("Cache loaded with %d orders", len(orders))
}

// Size возвращает количество заказов в кеше
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.orders)
}

// GetAll возвращает все заказы из кеша (для отладки)
func (c *Cache) GetAll() map[string]*models.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]*models.Order, len(c.orders))
	for k, v := range c.orders {
		result[k] = v
	}

	return result
}
