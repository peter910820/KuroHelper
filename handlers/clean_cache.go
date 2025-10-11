package handlers

import (
	"github.com/bwmarrin/discordgo"

	"kurohelper/cache"
	"kurohelper/utils"
)

func CleanCache(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cache.Clean()

	embed := &discordgo.MessageEmbed{
		Title:       "🔐管理員訊息",
		Color:       0xD0104C,
		Description: "刪除快取成功",
	}

	utils.InteractionEmbedRespondForSelf(s, i, embed, false)
}
