package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	internalerrors "kurohelper/errors"
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
	utils.InteractionEmbedRespond(s, i, embed, nil, false)
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
	utils.InteractionEmbedRespond(s, i, embed, nil, true)
}

func VndbFuzzySearchBrand(s *discordgo.Session, i *discordgo.InteractionCreate, cid *models.VndbInteractionCustomID) {
	var res *models.VndbProducerSearchResponse
	var component discordgo.ActionsRow
	var hasMore bool
	// 第一次查詢
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

		companyType, err := utils.GetOptions(i, "type")
		if err != nil && errors.Is(err, internalerrors.ErrOptionTranslateFail) {
			logrus.Error(err)
			utils.InteractionEmbedErrorRespond(s, i, "該功能目前異常，請稍後再嘗試", true)
			return
		}

		res, err = vndb.ProducerFuzzySearch(keyword, companyType)
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
		hasMore = pagination(&(res.Vn.Results), 0, false)

		if hasMore {
			component = discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "▶️",
						Style:    discordgo.PrimaryButton,
						CustomID: fmt.Sprintf("SearchBrand_1_%s", idStr),
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
		resValue := cacheValue.(models.VndbProducerSearchResponse)
		res = &resValue
		// 資料分頁
		hasMore = pagination(&(res.Vn.Results), cid.Page, true)
		if hasMore {
			if cid.Page == 0 {
				component = discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "▶️",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("SearchBrand_1_%s", cid.Key),
						},
					},
				}
			} else {
				component = discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "◀️",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("SearchBrand_%d_%s", cid.Page-1, cid.Key),
						},
						&discordgo.Button{
							Label:    "▶️",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("SearchBrand_%d_%s", cid.Page+1, cid.Key),
						},
					},
				}
			}
		} else {
			component = discordgo.ActionsRow{
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

	if cid == nil {
		utils.InteractionEmbedRespond(s, i, embed, &component, true)
	} else {
		utils.EditEmbedRespond(s, i, embed, &component)
	}

}

// 資料分頁
func pagination(result *[]models.VndbProducerSearchVnResponse, page int, useCache bool) bool {
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
