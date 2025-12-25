package bootstrap

import (
	"kurohelper/utils"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/peter910820/kurohelper-core/cache"
	"github.com/peter910820/kurohelper-core/erogs"
	"github.com/peter910820/kurohelper-core/seiya"
	"github.com/peter910820/kurohelper-core/ymgal"

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
	kurohelperdb.Migration(kurohelperdb.Dbs) // 選填

	// 將白名單存成快取
	cache.InitAllowList()

	// init ZhtwToJp var
	cache.InitZhtwToJp()

	cache.InitSeiyaCorrespond()

	// erogs rate limit init
	erogs.InitRateLimit(time.Duration(utils.GetEnvInt("EROGS_RATE_LIMIT_RESET_TIME", 10)))

	// seiya init
	err = seiya.Init()
	if err != nil {
		logrus.Fatal(err)
	}

	// ymgal init
	err = ymgalInit()
	if err != nil {
		logrus.Fatal(err)
	}

	cache.InitUser()
	// 掛載自動清除快取job
	go cache.CleanCacheJob(360, stopChan)
}

// ymgal init
func ymgalInit() error {
	// init config
	ymgal.Init(os.Getenv("YMGAL_ENDPOINT"), os.Getenv("YMGAL_CLIENT_ID"), os.Getenv("YMGAL_CLIENT_SECRET"))

	// init token
	// ymgal init token
	err := ymgal.GetToken()
	if err != nil {
		return err
	}
	return nil
}
