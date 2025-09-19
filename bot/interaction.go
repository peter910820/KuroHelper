package bot

import (
	"kurohelper/vndb"

	"github.com/bwmarrin/discordgo"
)

func OnInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "vndb統計資料":
		go vndb.GetStats(s, i)
	}
}
