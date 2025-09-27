package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"

	"kurohelper/models"
	vndbmodels "kurohelper/models/vndb"
	"kurohelper/utils"
	"kurohelper/vndb"
)

func SearchBrandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, cid *models.VndbInteractionCustomID) {
	var res *vndbmodels.ProducerSearchResponse
	var component *discordgo.ActionsRow
	var hasMore bool

	if cid == nil {
		// 長時間查詢
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

		// get options
		keyword, err := utils.GetOptions(i, "keyword")
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		brandType, err := utils.GetOptionOptional[string](i, "type")
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		source, err := utils.GetOptionOptional[string](i, "source")
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		if source == "erogs" {
			res, err = vndb.GetProducerByFuzzy(keyword, brandType) // 這邊先不接批評空間
			if err != nil {
				utils.HandleError(err, s, i)
				return
			}
		} else {
			res, err = vndb.GetProducerByFuzzy(keyword, brandType)
			if err != nil {
				utils.HandleError(err, s, i)
				return
			}
		}

		idStr := uuid.New().String()
		SetCache(idStr, *res)
		hasMore = pagination(&(res.Vn.Results), 0, false)

		if hasMore {
			component = &discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "▶️",
						Style:    discordgo.PrimaryButton,
						CustomID: fmt.Sprintf("SearchBrandNew_1_%s", idStr),
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
		resValue := cacheValue.(vndbmodels.ProducerSearchResponse)
		res = &resValue
		// 資料分頁
		hasMore = pagination(&(res.Vn.Results), cid.Page, true)
		if hasMore {
			if cid.Page == 0 {
				component = &discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "▶️",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("SearchBrandNew_1_%s", cid.Key),
						},
					},
				}
			} else {
				component = &discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "◀️",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("SearchBrandNew_%d_%s", cid.Page-1, cid.Key),
						},
						&discordgo.Button{
							Label:    "▶️",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("SearchBrandNew_%d_%s", cid.Page+1, cid.Key),
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
						CustomID: fmt.Sprintf("SearchBrandNew_%d_%s", cid.Page-1, cid.Key),
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
		utils.InteractionEmbedRespond(s, i, embed, component, true)
	} else {
		utils.EditEmbedRespond(s, i, embed, component)
	}

}
