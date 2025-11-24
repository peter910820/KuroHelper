package handlers

import (
	"github.com/bwmarrin/discordgo"

	"kurohelper/utils"

	"github.com/peter910820/kurohelper-core/cache"
)

// 清除快取Handler
func CleanCache(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cache.Clean()

	embed := &discordgo.MessageEmbed{
		Title:       "🔐管理員訊息",
		Color:       0xD0104C,
		Description: "刪除快取成功",
	}

	utils.InteractionEmbedRespondForSelf(s, i, embed, nil, false)
}
