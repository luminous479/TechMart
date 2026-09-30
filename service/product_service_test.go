package service

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/luminous479/TechMart/model"
)

type mockProductRepository struct {
	products []model.Product
}

func (m *mockProductRepository) GetProducts() ([]model.Product, error) {
	return m.products, nil
}

func (m *mockProductRepository) GetProduct(id int) (*model.Product, error) {
	for _, product := range m.products {
		if product.ID == id {
			return &product, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (m *mockProductRepository) CreateProduct(product model.Product) (int, error) {
	product.ID = len(m.products) + 1
	m.products = append(m.products, product)

	return product.ID, nil
}

func (m *mockProductRepository) UpdateProduct(
	id int,
	product model.Product,
) error {
	for i := range m.products {
		if m.products[i].ID == id {
			m.products[i] = product
			m.products[i].ID = id
			return nil
		}
	}

	return errors.New("product not found")
}

func (m *mockProductRepository) DeleteProduct(id int) error {
	for i := range m.products {
		if m.products[i].ID == id {
			m.products = append(m.products[:i], m.products[i+1:]...)
			return nil
		}
	}

	return errors.New("product not found")
}
func TestProductService_CreateProduct(t *testing.T) {
	repo := &mockProductRepository{}
	productValidator := BasicProductValidator{}

	service := NewProductService(repo,productValidator)

	product := model.Product{
		Name:     "Keyboard",
		SKU:      "KB-001",
		Price:    45.00,
		Quantity: 10,
	}

	id, err := service.CreateProduct(product)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if id != 1 {
		t.Fatalf("expected ID 1, got %d", id)
	}
}
