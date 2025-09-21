package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	internalerrors "kurohelper/errors"
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
	utils.InteractionEmbedRespond(s, i, embed, false)
}

func VndbSearchGameByID(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 長時間查詢
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	brandid, err := utils.GetOptions(i, "brandid")
	if err != nil {
		logrus.Error(err)
		utils.InteractionEmbedErrorRespond(s, i, "該功能目前異常，請稍後再嘗試", true)
		return
	}

	res, err := vndb.GetVnUseID(brandid)
	if err != nil {
		logrus.Error(err)
		if errors.Is(err, internalerrors.ErrVndbNoResult) {
			utils.InteractionEmbedErrorRespond(s, i, "找不到任何結果喔", true)
		} else {
			utils.InteractionEmbedErrorRespond(s, i, "該功能目前異常，請稍後再嘗試", true)
		}
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

	// staff block(待優化)
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

	// character block
	character := make([]string, 0, len(res.Results[0].Va))
	for _, va := range res.Results[0].Va {
		for _, vns := range va.Character.Vns {
			if vns.ID == brandid {
				if vns.Role == "primary" {
					character = append(character, fmt.Sprintf("**%s**(%s)", va.Character.Original, "主要角色"))
				} else {
					character = append(character, fmt.Sprintf("**%s**(%s)", va.Character.Original, "次要角色"))
				}
				break
			}
		}
	}

	// relations block
	relationsGame := make([]string, 0, len(res.Results[0].Relations))
	for _, rg := range res.Results[0].Relations {
		relationsGame = append(relationsGame, fmt.Sprintf("%s(%s)", rg.Titles[0].Title, rg.ID))
	}
	relationsGameDisplay := strings.Join(relationsGame, ", ")
	if strings.TrimSpace(relationsGameDisplay) == "" {
		relationsGameDisplay = "無"
	}

	embed := &discordgo.MessageEmbed{
		Title: gameTitle,
		Color: 0x04108e,
		Image: &discordgo.MessageEmbedImage{
			URL: res.Results[0].Image.Url,
		},
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
				Value:  strings.Join(character, ", "),
				Inline: false,
			},
			{
				Name:   "相關遊戲",
				Value:  relationsGameDisplay,
				Inline: false,
			},
		},
	}
	utils.InteractionEmbedRespond(s, i, embed, true)
}

func VndbFuzzySearchBrand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	companyType, err := utils.GetOptions(i, "type")
	if err != nil && errors.Is(err, internalerrors.ErrOptionTranslateFail) {
		logrus.Error(err)
		utils.InteractionEmbedErrorRespond(s, i, "該功能目前異常，請稍後再嘗試", true)
		return
	}

	res, err := vndb.ProducerFuzzySearch(keyword, companyType)
	if err != nil {
		logrus.Error(err)
		if errors.Is(err, internalerrors.ErrVndbNoResult) {
			utils.InteractionEmbedErrorRespond(s, i, "找不到任何結果喔", true)
		} else {
			utils.InteractionEmbedErrorRespond(s, i, "該功能目前異常，請稍後再嘗試", true)
		}
		return
	}

	// 資料分頁
	if len(res.Vn.Results) > 10 {
		res.Vn.Results = res.Vn.Results[:10]
	}

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

	utils.InteractionEmbedRespond(s, i, embed, true)
}
