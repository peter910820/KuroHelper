package handlers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	kurohelperdb "github.com/peter910820/kurohelper-db/v2"

	"gorm.io/gorm"

	kurohelpererrors "kurohelper/errors"
	"kurohelper/utils"

	"github.com/peter910820/kurohelper-core/cache"
	"github.com/peter910820/kurohelper-core/erogs"
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
			// 先檢查該筆以及操作是不是已經執行過
			gameRecord, err := kurohelperdb.GetUserGameErogs(userID, res.ID)
			if err != nil {
				// 沒有資料 開始新建
				if errors.Is(err, gorm.ErrRecordNotFound) {
					apiModel := kurohelperdb.UpsertUserGameErogsTXInput{
						UserID:         userID,
						UserName:       userName,
						ErogsBrandID:   res.BrandID,
						ErogsBrandName: res.BrandName,
						ErogsGameID:    res.ID,
						ErogsGamename:  res.Gamename,
						HasPlayed:      true,
						InWish:         false,
						CompleteDate:   completeDate,
					}
					err = kurohelperdb.UpsertUserGameErogsTransaction(apiModel)
					if err != nil {
						utils.HandleError(err, s, i)
						return
					}
				} else {
					utils.HandleError(err, s, i)
					return
				}
			} else if gameRecord.HasPlayed { // 有資料的狀況下，hasPlayed已被設定過
				embed := &discordgo.MessageEmbed{
					Title: "資料已建立，本次動作無效",
					Color: 0x7BA23F,
				}
				utils.InteractionEmbedRespondForSelf(s, i, embed, nil, true)
				return
			} else { // 有資料，hasPlayed還沒被設定過
				apiModel := kurohelperdb.UpsertUserGameErogsTXInput{
					UserID:         userID,
					UserName:       userName,
					ErogsBrandID:   res.BrandID,
					ErogsBrandName: res.BrandName,
					ErogsGameID:    res.ID,
					ErogsGamename:  res.Gamename,
					HasPlayed:      true,
					InWish:         gameRecord.InWish,
					CompleteDate:   completeDate,
				}
				err = kurohelperdb.UpsertUserGameErogsTransaction(apiModel)
				if err != nil {
					utils.HandleError(err, s, i)
					return
				}
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
