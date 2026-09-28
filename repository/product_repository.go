package repository

import (
	"database/sql"

	"github.com/luminous479/TechMart/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func getProductsFromDB(db *sql.DB) ([]model.Product, error) {
	rows, err := db.Query(`SELECT id, name, sku, price, quantity
		FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product

	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.SKU,
			&product.Price,
			&product.Quantity,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)

	}
	if err := rows.Err(); err != nil {

		return nil, err

	}

	return products, nil

}
func getProductFromDB(db *sql.DB, id int) (*model.Product, error) {
	var product model.Product

	err := db.QueryRow(`
		SELECT id, name, sku, price, quantity
		FROM products
		WHERE id = $1
	`, id).Scan(
		&product.ID,
		&product.Name,
		&product.SKU,
		&product.Price,
		&product.Quantity,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil
}
func createProductInDB(db *sql.DB, product model.Product) (int, error) {
	var id int

	err := db.QueryRow(`
		INSERT INTO products (name, sku, price, quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`,
		product.Name,
		product.SKU,
		product.Price,
		product.Quantity,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
func (r *ProductRepository) UpdateProduct(
	id int,
	product model.Product,
) error {
	result, err := r.db.Exec(`
		UPDATE products
		SET name = $1,
		    sku = $2,
		    price = $3,
		    quantity = $4
		WHERE id = $5
	`,
		product.Name,
		product.SKU,
		product.Price,
		product.Quantity,
	)
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

	return nil
}
