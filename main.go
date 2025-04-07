package main

import (
	"arbzoo/internal/sources/bybit"
	"encoding/json"
	"log"
)

func main() {
	// Initialize the Bybit WebSocket client
	client := bybit.NewClient([]string{
		bybit.GetOrderBookTopic(1, "FARTCOINUSDT"),
		bybit.GetOrderBookTopic(1, "WIFUSDT"),
	})
	client.RegisterHandler(bybit.GetOrderBookTopic(1, "FARTCOINUSDT"), func(topic string, msg []byte) {
		var orderBook bybit.OrderBookMessage
		err := json.Unmarshal(msg, &orderBook)
		if err != nil {
			log.Println("[FARTCOIN OrderBook] Unmarshal error:", err)
			return
		}
		log.Printf("[FARTCOIN OrderBook] %+v\n", orderBook)
	})

	client.RegisterHandler(bybit.GetOrderBookTopic(1, "WIFUSDT"), func(topic string, msg []byte) {
		var orderBook bybit.OrderBookMessage
		err := json.Unmarshal(msg, &orderBook)
		if err != nil {
			log.Println("[WIFUSDT OrderBook] Unmarshal error:", err)
			return
		}
		log.Printf("[WIFUSDT OrderBook] %+v\n", orderBook)
	})
	client.Start()
}
