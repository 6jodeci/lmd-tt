BEGIN;

-- CONFIG -- 
SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = ON;
SET check_function_bodies = FALSE;
SET client_min_messages = WARNING;
SET search_path = public, extensions;
SET default_tablespace = '';
SET default_with_oids = FALSE;

-- TABLES --
CREATE TABLE IF NOT EXISTS warehouses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    is_available BOOLEAN DEFAULT TRUE,
    CONSTRAINT unique_name UNIQUE (name) -- уникальный индекс на поле name
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    size VARCHAR(50),
    unique_code VARCHAR(100) NOT NULL,
    quantity INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS warehouse_products (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER NOT NULL REFERENCES warehouses(id),
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    CONSTRAINT unique_warehouse_product UNIQUE (warehouse_id, product_id), -- уникальный индекс на поля warehouse_id и product_id
    CONSTRAINT positive_quantity CHECK (quantity >= 0) -- констрэйнт на поле quantity, чтобы не было отрицательных значений
);

CREATE INDEX idx_warehouse_id ON warehouse_products (warehouse_id); -- индекс на поле warehouse_id таблицы warehouse_products
CREATE INDEX idx_product_id ON warehouse_products (product_id); -- индекс на поле product_id таблицы warehouse_products
COMMIT;