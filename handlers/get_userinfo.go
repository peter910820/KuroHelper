package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	kurohelperdb "github.com/peter910820/kurohelper-db/v2"
	"github.com/sirupsen/logrus"

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
	user, err := kurohelperdb.GetUser(userID)
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	// Game資料
	userGames, err := kurohelperdb.GetUserData(userID)
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	// Brand資料統計
	brandData, err := kurohelperdb.BrandCount(userID, true, false)
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	hasPlayedCount := 0
	inWishCount := 0
	listHasPlayed := make([]string, 0, 10)
	listInWish := make([]string, 0, 10)
	for _, r := range userGames {
		if r.HasPlayed {
			hasPlayedCount++
			if hasPlayedCount <= 10 {
				t := getUserPlayRecordTime(&r)
				if t != "" {
					listHasPlayed = append(listHasPlayed, fmt.Sprintf("%d. **%s**/*⏱️%s*", hasPlayedCount, r.GameErogs.Name, t))
				} else {
					listHasPlayed = append(listHasPlayed, fmt.Sprintf("%d. **%s**", hasPlayedCount, r.GameErogs.Name))
				}
			}
		}

		if r.InWish {
			inWishCount++
			if inWishCount <= 10 {
				listInWish = append(listInWish, fmt.Sprintf("%d. %s", inWishCount, r.GameErogs.Name))
			}
		}
	}

	listData := make([]string, 0, len(brandData))
	for i, b := range brandData {
		if i >= 5 { // 已經到第六筆，直接跳出
			break
		}
		if i <= 4 {
			star := strings.Repeat("⭐", 5-i)
			listData = append(listData, fmt.Sprintf("%s\n**%s: (%d)**", star, b.BrandName, b.Count))
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("**%s 的個人資料**", userName),
		Color:       0xB481BB,
		Description: fmt.Sprintf("資料建檔日期: %s", user.CreatedAt.Format("2006-01-02")),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("最近遊玩完畢(%d)", hasPlayedCount),
				Value:  strings.Join(listHasPlayed, "\n"),
				Inline: true,
			},
			{
				Name:   fmt.Sprintf("最近收藏(%d)", inWishCount),
				Value:  strings.Join(listInWish, "\n"),
				Inline: true,
			},
			{
				Name:   "玩過最多(公司品牌)",
				Value:  strings.Join(listData, "\n"),
				Inline: true,
			},
		},
	}
	utils.InteractionEmbedRespond(s, i, embed, nil, true)
}

func getUserPlayRecordTime(r *kurohelperdb.UserGameErogs) string {
	if r.CompletedAt != nil {
		return r.CompletedAt.Format("2006-01-02")
	}
	return ""
}
