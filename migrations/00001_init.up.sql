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
CREATE TABLE products (
  id serial PRIMARY KEY, 
  name text, 
  size text, 
  code text UNIQUE, 
  quantity integer NOT NULL, 
  warehouse_id integer NOT NULL REFERENCES warehouse (id)
);
CREATE TABLE warehouse (
  id serial PRIMARY KEY, 
  name text, 
  is_available boolean
);

-- СОЗДАНИЕ ИНДЕКСОВ --
CREATE INDEX idx_products_code ON products (code);
CREATE INDEX idx_products_quantity ON products (quantity);
CREATE INDEX idx_warehouse_name ON warehouse (name);
CREATE UNIQUE INDEX idx_products_warehouse_code ON products (warehouse_id, code);