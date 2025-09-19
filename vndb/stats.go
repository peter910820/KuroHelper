package vndb

import (
	"io"
	"kurohelper/utils"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func GetStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	resp, err := http.Get(os.Getenv("VNDB_ENDPOINT") + "/stats")
	if err != nil {
		logrus.Error(err)
		utils.SlashCommandRespond(s, i, "該功能目前異常，請稍後再嘗試")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		utils.SlashCommandRespond(s, i, "該功能目前異常，請稍後再嘗試")
	}

	if resp.StatusCode != 200 {
		logrus.Errorf("the server returned an error status code %d", resp.StatusCode)
		utils.SlashCommandRespond(s, i, "該功能目前異常，請稍後再嘗試")
	}

	embed := &discordgo.MessageEmbed{
		Title: "vndb統計資料",
		Color: 0x04108e,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "JSON內容",
				Value:  string(body),
				Inline: false,
			},
		},
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
