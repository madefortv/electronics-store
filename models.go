package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func (repository *ProductRepository) insertOffering(offering Offering) error {
	log.Printf("%d, %d, %s, %v", offering.ProductId, offering.DealId, offering.ModifedPrice, offering.Active)
	tx, _ := repository.database.Begin()
	stmt, _ := tx.Prepare(`INSERT INTO offerings (product_id, deal_id, modified_price, active) VALUES (?, ?, ?, ?);`)
	defer stmt.Close()
	_, err := stmt.Exec(offering.ProductId, offering.DealId, offering.ModifedPrice, offering.Active)
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

func (repository *ProductRepository) insertDeal(deal Deal) error {
	log.Printf("%s, %s, %s, %s, %d, %d, %v,", deal.Name, deal.Type, deal.Coupon, deal.Percent, deal.X, deal.Y, deal.Exclusive)
	tx, _ := repository.database.Begin()
	stmt, _ := tx.Prepare(`INSERT INTO deals (name, type, coupon, percent, x, y, exclusive) VALUES (?, ?, ?, ?, ?, ?, ?);`)
	defer stmt.Close()
	_, err := stmt.Exec(deal.Name, deal.Type, deal.Coupon, deal.Percent, deal.X, deal.Y, deal.Exclusive)
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

func (repository *ProductRepository) listDeals() []*Deal {
	rows, _ := repository.database.Query(`SELECT * FROM deals;`)
	defer rows.Close()

	deals := []*Deal{}

	for rows.Next() {
		var (
			id        int
			name      string
			btype     DealType
			coupon    string
			percent   string
			x         int
			y         int
			exclusive bool
		)

		rows.Scan(&id, &name, &btype, &coupon, &percent, &x, &y, &exclusive)

		deals = append(deals, &Deal{
			Id:        id,
			Name:      name,
			Type:      btype,
			Coupon:    coupon,
			Percent:   percent,
			X:         x,
			Y:         y,
			Exclusive: exclusive,
		})
	}

	return deals
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

func (repository *ProductRepository) listProducts() []*Product {
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

type DealType string

const (
	Retail       DealType = "Retail"
	Flat                  = "Flat"
	Percent               = "Percent"
	Bundle                = "Bundle"
	BuyXGetYFree          = "BuyXGetYFree"
	Other                 = "Other"
)

/*
   Helpful for signaling for delete or responses
   @id
*/

type ProductCode struct {
	Id int64 `json:"id"`
}

/*
   The Product model represents an item in stock or was in
   stock historically. Price is not the final price, but rather
   the set price at the time
   @Id the products primary key
   @Name product name
   @Description of the product
   @Price the retail price of an item
*/

type Product struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       string `json:"price"`
}

/*
   The deal struct is where we keep our data that drives the logic behind
   different types of offerings like bundles, or
   @Id is the deal id
   @Name refers to the bundle name
   @Type referes to the type of deal
   @Exclusive flag for whether this deal can work with other deals,
   @Coupon is a flat reduction in price from the msdrg price
   @X the first number of a Buy X Get Y Free modifier
   @Y the second number of a Buy X GEt Y Free modifier


   TODO: Add start and end timestamps
*/
type Deal struct {
	Id        int      `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Type      DealType `json:"type,omitempty"`
	Coupon    string   `json:"coupon,omitempty"`
	Percent   string   `json:"percent,omitempty"`
	X         int      `json:"x,omitempty"`
	Y         int      `json:"y,omitempty"`
	Exclusive bool     `json:"exclusive,omitempty"`
}

/* The offering model is a relationship between one or more products and
   zero or more deals. By default a product for sale is an offering with
   "retail" type deal modifying it.

   An Offering has Products and is modified by Deals. Users add Product Offerings
   to their cart, not Products

   ProductId, DealId -> primary key
   Each deal has a type and that type will influece the control flow for
   calculating the final price of a set of tiems.
   @Id is the primary key of this, although ProductId/DealId would work
   @ProductId is a product associated with this
   @DealId deal that modifies the product(s)
   @ModifiedPrice is the final price of an offering after being modified by
   a deal
   @Active flag determines whether this deal is active
*/

type Offering struct {
	Id           int    `json:"id,omitempty"`
	ProductId    int    `json:"product_id,omitempty"`
	DealId       int    `json:"deal_id,omitempty"`
	ModifedPrice string `json:"modified_price,omitempty"`
	Active       bool   `json:"active,omitempty"`
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
		"price" VARCHAR(8) NOT NULL,
		"name" VARCHAR(32) NOT NULL,
		"description" TEXT
	  );`

	statement, err := repository.database.Prepare(createProductsTableSQL)
	defer statement.Close()
	if err != nil {
		log.Fatalf("Failed to create table: %v", err.Error())
	}
	statement.Exec()
}

func (repository *ProductRepository) createDealsTable() {
	createDealsTableSQL := `CREATE TABLE deals (
	    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	    name VARCHAR(32) NOT NULL DEFAULT "Regular Price",
	    type VARCHAR(16) NOT NULL DEFAULT "Retail",
	    coupon VARCHAR(8) NOT NULL DEFAULT "0.00",
	    percent VARCHAR(8) NOT NULL DEFAULT "0.00",
	    x INTEGER NOT NULL DEFAULT 0,
	    y INTEGER NOT NULL DEFAULT 0,
	    exclusive BOOLEAN NOT NULL DEFAULT 1
	);`

	statement, err := repository.database.Prepare(createDealsTableSQL)
	defer statement.Close()
	if err != nil {
		log.Fatalf("Failed to create table: %v", err.Error())
	}
	statement.Exec()
}

func (repository *ProductRepository) createOfferingsTable() {
	createOfferingsTableSQL := `CREATE TABLE offerings (
	    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	    product_id INTEGER NOT NULL,
	    deal_id INTEGER NOT NULL,
	    modified_price VARCHAR(8) NOT NULL DEFAULT "NAN",
	    active BOOLEAN NOT NULL DEFAULT 1,
	    FOREIGN KEY (product_id) REFERENCES products (id),
	    FOREIGN KEY (deal_id) REFERENCES deals (id) );`

	statement, err := repository.database.Prepare(createOfferingsTableSQL)
	defer statement.Close()
	if err != nil {
		log.Fatalf("Failed to create table: %v", err.Error())
	}
	statement.Exec()
}
