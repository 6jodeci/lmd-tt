-- КОНФИГУРАЦИЯ -- 
SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = ON;
SET check_function_bodies = FALSE;
SET client_min_messages = WARNING;
SET search_path = public, extensions;
SET default_tablespace = '';
SET default_with_oids = FALSE;

-- ТАБЛИЦЫ --
CREATE TABLE warehouse (
  id SERIAL PRIMARY KEY, 
  name TEXT, 
  is_available BOOLEAN
);

CREATE TABLE products (
  id SERIAL PRIMARY KEY, 
  name TEXT, 
  size TEXT, 
  code TEXT UNIQUE, 
  quantity INTEGER NOT NULL, 
  warehouse_id INTEGER REFERENCES warehouse(id)
);

-- СОЗДАНИЕ ИНДЕКСОВ --
CREATE INDEX idx_products_code ON products (code);
CREATE INDEX idx_products_quantity ON products (quantity);
CREATE INDEX idx_warehouse_name ON warehouse (name);
