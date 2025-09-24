package handlers

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"kurohelper/erogs"
	internalerrors "kurohelper/errors"
	"kurohelper/models"
	erogsmodels "kurohelper/models/erogs"
	"kurohelper/utils"
)

func ErogsFuzzySearchCreator(s *discordgo.Session, i *discordgo.InteractionCreate, cid *models.VndbInteractionCustomID) {
	var res *erogsmodels.FuzzySearchCreatorResponse
	var component *discordgo.ActionsRow
	var hasMore bool
	var count int

	if cid == nil {
		// 長時間查詢
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

		keyword, err := utils.GetOptions(i, "keyword")
		if err != nil {
			logrus.Error(err)
			utils.InteractionEmbedErrorRespond(s, i, "該功能目前異常，請稍後再嘗試", true)
			return
		}

		res, err = erogs.GetCreatorByFuzzy(keyword)
		if err != nil {
			logrus.Error(err)
			if errors.Is(err, internalerrors.ErrVndbNoResult) {
				utils.InteractionEmbedErrorRespond(s, i, "找不到任何結果喔", true)
			} else {
				utils.InteractionEmbedErrorRespond(s, i, "該功能目前異常，請稍後再嘗試", true)
			}
			return
		}

		idStr := uuid.New().String()
		SetCache(idStr, *res)

		// 根據遊戲評價排序
		sort.Slice(res.Games, func(i, j int) bool {
			return res.Games[i].Median > res.Games[j].Median // 大到小排序
		})
		// 計算筆數
		for _, inner := range res.Games {
			count += len(inner.Shokushu)
		}

		hasMore = erogsPagination(&(res.Games), 0, false)

		if hasMore {
			component = &discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "▶️",
						Style:    discordgo.PrimaryButton,
						CustomID: fmt.Sprintf("ErogsFuzzySearchCreator_1_%s", idStr),
					},
				},
			}
		}
	} else {
		cacheValue, ok := GetCache(cid.Key)
		if !ok {
			utils.EmbedErrorRespond(s, i, "快取遺失，請重新查詢")
			return
		}
		resValue := cacheValue.(erogsmodels.FuzzySearchCreatorResponse)
		res = &resValue

		// 根據遊戲評價排序
		sort.Slice(res.Games, func(i, j int) bool {
			return res.Games[i].Median > res.Games[j].Median // 大到小排序
		})
		// 計算筆數
		for _, inner := range res.Games {
			count += len(inner.Shokushu)
		}

		// 資料分頁
		hasMore = erogsPagination(&(res.Games), cid.Page, true)
		if hasMore {
			if cid.Page == 0 {
				component = &discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "▶️",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("ErogsFuzzySearchCreator_1_%s", cid.Key),
						},
					},
				}
			} else {
				component = &discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "◀️",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("ErogsFuzzySearchCreator_%d_%s", cid.Page-1, cid.Key),
						},
						&discordgo.Button{
							Label:    "▶️",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("ErogsFuzzySearchCreator_%d_%s", cid.Page+1, cid.Key),
						},
					},
				}
			}
		} else {
			component = &discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					&discordgo.Button{
						Label:    "◀️",
						Style:    discordgo.PrimaryButton,
						CustomID: fmt.Sprintf("ErogsFuzzySearchCreator_%d_%s", cid.Page-1, cid.Key),
					},
				},
			}
		}
	}

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
			gameData = append(gameData, fmt.Sprintf("%d. **%s**  (%s)  %d分", i+1, g.Gamename, strings.Join(shokushu, ", "), g.Median))
		} else {
			gameData = append(gameData, fmt.Sprintf("%d. **%s**  (%s)  %d分", cid.Page*15+1, g.Gamename, strings.Join(shokushu, ", "), g.Median))
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s(%d  筆)", res.Name, count),
		Color:       0x04108e,
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
		utils.InteractionEmbedRespond(s, i, embed, component, true)
	} else {
		utils.EditEmbedRespond(s, i, embed, component)
	}

}

// 資料分頁
func erogsPagination(result *[]erogsmodels.Game, page int, useCache bool) bool {
	resultLen := len(*result)
	expectedMin := page * 15
	expectedMax := page*15 + 15

	if !useCache || page == 0 {
		if resultLen > 15 {
			*result = (*result)[:15]
			return true
		}
		return false
	} else {
		if resultLen > expectedMax {
			*result = (*result)[expectedMin:expectedMax]
			return true
		} else {
			*result = (*result)[expectedMin:]
			return false
		}
	}
}
