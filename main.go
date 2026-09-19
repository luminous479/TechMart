package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Product struct {
	Id       int  `json:"id"`
	Name     string `json:"name"`
	SKU      string `json:"sku"`
	Price    float64 `json:"price"`
	Quantity int  `json: "quantity"`
}

var products = []Product{
	{
		Id:       1,
		Name:     "Logitech Keyboard",
		SKU:      "KB-001",
		Price:    45,
		Quantity: 50,
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

	var product Product

	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadGateway)
	}
	if product.Name == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}

	if product.SKU == "" {
		http.Error(w, "Product SKU is required", http.StatusBadRequest)
		return
	}

	if product.Price <= 0 {
		http.Error(w, "Product price must be greater than 0", http.StatusBadRequest)
		return
	}

	if product.Quantity < 0 {
		http.Error(w, "Product quantity cannot be negative", http.StatusBadRequest)
		return
	}
	product.Id = len(products) + 1

	products = append(products, product)
    w.Header().Set("Content-Type","application/json")
	w.Header().Set("Location",fmt.Sprintf("/products/%d",product.Id))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)

}

func getProduct(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, "Invalid Product Id ", http.StatusBadRequest)
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

	

	id, err := strconv.Atoi(r.PathValue("id"))
      if err != nil{
				http.Error(w, "Invalid Request Body ", http.StatusBadRequest)
	  }
	  var updateProduct Product
     err = json.NewDecoder(r.Body).Decode(&updateProduct)

	if err != nil {
		http.Error(w, "Invalid Request Body ", http.StatusBadRequest)
		return
	}


	for i, product := range products{

		if product.Id == id {
			updateProduct.Id = id
          products[i] = updateProduct
		}
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updateProduct)
		return 
	}

	http.Error(w,"Product Not Found",http.StatusNotFound)


}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	for i, product := range products {
		if product.Id == id {
			products = append(products[:i], products[i+1:]...)

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}
