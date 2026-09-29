package repository

import (
	"database/sql"
	"errors"
)

type stockRepository struct {
	db *sql.DB
}

func NewstockRepository(db *sql.DB) *stockRepository {
	return &stockRepository{
		db: db,
	}
}

func (r *stockRepository) StockIn(
	productID int,
	quantity int,
	reason string,
) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	result, err := tx.Exec(`
		UPDATE products
		SET quantity = quantity + $1
		WHERE id = $2
	`, quantity, productID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	_, err = tx.Exec(`
		INSERT INTO stock_movements
			(product_id, type, quantity, reason)
		VALUES ($1, 'IN', $2, $3)
	`, productID, quantity, reason)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *stockRepository) StockOut(
	productID int,
	quantity int,
	reason string,
) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var currentQuantity int

	err = tx.QueryRow(`
		SELECT quantity
		FROM products
		WHERE id = $1
		FOR UPDATE
	`, productID).Scan(&currentQuantity)

	if err != nil {
		return err
	}

	if currentQuantity < quantity {
		return errors.New("insufficient stock")
	}

	_, err = tx.Exec(`
		UPDATE products
		SET quantity = quantity - $1
		WHERE id = $2
	`, quantity, productID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO stock_movements
			(product_id, type, quantity, reason)
		VALUES ($1, 'OUT', $2, $3)
	`, productID, quantity, reason)

	if err != nil {
		return err
	}

	return tx.Commit()
}
