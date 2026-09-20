package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	SKU      string  `json:"sku"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

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

type ProductResponse struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	SKU      string  `json:"sku"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var products = []Product{
	{
		ID:       1,
		Name:     "Logitech Keyboard",
		SKU:      "KB-001",
		Price:    45,
		Quantity: 50,
	},
	{
		ID:       2,
		Name:     "Dell Monitor",
		SKU:      "MN-001",
		Price:    220,
		Quantity: 20,
	},
	{
		ID:       3,
		Name:     "HP Mouse",
		SKU:      "MS-001",
		Price:    25,
		Quantity: 100,
	},
}

type StatusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (sr *StatusRecorder) WriteHeader(status int) {
	if sr.wroteHeader {
		return
	}

	sr.status = status
	sr.wroteHeader = true

	sr.ResponseWriter.WriteHeader(status)
}

func (sr *StatusRecorder) Write(data []byte) (int, error) {
	if !sr.wroteHeader {
		sr.WriteHeader(http.StatusOK)
	}

	return sr.ResponseWriter.Write(data)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", getProducts)
	mux.HandleFunc("POST /products", createProduct)
	mux.HandleFunc("GET /products/{id}", getProduct)
	mux.HandleFunc("PUT /products/{id}", updateProduct)
	mux.HandleFunc("DELETE /products/{id}", deleteProduct)

	fmt.Println("server is running at http://localhost:8080")

	handler := chainMiddleware(
		mux,
		recoveryMiddleware,
		loggingMiddleware,
	)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	responses := make([]ProductResponse, 0, len(products))

	for _, product := range products {
		responses = append(responses, toProductResponse(product))
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(responses)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var requestBody CreateProductRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		writeJSONError(
			w,
			"Invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	err = validateProduct(
		requestBody.Name,
		requestBody.SKU,
		requestBody.Price,
		requestBody.Quantity,
	)

	if err != nil {
		writeJSONError(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	product := Product{
		ID:       len(products) + 1,
		Name:     requestBody.Name,
		SKU:      requestBody.SKU,
		Price:    requestBody.Price,
		Quantity: requestBody.Quantity,
	}

	products = append(products, product)

	response := toProductResponse(product)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(
		"Location",
		fmt.Sprintf("/products/%d", product.ID),
	)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSONError(
			w,
			"Invalid product ID",
			http.StatusBadRequest,
		)
		return
	}

	for _, product := range products {
		if product.ID == id {
			response := toProductResponse(product)

			w.Header().Set("Content-Type", "application/json")

			json.NewEncoder(w).Encode(response)
			return
		}
	}

	writeJSONError(
		w,
		"Product not found",
		http.StatusNotFound,
	)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSONError(
			w,
			"Invalid product ID",
			http.StatusBadRequest,
		)
		return
	}

	var requestBody UpdateProductRequest

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		writeJSONError(
			w,
			"Invalid request body",
			http.StatusBadRequest,
		)
		return
	}
	err = validateProduct(
		requestBody.Name,
		requestBody.SKU,
		requestBody.Price,
		requestBody.Quantity,
	)

	if err != nil {
		writeJSONError(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	for i, product := range products {
		if product.ID == id {
			updatedProduct := Product{
				ID:       id,
				Name:     requestBody.Name,
				SKU:      requestBody.SKU,
				Price:    requestBody.Price,
				Quantity: requestBody.Quantity,
			}

			products[i] = updatedProduct

			response := toProductResponse(updatedProduct)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			json.NewEncoder(w).Encode(response)
			return
		}
	}

	writeJSONError(
		w,
		"Product not found",
		http.StatusNotFound,
	)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSONError(
			w,
			"Invalid product ID",
			http.StatusBadRequest,
		)
		return
	}

	for i, product := range products {
		if product.ID == id {
			products = append(
				products[:i],
				products[i+1:]...,
			)

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	writeJSONError(
		w,
		"Product not found",
		http.StatusNotFound,
	)
}

func toProductResponse(product Product) ProductResponse {
	return ProductResponse{
		ID:       product.ID,
		Name:     product.Name,
		SKU:      product.SKU,
		Price:    product.Price,
		Quantity: product.Quantity,
	}
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recorder := &StatusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start)

		fmt.Printf(
			"%s %s → %d → %v\n",
			r.Method,
			r.URL.Path,
			recorder.status,
			duration,
		)
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("PANIC:", err)

				writeJSONError(
					w,
					"Internal Server Error",
					http.StatusInternalServerError,
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func chainMiddleware(
	handler http.Handler,
	middlewares ...func(http.Handler) http.Handler,
) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
func validateProduct(
	name string,
	sku string,
	price float64,
	quantity int,
) error {
	if name == "" {
		return errors.New("product name is required")
	}

	if sku == "" {
		return errors.New("product SKU is required")
	}

	if price <= 0 {
		return errors.New("product price must be greater than 0")
	}

	if quantity < 0 {
		return errors.New("product quantity cannot be negative")
	}

	return nil
}
