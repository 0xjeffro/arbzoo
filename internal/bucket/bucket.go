package bucket

import (
	"sync"
	"time"
)

type PriceData struct {
	Exchange      string // Exchange name, e.g. Binance, Bybit
	Symbol        string // Trading pair symbol, e.g. BTCUSDT
	Token         string // Token name, e.g. BTC
	BestBidPrice  float64
	BestBidAmount float64
	BestAskPrice  float64
	BestAskAmount float64
	ReceivedAt    int64 // timestamp when the data was received
}

type DataStore struct {
	buckets       map[string]*TokenBucket
	mu            sync.RWMutex
	retentionSec  int64 // the retention period in seconds
	cleanupTicker *time.Ticker
}

type TokenBucket struct {
	Prices []PriceData
	mu     sync.RWMutex
}

func (store *DataStore) Insert(data PriceData) {
	store.mu.RLock()
	bucket, exists := store.buckets[data.Token]
	store.mu.RUnlock()

	if !exists {
		store.mu.Lock()
		if _, doubleCheck := store.buckets[data.Token]; !doubleCheck {
			store.buckets[data.Token] = &TokenBucket{Prices: make([]PriceData, 0)}
		}
		bucket = store.buckets[data.Token]
		store.mu.Unlock()
	}

	bucket.mu.Lock()
	defer bucket.mu.Unlock()
	bucket.Prices = append(bucket.Prices, data)
}

func (store *DataStore) StartAutoCleanup() {
	store.cleanupTicker = time.NewTicker(5 * time.Second) // 每5秒清理一次
	go func() {
		for range store.cleanupTicker.C {
			now := time.Now().Unix()
			store.mu.RLock()
			for _, bucket := range store.buckets {
				bucket.mu.Lock()
				// 只保留最近 retentionSec 秒内的数据
				filtered := bucket.Prices[:0]
				for _, p := range bucket.Prices {
					if now-p.ReceivedAt <= store.retentionSec {
						filtered = append(filtered, p)
					}
				}
				bucket.Prices = filtered
				bucket.mu.Unlock()
			}
			store.mu.RUnlock()
		}
	}()
}

func (store *DataStore) GetPrices(token string) []PriceData {
	store.mu.RLock()
	bucket, exists := store.buckets[token]
	store.mu.RUnlock()
	if !exists {
		return nil
	}
	bucket.mu.RLock()
	defer bucket.mu.RUnlock()
	return append([]PriceData(nil), bucket.Prices...)
}
