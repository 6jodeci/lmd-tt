package main

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestReserveProducts(t *testing.T) {
	// Open a test database connection
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=mydb sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Insert some test data
	_, err = db.Exec("INSERT INTO products (name, size, code, quantity) VALUES ('Product 1', 'Large', 'PRD1', 10)")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO products (name, size, code, quantity) VALUES ('Product 2', 'Small', 'PRD2', 5)")
	if err != nil {
		t.Fatal(err)
	}

	// Reserve some products
	err = ReserveProducts(db, []string{"PRD1", "PRD2"})
	if err != nil {
		t.Fatal(err)
	}

	// Check that the product quantities have been updated
	var p1Quantity int
	err = db.QueryRow("SELECT quantity FROM products WHERE code = 'PRD1'").Scan(&p1Quantity)
	if err != nil {
		t.Fatal(err)
	}
	if p1Quantity != 9 {
		t.Errorf("Expected product quantity to be 9, but got %d", p1Quantity)
	}

	var p2Quantity int
	err = db.QueryRow("SELECT quantity FROM products WHERE code = 'PRD2'").Scan(&p2Quantity)
	if err != nil {
		t.Fatal(err)
	}
	if p2Quantity != 4 {
		t.Errorf("Expected product quantity to be 4, but got %d", p2Quantity)
	}
}

func TestReleaseProducts(t *testing.T) {
	// Open a test database connection
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=mydb sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Insert some test data
	_, err = db.Exec("INSERT INTO products (name, size, code, quantity) VALUES ('Product 1', 'Large', 'PRD1', 10)")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO products (name, size, code, quantity) VALUES ('Product 2', 'Small', 'PRD2', 5)")
	if err != nil {
		t.Fatal(err)
	}

	// Reserve some products
	err = ReserveProducts(db, []string{"PRD1", "PRD2"})
	if err != nil {
		t.Fatal(err)
	}

	// Release some products
	err = ReleaseProducts(db, []string{"PRD1"})
	if err != nil {
		t.Fatal(err)
	}

	// Check that the product quantities have been updated
	var p1Quantity int
	err = db.QueryRow("SELECT quantity FROM products WHERE code = 'PRD1'").Scan(&p1Quantity)
	if err != nil {
		t.Fatal(err)
	}
	if p1Quantity != 10 {
		t.Errorf("Expected product quantity to be 10, but got %d", p1Quantity)
	}
}
