package bot

import (
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
		go handlers.SearchBrand(s, i, nil)
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
}

func newOnInteractionMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate, newCID utils.NewCID) {
	switch newCID.GetCommandName() {
	// case CustomIDTypeAddWish:
	case "查詢遊戲":
		go handlers.SearchGame(s, i, &newCID)
	case "查詢公司品牌":
		go handlers.SearchBrand(s, i, &newCID)
	case "查詢創作者":
		go handlers.SearchCreator(s, i, &newCID)
	case "查詢音樂":
		go handlers.SearchMusic(s, i, &newCID)
	case "查詢角色":
		go handlers.SearchCharacter(s, i, &newCID)
	case "加已玩":
		go handlers.AddHasPlayed(s, i, &newCID)
		// default:
		// 	logrus.Fatal(err)
	}
}
