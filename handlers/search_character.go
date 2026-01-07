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
	"github.com/sirupsen/logrus"

	"kurohelper/cache"
	kurohelperrerrors "kurohelper/errors"
	"kurohelper/utils"

	"github.com/kuro-helper/kurohelper-core/v3/bangumi"
	"github.com/kuro-helper/kurohelper-core/v3/erogs"
	"github.com/kuro-helper/kurohelper-core/v3/vndb"
)

// 查詢角色Handler
func SearchCharacter(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.NewCID) {
	// 長時間查詢
	if cid == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
	}

	if i.Type == discordgo.InteractionApplicationCommand {
		optList, err := utils.GetOptions(i, "列表搜尋")
		if err != nil && errors.Is(err, kurohelperrerrors.ErrOptionTranslateFail) {
			utils.HandleError(err, s, i)
			return
		}
		optDB, err := utils.GetOptions(i, "查詢資料庫選項")
		if err != nil && errors.Is(err, kurohelperrerrors.ErrOptionTranslateFail) {
			utils.HandleError(err, s, i)
			return
		}
		switch optDB {
		case "":
			fallthrough
		case "1":
			if optList == "" {
				vndbSearchCharacter(s, i)
			} else {
				vndbSearchCharacterList(s, i, cid)
			}
		case "2":
			if optList == "" {
				erogsSearchCharacter(s, i)
			} else {
				erogsSearchCharacterList(s, i, cid)
			}
		case "3":
			if optList == "" {
				bangumiSearchCharacter(s, i)
			} else {
				utils.HandleError(kurohelperrerrors.ErrBangumiCharacterListSearchNotSupported, s, i)
			}
		}
	} else {
		commandNameProvider := cid.GetCommandNameProvider()

		switch commandNameProvider {
		case "vndb":
			if !cid.GetCommandNameIsList() {
				vndbSearchCharacter(s, i)
			} else {
				vndbSearchCharacterList(s, i, cid)
			}
		case "erogs":
			if !cid.GetCommandNameIsList() {
				erogsSearchCharacter(s, i)
			} else {
				erogsSearchCharacterList(s, i, cid)
			}
		}
	}
}

// erogs查詢角色處理
func erogsSearchCharacter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	var res *erogs.FuzzySearchCharacterResponse
	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}
	idSearch, _ := regexp.MatchString(`^e\d+$`, keyword)
	res, err = erogs.GetCharacterByFuzzy(keyword, idSearch)
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	logrus.WithField("guildID", i.GuildID).Infof("erogs查詢角色: %s", keyword)

	if res.Birthday == "公式情報なし" {
		res.Birthday = "未收錄"
	}
	if res.BloodType == "公式情報なし" {
		res.BloodType = "未收錄"
	}
	bodyData := make([]string, 2)
	if res.Bust == "" && res.Waist == "" && res.Hip == "" {
		bodyData[0] = "未收錄"
	} else {
		bodyData[0] = res.Bust + "/" + res.Waist + "/" + res.Hip
	}
	bodyData[1] = res.Cup
	roleData := erogs.Role[res.Role]

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s (%s)", res.CharacterName, res.CreatorName),
		Color: 0xF8F8DF,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "登場於",
				Value:  fmt.Sprintf("%s (%s)", res.GameName, roleData),
				Inline: false,
			},
			{
				Name:   "性別",
				Value:  res.Sex,
				Inline: true,
			},
			{
				Name:   "年齡",
				Value:  res.Age,
				Inline: true,
			},
			{
				Name:   "身高/體重",
				Value:  fmt.Sprintf("%s/%s", res.Height, res.Weight),
				Inline: true,
			},
			{
				Name:   "生日",
				Value:  res.Birthday,
				Inline: true,
			},
			{
				Name:   "血型",
				Value:  res.BloodType,
				Inline: true,
			},
			{
				Name:   "三圍/罩杯",
				Value:  strings.Join(bodyData, "/"),
				Inline: true,
			},
			{
				Name:   "角色敘述",
				Value:  res.FormalExplain,
				Inline: true,
			},
		},
	}
	utils.InteractionEmbedRespond(s, i, embed, nil, true)
}

