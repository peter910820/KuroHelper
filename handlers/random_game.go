package handlers

import (
	"github.com/bwmarrin/discordgo"

	"kurohelper/cache"
	"kurohelper/provider/ymgal"
	"kurohelper/utils"
)

// 隨機遊戲Handler
func RandomGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 長時間查詢
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	game, err := ymgal.GetRandomGame()
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	title := game[0].Name
	if game[0].HaveChinese {
		title += "/" + game[0].ChineseName
	}

	var image *discordgo.MessageEmbedImage
	if i.GuildID != "" {
		// guild
		if _, ok := cache.GuildDiscordAllowList[i.GuildID]; ok {
			image = &discordgo.MessageEmbedImage{
				URL: "https://store.ymgal.games/" + game[0].MainImg,
			}
		}
	} else {
		// DM
		userID := utils.GetUserID(i)
		if _, ok := cache.GuildDiscordAllowList[userID]; ok {
			image = &discordgo.MessageEmbedImage{
				URL: "https://store.ymgal.games/" + game[0].MainImg,
			}
		}
	}

	embed := &discordgo.MessageEmbed{
		Title: title,
		Color: 0x261E47,
		Image: image,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "發售日",
				Value:  game[0].ReleaseDate,
				Inline: false,
			},
		},
	}

	utils.InteractionEmbedRespond(s, i, embed, nil, true)
}
