package bootstrap

import (
	"kurohelper/cache"
	"kurohelper/store"
	"kurohelper/utils"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/kuro-helper/kurohelper-core/v3/erogs"
	"github.com/kuro-helper/kurohelper-core/v3/seiya"
	corestore "github.com/kuro-helper/kurohelper-core/v3/store"
	"github.com/kuro-helper/kurohelper-core/v3/ymgal"

	kurohelperdb "github.com/kuro-helper/kurohelper-db/v3"
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
	store.InitAllowList()

	// init ZhtwToJp var
	corestore.InitZhtwToJp()

	corestore.InitSeiyaCorrespond()

	// erogs rate limit init
	erogs.InitRateLimit(time.Duration(utils.GetEnvInt("EROGS_RATE_LIMIT_RESET_TIME", 10)))

	// seiya init
	err = seiya.Init()
	if err != nil {
		logrus.Fatal(err)
	}

	// ymgal init
	if strings.EqualFold(os.Getenv("INIT_YMGAL"), "true") {
		err = ymgalInit()
		if err != nil {
			logrus.Fatal(err)
		}
	}

	store.InitUser()
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
