package bot

import (
	"github.com/bwmarrin/discordgo"
)

func OnInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	}
}
