package cache

import (
	"os"
	"sync"
	"time"

	"kurohelper/database"
	kurohelpererrors "kurohelper/errors"

	"github.com/sirupsen/logrus"
)

// cache struct
type Cache struct {
	Value    any
	ExpireAt time.Time
}

var (
	commonCache   = make(map[string]*Cache)
	commonCacheMu sync.RWMutex

	ZhtwToJp map[rune]rune
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
	for k := range commonCache {
		delete(commonCache, k)
	}
}

func InitZhtwToJp() {
	var entries []database.ZhtwToJp
	if err := database.Dbs[os.Getenv("DB_NAME")].Find(&entries).Error; err != nil {
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
