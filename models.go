package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

/*
   The Product model represents an item or bundle component. These are atomic.
*/

type Product struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       string `json:"price"`
}

/*
   Enum for deal type
*/
type DealType string

const (
	Retail   DealType = "Retail"
	Flat              = "Flat"
	Percent           = "Percent"
	Bundle            = "Bundle"
	BuyXGetY          = "BuyXGetY"
	Other             = "Other"
)

/*
   A deal struct holds data about the discounts applied to certain products
   they have different types and fields depending on the type

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
	ID        int      `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Type      DealType `json:"type,omitempty"`
	Coupon    string   `json:"coupon,omitempty"`
	Percent   string   `json:"percent,omitempty"`
	X         int      `json:"x,omitempty"`
	Y         int      `json:"y,omitempty"`
	Exclusive bool     `json:"exclusive,omitempty"`
}

/* The offering model is a relationship between one or more products and
   a deal. By default a product for sale is an offering with
   "retail" type deal modifying it.

   An Offering has Products and is modified by Deals. Users add Product Offerings
   to their cart, not Products

   ProductId, DealId -> primary key
   Each deal has a type and that type will influece the control flow for
   calculating the final price of a set of tiems. See utils.go

   @Id is the primary key of this, although ProductId/DealId would work
   @ProductId is a product associated with this
   @DealId deal that modifies the product(s)
   @ModifiedPrice is the Bundle Price
   @Active flag determines whether this deal is active
*/

type Offering struct {
	ID            int    `json:"id,omitempty"`
	ProductID     int    `json:"product_id,omitempty"`
	DealID        int    `json:"deal_id,omitempty"`
	ModifiedPrice string `json:"modified_price,omitempty"`
	Active        bool   `json:"active,omitempty"`
}

/*
   A helpful struct for unzipping joins into.
   After a join of offerings x products x deals, we get a product offering
*/
type ProductOffering struct {
	ProductID     int      `json:"product_id,omitempty"`
	DealID        int      `json:"deal_id,omitempty"`
	Quantity      int      `json:"quantity,omitempty"`
	ModifiedPrice string   `json:"modified_price,omitempty"`
	DealName      string   `json:"deal_name,omitempty"`
	Type          DealType `json:"type,omitempty"`
	Coupon        string   `json:"coupon,omitempty"`
	Percent       string   `json:"percent,omitempty"`
	X             int      `json:"x,omitempty"`
	Y             int      `json:"y,omitempty"`
	Exclusive     bool     `json:"exclusive,omitempty"`
	ProductName   string   `json:"product_name"`
	Description   string   `json:"description,omitempty"`
	Price         string   `json:"price"`
}

type ShoppingCart struct {
	Items []Item `json:"items"`
	Total string `json:"total"`
}

type Item struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

/* Database service */

func NewProductRepository(database *sql.DB) *ProductRepository {
	return &ProductRepository{database: database}
}

type ProductRepository struct {
	database *sql.DB
}

func ConnectDatabase(config *Config) (*sql.DB, error) {
	return sql.Open("sqlite3", config.DatabasePath)
}

/* Helpers */
func (repository *ProductRepository) createCartTable() {
	createProductsTableSQL := `CREATE TABLE cart (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"product_id" INTEGER NOT NULL,
		"quantity" INTEGER NOT NULL DEFAULT 1,
		FOREIGN KEY (product_id) REFERENCES products (id)
	  );`

	statement, err := repository.database.Prepare(createProductsTableSQL)
	if err != nil {
		log.Fatalf("Failed to open database connection")
	}

	defer statement.Close()
	if err != nil {
		log.Fatalf("Failed to create table: %v", err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalf("Database transaction failed: %v", err.Error())
	}

}

func (repository *ProductRepository) createProductsTable() {
	createProductsTableSQL := `CREATE TABLE products (
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		price VARCHAR(8) NOT NULL DEFAULT "NAN",
		name VARCHAR(32) NOT NULL DEFAULT "EMPTY",
		description TEXT
	  );`

	statement, err := repository.database.Prepare(createProductsTableSQL)
	if err != nil {
		log.Fatalf("Failed to open database connection")
	}

	defer statement.Close()
	if err != nil {
		log.Fatalf("Failed to create table: %v", err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalf("Database transaction failed: %v", err.Error())
	}

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
	if err != nil {
		log.Fatalf("Failed to open database connection")
	}

	defer statement.Close()
	if err != nil {
		log.Fatalf("Failed to create table: %v", err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalf("Database transaction failed: %v", err.Error())
	}

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
	if err != nil {
		log.Fatalf("Failed to open database connection")
	}
	defer statement.Close()
	if err != nil {
		log.Fatalf("Failed to create table: %v", err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalf("Database transaction failed: %v", err.Error())
	}

}
