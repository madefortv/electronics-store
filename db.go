package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
)

/* Get all the relevant deals and offerings that are also in the cart*/
func (repository *ProductRepository) getProductOfferings() []*ProductOffering {
	rows, _ := repository.database.Query(`
	    SELECT PID, DID, PNAME, DNAME, price, quantity, type, coupon, percent, x, y, modified_price
	    FROM (
		SELECT products.id AS PID, products.name AS PNAME,
		products.price, deals.id AS DID, deals.name AS DNAME,
		deals.type, deals.x, deals.y, deals.coupon, deals.percent,
		offerings.modified_price
		   FROM offerings
		   INNER JOIN products on products.id = offerings.product_id
		   INNER JOIN deals on deals.id = offerings.deal_id
		   WHERE active = 1
	       ) INNER JOIN cart on cart.product_id = pid WHERE cart.quantity > 0;`)
	defer rows.Close()

	var productOfferings []*ProductOffering
	for rows.Next() {
		var (
			pid           int
			did           int
			pname         string
			dname         string
			price         string
			quantity      int
			dtype         DealType
			coupon        string
			percent       string
			x             int
			y             int
			modifiedPrice string
		)
		err := rows.Scan(&pid, &did, &pname, &dname, &price, &quantity, &dtype, &coupon, &percent, &x, &y, &modifiedPrice)
		if err != nil {
			log.Fatalf("DB Scan error %v", err.Error())
		}

		productOfferings = append(productOfferings, &ProductOffering{
			ProductID:     pid,
			DealID:        did,
			ProductName:   pname,
			DealName:      dname,
			Type:          dtype,
			Price:         price,
			Quantity:      quantity,
			X:             x,
			Y:             y,
			Coupon:        coupon,
			Percent:       percent,
			ModifiedPrice: modifiedPrice,
		})

	}
	return productOfferings

}

func (repository *ProductRepository) listCart() []Item {
	rows, _ := repository.database.Query(`SELECT
		products.id,
		products.name,
		products.description,
		products.price,
		cart.quantity
		FROM cart INNER JOIN
		products ON products.id = cart.product_id;`)

	defer rows.Close()

	var items []Item

	for rows.Next() {
		var (
			id          int
			name        string
			price       string
			description string
			quantity    int
		)

		err := rows.Scan(&id, &name, &description, &price, &quantity)
		if err != nil {
			log.Fatalf("DB Scan error %v", err.Error())
		}

		items = append(items, Item{
			Product{
				ID:          id,
				Name:        name,
				Description: description,
				Price:       price,
			},
			quantity,
		})
	}

	return items
}

func (repository *ProductRepository) addToCart(product Product) error {
	tx, _ := repository.database.Begin()
	stmt, _ := tx.Prepare(`INSERT INTO cart (product_id) VALUES (?);`)
	defer stmt.Close()
	_, err := stmt.Exec(product.ID)
	if err != nil {
		err = tx.Rollback()

		log.Fatalf("Statement error %v", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("DB Commit error %v", err.Error())
	}

	return err
}

func (repository *ProductRepository) updateCart(item Item) error {
	tx, _ := repository.database.Begin()
	stmt, _ := tx.Prepare(`UPDATE cart SET quantity = ? WHERE product_id = ?;`)
	defer stmt.Close()
	_, err := stmt.Exec(item.Quantity, item.Product.ID)
	if err != nil {
		err = tx.Rollback()
		log.Fatalf("Statement error %v", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("DB Commit error %v", err.Error())
	}

	return err
}

func (repository *ProductRepository) removeFromCart(product Product) error {
	tx, _ := repository.database.Begin()

	stmt, _ := tx.Prepare(`DELETE FROM cart WHERE product_id = ?`)
	defer stmt.Close()
	_, err := stmt.Exec(product.ID)
	if err != nil {
		err = tx.Rollback()
		log.Fatalf("Statement error %v", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("DB Commit error %v", err.Error())
	}

	return err
}

/* Offerings */
func (repository *ProductRepository) insertOffering(offering Offering) error {
	tx, _ := repository.database.Begin()
	stmt, _ := tx.Prepare(`INSERT INTO offerings (product_id, deal_id, modified_price, active) VALUES (?, ?, ?, ?);`)
	defer stmt.Close()
	_, err := stmt.Exec(offering.ProductID, offering.DealID, offering.ModifiedPrice, offering.Active)
	if err != nil {
		err = tx.Rollback()
		log.Fatalf("Statement error %v", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("DB Commit error %v", err.Error())
	}

	return err
}

/* Lists all the product ID's in a given bundle */
func (repository *ProductRepository) getBundleComponents(dID int) []*Offering {

	rows, _ := repository.database.Query(`SELECT product_id, deal_id FROM offerings WHERE deal_id = ? ;`, dID)
	defer rows.Close()

	offerings := []*Offering{}

	for rows.Next() {
		var (
			productID int
			dealID    int
		)

		err := rows.Scan(&productID, &dealID)
		if err != nil {

			log.Fatalf("Error during scanning rows %v", err.Error())
		}

		offerings = append(offerings, &Offering{
			ProductID: productID,
			DealID:    dealID,
		})
	}

	return offerings
}

/* Deals */
func (repository *ProductRepository) insertDeal(deal Deal) error {
	tx, _ := repository.database.Begin()
	stmt, _ := tx.Prepare(`INSERT INTO deals (name, type, coupon, percent, x, y, exclusive) VALUES (?, ?, ?, ?, ?, ?, ?);`)
	defer stmt.Close()
	_, err := stmt.Exec(deal.Name, deal.Type, deal.Coupon, deal.Percent, deal.X, deal.Y, deal.Exclusive)
	if err != nil {
		err = tx.Rollback()
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

		err := rows.Scan(&id, &name, &btype, &coupon, &percent, &x, &y, &exclusive)
		if err != nil {

			log.Fatalf("Error during scanning rows %v", err.Error())
		}

		deals = append(deals, &Deal{
			ID:        id,
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
		err = tx.Rollback()
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
	_, err := stmt.Exec(product.Name, product.Description, product.Price, product.ID) //.Scan(&id, &name, &description, &price)
	if err != nil {
		err = tx.Rollback()
		log.Fatalf("Statement error %v", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("DB Commit error %v", err.Error())
	}
	return err
}

func (repository *ProductRepository) deleteProduct(product Product) error {
	tx, _ := repository.database.Begin()
	stmt, _ := tx.Prepare(`DELETE FROM products WHERE id = ?`)
	defer stmt.Close()

	_, err := stmt.Exec(product.ID)
	if err != nil {
		err = tx.Rollback()
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

		err := rows.Scan(&id, &name, &description, &price)
		if err != nil {

			log.Fatalf("Error during scanning rows %v", err.Error())
		}

		products = append(products, &Product{
			ID:          id,
			Name:        name,
			Description: description,
			Price:       price,
		})
	}

	return products
}

func (repository *ProductRepository) getProduct(product Product) (Product, error) {
	row := repository.database.QueryRow(`SELECT id, name, description, price FROM products WHERE id = ?;`, product.ID)

	var (
		id          int
		name        string
		description string
		price       string
	)

	err := row.Scan(&id, &name, &description, &price)
	if err != nil {
		log.Fatalf("Query error %v", err.Error())
	}

	product = Product{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
	}

	return product, err
}
