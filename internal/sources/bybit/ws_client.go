package bybit

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	wsURL          = "wss://stream.bybit.com/v5/public/linear"
	reconnectDelay = 5 * time.Second
	pingInterval   = 20 * time.Second
)

type MessageHandler func(topic string, msg []byte)

type Client struct {
	topics   []string
	conn     *websocket.Conn
	handlers map[string]MessageHandler
	done     chan struct{}
}

func NewClient(topics []string) *Client {
	return &Client{
		topics:   topics,
		handlers: make(map[string]MessageHandler),
		done:     make(chan struct{}),
	}
}

func (c *Client) RegisterHandler(topic string, handler MessageHandler) {
	c.handlers[topic] = handler
}

func (c *Client) Start() {
	for {
		err := c.connectAndServe()
		log.Println("[start] Reconnecting after delay...")
		time.Sleep(reconnectDelay)
		if err != nil {
			log.Println("[start] Reconnect error:", err)
		}
	}
}

func (c *Client) connectAndServe() error {
	log.Println("[connectAndServe] Connecting to Bybit WebSocket...")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("[connectAndServe] Error closing connection:", err)
		}
	}(conn)

	if err := c.subscribe(); err != nil {
		log.Println("[connectAndServe] Subscription error:", err)
		return err
	}

	go c.keepAlive()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("[connectAndServe] Read error:", err)
			close(c.done)
			return err
		}
		c.handleMessage(msg)
	}
}

func (c *Client) subscribe() error {
	req := map[string]interface{}{
		"op":   "subscribe",
		"args": c.topics,
	}
	return c.conn.WriteJSON(req)
}

func (c *Client) keepAlive() {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				log.Println("[keepAlive] Ping error:", err)
				return
			}
		case <-c.done:
			return
		}
	}
}

func (c *Client) handleMessage(msg []byte) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(msg, &raw); err != nil {
		log.Println("[handleMessage] Invalid JSON:", err)
		return
	}

	topicRaw, ok := raw["topic"]
	if !ok {
		log.Println("[handleMessage] No topic field, raw msg:", string(msg))
		return
	}

	var topic string
	err := json.Unmarshal(topicRaw, &topic)
	if err != nil {
		return
	}

	if handler, ok := c.handlers[topic]; ok {
		handler(topic, msg)
	} else {
		log.Println("[handleMessage] No handler for topic:", topic)
	}
}
