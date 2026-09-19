package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/products", productHandler)

	fmt.Println("server is running at 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}

func productHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		fmt.Fprintln(w, "GET: Give me the products")

	case http.MethodPost:
		fmt.Fprintln(w, "POST: Create a product")

	case http.MethodPut:
		fmt.Fprintln(w, "PUT: Update a product")

	case http.MethodDelete:
		fmt.Fprintln(w, "DELETE: Delete a product")

	default:
	   http.Error(w, "Methos not allowd", http.StatusMethodNotAllowed)

	}
}
