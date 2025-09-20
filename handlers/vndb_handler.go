package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"kurohelper/models"
	"kurohelper/utils"
	"kurohelper/vndb"
)

func VndbStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	r, err := vndb.GetStats()
	if err != nil {
		logrus.Error(err)
		utils.InteractionRespond(s, i, "該功能目前異常，請稍後再嘗試")
	}

	embed := &discordgo.MessageEmbed{
		Title: "VNDB統計資料",
		Color: 0x04108e,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "JSON內容",
				Value:  string(r),
				Inline: false,
			},
		},
	}
	utils.InteractionEmbedRespond(s, i, embed)
}

func VndbSearchGameByID(s *discordgo.Session, i *discordgo.InteractionCreate) {
	brandid, err := utils.GetOptions(i, "brandid")
	if err != nil {
		logrus.Error(err)
		utils.InteractionRespond(s, i, "該功能目前異常，請稍後再嘗試")
	}

	r, err := vndb.GetVnUseID(brandid)
	if err != nil {
		logrus.Error(err)
		utils.InteractionRespond(s, i, "該功能目前異常，請稍後再嘗試")
	}

	var res models.VndbResponse[models.VndbGetVnUseIDResponse]
	err = json.Unmarshal(r, &res)
	if err != nil {
		logrus.Error(err)
		utils.InteractionRespond(s, i, "該功能目前異常，請稍後再嘗試")
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

	var scenario string
	var art string
	var songs string
	var tmpAlias string
	for _, staff := range res.Results[0].Staff {
		if staff.Aliases != nil {
			tmpAlias = "("
			for _, alias := range staff.Aliases {
				tmpAlias += alias.Name + ", "
			}
			tmpAlias = tmpAlias[:len(tmpAlias)-2]
			tmpAlias += ")"

			if len(tmpAlias) == 2 {
				tmpAlias = ""
			}
		}

		switch staff.Role {
		case "scenario":
			scenario += fmt.Sprintf("*%s*%s\n", staff.Name, tmpAlias)
		case "art":
			art += fmt.Sprintf("*%s*%s\n", staff.Name, tmpAlias)
		case "songs":
			songs += fmt.Sprintf("*%s*%s\n", staff.Name, tmpAlias)
		}
	}

	embed := &discordgo.MessageEmbed{
		Title: gameTitle,
		Color: 0x04108e,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "公司名稱",
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
				Value:  fmt.Sprintf("%d(小時)/%d", res.Results[0].LengthMinutes/60, res.Results[0].LengthVotes),
				Inline: true,
			},
		},
	}
	utils.InteractionEmbedRespond(s, i, embed)
}
