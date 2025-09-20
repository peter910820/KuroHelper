package bot

import (
	"github.com/bwmarrin/discordgo"

	"kurohelper/handlers"
)

func OnInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "VNDB統計資料":
		go handlers.VndbStats(s, i)
	case "VNDB查詢遊戲":
		go handlers.VndbSearchGame(s, i)
	}
}
