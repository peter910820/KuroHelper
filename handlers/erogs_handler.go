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

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		logrus.Error(err)
		utils.InteractionEmbedErrorRespond(s, i, "該功能目前異常，請稍後再嘗試", true)
		return
	}

	if cid == nil {
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
		hasMore = erogsPagination(&(res.Creator[0].Games), 0, false)

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
		// 資料分頁
		hasMore = erogsPagination(&(res.Creator[0].Games), cid.Page, true)
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
						CustomID: fmt.Sprintf("SearchBrand_%d_%s", cid.Page-1, cid.Key),
					},
				},
			}
		}
	}

	link := ""
	if res.Creator[0].TwitterUsername != "" {
		link += fmt.Sprintf("[Twitter](https://x.com/%s) ", res.Creator[0].TwitterUsername)
	}
	if res.Creator[0].Blog != "" {
		link += fmt.Sprintf("[Blog](%s) ", res.Creator[0].Blog)
	}
	if res.Creator[0].Pixiv != nil {
		link += fmt.Sprintf("[Pixiv](https://www.pixiv.net/users/%d) ", *res.Creator[0].Pixiv)
	}

	// 根據遊戲評價排序
	sort.Slice(res.Creator[0].Games, func(i, j int) bool {
		return res.Creator[0].Games[i].Median > res.Creator[0].Games[j].Median // 大到小排序
	})

	gameData := make([]string, 0, len(res.Creator[0].Games))
	for i, g := range res.Creator[0].Games {
		shokushu := make([]string, 0, len(g.Shokushu))
		for _, s := range g.Shokushu {
			if s.Shubetu != 7 {
				shokushu = append(shokushu, fmt.Sprintf("*%s*", erogs.ShubetuMap[s.Shubetu]))
			} else {
				shokushu = append(shokushu, fmt.Sprintf("*%s*", s.ShubetuDetailName))
			}
		}
		gameData = append(gameData, fmt.Sprintf("%d. **%s**  (%s)", i+1, g.Gamename, strings.Join(shokushu, ", ")))
	}

	embed := &discordgo.MessageEmbed{
		Title:       res.Creator[0].Name,
		Color:       0x04108e,
		Description: link,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "歷代作品(評價排序)",
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
	expectedMin := page * 10
	expectedMax := page*10 + 9

	if !useCache || page == 0 {
		if resultLen > 10 {
			*result = (*result)[:10]
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
