package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type ProductCode struct {
	Id int64 `json:"id"`
}

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

func NewProductRepository(database *sql.DB) *ProductRepository {
	return &ProductRepository{database: database}
}

type ProductRepository struct {
	database *sql.DB
}

func ConnectDatabase(config *Config) (*sql.DB, error) {
	return sql.Open("sqlite3", config.DatabasePath)
}

func (repository *ProductRepository) createProductsTable() {
	createProductsTableSQL := `CREATE TABLE products (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"description" TEXT,
		"price" TEXT
	  );`

	statement, err := repository.database.Prepare(createProductsTableSQL)
	defer statement.Close()
	if err != nil {
		log.Fatalf("Failed to create table: %v", err.Error())
	}
	statement.Exec()
}

func (repository *ProductRepository) insertProduct(product Product) error {
	tx, _ := repository.database.Begin()
	stmt, _ := tx.Prepare(`INSERT INTO products (name, description, price) VALUES (?, ?, ?);`)
	defer stmt.Close()
	_, err := stmt.Exec(product.Name, product.Description, product.Price)

	if err != nil {
		tx.Rollback()
		log.Fatalf("Statement error %v", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("DB Commit error %v", err.Error())
	}

	return err
}

func (repository *ProductRepository) updateProduct(product Product) error {
	tx, _ := repository.database.Begin()
	stmt, _ := tx.Prepare(`UPDATE products SET name = ?, description = ?, price = ? WHERE id = ?;`)
	defer stmt.Close()
	_, err := stmt.Exec(product.Name, product.Description, product.Price, product.Id) //.Scan(&id, &name, &description, &price)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Statement error %v", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("DB Commit error %v", err.Error())
	}
	/*
		getProductSql := `SELECT id, name, description, price FROM products WHERE id = $1`
		err = repository.database.QueryRow(getProductSql, product.Id).Scan(&id, &name, &description, &price)

		if err != nil {
			log.Fatal(err.Error())
		}

		updatedProduct := Product{Id: id, Name: name, Description: description, Price: price}
	*/
	return err
}

func (repository *ProductRepository) deleteProduct(code ProductCode) error {
	tx, _ := repository.database.Begin()
	stmt, _ := tx.Prepare(`DELETE FROM products WHERE id = ?`)
	defer stmt.Close()

	_, err := stmt.Exec(code.Id)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Statement error %v", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("DB Commit error %v", err.Error())
	}
	return err
}

func (repository *ProductRepository) FindAll() []*Product {
	rows, _ := repository.database.Query(`SELECT id, name, description, price FROM products;`)
	defer rows.Close()

	products := []*Product{}

	for rows.Next() {
		var (
			id          int
			name        string
			description string
			price       string
		)

		rows.Scan(&id, &name, &description, &price)

		products = append(products, &Product{
			Id:          id,
			Name:        name,
			Description: description,
			Price:       price,
		})
	}

	return products
}
