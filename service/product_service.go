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

func (s *ProductService) UpdateProduct(
	id int,
	product model.Product,
) error {
	return s.repo.UpdateProduct(id, product)
}
