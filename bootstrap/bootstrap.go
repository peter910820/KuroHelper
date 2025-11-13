package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"kurohelper/cache"
	"kurohelper/provider/erogs"
	"kurohelper/provider/seiya"
	"kurohelper/provider/ymgal"

	kurohelperdb "github.com/peter910820/kurohelper-db/v2"
)

// 啟動函式
func Init(stopChan <-chan struct{}) {
	// logrus settings
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
	// load .env
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Fatal(err)
	}

	config := kurohelperdb.Config{
		DBOwner:    os.Getenv("DB_OWNER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
	}

	err = kurohelperdb.InitDsn(config)
	if err != nil {
		logrus.Fatal(err)
	}
	// kurohelperdb.Migration() // 選填

	// 將白名單存成快取
	cache.InitAllowList()

	// init ZhtwToJp var
	cache.InitZhtwToJp()

	cache.InitSeiyaCorrespond()

	// erogs rate limit init
	erogs.InitRateLimit()

	// seiya init
	err = seiya.Init()
	if err != nil {
		logrus.Fatal(err)
	}

	// ymgal init token
	err = ymgal.GetToken()
	if err != nil {
		logrus.Fatal(err)
	}

	cache.InitUser()
	// 掛載自動清除快取job
	go cache.CleanCacheJob(360, stopChan)
}
