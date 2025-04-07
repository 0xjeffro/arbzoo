package bybit

type OrderBookMessage struct {
	// https://bybit-exchange.github.io/docs/v5/websocket/public/orderbook
	Topic string        `json:"topic"` // Topic name, e.g. "orderbook.1.BTCUSDT"
	Type  string        `json:"type"`  // Data type: "snapshot" or "delta"
	TS    int64         `json:"ts"`    // Server timestamp in milliseconds
	Data  OrderBookData `json:"data"`  // Order book data
}

type OrderBookData struct {
	Symbol   string      `json:"s"`   // Trading pair symbol, e.g. "BTCUSDT"
	Bids     [][2]string `json:"b"`   // Bid levels: [price, size]
	Asks     [][2]string `json:"a"`   // Ask levels: [price, size]
	UpdateID int64       `json:"u"`   // Update ID for this snapshot or delta
	Sequence int64       `json:"seq"` // Cross-sequence number for orderbook comparison
	CTime    int64       `json:"cts"` // Matching engine timestamp (in ms)
}
