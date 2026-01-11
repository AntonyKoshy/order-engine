package main

import (
	"fmt"
	"net/http"
)

func main() {

	StartMatcher()

	http.HandleFunc("/order", orderHandler)
	http.HandleFunc("/orderbook", orderBookHandler)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)

}
