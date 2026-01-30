package cache

import (
	"sync"
	"time"

	"kurohelper-core/erogs"

	kurohelpercore "kurohelper-core"
)

// CID快取
//
// 對每一次查詢建立CID以及關鍵字的關聯
// 因為CID不允許過長字元，所以遇到很長的關鍵字時會直接丟錯，所以才需要這層快取
var (
	CIDStore = NewCacheStoreV2[string](time.Hour)
)

// 批評空間快取
var (
	// 使用關鍵字Base64作為鍵
	ErogsGameListStore = NewCacheStoreV2[[]erogs.FuzzySearchListResponse](2 * time.Hour)
	// 使用批評空間ID作為鍵
	ErogsGameStore = NewCacheStoreV2[*erogs.FuzzySearchGameResponse](2 * time.Hour)
	// 使用關鍵字Base64作為鍵
	ErogsMusicListStore = NewCacheStoreV2[[]erogs.MusicList](2 * time.Hour)
	// 使用批評空間ID作為鍵
	ErogsMusicStore = NewCacheStoreV2[*erogs.Music](2 * time.Hour)
)

// VNDB快取
// var (
// 	VndbGameListStore = NewCacheStoreV2[*vndb.ProducerSearchResponse](time.Hour)
// 	VndbGameStore     = NewCacheStoreV2[*vndb.BasicResponse[vndb.GetVnUseIDResponse]](time.Hour)
// )

// 月幕快取
// var (
// 	YmgalGame = NewCacheStoreV2[*ymgal.SearchGameResp](time.Hour)
// )

// cache struct
type CacheV2[T any] struct {
	Value    T
	ExpireAt time.Time
}

// CacheStoreV2 泛型快取儲存
// data的鍵值為UUID
type CacheStoreV2[T any] struct {
	data       map[string]*CacheV2[T]
	expireTime time.Duration
	mu         sync.RWMutex
}

// NewCacheStore 建立新的快取儲存
func NewCacheStoreV2[T any](expireTime time.Duration) *CacheStoreV2[T] {
	return &CacheStoreV2[T]{
		data:       make(map[string]*CacheV2[T]),
		expireTime: expireTime,
	}
}

// Set 設定快取
func (c *CacheStoreV2[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &CacheV2[T]{
		Value:    value,
		ExpireAt: time.Now().Add(c.expireTime),
	}
}

// Get 從快取中取得資料
func (c *CacheStoreV2[T]) Get(key string) (T, error) {
	c.mu.RLock()
	item, ok := c.data[key]
	c.mu.RUnlock()

	// 不存在或已過期
	if !ok || time.Now().After(item.ExpireAt) {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		var zero T
		return zero, kurohelpercore.ErrCacheLost
	}

	return item.Value, nil
}

// Clean 清除過期快取
func (c *CacheStoreV2[T]) Clean() (deleteCount int, total int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	total = len(c.data)
	for k, d := range c.data {
		if time.Now().After(d.ExpireAt) {
			delete(c.data, k)
			deleteCount++
		}
	}

	return
}

// CleanAll 清除所有快取
func (c *CacheStoreV2[T]) CleanAll() (count int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range c.data {
		delete(c.data, k)
		count++
	}

	return
}
