package handlers

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"

	"kurohelper/cache"
	"kurohelper/database"
	kurohelpererrors "kurohelper/errors"
	"kurohelper/provider/erogs"
	"kurohelper/provider/vndb"
	"kurohelper/utils"
)

func SearchBrand(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.NewCID) {
	// 長時間查詢
	if cid == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
	}

	if i.Type == discordgo.InteractionApplicationCommand {
		opt, err := utils.GetOptions(i, "查詢資料庫選項")
		if err != nil && errors.Is(err, kurohelpererrors.ErrOptionTranslateFail) {
			utils.HandleError(err, s, i)
			return
		}
		if opt == "" {
			erogsSearchBrand(s, i, cid)
		} else {
			vndbSearchBrand(s, i, cid)
		}
	} else {
		if cid.GetCommandNameProvider() == "erogs" {
			erogsSearchBrand(s, i, cid)
		} else {
			vndbSearchBrand(s, i, cid)
		}
	}
}

func erogsSearchBrand(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.NewCID) {
	var res *erogs.FuzzySearchBrandResponse
	var messageComponent []discordgo.MessageComponent
	var hasMore bool
	var count int
	var pageIndex int
	if cid == nil {
		keyword, err := utils.GetOptions(i, "keyword")
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		res, err = erogs.GetBrandByFuzzy(keyword)
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		idStr := uuid.New().String()
		cache.Set(idStr, *res)

		// 根據遊戲評價排序
		sort.Slice(res.GameList, func(i, j int) bool {
			if res.GameList[i].SellDay == "2050-01-01" {
				return false
			} else if res.GameList[j].SellDay == "2050-01-01" {
				return true
			} else {
				return res.GameList[i].SellDay > res.GameList[j].SellDay // 晚到早排序
			}
		})
		// 計算筆數
		count = len(res.GameList)

		hasMore = pagination(&(res.GameList), 0, false)

		if hasMore {
			messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("▶️", idStr, 1, false, i.ApplicationCommandData().Name, "erogs")}
		}
	} else {
		// 處理CID
		pageCID := utils.PageCID{
			NewCID: *cid,
		}
		cacheValue, err := cache.Get(pageCID.GetCacheID())
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		resValue := cacheValue.(erogs.FuzzySearchBrandResponse)
		res = &resValue
		// 資料分頁
		pageIndex, err = pageCID.GetPageIndex()
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		count = len(res.GameList)
		hasMore = pagination(&(res.GameList), pageIndex, true)
		if hasMore {
			if pageIndex == 0 {
				messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("▶️", pageCID.GetCacheID(), 1, false, cid.GetCommandName(), "erogs")}
			} else {
				messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("◀️", pageCID.GetCacheID(), pageIndex-1, false, cid.GetCommandName(), "erogs")}
				messageComponent = append(messageComponent, utils.MakeCIDPageComponent("▶️", pageCID.GetCacheID(), pageIndex+1, false, cid.GetCommandName(), "erogs"))
			}
		} else {
			messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("◀️", pageCID.GetCacheID(), pageIndex-1, false, cid.GetCommandName(), "erogs")}
		}
	}

	actionsRow := utils.MakeActionsRow(messageComponent)

	// 處理資料庫
	var userGameErogs []database.UserGameErogs
	status := make(map[int]byte)
	userID := utils.GetUserID(i)
	if strings.TrimSpace(userID) != "" {
		_, ok := cache.UserCache[userID]
		if ok {
			database.Dbs[os.Getenv("DB_NAME")].Where("user_id = ?", userID).Preload("GameErogs").Find(&userGameErogs)
			// 利用位元運算壓縮狀態
			for _, game := range userGameErogs {
				if game.HasPlayed {
					status[game.GameErogsID] |= Played
				}
				if game.InWish {
					status[game.GameErogsID] |= Wish
				}
			}
		}
	}

	gameData := make([]string, 0, len(res.GameList))
	for _, g := range res.GameList {
		var prefix string
		flags := status[g.ID]
		if flags&Played != 0 {
			prefix += "✅"
		}
		if flags&Wish != 0 {
			prefix += "❤️"
		}
		gameData = append(gameData, fmt.Sprintf("%s%s　%d(%d)　**%s** (%s)", prefix, g.SellDay, g.Median, g.Count2, g.GameName, g.Model))
	}

	if res.Lost {
		res.BrandName += " (解散)"
	}
	link := ""
	if res.URL != "" {
		link += fmt.Sprintf("[官網](%s) ", res.URL)
	}
	if res.Twitter != "" {
		link += fmt.Sprintf("[Twitter](https://x.com/%s) ", res.Twitter)
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s (%d筆)", res.BrandName, count),
		Color:       0xF8F8DF,
		Description: link,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "歷代作品(發售日排序)",
				Value:  strings.Join(gameData, "\n"),
				Inline: false,
			},
		},
	}

	if cid == nil {
		utils.InteractionEmbedRespond(s, i, embed, actionsRow, true)
	} else {
		utils.EditEmbedRespond(s, i, embed, actionsRow)
	}

}

