package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Product struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Size     string `db:"size"`
	Code     string `db:"code"`
	Quantity int    `db:"quantity"`
}

type Warehouse struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	IsAvailable bool   `db:"is_available"`
}

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

func main() {
	// Подключение к базе
	db, err := sql.Open("postgres", "dbname=your-db-name user=your-db-user password=your-db-password sslmode=require")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Инициализируем Роутер
	router := gin.Default()

	// Роут для резервирования продуктов
	router.POST("/products/reserve", func(c *gin.Context) {
		var productCodes []string
		if err := c.ShouldBindJSON(&productCodes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ReserveProducts(db, productCodes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Products reserved successfully"})
	})

	// Роут для релизов продуктов
	router.POST("/products/release", func(c *gin.Context) {
		var productCodes []string
		if err := c.ShouldBindJSON(&productCodes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ReleaseProducts(db, productCodes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Products released successfully"})
	})

	// Роут оставшихся продуктов
	router.GET("/products/remaining/:warehouseID", func(c *gin.Context) {
		warehouseIDStr := c.Param("warehouseID")
		warehouseID, err := strconv.Atoi(warehouseIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID"})
			return
		}

		products, err := GetRemainingProducts(db, warehouseID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get remain product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"products": products})
	})

	// Запуск сервера
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
