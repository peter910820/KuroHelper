package handlers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	kurohelperdb "github.com/peter910820/kurohelper-db/v2"

	"gorm.io/gorm"

	"kurohelper/utils"

	"github.com/peter910820/kurohelper-core/cache"
	"github.com/peter910820/kurohelper-core/erogs"
)

// 加收藏Handler
func AddInWish(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.NewCID) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	if cid != nil {
		// get cache
		cacheValue, err := cache.Get(cid.GetCacheID())
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		resValue := cacheValue.(erogs.FuzzySearchGameResponse)
		res := &resValue

		userID := utils.GetUserID(i)
		userName := utils.GetUsername(i)
		if strings.TrimSpace(userID) != "" && strings.TrimSpace(userName) != "" {
			err := kurohelperdb.Dbs.Transaction(func(tx *gorm.DB) error {
				// 1. 確保 User 存在
				if _, err := kurohelperdb.EnsureUserTx(tx, userID, userName); err != nil {
					return err
				}

				// 2. 確保 Brand 存在
				if _, err := kurohelperdb.EnsureBrandErogsTx(tx, res.BrandID, res.BrandName); err != nil {
					return err
				}

				// 3. 確保 Game 存在
				if _, err := kurohelperdb.EnsureGameErogsTx(tx, res.ID, res.Gamename, res.BrandID); err != nil {
					return err
				}

				// 4. 建立資料
				if err := kurohelperdb.CreateUserInWishTx(tx, userID, res.ID); err != nil {
					return err
				}

				return nil // commit
			})
			if err != nil {
				utils.HandleError(err, s, i)
				return
			}

			// 確保新建立的使用者有加入快取
			if _, ok := cache.UserCache[userID]; !ok {
				cache.UserCache[userID] = struct{}{}
			}

			embed := &discordgo.MessageEmbed{
				Title: "加入成功！",
				Color: 0x90B44B,
			}
			utils.InteractionEmbedRespondForSelf(s, i, embed, nil, true)
		} else { // 找不到使用者，此狀況應該會是Discord官方問題或是程式碼邏輯問題
			embed := &discordgo.MessageEmbed{
				Title: "找不到使用者！",
				Color: 0x90B44B,
			}
			utils.InteractionEmbedRespondForSelf(s, i, embed, nil, true)
		}
	} else {
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

		cidCommandName := utils.MakeCIDCommandName(i.ApplicationCommandData().Name, false, "")
		messageComponent := []discordgo.MessageComponent{utils.MakeCIDCommonComponent("✅", idStr, cidCommandName)}
		actionsRow := utils.MakeActionsRow(messageComponent)

		image := generateImage(i, res.BannerUrl)

		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name: res.BrandName,
			},
			Title: fmt.Sprintf("**%s(%s)**", res.Gamename, res.SellDay),
			URL:   res.Shoukai,
			Color: 0x90B44B,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "發行機種",
					Value:  res.Model,
					Inline: false,
				},
				{
					Name:   "確認",
					Value:  "你確定要加入收藏嗎?",
					Inline: false,
				},
			},
			Image: image,
		}
		utils.InteractionEmbedRespondForSelf(s, i, embed, actionsRow, true)
	}
}
