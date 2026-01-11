package main

import (
	"time"
)

type OrderBook struct {
	BuyOrders  []Order
	SellOrders []Order
}

type Order struct {
	ID        string `json:"id"`
	Symbol    string `json:"symbol"`
	Side      string `json:"side"`
	Price     int    `json:"price"`
	Quantity  int    `json:"quantity"`
	TimeStamp time.Time
}
