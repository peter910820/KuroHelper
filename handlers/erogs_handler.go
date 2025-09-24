package handlers

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"kurohelper/erogs"
	internalerrors "kurohelper/errors"
	"kurohelper/utils"
)

func ErogsFuzzySearchCreator(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		logrus.Error(err)
		utils.InteractionEmbedErrorRespond(s, i, "該功能目前異常，請稍後再嘗試", true)
		return
	}

	res, err := erogs.GetCreatorByFuzzy(keyword)
	if err != nil {
		logrus.Error(err)
		if errors.Is(err, internalerrors.ErrVndbNoResult) {
			utils.InteractionEmbedErrorRespond(s, i, "找不到任何結果喔", true)
		} else {
			utils.InteractionEmbedErrorRespond(s, i, "該功能目前異常，請稍後再嘗試", true)
		}
		return
	}

	link := ""
	if res.Creator[0].TwitterUsername != "" {
		link += fmt.Sprintf("[Twitter](https://x.com/%s) ", res.Creator[0].TwitterUsername)
	}
	if res.Creator[0].Blog != "" {
		link += fmt.Sprintf("[Blog](%s) ", res.Creator[0].Blog)
	}
	if res.Creator[0].Pixiv != nil {
		link += fmt.Sprintf("[Pixiv](https://www.pixiv.net/users/%d) ", *res.Creator[0].Pixiv)
	}

	// 根據遊戲評價排序
	sort.Slice(res.Creator[0].Games, func(i, j int) bool {
		return res.Creator[0].Games[i].Median > res.Creator[0].Games[j].Median // 大到小排序
	})

	gameData := make([]string, 0, len(res.Creator[0].Games))
	for i, g := range res.Creator[0].Games {
		shokushu := make([]string, 0, len(g.Shokushu))
		for _, s := range g.Shokushu {
			if s.Shubetu != 7 {
				shokushu = append(shokushu, fmt.Sprintf("*%s*", erogs.ShubetuMap[s.Shubetu]))
			} else {
				shokushu = append(shokushu, fmt.Sprintf("*%s*", s.ShubetuDetailName))
			}
		}
		gameData = append(gameData, fmt.Sprintf("%d. **%s**  (%s)", i+1, g.Gamename, strings.Join(shokushu, ", ")))
	}

	embed := &discordgo.MessageEmbed{
		Title:       res.Creator[0].Name,
		Color:       0x04108e,
		Description: link,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "歷代作品(評價排序)",
				Value:  strings.Join(gameData, "\n"),
				Inline: false,
			},
		},
	}

	utils.InteractionEmbedRespond(s, i, embed, nil, true)

}
