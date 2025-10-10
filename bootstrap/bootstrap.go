package bootstrap

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"kurohelper/cache"
	"kurohelper/database"
	"kurohelper/provider/erogs"
	"kurohelper/provider/seiya"
	"kurohelper/provider/ymgal"
)

// 啟動函式
func Init() {
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

	// init migration
	for i := range 1 {
		dbName, db := database.InitDsn(i)
		database.Dbs[dbName] = db
		database.Migration(dbName, database.Dbs[dbName])
	}

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
}
