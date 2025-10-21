package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	kurohelperdb "github.com/peter910820/kurohelper-db"
	"github.com/peter910820/kurohelper-db/models"

	"kurohelper/utils"
)

type BrandCount struct {
	BrandID   int
	BrandName string
	Count     int
}

func GetUserinfo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 長時間查詢
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	userID := utils.GetUserID(i)
	userName := utils.GetUsername(i)

	// User資料
	var user models.User
	err := kurohelperdb.Dbs.First(&user, "id = ?", userID).Error
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	// Brand資料統計
	var brandData []BrandCount
	err = kurohelperdb.Dbs.
		Table("user_game_erogs AS uge").
		Select("b.id AS brand_id, b.name AS brand_name, COUNT(*) AS count").
		Joins("JOIN game_erogs AS g ON uge.game_erogs_id = g.id").
		Joins("JOIN brand_erogs AS b ON g.brand_erogs_id = b.id").
		Where("uge.user_id = ?", userID).
		Group("b.id, b.name").
		Order("count DESC").
		Scan(&brandData).Error
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	listData := make([]string, 0, len(brandData))
	for i, b := range brandData {
		if i <= 2 {
			star := strings.Repeat("⭐", 3-i)
			listData = append(listData, fmt.Sprintf("%s**%s: %d部**", star, b.BrandName, b.Count))
		} else {
			listData = append(listData, fmt.Sprintf("**%s: %d部**", b.BrandName, b.Count))
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("**%s 的個人資料**", userName),
		Color:       0xB481BB,
		Description: fmt.Sprintf("資料建檔日期: %s", user.CreatedAt.Format("2006-01-02")),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "玩過的作品(公司品牌)統計",
				Value:  strings.Join(listData, "\n"),
				Inline: false,
			},
		},
	}
	utils.InteractionEmbedRespond(s, i, embed, nil, true)
}
