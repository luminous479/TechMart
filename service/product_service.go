package service

import (
	"errors"

	"github.com/luminous479/TechMart/model"
	
)

var ErrDuplicateSKU = errors.New("product SKU already exists")


type ProductRepository interface {
	GetProducts() ([]model.Product, error)
	GetProduct(id int) (*model.Product, error)
	CreateProduct(product model.Product) (int, error)
	UpdateProduct(id int, product model.Product) error
	DeleteProduct(id int) error
}

type ProductService struct {
	repo ProductRepository
	validator ProductValidator
}
func NewProductService(repo ProductRepository, validator ProductValidator) *ProductService {
	return &ProductService{
		repo: repo,
		validator: validator,
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
	if err := s.validator.validate(product); err != nil {
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
	if err := s.validator.validate(product); err != nil {
		return err
	}

	return s.repo.UpdateProduct(id, product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}

