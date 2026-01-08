package main

import (
	"kurobidder/letao"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"kurohelper/bootstrap"
	"kurohelper/bot"
)

func main() {
	// 基本初始化
	bootstrap.BasicInit()

	token := os.Getenv("BOT_TOKEN")
	kuroHelper, err := discordgo.New("Bot " + token)
	if err != nil {
		logrus.Fatal(err)
	}

	kuroHelper.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	logrus.Info("KuroHelper is now running. Press CTRL+C to exit.")

	kuroHelper.AddHandler(bot.Ready)
	kuroHelper.AddHandler(bot.OnInteraction)

	err = kuroHelper.Open() // websocket connect
	if err != nil {
		logrus.Fatal(err)
	}

	// 初始化專案作業
	stopChan := make(chan struct{})
	kurobidderDataChan := make(chan []letao.AuctionItem, 2)
	bootstrap.Init(stopChan, kurobidderDataChan)

	// 監聽 kurobidder 爬蟲數據並發送到 Discord
	go func() {
		for {
			select {
			case items := <-kurobidderDataChan:
				if len(items) > 0 {
					bootstrap.SendKurobidderDataToDiscord(kuroHelper, items)
				}
			case <-stopChan:
				return
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	interruptSignal := <-c
	logrus.Debug(interruptSignal)

	// 關閉 jobs
	close(stopChan)

	kuroHelper.Close() // websocket disconnect
}
