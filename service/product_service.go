package service

import (
	"github.com/luminous479/TechMart/model"
	repos "github.com/luminous479/TechMart/repository"
)

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

	// Business rules can go here.

	return products, nil
}
func (s *ProductService) GetProduct(id int) (*model.Product, error) {
	return s.repo.GetProduct(id)
}

func (s *ProductService) UpdateProduct(id int, product model.Product) error {
	return s.repo.UpdateProduct(id, product)
}
func (s *ProductService) CreateProduct(product model.Product,
) (int, error) {
	return s.repo.CreateProduct(product)
}
func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}
