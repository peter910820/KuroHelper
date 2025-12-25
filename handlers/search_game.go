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
	kurohelperdb "github.com/peter910820/kurohelper-db/v2"
	"github.com/siongui/gojianfan"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	kurohelpererrors "kurohelper/errors"
	"kurohelper/utils"

	kurohelpercore "github.com/peter910820/kurohelper-core"
	"github.com/peter910820/kurohelper-core/cache"
	"github.com/peter910820/kurohelper-core/erogs"
	"github.com/peter910820/kurohelper-core/seiya"
	"github.com/peter910820/kurohelper-core/vndb"
	"github.com/peter910820/kurohelper-core/ymgal"
)

// 查詢遊戲Handler
func SearchGame(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.NewCID) {
	// 長時間查詢
	if cid == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
	}

	if i.Type == discordgo.InteractionApplicationCommand {
		optList, err := utils.GetOptions(i, "列表搜尋")
		if err != nil && errors.Is(err, kurohelpererrors.ErrOptionTranslateFail) {
			utils.HandleError(err, s, i)
			return
		}
		optSource, err := utils.GetOptions(i, "查詢資料庫選項")
		if err != nil && errors.Is(err, kurohelpererrors.ErrOptionTranslateFail) {
			utils.HandleError(err, s, i)
			return
		}
		switch optSource {
		case "":
			fallthrough
		case "2":
			if optList == "" {
				erogsSearchGame(s, i)
			} else {
				erogsSearchGameList(s, i, cid)
			}
		case "1":
			if optList == "" {
				vndbSearchGame(s, i)
			} else {
				vndbSearchGameList(s, i, cid)
			}
		}
	} else {
		commandNameProvider := cid.GetCommandNameProvider()
		switch commandNameProvider {
		case "erogs":
			erogsSearchGameList(s, i, cid)
		case "vndb":
			vndbSearchGameList(s, i, cid)
		}
	}
}

// erogs查詢遊戲處理
func erogsSearchGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var res *erogs.FuzzySearchGameResponse
	var resVndb *vndb.BasicResponse[vndb.GetVnUseIDResponse]

	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	idSearch, _ := regexp.MatchString(`^e\d+$`, keyword)

	// 條件符合就用月幕做跳板
	if !idSearch && utils.IsAllHanziOrDigit(keyword) {
		ymgalKeyword, err := ymgalGetGameString(keyword)
		if err != nil {
			logrus.Warn(err)
		}

		if strings.TrimSpace(ymgalKeyword) != "" {
			keyword = ymgalKeyword
		}
	}

	res, err = erogs.GetGameByFuzzy(keyword, idSearch)
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	logrus.Printf("erogs查詢遊戲: %s", keyword)

	// 處理使用者資訊
	userID := utils.GetUserID(i)
	var userData string
	userGameErogs, err := kurohelperdb.GetUserGameErogs(userID, res.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(err, s, i)
			return
		}
	} else { // 有找到資料
		if userGameErogs.HasPlayed {
			userData += "✅"
		}

		if userGameErogs.InWish {
			userData += "❤️"
		}
	}

	vndbRating := 0.0
	vndbVotecount := 0
	if strings.TrimSpace(res.VndbId) != "" {
		resVndb, err = vndb.GetVNByID(res.VndbId)
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

	image := generateImage(i, res.BannerUrl)

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: res.BrandName,
		},
		Title:       fmt.Sprintf("%s**%s(%s)**", userData, res.Gamename, res.SellDay),
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
		Image: image,
	}
	utils.InteractionEmbedRespond(s, i, embed, nil, true)
}

