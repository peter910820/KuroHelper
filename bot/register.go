package bot

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// 註冊命令
func RegisterCommand(s *discordgo.Session) {
	var guildCmds []*discordgo.ApplicationCommand
	var globalCmds []*discordgo.ApplicationCommand
	guildCmds = append(guildCmds, managementCommands()...)
	globalCmds = append(globalCmds, galgameCommands()...)
	globalCmds = append(globalCmds, vndbCommands()...)

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

// 群組專用管理指令，要使用群組內部整合管理複寫權限，預設是全部可見(私有)
func managementCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "清除快取",
			Description: "清除搜尋資料快取(管理員專用)",
		},
	}
}

// 主要專用指令(全域)
func galgameCommands() []*discordgo.ApplicationCommand {
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
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "列表搜尋",
					Description: "是否啟用列表搜尋",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "啟用",
							Value: "1",
						},
					},
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
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "列表搜尋",
					Description: "是否啟用列表搜尋",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "啟用",
							Value: "1",
						},
					},
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
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "列表搜尋",
					Description: "是否啟用列表搜尋",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "啟用",
							Value: "1",
						},
					},
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
		{
			Name:        "加已玩",
			Description: "把遊戲加到已玩(ErogameScape)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "keyword",
					Description: "關鍵字",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "complete_date",
					Description: "遊玩結束日期",
					Required:    false,
				},
			},
		},
		{
			Name:        "查詢角色",
			Description: "根據關鍵字查詢角色資料(ErogameScape)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "keyword",
					Description: "關鍵字",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "列表搜尋",
					Description: "是否啟用列表搜尋",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "啟用",
							Value: "1",
						},
					},
				},
			},
		},
		{
			Name:        "隨機遊戲",
			Description: "隨機一部Galgame(ymgal)",
		},
		{
			Name:        "個人資料",
			Description: "取得自己的個人資料(KuroHelper)",
		},
	}
}

// vndb專用指令(全域)
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
	}
}
