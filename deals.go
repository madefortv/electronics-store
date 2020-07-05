package store

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func (repository *ProductRepository) insertDeal(deal Deal) error {
	tx, _ := repository.database.Begin()
	// logic for type needs to be put here
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

/*
func (repository *ProductRepository) updateDeal(deal Deal) error {
	tx, _ := repository.database.Begin()
	stmt, _ := tx.Prepare(`UPDATE deals SET (name, type, coupon, percent, x, y, exclusive) Vname = ?, description = ?, price = ? WHERE id = ?;`)
	defer stmt.Close()
	_, err := stmt.Exec(deal.Name, deal.Description, deal.Price, deal.Id) //.Scan(&id, &name, &description, &price)
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

func (repository *ProductRepository) deleteDeal(deal Deal) error {
	tx, _ := repository.database.Begin()
	stmt, _ := tx.Prepare(`DELETE FROM deals WHERE id = ?`)
	defer stmt.Close()

	_, err := stmt.Exec(deal.Id)
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

func (repository *ProductRepository) FindAll() []*Deal {
	rows, _ := repository.database.Query(`SELECT id, name, description, price FROM deals;`)
	defer rows.Close()

	deals := []*Deal{}

	for rows.Next() {
		var (
			id          int
			name        string
			description string
			price       string
		)

		rows.Scan(&id, &name, &description, &price)

		deals = append(deals, &Deal{
			Id:          id,
			Name:        name,
			Description: description,
			Price:       price,
		})
	}

	return deals
}
*/
