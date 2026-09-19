package main

import (
	"fmt"
	"net/http"
)

func main() {
	
	http.HandleFunc("/products", productHandeler)

	fmt.Println("server is running at 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}

func productHandeler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Inventory products will appear here")
}