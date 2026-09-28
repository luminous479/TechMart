package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

func (h *ProductHandler) GetProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSONError(
			w,
			"Invalid product ID",
			http.StatusBadRequest,
		)
		return
	}

	product, err := h.service.GetProduct(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(
				w,
				"Product not found",
				http.StatusNotFound,
			)
			return
		}

		writeJSONError(
			w,
			"Failed to get product",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(product)
}
