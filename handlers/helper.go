package handlers

import (
	"discordbot/utils"

	"github.com/bwmarrin/discordgo"
)

func Helper(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title: "幫助",
		Color: 0xF19483,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "使用說明/文檔",
				Value:  "https://docs.kurohelper.com/docs",
				Inline: true,
			},
			{
				Name:   "邀請至伺服器",
				Value:  "https://discord.com/oauth2/authorize?client_id=1418225729241612298&permissions=3941734153714752&integration_type=0&scope=bot",
				Inline: true,
			},
			{
				Name:   "聯繫我們/加入群組",
				Value:  "https://discord.gg/6rkrm7tsXr",
				Inline: true,
			},
		},
	}
	utils.InteractionEmbedRespond(s, i, embed, nil, false)
}
