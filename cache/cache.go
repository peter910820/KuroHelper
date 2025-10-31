package cache

import (
	"sync"
	"time"

	kurohelperdb "github.com/peter910820/kurohelper-db"
	"github.com/peter910820/kurohelper-db/models"
	"github.com/sirupsen/logrus"

	kurohelpererrors "kurohelper/errors"
)

// cache struct
type Cache struct {
	Value    any
	ExpireAt time.Time
}

var (
	commonCache   = make(map[string]*Cache)
	commonCacheMu sync.RWMutex

	ZhtwToJp        map[rune]rune
	SeiyaCorrespond map[string]string
)

func Set(key string, value any) {
	commonCacheMu.Lock()
	defer commonCacheMu.Unlock()
	commonCache[key] = &Cache{
		Value:    value,
		ExpireAt: time.Now().Add(10 * time.Minute),
	}
}

func Get(key string) (any, error) {
	commonCacheMu.RLock()
	item, ok := commonCache[key]
	commonCacheMu.RUnlock()
	// 檢查過期或不存在
	if !ok || time.Now().After(item.ExpireAt) {
		commonCacheMu.Lock()
		delete(commonCache, key)
		commonCacheMu.Unlock()
		return nil, kurohelpererrors.ErrCacheLost
	}
	return item.Value, nil
}

func Clean() {
	commonCacheMu.Lock()
	defer commonCacheMu.Unlock()
	count := 0
	for k := range commonCache {
		delete(commonCache, k)
		count++
	}
	logrus.Printf("%d筆快取已被清除", count)
}

func InitZhtwToJp() {
	var entries []models.ZhtwToJp
	if err := kurohelperdb.Dbs.Find(&entries).Error; err != nil {
		logrus.Fatal(err)
	}

	// 轉換
	ZhtwToJp = make(map[rune]rune, len(entries))
	for _, e := range entries {
		keyRunes := []rune(e.ZhTw)
		valRunes := []rune(e.Jp)

		// 確保都是單一字
		if len(keyRunes) == 1 && len(valRunes) == 1 {
			ZhtwToJp[keyRunes[0]] = valRunes[0]
		}
	}
}

func InitSeiyaCorrespond() {
	var entries []models.SeiyaCorrespond
	if err := kurohelperdb.Dbs.Find(&entries).Error; err != nil {
		logrus.Fatal(err)
	}

	// 轉換
	SeiyaCorrespond = make(map[string]string, len(entries))
	for _, e := range entries {
		SeiyaCorrespond[e.GameName] = e.SeiyaURL
	}
}