func vndbSearchBrand(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.NewCID) {
	var res *vndb.ProducerSearchResponse
	var messageComponent []discordgo.MessageComponent
	var hasMore bool
	var pageIndex int
	// 第一次查詢
	if cid == nil {
		keyword, err := utils.GetOptions(i, "keyword")
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		// companyType, err := utils.GetOptions(i, "type")
		// if err != nil && errors.Is(err, kurohelpererrors.ErrOptionTranslateFail) {
		// 	utils.HandleError(err, s, i)
		// 	return
		// }

		res, err = vndb.GetProducerByFuzzy(keyword, "")
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		idStr := uuid.New().String()
		cache.Set(idStr, *res)
		hasMore = pagination(&(res.Vn.Results), 0, false)

		if hasMore {
			messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("▶️", idStr, 1, false, i.ApplicationCommandData().Name, "vndb")}
		}
	} else {
		// 處理CID
		pageCID := utils.PageCID{
			NewCID: *cid,
		}
		cacheValue, err := cache.Get(pageCID.GetCacheID())
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		resValue := cacheValue.(vndb.ProducerSearchResponse)
		res = &resValue
		// 資料分頁
		pageIndex, err = pageCID.GetPageIndex()
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		hasMore = pagination(&(res.Vn.Results), pageIndex, true)
		if hasMore {
			if pageIndex == 0 {
				messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("▶️", pageCID.GetCacheID(), 1, false, cid.GetCommandName(), "vndb")}
			} else {
				messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("◀️", pageCID.GetCacheID(), pageIndex-1, false, cid.GetCommandName(), "vndb")}
				messageComponent = append(messageComponent, utils.MakeCIDPageComponent("▶️", pageCID.GetCacheID(), pageIndex+1, false, cid.GetCommandName(), "vndb"))
			}
		} else {
			messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("◀️", pageCID.GetCacheID(), pageIndex-1, false, cid.GetCommandName(), "vndb")}
		}
	}

	actionsRow := utils.MakeActionsRow(messageComponent)

	/* 處理回傳結構 */
	title := res.Producer.Results[0].Original
	if len(res.Producer.Results[0].Aliases) != 0 {
		allAlias := make([]string, 0, len(res.Producer.Results[0].Aliases))
		allAlias = append(allAlias, res.Producer.Results[0].Aliases...)

		if strings.TrimSpace(title) != "" {
			title += fmt.Sprintf("%s(%s)", allAlias[0], strings.Join(allAlias[1:], "), ("))
		} else {
			if len(allAlias) > 1 {
				title = fmt.Sprintf("%s(%s)", allAlias[0], strings.Join(allAlias[1:], "), ("))
			} else {
				title = allAlias[0]
			}
		}

	}

	if strings.TrimSpace(title) == "" {
		title = res.Producer.Results[0].Name
	}

	gameData := make([]string, 0, len(res.Vn.Results))
	for _, game := range res.Vn.Results {
		if strings.TrimSpace(game.Alttitle) != "" {
			gameData = append(gameData, fmt.Sprintf("%.1f/%.1f/%03d　%02d(H)/%03d　**%s**", game.Average, game.Rating, game.Votecount, game.LengthMinutes/60, game.LengthVotes, game.Alttitle))
		} else {
			gameData = append(gameData, fmt.Sprintf("%.1f/%.1f/%03d　%02d(H)/%03d　**%s**", game.Average, game.Rating, game.Votecount, game.LengthMinutes/60, game.LengthVotes, game.Title))
		}
	}

	embed := &discordgo.MessageEmbed{
		Title: title,
		Color: 0x04108e,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "品牌(公司)名稱",
				Value:  title,
				Inline: false,
			},
			{
				Name:   "遊戲列表",
				Value:  strings.Join(gameData, "\n"),
				Inline: false,
			},
		},
	}

	if cid == nil {
		utils.InteractionEmbedRespond(s, i, embed, actionsRow, true)
	} else {
		utils.EditEmbedRespond(s, i, embed, actionsRow)
	}

}
