package bot

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"

	"kurohelper/handlers"
	"kurohelper/utils"
)

func OnInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		onInteractionApplicationCommand(s, i)
	case discordgo.InteractionMessageComponent:
		onInteractionMessageComponent(s, i)
	}
}

func onInteractionApplicationCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "vndb統計資料":
		go handlers.VndbStats(s, i)
	case "查詢指定遊戲":
		go handlers.VndbSearchGameByID(s, i)
	case "查詢公司品牌":
		go handlers.SearchBrand(s, i, nil, "")
	case "查詢創作者":
		go handlers.SearchCreator(s, i, nil)
	case "查詢音樂":
		go handlers.SearchMusic(s, i, nil)
	case "查詢遊戲":
		go handlers.SearchGame(s, i, nil)
	case "隨機遊戲":
		go handlers.RandomGame(s, i)
	case "加已玩":
		go handlers.AddHasPlayed(s, i, nil)
	case "清除快取":
		go handlers.CleanCache(s, i)
	case "查詢角色":
		go handlers.SearchCharacter(s, i, nil)
	}
}

func onInteractionMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	newCID := strings.Split(i.MessageComponentData().CustomID, "|")
	if len(newCID) > 1 {
		newOnInteractionMessageComponent(s, i, utils.NewCID(newCID))
		return
	}
	cid := strings.SplitN(i.MessageComponentData().CustomID, "::", 4)
	value, err := strconv.Atoi(cid[3])
	if err != nil {
		if cid[3] == "true" {
			value = 1
		} else {
			value = 0
		}
	}
	cidStruct := handlers.CustomID{
		ID:          cid[1],
		CommandName: cid[0],
		Type:        cid[2],
		Value:       value,
	}

	switch cidStruct.CommandName {
	case "查詢公司品牌(vndb)":
		go handlers.SearchBrand(s, i, &cidStruct, "vndb")
	case "查詢公司品牌(erogs)":
		go handlers.SearchBrand(s, i, &cidStruct, "erogs")
	case "查詢創作者":
		go handlers.SearchCreator(s, i, &cidStruct)
	case "查詢遊戲列表":
		go handlers.SearchGame(s, i, &cidStruct)
	case "查詢音樂列表":
		go handlers.SearchMusic(s, i, &cidStruct)
	case "查詢創作者列表":
		go handlers.SearchCreator(s, i, &cidStruct)
	case "查詢角色列表":
		go handlers.SearchCharacter(s, i, &cidStruct)
	}
}

func newOnInteractionMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate, newCID utils.NewCID) {
	switch newCID.GetCommandName() {
	// case CustomIDTypeAddWish:
	case "加已玩":
		go handlers.AddHasPlayed(s, i, newCID)
		// default:
		// 	logrus.Fatal(err)
	}
}
