package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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
	return nil
}
