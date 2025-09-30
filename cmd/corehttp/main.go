package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
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

	kuroHelper.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	logrus.Info("KuroHelper is now running. Press CTRL+C to exit.")

	kuroHelper.AddHandler(bot.Ready)
	kuroHelper.AddHandler(bot.OnInteraction)

	err = kuroHelper.Open() // websocket connect
	if err != nil {
		logrus.Fatal(err)
	}
	defer kuroHelper.Close() // defer websocket disconnect

	// Fiber server
	app := fiber.New()
	app.Post("/github-actions", func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth != "Bearer "+token {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		eventName := c.FormValue("event_name")
		if eventName == "push" {
			_ = PushSend(kuroHelper, c)
			// if err != nil {
			// 	return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			// }
		}

		return c.JSON(fiber.Map{"status": "ok"})
	})

	fiberDone := make(chan struct{})

	go func() {
		if err := app.Listen(":3000"); err != nil {
			logrus.Println("Fiber shutdown:", err)
		}
		close(fiberDone)
		logrus.Println("Fiber close success")
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	interruptSignal := <-c
	logrus.Debug(interruptSignal)

	// 優雅關閉 Fiber server
	if err := app.Shutdown(); err != nil {
		logrus.Println("Fiber shutdown error:", err)
	}

	// 等 Fiber goroutine 關閉
	<-fiberDone
}

func PushSend(kuroHelper *discordgo.Session, c *fiber.Ctx) error {
	branch := c.FormValue("branch")
	hash := c.FormValue("hash")
	fullHash := c.FormValue("full_hash")
	authorEmail := c.FormValue("author_email")
	authorName := c.FormValue("author_name")
	date := c.FormValue("date")
	subject := c.FormValue("subject")
	body := c.FormValue("body")

	embed := &discordgo.MessageEmbed{
		Title:       "KuroHelper Push Event",
		Color:       0xaf5f3c,
		Description: fmt.Sprintf("[%s](https://github.com/peter910820/KuroHelper/commit/%s)  %s\n%s", hash, fullHash, branch, date),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "發送人",
				Value:  authorName + "/" + authorEmail,
				Inline: false,
			},
			{
				Name:   "主旨",
				Value:  subject,
				Inline: false,
			},
			{
				Name:   "內容",
				Value:  body,
				Inline: false,
			},
		},
	}

	_, err := kuroHelper.ChannelMessageSendEmbed(os.Getenv("BOT_CHANNEL_ID"), embed)
	if err != nil {
		return err
	}
	return nil
}
