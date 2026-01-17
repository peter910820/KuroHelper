package cache

import (
	"time"

	"github.com/sirupsen/logrus"
)

func CleanCacheJob(minute time.Duration, stopChan <-chan struct{}) {
	logrus.Print("CleanCacheJob 正在啟動...")
	ticker := time.NewTicker(minute * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 清除Cache
			logrus.Printf("%d筆SearchCache快取已被清除", SearchCache.Clean())
			logrus.Printf("%d筆UserInfoCache快取已被清除", UserInfoCache.Clean())
			// logrus.Printf("%d筆SubmitDataCache快取已被清除", SubmitDataCache.Clean())
			logrus.Printf("%d筆SearchBrandCache快取已被清除", SearchBrandCache.Clean())
		case <-stopChan:
			logrus.Println("CleanCacheJob 正在關閉...")
			return
		}
	}
}
