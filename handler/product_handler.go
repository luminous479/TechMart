package handler

import (
	"encoding/json"
	"net/http"

	"github.com/luminous479/TechMart/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(ser *service.ProductService) *ProductHandler {

	return &ProductHandler{
		service: ser,
	}

}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {

	products, err := h.service.GetProducts()
	if err != nil {
		http.Error(
			w,
			"Failed to get products",
			http.StatusInternalServerError,
		)
		return
	}


	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(products)
}
