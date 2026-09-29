package service

import (
	"errors"

	"github.com/luminous479/TechMart/model"
	repos "github.com/luminous479/TechMart/repository"
)

var ErrDuplicateSKU = errors.New("product SKU already exists")

type ProductService struct {
	repo *repos.ProductRepository
}

func NewProductService(repo repos.ProductRepository) *ProductService {

	return &ProductService{
		repo: &repo,
	}
}
func (s *ProductService) GetProducts() ([]model.Product, error) {

	products, err := s.repo.GetProducts()
	if err != nil {
		return nil, err
	}

	return products, nil
}
func (s *ProductService) GetProduct(id int) (*model.Product, error) {
	return s.repo.GetProduct(id)
}
func (s *ProductService) CreateProduct(product model.Product,
) (int, error) {
	if err := validateProduct(product); err != nil {
		return 0, err
	}

	id, err := s.repo.CreateProduct(product)
	if err != nil {

		return 0, err

	}
	return id, nil
}
func (s *ProductService) UpdateProduct(
	id int,
	product model.Product,
) error {
	if err := validateProduct(product); err != nil {
		return err
	}

	return s.repo.UpdateProduct(id, product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}
func validateProduct(product model.Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}

	if product.SKU == "" {
		return errors.New("product SKU is required")
	}

	if product.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}

	if product.Quantity < 0 {
		return errors.New("product quantity cannot be negative")
	}

	return nil
}
