package bybit

import (
	"strconv"
)

func GetOrderBookTopic(depth int, symbol string) string {
	// https://bybit-exchange.github.io/docs/v5/websocket/public/orderbook
	//Depths
	//	Linear & inverse:
	//	Level 1 data, push frequency: 10ms
	//	Level 50 data, push frequency: 20ms
	//	Level 200 data, push frequency: 100ms
	//	Level 500 data, push frequency: 100ms
	//
	//Spot:
	//	Level 1 data, push frequency: 10ms
	//	Level 50 data, push frequency: 20ms
	//	Level 200 data, push frequency: 200ms
	//
	//Option:
	//	Level 25 data, push frequency: 20ms
	//	Level 100 data, push frequency: 100ms
	//
	//Topic:
	//	orderbook.{depth}.{symbol} e.g., orderbook.1.BTCUSDT

	return "orderbook." + strconv.Itoa(depth) + "." + symbol
}
