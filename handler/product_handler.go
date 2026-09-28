package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/luminous479/TechMart/model"
	"github.com/luminous479/TechMart/service"
)

type CreateProductRequest struct {
	Name     string  `json:"name"`
	SKU      string  `json:"sku"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
type UpdateProductRequest struct {
	Name     string  `json:"name"`
	SKU      string  `json:"sku"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}


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

func (h *ProductHandler) CreateProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request CreateProductRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeJSONError(
			w,
			"Invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if request.Name == "" {
		writeJSONError(
			w,
			"Product name is required",
			http.StatusBadRequest,
		)
		return
	}

	if request.SKU == "" {
		writeJSONError(
			w,
			"Product SKU is required",
			http.StatusBadRequest,
		)
		return
	}

	if request.Price <= 0 {
		writeJSONError(
			w,
			"Product price must be greater than 0",
			http.StatusBadRequest,
		)
		return
	}

	if request.Quantity < 0 {
		writeJSONError(
			w,
			"Product quantity cannot be negative",
			http.StatusBadRequest,
		)
		return
	}

	product := model.Product{
		Name:     request.Name,
		SKU:      request.SKU,
		Price:    request.Price,
		Quantity: request.Quantity,
	}

	id, err := h.service.CreateProduct(product)
	if err != nil {
		writeJSONError(
			w,
			"Failed to create product",
			http.StatusInternalServerError,
		)
		return
	}

	product.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(
		"Location",
		fmt.Sprintf("/products/%d", product.ID),
	)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(
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

	var request UpdateProductRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeJSONError(
			w,
			"Invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if request.Name == "" {
		writeJSONError(
			w,
			"Product name is required",
			http.StatusBadRequest,
		)
		return
	}

	if request.SKU == "" {
		writeJSONError(
			w,
			"Product SKU is required",
			http.StatusBadRequest,
		)
		return
	}

	if request.Price <= 0 {
		writeJSONError(
			w,
			"Product price must be greater than 0",
			http.StatusBadRequest,
		)
		return
	}

	if request.Quantity < 0 {
		writeJSONError(
			w,
			"Product quantity cannot be negative",
			http.StatusBadRequest,
		)
		return
	}

	product := model.Product{
		ID:       id,
		Name:     request.Name,
		SKU:      request.SKU,
		Price:    request.Price,
		Quantity: request.Quantity,
	}

	err = h.service.UpdateProduct(id, product)
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
			"Failed to update product",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(product)
}
