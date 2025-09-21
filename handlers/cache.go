package handlers

import (
	"sync"
	"time"

	"kurohelper/models"
)

var (
	vndbCache   = make(map[string]*models.Cache)
	vndbCacheMu sync.RWMutex
)

func SetCache(key string, value interface{}) {
	vndbCacheMu.Lock()
	defer vndbCacheMu.Unlock()
	vndbCache[key] = &models.Cache{
		Value:    value,
		ExpireAt: time.Now().Add(3),
	}
}

func GetCache(key string) (interface{}, bool) {
	vndbCacheMu.RLock()
	item, ok := vndbCache[key]
	vndbCacheMu.RUnlock()
	// 檢查過期或不存在
	if !ok || time.Now().After(item.ExpireAt) {
		vndbCacheMu.Lock()
		delete(vndbCache, key)
		vndbCacheMu.Unlock()
		return nil, false
	}
	return item.Value, true
}
