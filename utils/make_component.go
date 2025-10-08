package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func MakeActionsRow(messageComponent []discordgo.MessageComponent) *discordgo.ActionsRow {
	if len(messageComponent) != 0 {
		return &discordgo.ActionsRow{
			Components: messageComponent,
		}
	} else {
		return nil
	}

}

func MakePageComponent(label string, commandName string, id string, value int) *discordgo.Button {
	return &discordgo.Button{
		Label:    label,
		Style:    discordgo.PrimaryButton,
		CustomID: fmt.Sprintf("%s::%s::Page::%d", commandName, id, value),
	}
}

func MakeAddHasPlayedComponent(label string, commandName string, id int, value bool) *discordgo.Button {
	return &discordgo.Button{
		Label:    label,
		Style:    discordgo.PrimaryButton,
		CustomID: fmt.Sprintf("%s::%d::_::%t", commandName, id, value),
	}
}
