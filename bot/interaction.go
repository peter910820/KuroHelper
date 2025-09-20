package bot

import (
	"github.com/bwmarrin/discordgo"

	"kurohelper/handlers"
)

func OnInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "vndb統計資料":
		go handlers.VndbStats(s, i)
	case "vndb查詢指定遊戲":
		go handlers.VndbSearchGameByID(s, i)
	case "vndb模糊查詢品牌":
		go handlers.VndbFuzzySearchBrand(s, i)
	}
}
