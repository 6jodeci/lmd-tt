package controller

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestCreateWarehouse(t *testing.T) {
	// Коннектимся к базе
	db, err := sql.Open("postgres", "host=localhost port=5432 user=root password=secret dbname=lamoda_db sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создам новый склад
	w := &Warehouse{
		Name:        "Test Warehouse",
		IsAvailable: true,
	}

	// Вызываем функцию создания нового склада
	err = CreateWarehouse(db, w)
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем успешно ли был создан склад
	if w.ID == 0 {
		t.Errorf("Expected warehouse ID to be non-zero, got %d", w.ID)
	}
}

func TestCreateProduct(t *testing.T) {
	// Подключаемся к базе
	db, err := sql.Open("postgres", "host=localhost port=5432 user=root password=secret dbname=lamoda_db sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем новый склад (снова т.к мы не знаем айдишник и название склада заренее)
	w := &Warehouse{
		Name:        "Test Warehouse",
		IsAvailable: true,
	}

	// Вызываем функцию для создания склада
	err = CreateWarehouse(db, w)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем продукт
	p := &Product{
		Name:        "Test Product",
		Size:        5,
		Code:        "TP0111",
		Quantity:    10,
		WarehouseID: w.ID,
	}

	// Вызываем функцию создания продукта
	err = CreateProduct(db, p)
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем успешно ли создан продукт
	if p.ID == 0 {
		t.Errorf("Expected product ID to be non-zero, got %d", p.ID)
	}
}
