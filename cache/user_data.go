package cache

import (
	"kurohelper/database"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	UserCache = make(map[string]struct{})
)

// 把有存在的User從資料庫仔入快取
//
// 目的是檢查使用者的時候不用先檢查他是否在資料庫，可以直接決定要產生User紀錄還是直接抓出資料
func InitUser() {
	var entries []database.User
	if err := database.Dbs[os.Getenv("DB_NAME")].Find(&entries).Error; err != nil {
		logrus.Fatal(err)
	}

	// 存進快取
	for _, e := range entries {
		UserCache[e.ID] = struct{}{}
	}
}
