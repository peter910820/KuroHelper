package handlers

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"kurohelper/cache"
	"kurohelper/database"
	"kurohelper/provider/erogs"
	"kurohelper/utils"
)

func AddHasPlayed(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.NewCID) {
	if cid != nil {
		addWishCID := utils.AddWishCID{
			NewCID: *cid,
		}

		confirmMark, err := addWishCID.GetConfirmMark()
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		if !confirmMark {
			embed := &discordgo.MessageEmbed{
				Title: "操作已取消",
				Color: 0x7BA23F,
			}
			utils.EditEmbedRespond(s, i, embed, nil)
			return
		}
		// get cache
		cacheValue, err := cache.Get(addWishCID.GetCacheID())
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		resValue := cacheValue.(erogs.FuzzySearchGameResponse)
		res := &resValue

		userID := utils.GetUserID(i)
		userName := utils.GetUsername(i)
		if strings.TrimSpace(userID) != "" && strings.TrimSpace(userName) != "" {
			var user database.User
			var gameErogs database.GameErogs
			var brandErogs database.BrandErogs
			err := database.Dbs[os.Getenv("DB_NAME")].Transaction(func(tx *gorm.DB) error {
				// 1. 確保 User 存在
				if err := tx.Where("id = ?", userID).FirstOrCreate(&user, database.User{ID: userID, Name: userName}).Error; err != nil {
					return err
				}

				// 2. 確保 Brand 存在
				if err := tx.Where("id = ?", res.BrandID).FirstOrCreate(&brandErogs, database.BrandErogs{ID: res.BrandID, Name: res.BrandName}).Error; err != nil {
					return err
				}

				// 3. 確保 Game 存在
				if err := tx.Where("id = ?", res.ID).FirstOrCreate(&gameErogs, database.GameErogs{ID: res.ID, Name: res.Gamename, BrandErogsID: res.BrandID}).Error; err != nil {
					return err
				}

				// 4. 建立 UserGame
				ug := database.UserGameErogs{UserID: user.ID, GameErogsID: res.ID, HasPlayed: true, InWish: false}
				if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&ug).Error; err != nil {
					return err
				}

				return nil // commit
			})
			if err != nil {
				logrus.Fatal(err)
			}

			if _, ok := cache.UserCache[userID]; !ok {
				cache.UserCache[userID] = struct{}{}
			}

			embed := &discordgo.MessageEmbed{
				Title: "加入成功！",
				Color: 0x7BA23F,
			}
			utils.EditEmbedRespond(s, i, embed, nil)
			return
		} else {
			embed := &discordgo.MessageEmbed{
				Title: "找不到使用者！",
				Color: 0x7BA23F,
			}
			utils.EditEmbedRespond(s, i, embed, nil)
			return
		}
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

	idStr := uuid.New().String()
	cache.Set(idStr, *res)

	messageComponent := []discordgo.MessageComponent{utils.MakeCIDAddHasPlayedComponent("✅", utils.AddHasPlayedArgs{CacheID: idStr, ConfirmMark: true}, i)}
	messageComponent = append(messageComponent, utils.MakeCIDAddHasPlayedComponent("❌", utils.AddHasPlayedArgs{CacheID: idStr, ConfirmMark: false}, i))
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
