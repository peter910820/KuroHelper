package main

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"kurohelper/bot"
	"kurohelper/config"
)

func main() {
	// 初始化config
	if err := config.Init(); err != nil {
		logrus.Fatal(err)
	}

	token := os.Getenv("BOT_TOKEN")
	kuroHelper, err := discordgo.New("Bot " + token)
	if err != nil {
		logrus.Fatal(err)
	}

	err = kuroHelper.Open() // websocket connect
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("KuroHelper is now running. Press CTRL+C to exit.")

	kuroHelper.AddHandler(bot.Ready)
	kuroHelper.AddHandler(bot.OnInteraction)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	interruptSignal := <-c
	kuroHelper.Close() // websocket disconnect
	logrus.Debug(interruptSignal)
}
