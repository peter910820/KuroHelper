package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

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
		Title: "vndb統計資料",
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
