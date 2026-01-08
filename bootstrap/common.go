package bootstrap

import (
	"fmt"
	"kurobidder/letao"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func SendKurobidderDataToDiscord(s *discordgo.Session, items []letao.AuctionItem) {
	channelID := os.Getenv("LETAO_CHANNEL_ID")
	if len(items) == 0 {
		return
	}

	// 一筆一筆發送
	for _, item := range items {
		itemInfo := fmt.Sprintf("**目前出價: %s/%s**\n**出價次數:** %s\n**剩餘時間:** %s\n**商品連結:** [點擊查看](%s)",
			item.PriceMP, item.PriceM, item.BidsInfo, item.TimeInfo, item.URL)

		embed := &discordgo.MessageEmbed{
			Title: item.Title,
			Color: 0x261E47,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "商品資訊:",
					Value:  itemInfo,
					Inline: false,
				},
			},
		}

		// 使用商品圖片
		if item.ImageURL != "" {
			embed.Image = &discordgo.MessageEmbedImage{
				URL: item.ImageURL,
			}
		}

		_, err := s.ChannelMessageSendEmbed(channelID, embed)
		if err != nil {
			logrus.Error(err)
		}
	}
}
