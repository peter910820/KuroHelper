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
		go handlers.VndbFuzzySearchBrand(s, i, nil)
	// case "vndb模糊查詢創作家":
	// 	go handlers.VndbFuzzySearchStaff(s, i, nil)
	case "查詢創作者":
		go handlers.ErogsFuzzySearchCreator(s, i, nil)
	case "查詢音樂":
		go handlers.ErogsFuzzySearchMusic(s, i, nil)
	case "查詢遊戲":
		go handlers.ErogsFuzzySearchGame(s, i, nil)
	case "隨機遊戲":
		go handlers.RandomGameHandler(s, i)
	case "清除快取":
		go handlers.CleanCache(s, i)
	}
}

func onInteractionMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cid := strings.SplitN(i.MessageComponentData().CustomID, "::", 4)
	value, err := strconv.Atoi(cid[3])
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}
	cidStruct := handlers.CustomID{
		ID:          cid[1],
		CommandName: cid[0],
		Type:        cid[2],
		Value:       value,
	}

	switch cidStruct.CommandName {
	case "查詢公司品牌":
		go handlers.VndbFuzzySearchBrand(s, i, &cidStruct)
	case "查詢創作者":
		go handlers.ErogsFuzzySearchCreator(s, i, &cidStruct)
	}

}
