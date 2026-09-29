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

	handler "github.com/luminous479/TechMart/handler"
	model "github.com/luminous479/TechMart/model"
	repo "github.com/luminous479/TechMart/repository"
	service "github.com/luminous479/TechMart/service"
	helper "github.com/luminous479/TechMart/helper"
)


type ProductResponse struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	SKU      string  `json:"sku"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type StatusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}
type StockRequest struct {
	Quantity int    `json:"quantity"`
	Reason   string `json:"reason"`
}
type StockMovement struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	Type      string    `json:"type"`
	Quantity  int       `json:"quantity"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
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

	productRepository := repo.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", productHandler.GetProducts)
	mux.HandleFunc("POST /products", productHandler.CreateProduct)
	mux.HandleFunc("GET /products/{id}", productHandler.GetProduct)
	mux.HandleFunc("PUT /products/{id}", productHandler.UpdateProduct)
	mux.HandleFunc("DELETE /products/{id}", productHandler.DeleteProduct)
	mux.HandleFunc("POST /products/{id}/stock-in", stockInHandler(db))
	mux.HandleFunc("POST /products/{id}/stock-out", stockOutHandler(db))
	mux.HandleFunc("GET /products/{id}/movements", getStockMovements(db))

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

func toProductResponse(product model.Product) ProductResponse {
	return ProductResponse{
		ID:       product.ID,
		Name:     product.Name,
		SKU:      product.SKU,
		Price:    product.Price,
		Quantity: product.Quantity,
	}
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

		helper.WriteJSONError(
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

func stockInHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
		helper.WriteJSONError(
				w,
				"Invalid product ID",
				http.StatusBadRequest,
			)
			return
		}

		var request StockRequest

		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
		helper.WriteJSONError(
				w,
				"Invalid request body",
				http.StatusBadRequest,
			)
			return
		}

		err = stockIn(
			db,
			id,
			request.Quantity,
			request.Reason,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
			helper.WriteJSONError(
					w,
					"Product not found",
					http.StatusNotFound,
				)
				return
			}

	helper.WriteJSONError(
				w,
				"Failed to process stock-in",
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Stock added successfully",
		})
	}
}
func stockOut(db *sql.DB, productID int, quantity int, reason string) error {
	if quantity <= 0 {
		return errors.New("stock-out quantity must be greater than 0")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.Exec(`
		UPDATE products
		SET quantity = quantity - $1
		WHERE id = $2
		  AND quantity >= $1
	`, quantity, productID)

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

	_, err = tx.Exec(`
		INSERT INTO stock_movements
		(product_id, type, quantity, reason)
		VALUES ($1, $2, $3, $4)
	`, productID, "OUT", quantity, reason)

	if err != nil {

		return err
	}

	return tx.Commit()
}
func stockOutHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
	helper.WriteJSONError(
				w,
				"Invalid product ID",
				http.StatusBadRequest,
			)
			return
		}

		var request StockRequest

		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
		helper.WriteJSONError(
				w,
				"Invalid request body",
				http.StatusBadRequest,
			)
			return
		}

		err = stockOut(
			db,
			id,
			request.Quantity,
			request.Reason,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
			helper.WriteJSONError(
					w,
					"Product not found or insufficient stock",
					http.StatusBadRequest,
				)
				return
			}

		helper.WriteJSONError(
				w,
				"Failed to process stock-out",
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Stock removed successfully",
		})
	}
}
func getStockMovementsFromDB(
	db *sql.DB,
	productID int,
) ([]StockMovement, error) {

	rows, err := db.Query(`
		SELECT id, product_id, type, quantity, reason, created_at
		FROM stock_movements
		WHERE product_id = $1
		ORDER BY created_at DESC
	`, productID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var movements []StockMovement

	for rows.Next() {
		var movement StockMovement

		err := rows.Scan(
			&movement.ID,
			&movement.ProductID,
			&movement.Type,
			&movement.Quantity,
			&movement.Reason,
			&movement.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		movements = append(movements, movement)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movements, nil
}
func getStockMovements(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
		helper.WriteJSONError(
				w,
				"Invalid product ID",
				http.StatusBadRequest,
			)
			return
		}

		movements, err := getStockMovementsFromDB(db, id)
		if err != nil {
		helper.WriteJSONError(
				w,
				"Failed to get stock movements",
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(movements)
	}
}
