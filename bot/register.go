package bot

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func RegisterCommand(s *discordgo.Session) {
	var guildCmds []*discordgo.ApplicationCommand
	var globalCmds []*discordgo.ApplicationCommand
	guildCmds = append(guildCmds, managementCommands()...)
	globalCmds = append(globalCmds, vndbCommands()...)
	globalCmds = append(globalCmds, erogsCommands()...)

	// guild commands
	// main mangement group ID
	guildID := os.Getenv("SELF_GROUP_ID")
	for _, cmd := range guildCmds {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
		if err != nil {
			logrus.Errorf("register guild slash command failed: %s", err)
		}
	}
	// global commands
	for _, cmd := range globalCmds {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			logrus.Errorf("register global slash command failed: %s", err)
		}
	}

}

// 群組專用管理指令，要使用群組內部整合管理複寫權限，預設是全部可見
func managementCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "清除快取",
			Description: "清除搜尋資料快取(管理員專用)",
		},
	}
}

func vndbCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "vndb統計資料",
			Description: "取得VNDB統計資料(VNDB)",
		},
		{
			Name:        "查詢指定遊戲",
			Description: "根據VNDB ID查詢指定遊戲資料(VNDB)",
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
			Name:        "查詢公司品牌",
			Description: "根據關鍵字查詢公司品牌資料(VNDB)",
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
		// {
		// 	Name:        "vndb模糊查詢創作家",
		// 	Description: "根據關鍵字查詢創作家資料",
		// 	Options: []*discordgo.ApplicationCommandOption{
		// 		{
		// 			Type:        discordgo.ApplicationCommandOptionString,
		// 			Name:        "keyword",
		// 			Description: "關鍵字",
		// 			Required:    true,
		// 		},
		// 		{
		// 			Type:        discordgo.ApplicationCommandOptionString,
		// 			Name:        "role",
		// 			Description: "角色過濾",
		// 			Required:    false,
		// 			Choices: []*discordgo.ApplicationCommandOptionChoice{
		// 				{
		// 					Name:  "腳本家",
		// 					Value: "scenario",
		// 				},
		// 				{
		// 					Name:  "角色設計",
		// 					Value: "chardesign",
		// 				},
		// 				{
		// 					Name:  "美術",
		// 					Value: "art",
		// 				},
		// 				{
		// 					Name:  "音樂",
		// 					Value: "music",
		// 				},
		// 				{
		// 					Name:  "歌曲",
		// 					Value: "songs",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
	}
}

// 主要的Galgame搜尋來源，之後會整合指令變成可選搜尋來源，現在先分開
func erogsCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "查詢創作者",
			Description: "根據關鍵字查詢創作者資料(ErogameScape)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "keyword",
					Description: "關鍵字",
					Required:    true,
				},
			},
		},
		{
			Name:        "查詢音樂",
			Description: "根據關鍵字查詢音樂資料(ErogameScape)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "keyword",
					Description: "關鍵字",
					Required:    true,
				},
			},
		},
		{
			Name:        "查詢遊戲",
			Description: "根據關鍵字查詢遊戲資料(ErogameScape)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "keyword",
					Description: "關鍵字",
					Required:    true,
				},
			},
		},
		{
			Name:        "查詢公司品牌",
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
					Name:        "查詢資料庫選項",
					Description: "選擇查詢的資料庫",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "VNDB",
							Value: "1",
						},
					},
				},
			},
		},
	}
}
