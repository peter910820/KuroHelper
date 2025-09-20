package handlers

import (
	"errors"

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
	utils.InteractionEmbedRespond(s, i, embed)
}

func VndbSearchGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		logrus.Error(err)
		utils.InteractionRespond(s, i, "該功能目前異常，請稍後再嘗試")
	}
	brandid, err := utils.GetOptions(i, "brandid")
	if err != nil && errors.Is(err, internalerrors.ErrOptionTranslateFail) {
		logrus.Error(err)
		utils.InteractionRespond(s, i, "該功能目前異常，請稍後再嘗試")
	}
}
