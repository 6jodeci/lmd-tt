package controller

// @title Warehouse API Documentation
// @description This is a sample API for a warehouse application
// @version 1
// @host localhost:8080
// @BasePath /api/v1

import (
	"database/sql"
	"errors"
)

// Product структура продукта
type Product struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Size     string `json:"size" db:"size"`
	Code     string `json:"code" db:"code"`
	Quantity int    `json:"quantity" db:"quantity"`
}

// Warehouse структура склада
type Warehouse struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	IsAvailable bool   `json:"is_available" db:"is_available"`
}

// ReserveProducts reserves products
// @Summary Reserves products
// @Description Reserves products and updates their quantities
// @Tags products
// @Accept json
// @Produce json
// @Param productCodes query []string true "Product codes"
// @Success 204 {string} string ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/reserve [post]
// ReserveProducts резервирует продукты
func ReserveProducts(db *sql.DB, productCodes []string) error {
	if len(productCodes) == 0 {
		return errors.New("empty product codes")
	}

	// Начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Подготавливаем SQL-запрос для обновления количества продуктов
	updateStmt, err := tx.Prepare("UPDATE products SET quantity = quantity - 1 WHERE code = $1")
	if err != nil {
		tx.Rollback()
		return err
	}
	// Закрываем подготовленный запрос при выходе из функции
	defer updateStmt.Close()

	// Зарезервируем каждый продукт в цикле
	for _, code := range productCodes {
		// Проверяем, существует ли продукт
		var p Product
		err := db.QueryRow("SELECT id, name, size, code, quantity FROM products WHERE code = $1", code).Scan(&p.ID, &p.Name, &p.Size, &p.Code, &p.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Проверяем, доступен ли продукт для бронирования
		if p.Quantity < 1 {
			tx.Rollback()
			return errors.New("product is out of stock")
		}

		// Обновление количества продукта
		_, err = updateStmt.Exec(p.Code)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// ReleaseProducts releases products
// @Summary Releases products
// @Description Releases reserved products and updates their quantities
// @Tags products
// @Accept json
// @Produce json
// @Param productCodes query []string true "Product codes"
// @Success 204 {string} string ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/release [post]
// ReleaseProducts отменяет резервирование товаров.
func ReleaseProducts(db *sql.DB, productCodes []string) error {
	// Проверяем массив на пустоту массива кодов
	if len(productCodes) == 0 {
		return errors.New("empty product codes")
	}

	// Начинаем новую транзакцию
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Создаем новую транзакцию
	updateStmt, err := tx.Prepare("UPDATE products SET quantity = quantity + 1 WHERE code = $1")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer updateStmt.Close()

	// Проходимся по каждому продукту
	for _, code := range productCodes {
		// Проверяем существует ли продукт
		var p Product
		err := db.QueryRow("SELECT id, name, size, code, quantity FROM products WHERE code = $1", code).Scan(&p.ID, &p.Name, &p.Size, &p.Code, &p.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Обновляем количество продуктов
		_, err = updateStmt.Exec(p.Code)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// GetRemainingProducts returns remaining products
// @Summary Returns remaining products
// @Description Returns the remaining products in the warehouse
// @Tags products
// @Accept json
// @Produce json
// @Param warehouseID query int true "Warehouse ID"
// @Success 200 {array} Product
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/remaining [get]
// GetRemainingProducts возвращает оставшееся количество продуктов на складе
func GetRemainingProducts(db *sql.DB, warehouseID int) ([]Product, error) {
	rows, err := db.Query("SELECT code, quantity FROM products WHERE warehouse_id=$1", warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Создаем пустой слайс для хранения результатов
	var products []Product
	// Проходимся по строкам, возвращенным запросом, и добавляем каждую строку к слайсу продуктов.
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Code, &p.Quantity); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
