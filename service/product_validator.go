package service

import (
	"errors"

	"github.com/luminous479/TechMart/model"
)

type ProductValidator interface {
	validate(product model.Product) error
}

type BasicProductValidator struct{}

func (BasicProductValidator) validate(product model.Product) error {

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
