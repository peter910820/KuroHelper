package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"kurohelper/database"
	"kurohelper/models"
)

var (
	Dbs      = make(map[string]*gorm.DB)
	ZhtwToJp map[rune]rune
)

func Init() error {
	// logrus settings
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
	// load .env
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	// init migration
	for i := range 1 {
		dbName, db := database.InitDsn(i)
		Dbs[dbName] = db
		database.Migration(dbName, Dbs[dbName])
	}

	var entries []models.ZhtwToJp
	if err := Dbs[os.Getenv("DB_NAME")].Find(&entries).Error; err != nil {
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

	return nil
}
