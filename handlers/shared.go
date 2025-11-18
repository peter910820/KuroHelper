package handlers

import (
	"kurohelper/cache"
	"kurohelper/utils"

	"github.com/bwmarrin/discordgo"
)

type CharacterData struct {
	Name string
	Role string
}

const (
	Played = 1 << iota
	Wish
)

// 資料分頁
func pagination[T any](result *[]T, page int, useCache bool) bool {
	resultLen := len(*result)
	expectedMin := page * 10
	expectedMax := page*10 + 10

	if !useCache || page == 0 {
		if resultLen > 10 {
			*result = (*result)[:10]
			return true
		}
		return false
	} else {
		if resultLen > expectedMax {
			*result = (*result)[expectedMin:expectedMax]
			return true
		} else {
			*result = (*result)[expectedMin:]
			return false
		}
	}
}

// 產生顯示圖片，會檢查白名單來判斷要不要顯示
func generateImage(i *discordgo.InteractionCreate, url string) *discordgo.MessageEmbedImage {
	var image *discordgo.MessageEmbedImage
	if i.GuildID != "" {
		// guild
		if _, ok := cache.GuildDiscordAllowList[i.GuildID]; ok {
			image = &discordgo.MessageEmbedImage{
				URL: url,
			}
		}
	} else {
		// DM
		userID := utils.GetUserID(i)
		if _, ok := cache.GuildDiscordAllowList[userID]; ok {
			image = &discordgo.MessageEmbedImage{
				URL: url,
			}
		}
	}
	return image
}
