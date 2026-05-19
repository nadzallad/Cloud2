-- =========================
-- CREATE DATABASES
-- =========================
CREATE DATABASE order_db;
CREATE DATABASE payment_db;
CREATE DATABASE pickup_db;
CREATE DATABASE warehouse_db;
CREATE DATABASE shipment_db;
CREATE DATABASE tracking_db;
CREATE DATABASE delivery_db;
CREATE DATABASE notification_db;

-- =========================
-- ORDER SERVICE
-- =========================
USE order_db;

CREATE TABLE orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,

    sender_name VARCHAR(100),
    sender_address TEXT,

    receiver_name VARCHAR(100),
    receiver_address TEXT,

    weight_kg DECIMAL(10,2),
    distance_km DECIMAL(10,2),

    base_price DECIMAL(10,2),
    shipping_cost DECIMAL(10,2),
    total_price DECIMAL(10,2),

    status VARCHAR(50) DEFAULT 'CREATED',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- PAYMENT SERVICE
-- =========================
USE payment_db;

CREATE TABLE payments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT NOT NULL,

    amount DECIMAL(10,2),
    payment_method VARCHAR(50),
    status VARCHAR(50) DEFAULT 'UNPAID',

    paid_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- PICKUP SERVICE
-- =========================
USE pickup_db;

CREATE TABLE pickups (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT NOT NULL,

    courier_id BIGINT,
    pickup_address TEXT,
    scheduled_time DATETIME,

    status VARCHAR(50) DEFAULT 'SCHEDULED'
);

-- =========================
-- WAREHOUSE SERVICE
-- =========================
USE warehouse_db;

CREATE TABLE warehouse_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT NOT NULL,

    warehouse_location VARCHAR(100),
    status VARCHAR(50), 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- SHIPMENT SERVICE
-- =========================
USE shipment_db;

CREATE TABLE shipments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT NOT NULL,
    courier_id BIGINT,

    tracking_number VARCHAR(100) UNIQUE NOT NULL,

    status VARCHAR(50) DEFAULT 'IN_TRANSIT',
    shipped_at TIMESTAMP NULL
);

-- =========================
-- TRACKING SERVICE
-- =========================
USE tracking_db;

CREATE TABLE tracking_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,

    shipment_id BIGINT NOT NULL,
    tracking_number VARCHAR(100) NOT NULL,

    status VARCHAR(100),
    location VARCHAR(100),
    note TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- DELIVERY SERVICE
-- =========================
USE delivery_db;

CREATE TABLE deliveries (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    shipment_id BIGINT NOT NULL,

    receiver_name VARCHAR(100),

    status VARCHAR(50) DEFAULT 'PENDING',
    delivered_at TIMESTAMP NULL
);

-- =========================
-- NOTIFICATION SERVICE
-- =========================
USE notification_db;

CREATE TABLE notifications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT,

    type VARCHAR(50),
    message TEXT,

    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);