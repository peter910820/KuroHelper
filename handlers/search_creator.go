package handlers

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"

	"kurohelper/cache"
	kurohelpererrors "kurohelper/errors"
	"kurohelper/provider/erogs"
	"kurohelper/utils"
)

func SearchCreator(s *discordgo.Session, i *discordgo.InteractionCreate, cid *CustomID) {
	// 長時間查詢
	if cid == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
	}

	if i.Type == discordgo.InteractionApplicationCommand {
		opt, err := utils.GetOptions(i, "列表搜尋")
		if err != nil && errors.Is(err, kurohelpererrors.ErrOptionTranslateFail) {
			utils.HandleError(err, s, i)
			return
		}
		if opt == "" {
			erogsSearchCreator(s, i, cid)
		} else {
			erogsSearchCreatorList(s, i, cid)
		}
	} else {
		erogsSearchCreatorList(s, i, cid)
	}
}

func erogsSearchCreator(s *discordgo.Session, i *discordgo.InteractionCreate, cid *CustomID) {
	var res *erogs.FuzzySearchCreatorResponse
	var messageComponent []discordgo.MessageComponent
	var hasMore bool
	var count int
	var countInner int

	if cid == nil {
		keyword, err := utils.GetOptions(i, "keyword")
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		idSearch, _ := regexp.MatchString(`^e\d+$`, keyword)
		res, err = erogs.GetCreatorByFuzzy(keyword, idSearch)
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		idStr := uuid.New().String()
		cache.Set(idStr, *res)

		// 根據遊戲評價排序
		sort.Slice(res.Games, func(i, j int) bool {
			return res.Games[i].Median > res.Games[j].Median // 大到小排序
		})
		// 計算筆數
		for _, inner := range res.Games {
			countInner += len(inner.Shokushu)
		}
		count = len(res.Games)

		hasMore = pagination(&(res.Games), 0, false)

		if hasMore {
			messageComponent = []discordgo.MessageComponent{utils.MakePageComponent("▶️", i.ApplicationCommandData().Name, idStr, 1)}
		}
	} else {
		cacheValue, err := cache.Get(cid.ID)
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		resValue := cacheValue.(erogs.FuzzySearchCreatorResponse)
		res = &resValue

		// 根據遊戲評價排序
		sort.Slice(res.Games, func(i, j int) bool {
			return res.Games[i].Median > res.Games[j].Median // 大到小排序
		})
		// 計算筆數
		for _, inner := range res.Games {
			countInner += len(inner.Shokushu)
		}
		count = len(res.Games)

		// 資料分頁
		hasMore = pagination(&(res.Games), cid.Value, true)
		if hasMore {
			if cid.Value == 0 {
				messageComponent = []discordgo.MessageComponent{utils.MakePageComponent("▶️", cid.CommandName, cid.ID, 1)}
			} else {
				messageComponent = []discordgo.MessageComponent{utils.MakePageComponent("◀️", cid.CommandName, cid.ID, cid.Value-1)}
				messageComponent = append(messageComponent, utils.MakePageComponent("▶️", cid.CommandName, cid.ID, cid.Value+1))
			}
		} else {
			messageComponent = []discordgo.MessageComponent{utils.MakePageComponent("◀️", cid.CommandName, cid.ID, cid.Value-1)}
		}
	}

	actionsRow := utils.MakeActionsRow(messageComponent)

	link := ""
	if res.TwitterUsername != "" {
		link += fmt.Sprintf("[Twitter](https://x.com/%s) ", res.TwitterUsername)
	}
	if res.Blog != "" {
		link += fmt.Sprintf("[Blog](%s) ", res.Blog)
	}
	if res.Pixiv != nil {
		link += fmt.Sprintf("[Pixiv](https://www.pixiv.net/users/%d) ", *res.Pixiv)
	}

	gameData := make([]string, 0, len(res.Games))
	for i, g := range res.Games {
		shokushu := make([]string, 0, len(g.Shokushu))
		for _, s := range g.Shokushu {
			if s.Shubetu != 7 {
				shokushu = append(shokushu, fmt.Sprintf("*%s*", erogs.ShubetuMap[s.Shubetu]))
			} else {
				shokushu = append(shokushu, fmt.Sprintf("*%s*", s.ShubetuDetailName))
			}
		}

		if cid == nil {
			gameData = append(gameData, fmt.Sprintf("%d. **%s**  (%s) / %d分 / %s", i+1, g.Gamename, strings.Join(shokushu, ", "), g.Median, g.SellDay))
		} else {
			gameData = append(gameData, fmt.Sprintf("%d. **%s**  (%s) / %d分 / %s", cid.Value*10+i+1, g.Gamename, strings.Join(shokushu, ", "), g.Median, g.SellDay))
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s(%d/%d  筆)", res.Name, count, countInner),
		Color:       0xF8F8DF,
		Description: link,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "歷代作品(遊戲評價排序)",
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

func erogsSearchCreatorList(s *discordgo.Session, i *discordgo.InteractionCreate, cid *CustomID) {
	var res *[]erogs.FuzzySearchListResponse
	var messageComponent []discordgo.MessageComponent
	var hasMore bool
	var count int
	if cid == nil {
		keyword, err := utils.GetOptions(i, "keyword")
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		res, err = erogs.GetCreatorListByFuzzy(keyword)
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		idStr := uuid.New().String()
		cache.Set(idStr, *res)

		// 計算筆數
		count = len(*res)

		hasMore = pagination(res, 0, false)

		if hasMore {
			messageComponent = []discordgo.MessageComponent{utils.MakePageComponent("▶️", "查詢創作者列表", idStr, 1)}
		}
	} else {
		cacheValue, err := cache.Get(cid.ID)
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		resValue := cacheValue.([]erogs.FuzzySearchListResponse)
		res = &resValue

		// 計算筆數
		count = len(*res)

		// 資料分頁
		hasMore = pagination(res, cid.Value, true)
		if hasMore {
			if cid.Value == 0 {
				messageComponent = []discordgo.MessageComponent{utils.MakePageComponent("▶️", cid.CommandName, cid.ID, 1)}
			} else {
				messageComponent = []discordgo.MessageComponent{utils.MakePageComponent("◀️", cid.CommandName, cid.ID, cid.Value-1)}
				messageComponent = append(messageComponent, utils.MakePageComponent("▶️", cid.CommandName, cid.ID, cid.Value+1))
			}
		} else {
			messageComponent = []discordgo.MessageComponent{utils.MakePageComponent("◀️", cid.CommandName, cid.ID, cid.Value-1)}
		}
	}
	actionsRow := utils.MakeActionsRow(messageComponent)

	listData := make([]string, 0, len(*res))
	for _, r := range *res {
		listData = append(listData, fmt.Sprintf("e%-5s　%s", strconv.Itoa(r.ID), r.Name))
	}
	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("創作者列表搜尋 (%d筆)", count),
		Color: 0xF8F8DF,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "ID/名稱",
				Value:  strings.Join(listData, "\n"),
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
