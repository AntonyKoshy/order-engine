package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// incoming order is queued thats it
func orderHandler(w http.ResponseWriter, r *http.Request) {
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid Order", http.StatusBadRequest)
	}
	order.TimeStamp = time.Now()
	//queuing the order
	wg.Add(1)
	orderCh <- order
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("order received"))

}

// in-memory orders for that symbol are returned thats it
func orderBookHandler(w http.ResponseWriter, r *http.Request) {

	mu.RLock()
	defer mu.RUnlock()
	symbol := r.URL.Query().Get("symbol")

	book, exists := orderBook[symbol]
	if !exists {
		http.Error(w, "Symbol not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(book)

}
