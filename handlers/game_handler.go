package handlers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"

	"kurohelper/cache"
	kurohelpererrors "kurohelper/errors"
	"kurohelper/provider/erogs"
	"kurohelper/provider/seiya"
	"kurohelper/provider/vndb"
	"kurohelper/utils"
)

func FuzzySearchGame(s *discordgo.Session, i *discordgo.InteractionCreate, cid *CustomID) {
	if i.Type == discordgo.InteractionApplicationCommand {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

		opt, err := utils.GetOptions(i, "列表搜尋")
		if err != nil && errors.Is(err, kurohelpererrors.ErrOptionTranslateFail) {
			utils.HandleError(err, s, i)
			return
		}
		if opt == "" {
			erogsFuzzySearchGame(s, i)
		} else {
			erogsFuzzySearchGameList(s, i, cid)
		}
	} else {
		erogsFuzzySearchGameList(s, i, cid)
	}

}

func erogsFuzzySearchGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	var res *erogs.FuzzySearchGameResponse
	var resVndb *vndb.BasicResponse[vndb.GetVnUseIDResponse]

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

	vndbRating := 0.0
	vndbVotecount := 0
	if strings.TrimSpace(res.VndbId) != "" {
		resVndb, err = vndb.GetVnUseID(res.VndbId)
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		vndbRating = resVndb.Results[0].Rating
		vndbVotecount = resVndb.Results[0].Votecount
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
	if res.Junni == 0 || res.Junni > 500 {
		junni = 0x04108e // Default
	} else if res.Junni <= 50 {
		junni = 0xFFD700 // Gold
		rank = "批評空間 TOP 50"
	} else if res.Junni <= 100 {
		junni = 0xC0C0C0 // Silver
		rank = "批評空間 TOP 100"
	} else {
		junni = 0xCD7F32 // Bronze
		rank = "批評空間 TOP 500"
	}

	// 用批評空間回來的遊戲名對誠也做模糊搜尋
	seiyaURL := seiya.GetGuideURL(res.Gamename)
	if seiyaURL != "" {
		rank += "  " + fmt.Sprintf("[誠也攻略](%s)", seiyaURL)
	}
	erogsURL := "https://erogamescape.dyndns.org/~ap2/ero/toukei_kaiseki/game.php?game=" + fmt.Sprint(res.ID)
	rank += "  " + fmt.Sprintf("[批評空間](%s)", erogsURL)
	if res.VndbId != "" {
		vndbURL := "https://vndb.org/" + res.VndbId
		rank += "  " + fmt.Sprintf("[VNDB](%s)", vndbURL)
	}
	vndbData := "無"
	if vndbVotecount != 0 {
		vndbData = fmt.Sprintf("%.1f/%d", vndbRating, vndbVotecount)
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
				Value:  vndbData,
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
				Inline: true,
			},
			{
				Name:   "發行機種",
				Value:  res.Model,
				Inline: true,
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

func erogsFuzzySearchGameList(s *discordgo.Session, i *discordgo.InteractionCreate, cid *CustomID) {
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

		res, err = erogs.GetGameListByFuzzy(keyword)
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
			messageComponent = []discordgo.MessageComponent{utils.MakePageComponent("▶️", "查詢遊戲列表", idStr, 1)}
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
		listData = append(listData, fmt.Sprintf("e%-5s　%s (%s)", strconv.Itoa(r.ID), r.Name, r.Category))
	}
	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("遊戲列表搜尋 (%d筆)", count),
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
