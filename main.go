package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
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

	db, err := connectDB()
	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}

	defer db.Close()

	fmt.Println("Connected to PostgreSQL")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", getProducts(db))
	mux.HandleFunc("POST /products", createProduct(db))
	mux.HandleFunc("GET /products/{id}", getProduct(db))
	mux.HandleFunc("PUT /products/{id}", updateProduct(db))
	mux.HandleFunc("DELETE /products/{id}", deleteProduct(db))

	fmt.Println("server is running at http://localhost:8080")

	handler := chainMiddleware(
		mux,
		recoveryMiddleware,
		loggingMiddleware,
	)

	err = http.ListenAndServe(":8080", handler)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func getProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		products, err := getProductsFromDB(db)
		if err != nil {
			writeJSONError(
				w,
				"Failed to get products",
				http.StatusInternalServerError,
			)
			return
		}

		responses := make([]ProductResponse, 0, len(products))

		for _, product := range products {
			responses = append(responses, toProductResponse(product))
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(responses)
	}
}

func createProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			Name:     requestBody.Name,
			SKU:      requestBody.SKU,
			Price:    requestBody.Price,
			Quantity: requestBody.Quantity,
		}

		id, err := createProductInDB(db, product)
		if err != nil {
			writeJSONError(
				w,
				"Failed to create product",
				http.StatusInternalServerError,
			)
			return
		}

		product.ID = id

		response := toProductResponse(product)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set(
			"Location",
			fmt.Sprintf("/products/%d", product.ID),
		)
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(response)
	}
}

func getProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeJSONError(
				w,
				"Invalid product ID",
				http.StatusBadRequest,
			)
			return
		}

		product, err := getProductFromDB(db, id)
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

		response := toProductResponse(*product)

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(response)
	}
}

func updateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		product := Product{
			ID:       id,
			Name:     requestBody.Name,
			SKU:      requestBody.SKU,
			Price:    requestBody.Price,
			Quantity: requestBody.Quantity,
		}

		err = updateProductInDB(db, id, product)
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

		response := toProductResponse(product)

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(response)
	}
}

func deleteProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeJSONError(
				w,
				"Invalid product ID",
				http.StatusBadRequest,
			)
			return
		}

		err = deleteProductFromDB(db, id)
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
				"Failed to delete product",
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
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

// database connection

func connectDB() (*sql.DB, error) {

	dsn := "postgres://postgres@localhost/inventory_db?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil

}

func getProductsFromDB(db *sql.DB) ([]Product, error) {
	rows, err := db.Query(`SELECT id, name, sku, price, quantity
		FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.SKU,
			&product.Price,
			&product.Quantity,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)

	}
	if err := rows.Err(); err != nil {

		return nil, err

	}

	return products, nil

}
func getProductFromDB(db *sql.DB, id int) (*Product, error) {
	var product Product

	err := db.QueryRow(`
		SELECT id, name, sku, price, quantity
		FROM products
		WHERE id = $1
	`, id).Scan(
		&product.ID,
		&product.Name,
		&product.SKU,
		&product.Price,
		&product.Quantity,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil
}
func createProductInDB(db *sql.DB, product Product) (int, error) {
	var id int

	err := db.QueryRow(`
		INSERT INTO products (name, sku, price, quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`,
		product.Name,
		product.SKU,
		product.Price,
		product.Quantity,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
func updateProductInDB(
	db *sql.DB,
	id int,
	product Product,
) error {
	result, err := db.Exec(`
		UPDATE products
		SET name = $1,
		    sku = $2,
		    price = $3,
		    quantity = $4
		WHERE id = $5
	`,
		product.Name,
		product.SKU,
		product.Price,
		product.Quantity,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
func deleteProductFromDB(db *sql.DB, id int) error {
	result, err := db.Exec(`
		DELETE FROM products
		WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
