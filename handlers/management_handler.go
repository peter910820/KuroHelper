package handlers

import (
	"github.com/bwmarrin/discordgo"

	"kurohelper/utils"
)

func CleanCache(s *discordgo.Session, i *discordgo.InteractionCreate) {
	vndbCacheMu.Lock()
	defer vndbCacheMu.Unlock()
	for k := range vndbCache {
		delete(vndbCache, k)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🔐管理員訊息",
		Color:       0xD0104C,
		Description: "刪除快取成功",
	}

	utils.InteractionEmbedRespondForSelf(s, i, embed, false)
}
