package main

import (
	"testing"
	"time"
)

func TestMatchingScenarios(t *testing.T) {
	tests := []struct {
		Name          string
		Book          *OrderBook
		Incoming      Order
		ExpectedBuys  int
		ExpectedSells int
	}{
		{
			Name: "partial fills sell",
			Book: &OrderBook{
				SellOrders: []Order{
					{ID: "S1", Symbol: "TCS", Side: "sell", Quantity: 10, Price: 100},
				},
			},
			Incoming:      Order{ID: "S1", Symbol: "TCS", Side: "buy", Quantity: 5, Price: 105},
			ExpectedBuys:  0,
			ExpectedSells: 1,
		}, {
			Name: "full match sell",
			Book: &OrderBook{
				SellOrders: []Order{
					{ID: "S1", Symbol: "TCS", Side: "sell", Quantity: 5, Price: 100},
				},
			}, Incoming: Order{ID: "S1", Symbol: "TCS", Side: "buy", Quantity: 10, Price: 105},
			ExpectedBuys:  1,
			ExpectedSells: 0,
		},
		{
			Name: "no cross sell",
			Book: &OrderBook{
				SellOrders: []Order{
					{ID: "S1", Symbol: "TCS", Side: "buy", Quantity: 5, Price: 120},
				},
			},
			Incoming:      Order{ID: "S1", Symbol: "TCS", Side: "buy", Quantity: 10, Price: 105},
			ExpectedBuys:  1,
			ExpectedSells: 1,
		},
		{
			Name: "Sell with multiple buys",
			Book: &OrderBook{
				BuyOrders: []Order{
					{ID: "B1", Symbol: "TCS", Price: 110, Quantity: 5},
					{ID: "B2", Symbol: "TCS", Price: 108, Quantity: 5},
					{ID: "B3", Symbol: "TCS", Price: 105, Quantity: 5},
				},
			}, Incoming: Order{ID: "S1", Symbol: "TCS", Price: 100, Quantity: 12, Side: "sell"},
			ExpectedBuys:  1,
			ExpectedSells: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			matchOrders(tt.Book, tt.Incoming)
			if len(tt.Book.BuyOrders) != tt.ExpectedBuys || len(tt.Book.SellOrders) != tt.ExpectedSells {
				t.Fatal("Unexpected Result")
			}
		})
	}

}

func TestTimePriorityInserting(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Hour)

	book := &OrderBook{
		BuyOrders: []Order{
			{ID: "B1", Symbol: "NIFTY", Price: 110, Quantity: 5, TimeStamp: t1},
			{ID: "B2", Symbol: "NIFTY", Price: 110, Quantity: 5, TimeStamp: t2},
		},
	}

	incoming := Order{
		ID: "B3", Symbol: "NIFTY", Price: 110, Quantity: 5, TimeStamp: t1.Add(time.Second),
	}
	matchBuy(book, incoming)

	if book.BuyOrders[1].ID != "B3" {
		t.Fatal("Time Priority inserting logic failed")
	}

}

//test insertBuyOrders
//case 1: same price different times - same price FIFO
//case 2: different price  - higher price priority

func TestInsertBuyOrderScenarios(t *testing.T) {

	t1 := time.Now()
	t2 := t1.Add(time.Second)
	//define a struct for test
	tests := []struct {
		Name     string
		Orders   []Order
		Expected []time.Time
	}{
		{
			Name: "same price FIFO",
			Orders: []Order{
				{Price: 150, TimeStamp: t2},
				{Price: 150, TimeStamp: t1},
			},
			Expected: []time.Time{t1, t2},
		},
		{
			Name: "higher price priority",
			Orders: []Order{
				{Price: 130, TimeStamp: t2},
				{Price: 150, TimeStamp: t1},
			},
			Expected: []time.Time{t1, t2},
		},
	}

	//function to call -  insertBuyOrders
	for _, tt := range tests {

		t.Run(tt.Name, func(t *testing.T) {
			book := &OrderBook{}

			for _, order := range tt.Orders {
				insertBuyOrder(book, order)

			}

			for idx, ts := range tt.Expected {
				if !book.BuyOrders[idx].TimeStamp.Equal(ts) {
					t.Fatalf("order mismatched at idx %d", idx)
				}
			}

		})

	}

	//
}

func TestInsertSellOrderScenarios(t *testing.T) {

	t1 := time.Now()
	t2 := t1.Add(time.Second)
	//define a struct for test
	tests := []struct {
		Name     string
		Orders   []Order
		Expected []time.Time
	}{
		{
			Name: "same price FIFO",
			Orders: []Order{
				{Price: 150, TimeStamp: t2},
				{Price: 150, TimeStamp: t1},
			},
			Expected: []time.Time{t1, t2},
		},
		{
			Name: "higher price priority",
			Orders: []Order{
				{Price: 130, TimeStamp: t2},
				{Price: 150, TimeStamp: t1},
			},
			Expected: []time.Time{t2, t1},
		},
	}

	//function to call -  insertBuyOrders
	for _, tt := range tests {

		t.Run(tt.Name, func(t *testing.T) {
			book := &OrderBook{}

			for _, order := range tt.Orders {
				insertSellOrder(book, order)

			}

			for idx, ts := range tt.Expected {
				if !book.SellOrders[idx].TimeStamp.Equal(ts) {
					t.Fatalf("order mismatched at idx %d", idx)
				}
			}

		})

	}

	//
}

func BenchmarkMatchBuy(b *testing.B) {
	book := &OrderBook{}
	buy := Order{
		Symbol:   "TCS",
		Quantity: 100,
		Price:    243,
		Side:     "buy",
	}
	for i := 0; i < b.N; i++ {
		matchBuy(book, buy)
	}
}
