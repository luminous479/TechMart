package main

import (
	"fmt"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", getProducts)
	mux.HandleFunc("POST /products", createProduct)

	mux.HandleFunc("GET /products/{id}", getProduct)
	mux.HandleFunc("PUT /products/{id}", updateProduct)
	mux.HandleFunc("DELETE /products/{id}", deleteProduct)

	fmt.Println("server is running at http://localhost:8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func getProducts(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "GET: All products")
}

func createProduct(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "POST: Create a product")
}

func getProduct(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	fmt.Fprintln(w, "GET: Product ID:", id)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	fmt.Fprintln(w, "PUT: Update product ID:", id)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	fmt.Fprintln(w, "DELETE: Delete product ID:", id)
}