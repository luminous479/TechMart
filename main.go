package main

import (
	"fmt"
	"net/http"

)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", productHandler)

	mux.HandleFunc("GET /products/{id}", productByIdHandler)

	fmt.Println("server is running at 8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}

func productHandler(w http.ResponseWriter, r *http.Request) {


	fmt.Fprintln(w,"View All Products")
}

func productByIdHandler(w http.ResponseWriter, r *http.Request){

	fmt.Fprintln(w, "Product ID : ", r.PathValue("id"))
}
