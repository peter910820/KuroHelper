package handlers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"

	"kurohelper/erogs"
	"kurohelper/models"
	erogsmodels "kurohelper/models/erogs"
	"kurohelper/utils"
	"kurohelper/vndb"
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
			utils.HandleError(err, s, i)
			return
		}

		res, err = erogs.GetCreatorByFuzzy(keyword)
		if err != nil {
			utils.HandleError(err, s, i)
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

func ErogsFuzzySearchMusic(s *discordgo.Session, i *discordgo.InteractionCreate, cid *models.VndbInteractionCustomID) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	var res *erogsmodels.FuzzySearchMusicResponse
	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	res, err = erogs.GetMusicByFuzzy(keyword)
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	musicData := make([]string, 0, len(res.GameCategories))
	for _, m := range res.GameCategories {
		musicData = append(musicData, m.GameName+" ("+m.Category+")")
	}

	singerList := strings.Split(res.Singers, ",")
	arrangementList := strings.Split(res.Arrangments, ",")
	lyricList := strings.Split(res.Lyrics, ",")
	compositionList := strings.Split(res.Compositions, ",")
	albumList := strings.Split(res.Album, ",")
	if res.PlayTime == "00:00:00" {
		res.PlayTime = "未收錄"
	}
	if res.ReleaseDate == "0001-01-01" {
		res.ReleaseDate = "未收錄"
	}

	embed := &discordgo.MessageEmbed{
		Title: res.MusicName,
		Color: 0x04108e,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "音樂時長",
				Value:  res.PlayTime,
				Inline: true,
			},
			{
				Name:   "發行日期",
				Value:  res.ReleaseDate,
				Inline: true,
			},
			{
				Name:   "平均分數/樣本數",
				Value:  fmt.Sprintf("%.2f / %d", res.AvgTokuten, res.TokutenCount),
				Inline: true,
			},
			{
				Name:   "歌手",
				Value:  strings.Join(singerList, "\n"),
				Inline: false,
			},
			{
				Name:   "作詞",
				Value:  strings.Join(lyricList, "\n"),
				Inline: true,
			},
			{
				Name:   "作曲",
				Value:  strings.Join(compositionList, "\n"),
				Inline: true,
			},
			{
				Name:   "編曲",
				Value:  strings.Join(arrangementList, "\n"),
				Inline: true,
			},
			{
				Name:   "遊戲收錄",
				Value:  strings.Join(musicData, "\n"),
				Inline: false,
			},
			{
				Name:   "專輯",
				Value:  strings.Join(albumList, "\n"),
				Inline: false,
			},
		},
	}
	utils.InteractionEmbedRespond(s, i, embed, nil, true)
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

func ErogsFuzzySearchGame(s *discordgo.Session, i *discordgo.InteractionCreate, cid *models.VndbInteractionCustomID) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	var res *erogsmodels.FuzzySearchGameResponse
	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	res, err = erogs.GetGameByFuzzy(keyword)
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	resVndb, err := vndb.GetVnUseID(res.VndbId)
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	shubetuData := make(map[int]map[int][]string) // map[shubetu_type]map[shubetu_detail]][]creator name + shube1tu_detail_name

	for typeIdx := 1; typeIdx <= 6; typeIdx++ {
		shubetuData[typeIdx] = make(map[int][]string)
		for detailIdx := 1; detailIdx <= 3; detailIdx++ {
			shubetuData[typeIdx][detailIdx] = make([]string, 0)
		}
	}
	for _, shubetu := range res.CreatorShubetu { // 遍歷shubetu_detail
		shubetuType := shubetu.ShubetuType
		detailType := shubetu.ShubetuDetailType
		creatorData := ""
		if shubetu.ShubetuDetailName == "" {
			creatorData = shubetu.CreatorName
		} else {
			creatorData = shubetu.CreatorName + " (" + shubetu.ShubetuDetailName + ")"
		}
		if shubetu.ShubetuType != 5 {
			shubetuData[shubetuType][1] = append(shubetuData[shubetuType][1], creatorData)
		} else {
			if detailType == 2 || detailType == 3 {
				shubetuData[shubetuType][2] = append(shubetuData[shubetuType][2], creatorData)
			} else {
				shubetuData[shubetuType][1] = append(shubetuData[shubetuType][1], creatorData)
			}
		}
	}

	switch res.Okazu {
	case "true":
		res.Okazu = "拔作"
	case "false":
		res.Okazu = "非拔作"
	default:
		res.Okazu = ""
	}

	switch res.Erogame {
	case "true":
		res.Erogame = "18禁"
	case "false":
		res.Erogame = "全年齡"
	default:
		res.Erogame = ""
	}

	otherInfo := ""
	if res.Erogame == "" && res.Okazu == "" {
		otherInfo = "無"
	} else if res.Erogame == "" || res.Okazu == "" {
		otherInfo = res.Erogame + res.Okazu
	} else {
		otherInfo = res.Okazu + " / " + res.Erogame
	}

	junni := 0x04108e
	rank := ""
	if res.Junni <= 50 {
		junni = 0xFFD700 // Gold
		rank = "批評空間 TOP 50"
	} else if res.Junni <= 100 {
		junni = 0xC0C0C0 // Silver
		rank = "批評空間 TOP 100"
	} else if res.Junni <= 500 {
		junni = 0xCD7F32 // Bronze
		rank = "批評空間 TOP 500"
	} else {
		junni = 0x04108e // Default
	}
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: res.BrandName,
		},
		Title:       fmt.Sprintf("**%s(%s)**", res.Gamename, res.SellDay),
		URL:         res.Shoukai,
		Color:       junni,
		Description: rank,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "劇本",
				Value:  strings.Join(shubetuData[2][1], " / "),
				Inline: false,
			},
			{
				Name:   "原畫",
				Value:  strings.Join(shubetuData[1][1], " / "),
				Inline: false,
			},
			{
				Name:   "主角群CV",
				Value:  strings.Join(shubetuData[5][1], " / "),
				Inline: false,
			},
			{
				Name:   "配角群CV",
				Value:  strings.Join(shubetuData[5][2], " / "),
				Inline: false,
			},
			{
				Name:   "歌手",
				Value:  strings.Join(shubetuData[6][1], " / "),
				Inline: false,
			},
			{
				Name:   "音樂",
				Value:  strings.Join(shubetuData[3][1], " / "),
				Inline: false,
			},
			{
				Name:   "批評空間分數/樣本數",
				Value:  res.Median + " / " + res.TokutenCount,
				Inline: true,
			},
			{
				Name:   "vndb分數/樣本數",
				Value:  fmt.Sprintf("%.1f/%d", resVndb.Results[0].Rating, resVndb.Results[0].Votecount),
				Inline: true,
			},
			{
				Name:   "遊玩時數",
				Value:  res.TotalPlayTimeMedian,
				Inline: true,
			},
			{
				Name:   "開始理解遊戲樂趣時數",
				Value:  res.TimeBeforeUnderstandingFunMedian,
				Inline: false,
			},
			{
				Name:   "類型",
				Value:  res.Genre,
				Inline: true,
			},
			{
				Name:   "其他資訊",
				Value:  otherInfo,
				Inline: true,
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: res.BannerUrl,
		},
	}
	utils.InteractionEmbedRespond(s, i, embed, nil, true)
}
