package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"kurohelper/provider/vndb"
	"kurohelper/utils"
)

// vndb統計資料Handler
func VndbStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	r, err := vndb.GetStats()
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "VNDB統計資料",
		Color: 0x04108e,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "角色收錄數量",
				Value:  fmt.Sprintf("%d筆", r.Chars),
				Inline: true,
			},
			{
				Name:   "公司/品牌收錄數量",
				Value:  fmt.Sprintf("%d筆", r.Producers),
				Inline: true,
			},
			{
				Name:   "發行版本收錄數量",
				Value:  fmt.Sprintf("%d筆", r.Releases),
				Inline: true,
			},
			{
				Name:   "標籤收錄數量",
				Value:  fmt.Sprintf("%d筆", r.Tags),
				Inline: true,
			},
			{
				Name:   "角色特徵收錄數量",
				Value:  fmt.Sprintf("%d筆", r.Traits),
				Inline: true,
			},
			{
				Name:   "視覺小說收錄數量",
				Value:  fmt.Sprintf("%d筆", r.VN),
				Inline: true,
			},
		},
	}
	utils.InteractionEmbedRespond(s, i, embed, nil, false)
}

// 查詢指定遊戲Handler
func VndbSearchGameByID(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 長時間查詢
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	brandid, err := utils.GetOptions(i, "brandid")
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	res, err := vndb.GetVnUseID(brandid)
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
