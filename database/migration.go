package database

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

func Migration(dbName string, db *gorm.DB) {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf(".env file error: %v", err)
	}

	switch dbName {
	case os.Getenv("DB_NAME"):
		db.AutoMigrate(&ZhtwToJp{})
	default:
		logrus.Fatal("error in migration function")
	}
}
