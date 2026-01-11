package main

func StartMatcher() {
	go func() {
		for order := range orderCh {
			mu.Lock()

			//save in memory to orderBook
			//check if symbol exists in orderBook

			book, exists := orderBook[order.Symbol]
			if !exists {
				book = &OrderBook{}
				orderBook[order.Symbol] = book
			}
			matchOrders(book, order)

			mu.Unlock()
			wg.Done()

		}
	}()
}

//next step - matching logic - buy and sell

func matchOrders(book *OrderBook, incoming Order) {
	switch incoming.Side {
	case "buy":
		matchBuy(book, incoming)
	case "sell":
		matchSell(book, incoming)
	}
}

func matchBuy(book *OrderBook, buy Order) {

	for i := 0; i < len(book.SellOrders) && buy.Quantity > 0; {

		sell := &book.SellOrders[i]

		if buy.Price < sell.Price {
			break
		}

		tradeQty := min(buy.Quantity, sell.Quantity)

		// fmt.Printf("TRADE: BUY %d %v @ %d\n", tradeQty, buy.Symbol, sell.Price)
		//if tq < order, then keep going but remove that sell order from list
		buy.Quantity -= tradeQty
		sell.Quantity -= tradeQty

		if sell.Quantity == 0 {
			book.SellOrders = append(book.SellOrders[:i], book.SellOrders[i+1:]...)
		} else {
			i++
		}

	}

	if buy.Quantity > 0 {
		//append left over to buy orders
		// book.BuyOrders = append(book.BuyOrders, buy)
		insertBuyOrder(book, buy)
	}

}

func matchSell(book *OrderBook, sell Order) {
	for i := 0; i < len(book.BuyOrders) && sell.Quantity > 0; {
		buy := &book.BuyOrders[i]

		if sell.Price > buy.Price {
			break
		}

		tradeQty := min(sell.Quantity, buy.Quantity)
		// fmt.Printf("TRADE: SELL %d %v @ %d\n", tradeQty, sell.Symbol, buy.Price)

		sell.Quantity -= tradeQty
		buy.Quantity -= tradeQty

		if buy.Quantity == 0 {
			book.BuyOrders = append(book.BuyOrders[:i], book.BuyOrders[i+1:]...)
		} else {
			i++
		}

	}

	if sell.Quantity > 0 {
		// book.SellOrders = append(book.SellOrders, sell)
		insertSellOrder(book, sell)
	}
}