// erogs查詢遊戲列表搜尋處理
func erogsSearchGameList(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.NewCID) {
	var res *[]erogs.FuzzySearchListResponse
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
			cidCommandName := utils.MakeCIDCommandName(i.ApplicationCommandData().Name, true, "")
			messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("▶️", idStr, 1, cidCommandName)}
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
		resValue := cacheValue.([]erogs.FuzzySearchListResponse)
		res = &resValue

		// 計算筆數
		count = len(*res)

		// 資料分頁
		pageIndex, err = pageCID.GetPageIndex()
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		hasMore = pagination(res, pageIndex, true)
		cidCommandName := utils.MakeCIDCommandName(cid.GetCommandName(), true, "")
		if hasMore {
			if pageIndex == 0 {
				messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("▶️", pageCID.GetCacheID(), 1, cidCommandName)}
			} else {
				messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("◀️", pageCID.GetCacheID(), pageIndex-1, cidCommandName)}
				messageComponent = append(messageComponent, utils.MakeCIDPageComponent("▶️", pageCID.GetCacheID(), pageIndex+1, cidCommandName))
			}
		} else {
			messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("◀️", pageCID.GetCacheID(), pageIndex-1, cidCommandName)}
		}
	}
	actionsRow := utils.MakeActionsRow(messageComponent)
	listData := make([]string, 0, len(*res))
	for _, r := range *res {
		listData = append(listData, fmt.Sprintf("**e%-5s**　%s (%s)", strconv.Itoa(r.ID), r.Name, r.Category))
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

// vndb查詢遊戲處理
func vndbSearchGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}
	var res *vndb.BasicResponse[vndb.GetVnUseIDResponse]
	idSearch, _ := regexp.MatchString(`^v\d+$`, keyword)
	if idSearch {
		res, err = vndb.GetVNByID(keyword)
		logrus.Printf("vndb搜尋遊戲ID: %s", keyword)
	} else {
		res, err = vndb.GetVNByFuzzy(keyword)
		logrus.Printf("vndb搜尋遊戲: %s", keyword)
	}

	if err != nil {
		utils.HandleError(err, s, i)
		return
	}
	/* 處理回傳結構 */

	gameTitle := res.Results[0].Alttitle
	if strings.TrimSpace(gameTitle) == "" {
		gameTitle = res.Results[0].Title
	}
	brandTitle := res.Results[0].Developers[0].Original
	if strings.TrimSpace(brandTitle) != "" {
		brandTitle += fmt.Sprintf("(%s)", res.Results[0].Developers[0].Name)
	} else {
		brandTitle = res.Results[0].Developers[0].Name
	}

	// staff block
	var scenario string
	var art string
	var songs string
	var tmpAlias string
	for _, staff := range res.Results[0].Staff {
		staffName := staff.Original
		if staffName == "" {
			staffName = staff.Name
		}
		if len(staff.Aliases) > 0 {
			aliases := make([]string, 0, len(staff.Aliases))
			for _, alias := range staff.Aliases {
				if alias.IsMain {
					staffName = alias.Name
				} else {
					aliases = append(aliases, alias.Name)
				}
			}
			tmpAlias = "(" + strings.Join(aliases, ", ") + ")"
			if len(aliases) == 0 {
				tmpAlias = ""
			}
		}

		switch staff.Role {
		case "scenario":
			scenario += fmt.Sprintf("%s %s\n", staffName, tmpAlias)
		case "art":
			art += fmt.Sprintf("%s %s\n", staffName, tmpAlias)
		case "songs":
			songs += fmt.Sprintf("%s %s\n", staffName, tmpAlias)
		}
	}

	// character block

	characterMap := make(map[string]CharacterData) // map[characterID]CharacterData
	for _, va := range res.Results[0].Va {
		characterName := va.Character.Original
		if characterName == "" {
			characterName = va.Character.Name
		}
		for _, vn := range va.Character.Vns {
			if vn.ID == res.Results[0].ID {
				characterMap[va.Character.ID] = CharacterData{
					Name: characterName,
					Role: vn.Role,
				}
				break
			}
		}
	}

	// 將 map 轉為 slice 並排序
	characterList := make([]CharacterData, 0, len(characterMap))
	for _, character := range characterMap {
		characterList = append(characterList, character)
	}
	sort.Slice(characterList, func(i, j int) bool {
		return characterList[i].Role < characterList[j].Role
	})

	// 格式化輸出
	characters := make([]string, 0, len(characterList))
	for _, character := range characterList {
		characters = append(characters, fmt.Sprintf("**%s** (%s)", character.Name, vndb.Role[character.Role]))
	}

	// relations block
	relationsGame := make([]string, 0, len(res.Results[0].Relations))
	for _, rg := range res.Results[0].Relations {
		titleName := ""
		for _, title := range rg.Titles {
			if title.Main {
				titleName = title.Title
			}
		}
		relationsGame = append(relationsGame, fmt.Sprintf("%s(%s)", titleName, rg.ID))
	}
	relationsGameDisplay := strings.Join(relationsGame, ", ")
	if strings.TrimSpace(relationsGameDisplay) == "" {
		relationsGameDisplay = "無"
	}

	image := generateImage(i, res.Results[0].Image.Url)
	embed := &discordgo.MessageEmbed{
		Title: gameTitle,
		Color: 0x04108e,
		Image: image,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "品牌(公司)名稱",
				Value:  brandTitle,
				Inline: false,
			},
			{
				Name:   "劇本",
				Value:  scenario,
				Inline: false,
			},
			{
				Name:   "美術",
				Value:  art,
				Inline: false,
			},
			{
				Name:   "音樂",
				Value:  songs,
				Inline: false,
			},
			{
				Name:   "評價(平均/貝式平均/樣本數)",
				Value:  fmt.Sprintf("%.1f/%.1f/%d", res.Results[0].Average, res.Results[0].Rating, res.Results[0].Votecount),
				Inline: true,
			},
			{
				Name:   "平均遊玩時數/樣本數",
				Value:  fmt.Sprintf("%d(H)/%d", res.Results[0].LengthMinutes/60, res.Results[0].LengthVotes),
				Inline: true,
			},
			{
				Name:   "角色列表",
				Value:  strings.Join(characters, " / "),
				Inline: false,
			},
			{
				Name:   "相關遊戲",
				Value:  relationsGameDisplay,
				Inline: false,
			},
		},
	}
	utils.InteractionEmbedRespond(s, i, embed, nil, true)
}

