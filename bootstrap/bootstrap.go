package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"kurohelper/cache"
	"kurohelper/provider/erogs"
	"kurohelper/provider/seiya"
	"kurohelper/provider/ymgal"

	kurohelperdb "github.com/peter910820/kurohelper-db"
	"github.com/peter910820/kurohelper-db/models"
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

	config := models.Config{
		DBOwner:    os.Getenv("DB_OWNER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
	}

	kurohelperdb.InitDsn(config)
	kurohelperdb.Migration()

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
	go cache.CleanCacheJob(240, stopChan)
}
