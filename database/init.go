package database

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDsn(choiceDb int) (string, *gorm.DB) {
	var dbname string
	switch choiceDb {
	case 0:
		dbname = os.Getenv("DB_NAME")
	default:
		logrus.Fatal("error in init dsn function")
	}

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_OWNER"),
		os.Getenv("DB_PASSWORD"),
		dbname,
		os.Getenv("DB_PORT"))

	// get connect db variable
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("連接資料庫失敗: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatalf("無法取得 sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return dbname, db
}
