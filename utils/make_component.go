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

func MakeCIDAddHasPlayedComponent(label string, data AddHasPlayedArgs, i *discordgo.InteractionCreate) *discordgo.Button {
	return &discordgo.Button{
		Label:    label,
		Style:    discordgo.PrimaryButton,
		CustomID: fmt.Sprintf("%s|%d|%s|%t", i.ApplicationCommandData().Name, CustomIDTypeAddHasPlayed, data.CacheID, data.ConfirmMark),
	}
}

func MakeCIDPageComponent(label string, id string, value int, isList bool, commandName string) *discordgo.Button {
	listString := ""
	if isList {
		listString = "list"
	}
	return &discordgo.Button{
		Label:    label,
		Style:    discordgo.PrimaryButton,
		CustomID: fmt.Sprintf("%s|%d|%s|%d", commandName+"/"+listString, CustomIDTypePage, id, value),
	}
}
