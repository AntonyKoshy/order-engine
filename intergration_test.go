package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrderFlowIntegration(t *testing.T) {
	resetState()
	StartMatcher()

	sell := Order{
		ID:       "S1",
		Symbol:   "TCS",
		Side:     "sell",
		Price:    100,
		Quantity: 120,
	}

	body, _ := json.Marshal(sell)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	orderHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status code 200 got %d", rec.Code)
	}
	buy := Order{
		ID:       "B1",
		Symbol:   "TCS",
		Side:     "buy",
		Price:    105,
		Quantity: 60,
	}
	body, _ = json.Marshal(buy)
	req = httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body))
	rec = httptest.NewRecorder()

	orderHandler(rec, req)

	wg.Wait()
	//calling orderbook
	req = httptest.NewRequest(http.MethodGet, "/orderbook?symbol=TCS", nil)
	rec = httptest.NewRecorder()

	orderBookHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status code 200 got %d", rec.Code)
	}

	var book OrderBook

	if err := json.NewDecoder(rec.Body).Decode(&book); err != nil {
		log.Fatal("failed to decode orderbook")
	}

	if len(book.SellOrders) != 1 {
		t.Fatalf("expected 1 sell order got %d", len(book.SellOrders))
	}
	if book.SellOrders[0].Quantity != 60 {
		t.Fatalf("expected qty 60 found %d", book.SellOrders[0].Quantity)
	}
	if len(book.BuyOrders) != 0 {
		t.Fatalf("expected 0 sell order got %d", len(book.BuyOrders))
	}

}

func resetState() {
	initState()
}
