package utils

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	internalerrors "kurohelper/errors"
)

// handle interaction command common respond
func InteractionRespond(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

// handle interaction command embed respond
func InteractionEmbedRespond(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed, editFlag bool) {
	if editFlag {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		})
		if err != nil {
			logrus.Error(err)
		}
		InteractionRespond(s, i, "該功能目前異常，請稍後再嘗試")
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
	}

}

// get slash command options
func GetOptions(i *discordgo.InteractionCreate, name string) (string, error) {
	for _, v := range i.ApplicationCommandData().Options {
		if v.Name == name {
			value, ok := v.Value.(string) // type assertion
			if !ok {
				return "", internalerrors.ErrOptionTranslateFail
			} else {
				return value, nil
			}
		}
	}
	return "", internalerrors.ErrOptionNotFound
}
