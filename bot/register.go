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
			Name:        "VNDB統計資料",
			Description: "取得vndb統計資料",
		},
		{
			Name:        "VNDB查詢遊戲",
			Description: "根據關鍵字取得vndb遊戲資料",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "keyword",
					Description: "關鍵字",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "brandid",
					Description: "指定的品牌ID(VNDB ID)",
					Required:    false,
				},
			},
		},
	}
}

func eroscapeCommands() {

}
