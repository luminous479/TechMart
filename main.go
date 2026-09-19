package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type Product struct{
	Id int
	Name string
	SKU string
	Price float64
	Quantity int

}

var products = []Product{
    {
		Id : 1,
		Name :"Logitech Keyboard",
		SKU : "KB-001",
		Price : 45,
		Quantity : 50,
	},
	{
		Id:       2,
		Name:     "Dell Monitor",
		SKU:      "MN-001",
		Price:    220,
		Quantity: 20,
	},
	{
		Id:       3,
		Name:     "HP Mouse",
		SKU:      "MS-001",
		Price:    25,
		Quantity: 100,
	},
}
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

	for _, product := range products {
		
		fmt.Fprintf(
			w,
			"ID: %d | Name: %s | SKU: %s | Price: $%.2f | Quantity: %d\n",
			product.Id,
			product.Name,
			product.SKU,
			product.Price,
			product.Quantity,
		)
		

	}
}

func createProduct(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "POST: Create a product")
}

func getProduct(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil{
		http.Error(w,"Invalid Product Id ",http.StatusBadRequest)
		return
	}
	
	for _, product := range products {
 
		if product.Id == id {

		fmt.Fprintf(
			w,
			"ID: %d | Name: %s | SKU: %s | Price: $%.2f | Quantity: %d\n",
			product.Id,
			product.Name,
			product.SKU,
			product.Price,
			product.Quantity,
		)
		return 
		}
		


	}

	http.Error(w, "Product not found", http.StatusNotFound)

}

func updateProduct(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	fmt.Fprintln(w, "PUT: Update product ID:", id)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	fmt.Fprintln(w, "DELETE: Delete product ID:", id)
}