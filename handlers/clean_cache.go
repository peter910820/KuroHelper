package handlers

import (
	"github.com/bwmarrin/discordgo"

	"kurohelper/utils"

	"kurohelper/cache"
)

// 清除快取Handler
func CleanCache(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cache.UserInfoCache.Clean()

	embed := &discordgo.MessageEmbed{
		Title:       "🔐管理員訊息",
		Color:       0xD0104C,
		Description: "刪除快取成功",
	}

	utils.InteractionEmbedRespondForSelf(s, i, embed, nil, false)
}
