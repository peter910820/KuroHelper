package bot

import (
	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateGameStatus(0, "さくら、もゆ。-as the Night's, Reincarnation-")
	RegisterCommand(s)
}
