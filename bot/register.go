package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func RegisterCommand(s *discordgo.Session) {
	var cmds []*discordgo.ApplicationCommand
	cmds = append(cmds, vndbCommands()...)

	for _, cmd := range cmds {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			logrus.Errorf("register slash command failed: %s", err)
		}
	}
}

func vndbCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "vndb統計資料",
			Description: "取得VNDB統計資料",
		},
		{
			Name:        "vndb查詢指定遊戲",
			Description: "根據VNDB ID查詢指定遊戲資料",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "brandid",
					Description: "VNDB ID",
					Required:    true,
				},
			},
		},
		{
			Name:        "vndb模糊查詢品牌",
			Description: "根據關鍵字查詢公司品牌資料",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "keyword",
					Description: "關鍵字",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "type",
					Description: "公司性質",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "公司",
							Value: "company",
						},
						{
							Name:  "個人",
							Value: "individual",
						},
						{
							Name:  "同人社團",
							Value: "amateur",
						},
					},
				},
			},
		},
	}
}

func eroscapeCommands() {

}
