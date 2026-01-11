# Concurrent Order Matching Engine (Go)

A concurrent, in-memory order matching engine written in Go.  
The system supports buy/sell order ingestion via HTTP, matches orders using **price–time priority**, and exposes the current order book per symbol.

This project focuses on **correctness, concurrency, testability, and clear trade-offs**, rather than infrastructure complexity.

---

## Features

- In-memory order book per symbol
- Buy/Sell order matching
- **Price–time priority**
- Concurrent request handling
- Deterministic testing (no sleeps)
- Benchmarked critical path

---

## High - Level Architecture

Orders are never applied directly from HTTP handlers. Instead, handlers enqueueorders into a buffered channel. A single background goroutine owns all mutations to the order book, which makes reasoning about correctness straightforward.

```
Client
  │
  │  HTTP POST /order
  │
  ▼
HTTP Handler (orderHandler)
  │
  │  enqueue order
  ▼
Buffered Channel (orderCh)
  │
  │  single consumer
  ▼
Matcher Goroutine
  │
  │  match + insert (price–time priority)
  ▼
In-Memory Order Book (per symbol)
  │
  │  protected by RWMutex
  ▼
HTTP Handler (orderBookHandler)
  │
  │  HTTP GET /orderbook
  ▼
Client
```


---

### Key Ideas

- **HTTP handlers do not modify shared state**
- Orders are **enqueued** into a channel
- A **single matcher goroutine** owns all writes to the order book
- Concurrent reads are protected with `sync.RWMutex`

---

## Concurrency Model

- **One writer**: matcher goroutine
- **Many readers**: order book HTTP handler
- Channel decouples ingestion from processing
- `RWMutex` allows safe concurrent reads
- Verified using `go test -race`

This design avoids complex locking and prevents race conditions by construction.

---

## Order Matching Logic

- Buy orders match against the lowest-priced sell orders
- Sell orders match against the highest-priced buy orders
- Matching continues until:
  - Quantity is exhausted, or
  - Price condition no longer holds
- Remaining quantity is inserted using **price–time priority**

### Price–Time Priority

- Better price always wins
- For equal prices, earlier timestamp wins
- Orders are inserted into sorted slices

---

## Testing Strategy

### Unit Tests
- Matching logic
- Partial fills
- Full fills
- Price–time priority ordering

### Integration Test
- End-to-end flow:
  - `POST /order`
  - Matcher goroutine processes order
  - `GET /orderbook`
- Uses `httptest`
- Synchronization via `sync.WaitGroup` (no `time.Sleep`)

Run tests:

```bash
go test ./...
go test -race ./...
```

## Running the Server

Ensure Go is installed (Go 1.25+ recommended).

From the project root:

```bash
go run main.go




