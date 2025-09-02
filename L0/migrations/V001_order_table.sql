-- Создание таблицы заказов
CREATE TABLE IF NOT EXISTS orders (
	order_uid VARCHAR(255) PRIMARY KEY,
	track_number VARCHAR(255) NOT NULL,
	entry VARCHAR(255),
	locale VARCHAR(10),
	internal_signature VARCHAR(255),
	customer_id VARCHAR(255),
	delivery_service VARCHAR(255),
	shardkey VARCHAR(10),
	sm_id INTEGER,
	date_created TIMESTAMP,
	oof_shard VARCHAR(10),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы доставки
CREATE TABLE IF NOT EXISTS deliveries (
	id SERIAL PRIMARY KEY,
	order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
	name VARCHAR(255),
	phone VARCHAR(50),
	zip VARCHAR(20),
	city VARCHAR(255),
	address TEXT,
	region VARCHAR(255),
	email VARCHAR(255),
	CONSTRAINT uq_deliveries_order_uid UNIQUE (order_uid)
);

-- Создание таблицы платежей
CREATE TABLE IF NOT EXISTS payments (
	id SERIAL PRIMARY KEY,
	order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
	transaction VARCHAR(255),
	request_id VARCHAR(255),
	currency VARCHAR(10),
	provider VARCHAR(255),
	amount INTEGER,
	payment_dt BIGINT,
	bank VARCHAR(255),
	delivery_cost INTEGER,
	goods_total INTEGER,
	custom_fee INTEGER,
	CONSTRAINT uq_payments_order_uid UNIQUE (order_uid)
);

-- Создание таблицы товаров
CREATE TABLE IF NOT EXISTS items (
	id SERIAL PRIMARY KEY,
	order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
	chrt_id INTEGER,
	track_number VARCHAR(255),
	price INTEGER,
	rid VARCHAR(255),
	name VARCHAR(255),
	sale INTEGER,
	size VARCHAR(50),
	total_price INTEGER,
	nm_id INTEGER,
	brand VARCHAR(255),
	status INTEGER
);

-- Создание индексов для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_orders_track_number ON orders(track_number);
CREATE INDEX IF NOT EXISTS idx_orders_date_created ON orders(date_created);
CREATE INDEX IF NOT EXISTS idx_deliveries_order_uid ON deliveries(order_uid);
CREATE INDEX IF NOT EXISTS idx_payments_order_uid ON payments(order_uid);
CREATE INDEX IF NOT EXISTS idx_items_order_uid ON items(order_uid);