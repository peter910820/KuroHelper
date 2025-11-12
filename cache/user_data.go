package cache

import (
	"github.com/sirupsen/logrus"

	kurohelperdb "github.com/peter910820/kurohelper-db/v2"
)

var (
	UserCache = make(map[string]struct{})
)

// 把有存在的User從資料庫載入快取
//
// 目的是檢查使用者的時候不用先檢查他是否在資料庫，可以直接決定要產生User紀錄還是直接抓出資料
func InitUser() {
	user, err := kurohelperdb.GetAllUser()
	if err != nil {
		logrus.Fatal(err)
	}

	// 存進快取
	for _, e := range user {
		UserCache[e.ID] = struct{}{}
	}
}
