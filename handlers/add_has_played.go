package handlers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	kurohelperdb "github.com/peter910820/kurohelper-db"
	"github.com/peter910820/kurohelper-db/models"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"kurohelper/cache"
	kurohelpererrors "kurohelper/errors"
	"kurohelper/provider/erogs"
	"kurohelper/utils"
)

// 加已玩Handler
func AddHasPlayed(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.NewCID) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	if cid != nil {
		addHasPlayedCID := utils.AddHasPlayedCID{
			NewCID: *cid,
		}

		completeDate, err := addHasPlayedCID.GetCompleteDate()
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		// get cache
		cacheValue, err := cache.Get(addHasPlayedCID.GetCacheID())
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		resValue := cacheValue.(erogs.FuzzySearchGameResponse)
		res := &resValue

		userID := utils.GetUserID(i)
		userName := utils.GetUsername(i)
		if strings.TrimSpace(userID) != "" && strings.TrimSpace(userName) != "" {
			var msg string
			var user models.User
			var gameErogs models.GameErogs
			var brandErogs models.BrandErogs
			err := kurohelperdb.Dbs.Transaction(func(tx *gorm.DB) error {
				// 1. 確保 User 存在
				if err := tx.Where("id = ?", userID).FirstOrCreate(&user, models.User{ID: userID, Name: userName}).Error; err != nil {
					return err
				}

				// 2. 確保 Brand 存在
				if err := tx.Where("id = ?", res.BrandID).FirstOrCreate(&brandErogs, models.BrandErogs{ID: res.BrandID, Name: res.BrandName}).Error; err != nil {
					return err
				}

				// 3. 確保 Game 存在
				if err := tx.Where("id = ?", res.ID).FirstOrCreate(&gameErogs, models.GameErogs{ID: res.ID, Name: res.Gamename, BrandErogsID: res.BrandID}).Error; err != nil {
					return err
				}

				// 4. 建立 UserGame
				if completeDate.IsZero() {
					ug := models.UserGameErogs{UserID: user.ID, GameErogsID: res.ID, HasPlayed: true, InWish: false}
					result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&ug)
					if result.Error != nil {
						return result.Error
					}
					if result.RowsAffected == 0 {
						msg = "資料已建立，本次動作無效"
					}
				} else {
					ug := models.UserGameErogs{UserID: user.ID, GameErogsID: res.ID, HasPlayed: true, InWish: false, CompletedAt: &completeDate}
					result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&ug)
					if result.Error != nil {
						return result.Error
					}
					if result.RowsAffected == 0 {
						msg = "資料已建立，本次動作無效"
					}
				}

				return nil // commit
			})
			if err != nil {
				logrus.Fatal(err)
			}

			if _, ok := cache.UserCache[userID]; !ok {
				cache.UserCache[userID] = struct{}{}
			}

			if msg == "" {
				msg = "加入成功！"
			}

			embed := &discordgo.MessageEmbed{
				Title: msg,
				Color: 0x7BA23F,
			}
			utils.InteractionEmbedRespondForSelf(s, i, embed, nil, true)
			return
		} else {
			embed := &discordgo.MessageEmbed{
				Title: "找不到使用者！",
				Color: 0x7BA23F,
			}
			utils.InteractionEmbedRespondForSelf(s, i, embed, nil, true)
			return
		}
	}

	var res *erogs.FuzzySearchGameResponse

	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	completeDate, err := utils.GetOptions(i, "complete_date")
	if err != nil && !errors.Is(err, kurohelpererrors.ErrOptionNotFound) {
		utils.HandleError(err, s, i)
		return
	}

	var t time.Time
	if completeDate != "" {
		t, err = utils.ParseYYYYMMDD(completeDate)
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		if t.After(time.Now().AddDate(0, 0, 1)) {
			utils.HandleError(err, s, i)
			return
		}
	}

	idSearch, _ := regexp.MatchString(`^e\d+$`, keyword)
	res, err = erogs.GetGameByFuzzy(keyword, idSearch)
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	idStr := uuid.New().String()
	cache.Set(idStr, *res)

	cidCommandName := utils.MakeCIDCommandName(i.ApplicationCommandData().Name, false, "")
	messageComponent := []discordgo.MessageComponent{utils.MakeCIDAddHasPlayedComponent("✅", idStr, t, cidCommandName)}
	actionsRow := utils.MakeActionsRow(messageComponent)

	image := generateImage(i, res.BannerUrl)

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
		Image: image,
	}
	utils.InteractionEmbedRespondForSelf(s, i, embed, actionsRow, true)
}
