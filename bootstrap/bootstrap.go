package bootstrap

import (
	"kurohelper/cache"
	"kurohelper/database"
	"kurohelper/erogs"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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

	// 初始化erogs速率鎖
	erogs.InitRateLimit()
}
