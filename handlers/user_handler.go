package handlers

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"

	"kurohelper/database"
	"kurohelper/provider/erogs"
	"kurohelper/utils"
)

func AddHasPlayedHandler(s *discordgo.Session, i *discordgo.InteractionCreate, cid *CustomID) {
	if cid != nil {
		userID := utils.GetUserID(i)
		if strings.TrimSpace(userID) != "" {
			// 先假設使用者跟games都存在

			i, err := strconv.Atoi(cid.ID)
			if err != nil {
				logrus.Fatalf("Translate failed: %v", err)
			}
			userGameErogs := database.UserGameErogs{
				UserID:      userID,
				GameErogsID: i,
				HasPlayed:   true,
				InWish:      false,
			}

			database.Dbs[os.Getenv("DB_NAME")].Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "user_id"},
					{Name: "game_erogs_id"},
				}, // composite primary key
				DoNothing: true,
			}).Create(&userGameErogs)
			if err != nil {
				log.Fatalf("failed to insert user game: %v", err)
			}
		}
		embed := &discordgo.MessageEmbed{
			Title: "加入成功！",
			Color: 0x7BA23F,
		}
		utils.EditEmbedRespond(s, i, embed, nil)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	var res *erogs.FuzzySearchGameResponse

	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	idSearch, _ := regexp.MatchString(`^e\d+$`, keyword)
	res, err = erogs.GetGameByFuzzy(keyword, idSearch)
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	messageComponent := []discordgo.MessageComponent{utils.MakeAddHasPlayedComponent("✅", utils.AddHasPlayedArgs{GameID: res.ID, BrandID: 0, ConfirmMark: true}, i)}
	messageComponent = append(messageComponent, utils.MakeAddHasPlayedComponent("❌", utils.AddHasPlayedArgs{GameID: res.ID, BrandID: 0, ConfirmMark: false}, i))
	actionsRow := utils.MakeActionsRow(messageComponent)

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: res.BrandName,
		},
		Title: fmt.Sprintf("**%s(%s)**", res.Gamename, res.SellDay),
		URL:   res.Shoukai,
		Color: 0x7BA23F,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "發行機種",
				Value:  res.Model,
				Inline: false,
			},
			{
				Name:   "確認",
				Value:  "你確定要加入已玩嗎?",
				Inline: false,
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: res.BannerUrl,
		},
	}
	utils.InteractionEmbedRespond(s, i, embed, actionsRow, true)
}