// erogs查詢角色列表搜尋處理
func erogsSearchCharacterList(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.NewCID) {
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

		res, err = erogs.GetCharacterListByFuzzy(keyword)
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		idStr := uuid.New().String()
		cache.SearchCache.Set(idStr, *res)

		// 計算筆數
		count = len(*res)

		hasMore = pagination(res, 0, false)

		if hasMore {
			cidCommandName := utils.MakeCIDCommandName(i.ApplicationCommandData().Name, true, "erogs")
			messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("▶️", idStr, 1, cidCommandName)}
		}
	} else {
		// 處理CID
		pageCID := utils.PageCID{
			NewCID: *cid,
		}
		cacheValue, err := cache.SearchCache.Get(pageCID.GetCacheID())
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
		listData = append(listData, fmt.Sprintf("e%-5s　%s (%s)(%s)", strconv.Itoa(r.ID), r.Name, r.Category, r.Model))
	}
	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("角色列表搜尋 (%d筆)", count),
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

// Bangumi查詢角色處理
func bangumiSearchCharacter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	var res *bangumi.Character
	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}
	res, err = bangumi.GetCharacterByFuzzy(keyword)
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}
	logrus.WithField("guildID", i.GuildID).Infof("Bangumi查詢角色: %s", keyword)
	nameData := ""
	if res.NameCN == "" {
		nameData = res.Name
	} else {
		nameData = fmt.Sprintf("%s (%s)", res.Name, res.NameCN)
	}
	image := generateImage(i, res.Image)
	embed := &discordgo.MessageEmbed{
		Title:       nameData,
		Description: res.Summary,
		Color:       0xF8F8DF,
		Image:       image,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "別名",
				Value:  strings.Join(res.Aliases, "/"),
				Inline: false,
			},
			{
				Name:   "性別",
				Value:  res.Gender,
				Inline: true,
			},
			{
				Name:   "年齡",
				Value:  res.Age,
				Inline: true,
			},
			{
				Name:   "身高/體重",
				Value:  fmt.Sprintf("%s/%s", res.Height, res.Weight),
				Inline: true,
			},
			{
				Name:   "生日",
				Value:  res.BirthDay,
				Inline: true,
			},
			{
				Name:   "血型",
				Value:  res.BloodType,
				Inline: true,
			},
			{
				Name:   "三圍",
				Value:  res.BWH,
				Inline: true,
			},
			{
				Name:   "CV",
				Value:  strings.Join(res.CV, "/"),
				Inline: false,
			},
			{
				Name:   "登場於",
				Value:  strings.Join(res.Game, "\n"),
				Inline: true,
			},
			{
				Name:   "其他",
				Value:  strings.Join(res.Other, "\n"),
				Inline: true,
			},
		},
	}
	utils.InteractionEmbedRespond(s, i, embed, nil, true)
}

