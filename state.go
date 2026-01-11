package main

import "sync"

var (
	orderCh   chan Order
	orderBook map[string]*OrderBook
	mu        sync.RWMutex
	wg        sync.WaitGroup
)

func initState() {
	orderCh = make(chan Order, 100)
	orderBook = make(map[string]*OrderBook)

}
