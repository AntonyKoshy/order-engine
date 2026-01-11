package main

//logic to insert in correct order

// for buy orders in descending
func insertBuyOrder(book *OrderBook, order Order) {
	idx := 0
	for idx < len(book.BuyOrders) {
		existing := book.BuyOrders[idx]
		//enforce time priority
		if order.Price < existing.Price {
			idx++
			continue
		}
		if order.Price == book.BuyOrders[idx].Price &&
			existing.TimeStamp.Before(order.TimeStamp) {
			idx++
			continue
		}
		break
	}

	book.BuyOrders = append(book.BuyOrders, Order{})
	copy(book.BuyOrders[idx+1:], book.BuyOrders[idx:])
	book.BuyOrders[idx] = order
}

// for sell orders in ascending
func insertSellOrder(book *OrderBook, order Order) {
	idx := 0
	for idx < len(book.SellOrders) {
		existing := book.SellOrders[idx]
		if order.Price > existing.Price {
			idx++
			continue
		}
		if order.Price == existing.Price && existing.TimeStamp.Before(order.TimeStamp) {
			idx++
			continue
		}
		break
	}
	book.SellOrders = append(book.SellOrders, Order{})
	copy(book.SellOrders[idx+1:], book.SellOrders[idx:])
	book.SellOrders[idx] = order

}