func vndbSearchCharacter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	var res *vndb.CharacterSearchResponse
	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	idSearch, _ := regexp.MatchString(`^c\d+$`, keyword)
	if idSearch {
		logrus.WithField("guildID", i.GuildID).Infof("vndb查詢角色ID: %s", keyword)
		res, err = vndb.GetCharacterByID(keyword)
	} else {
		logrus.WithField("guildID", i.GuildID).Infof("vndb查詢角色: %s", keyword)
		res, err = vndb.GetCharacterByFuzzy(keyword)
	}

	if err != nil {
		utils.HandleError(err, s, i)
		return
	}
	nameData := res.Name
	heightData := "未收錄"
	weightData := "未收錄"
	BWHData := "未收錄"
	ageData := "未收錄"
	birthDayData := "未收錄"
	sexData := ""
	genderData := ""
	vnData := make([]string, 0, len(res.VNs))
	if res.Original != "" {
		nameData = fmt.Sprintf("%s (%s)", res.Original, res.Name)
	}
	if len(res.Aliases) == 0 {
		res.Aliases = []string{"未收錄"}
	}
	if res.Description == "" {
		res.Description = "無角色敘述"
	}
	if res.BloodType == "" {
		res.BloodType = "未收錄"
	}
	if res.Height != 0 {
		heightData = strconv.Itoa(res.Height) + "cm"
	}
	if res.Weight != 0 {
		weightData = strconv.Itoa(res.Weight) + "kg"
	}
	if res.Bust != 0 && res.Waist != 0 && res.Hips != 0 {
		BWHData = fmt.Sprintf("%d/%d/%d", res.Bust, res.Waist, res.Hips)
	}
	if res.Cup == "" {
		res.Cup = "未收錄"
	}
	if res.Age != nil {
		ageData = strconv.Itoa(*res.Age)
	}
	if res.Birthday != [2]int{} {
		birthDayData = fmt.Sprintf("%d月%d號", res.Birthday[0], res.Birthday[1])
	}
	if res.Sex == [2]string{} {
		sexData = "未收錄"
	} else if res.Sex[0] != res.Sex[1] {
		sexData = fmt.Sprintf("%s/||%s||", vndb.Sex[res.Sex[0]], vndb.Sex[res.Sex[1]])
	} else {
		sexData = vndb.Sex[res.Sex[0]]
	}
	if res.Gender == [2]string{} {
		genderData = "未收錄"
	} else if res.Gender[0] != res.Gender[1] {
		genderData = fmt.Sprintf("%s/||%s||", vndb.Gender[res.Gender[0]], vndb.Gender[res.Gender[1]])
	} else {
		genderData = vndb.Gender[res.Gender[0]]
	}
	if len(res.VNs) == 0 {
		vnData = append(vnData, "未收錄")
	} else {
		sort.Slice(res.VNs, func(i, j int) bool { // 依照角色定位排序
			return vndb.RolePriority[res.VNs[i].Role] < vndb.RolePriority[res.VNs[j].Role]
		})
		for _, vn := range res.VNs {
			titleData := vn.Title
			for _, title := range vn.Titles {
				if title.Main { // 抓取原文標題
					titleData = title.Title
					break
				}
			}
			if vn.Spoiler == 0 {
				vnData = append(vnData, fmt.Sprintf("%s (%s)", titleData, vndb.Role[vn.Role]))
			} else {
				vnData = append(vnData, fmt.Sprintf("||%s (%s)||", titleData, vndb.Role[vn.Role]))
			}
		}
	}

	res.Description = vndb.ConvertBBCodeToMarkdown(res.Description)
	image := generateImage(i, res.Image.URL)
	embed := &discordgo.MessageEmbed{
		Title:       nameData,
		Description: res.Description, // 敘述放在Description欄位以避免超過字數限制
		Color:       0xF8F8DF,
		Image:       image,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "別名",
				Value:  strings.Join(res.Aliases, "/"),
				Inline: true,
			},
			{
				Name:   "CV",
				Value:  strings.Join(res.Vas, "/"),
				Inline: true,
			},
			{
				Name:   "生日",
				Value:  birthDayData,
				Inline: true,
			},
			{
				Name:   "生理性別",
				Value:  sexData,
				Inline: true,
			},
			{
				Name:   "性別認同",
				Value:  genderData,
				Inline: true,
			},
			{
				Name:   "身高/體重",
				Value:  fmt.Sprintf("%s/%s", heightData, weightData),
				Inline: true,
			},
			{
				Name:   "年齡",
				Value:  ageData,
				Inline: true,
			},
			{
				Name:   "血型",
				Value:  res.BloodType,
				Inline: true,
			},
			{
				Name:   "三圍",
				Value:  BWHData,
				Inline: true,
			},
			{
				Name:   "登場於",
				Value:  strings.Join(vnData, "\n"),
				Inline: true,
			},
		},
	}
	utils.InteractionEmbedRespond(s, i, embed, nil, true)
}

// vndb查詢角色列表搜尋處理
func vndbSearchCharacterList(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.NewCID) {
	var res *[]vndb.CharacterSearchResponse
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

		res, err = vndb.GetCharacterListByFuzzy(keyword)
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		idStr := uuid.New().String()
		cache.SearchCache.Set(idStr, *res)

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
		cacheValue, err := cache.SearchCache.Get(pageCID.GetCacheID())
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		resValue := cacheValue.([]vndb.CharacterSearchResponse)
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
		nameData := r.Name
		vnData := []string{}
		if r.Original != "" { // 有原文名稱則顯示原文名稱
			nameData = r.Original
		}
		sort.Slice(r.VNs, func(i, j int) bool { // 依照角色定位排序
			return vndb.RolePriority[r.VNs[i].Role] < vndb.RolePriority[r.VNs[j].Role]
		})
		for _, vn := range r.VNs {
			vnNameData := vn.Title
			for _, title := range vn.Titles {
				if title.Main { // 抓取原文標題
					vnNameData = title.Title
					break
				}
			}
			vnData = append(vnData, fmt.Sprintf("%s (%s)", vnNameData, vndb.Role[vn.Role]))
			if len(vnData) >= 2 { // 一個角色最多顯示兩個遊戲
				break
			}
		}
		listData = append(listData, fmt.Sprintf("%-6s　%s  ( %s )", r.ID, nameData, strings.Join(vnData, "\n")))
	}
	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("角色列表搜尋 (%d筆)", count),
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