// vndb查詢遊戲列表搜尋處理
func vndbSearchGameList(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.NewCID) {
	var res *[]vndb.GetVnIDUseListResponse
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

		res, err = vndb.GetVnID(keyword)
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
			cidCommandName := utils.MakeCIDCommandName(i.ApplicationCommandData().Name, true, "vndb")
			messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("▶️", idStr, 1, cidCommandName)}
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
		resValue := cacheValue.([]vndb.GetVnIDUseListResponse)
		res = &resValue

		// 計算筆數
		count = len(*res)

		// 資料分頁
		pageIndex, err = pageCID.GetPageIndex()
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		hasMore = pagination(res, pageIndex, true)
		cidCommandName := utils.MakeCIDCommandName(cid.GetCommandName(), true, "vndb")
		if hasMore {
			if pageIndex == 0 {
				messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("▶️", pageCID.GetCacheID(), 1, cidCommandName)}
			} else {
				messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("◀️", pageCID.GetCacheID(), pageIndex-1, cidCommandName)}
				messageComponent = append(messageComponent, utils.MakeCIDPageComponent("▶️", pageCID.GetCacheID(), pageIndex+1, cidCommandName))
			}
		} else {
			messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("◀️", pageCID.GetCacheID(), pageIndex-1, cidCommandName)}
		}
	}
	actionsRow := utils.MakeActionsRow(messageComponent)
	listData := make([]string, 0, len(*res))
	for _, r := range *res {
		listData = append(listData, fmt.Sprintf("**%s**　%s (%s)", r.ID, r.Title, r.Alttitle))
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

// 月幕查詢遊戲名稱處理
func ymgalGetGameString(keyword string) (string, error) {
	logrus.Printf("ymgal查詢遊戲: %s", keyword)

	searchGameRes, err := ymgal.SearchGame(gojianfan.T2S(keyword))
	if err != nil {
		return "", err
	}

	if len(searchGameRes.Result) == 0 {
		return "", kurohelpercore.ErrSearchNoContent
	}

	sort.Slice(searchGameRes.Result, func(i, j int) bool {
		return searchGameRes.Result[i].Weights > searchGameRes.Result[j].Weights
	})

	return searchGameRes.Result[0].Name, nil
}
