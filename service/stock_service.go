package service

import (
	"errors"

	"github.com/luminous479/TechMart/model"
)

type StockRepository interface {
	StockIn(productID int, quantity int, reason string) error
}

type StockService struct {
	repo StockRepository
}

func NewStockService(repo StockRepository) *StockService {
	return &StockService{
		repo: repo,
	}
}

func (s *StockService) StockIn(
	productID int,
	quantity int,
	reason string,
) error {

	if quantity <= 0 {
		return errors.New("stock-in quantity must be greater than 0")
	}

	return s.repo.StockIn(productID, quantity, reason)
}