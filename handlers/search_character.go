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
	"kurohelper/utils"
)

func SearchCharacter(s *discordgo.Session, i *discordgo.InteractionCreate, cid *CustomID) {
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
			erogsSearchCharacter(s, i)
		} else {
			erogsSearchCharacterList(s, i, cid)
		}
	} else {
		erogsSearchCharacterList(s, i, cid)
	}
}

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

func erogsSearchCharacterList(s *discordgo.Session, i *discordgo.InteractionCreate, cid *CustomID) {
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

		res, err = erogs.GetCharacterListByFuzzy(keyword)
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
			messageComponent = []discordgo.MessageComponent{utils.MakePageComponent("▶️", "查詢角色列表", idStr, 1)}
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
