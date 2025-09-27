package bot

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"

	"kurohelper/handlers"
	"kurohelper/models"
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
	case "vndb查詢指定遊戲":
		go handlers.VndbSearchGameByID(s, i)
	case "vndb模糊查詢品牌":
		go handlers.VndbFuzzySearchBrand(s, i, nil)
	// case "vndb模糊查詢創作家":
	// 	go handlers.VndbFuzzySearchStaff(s, i, nil)
	case "erogs模糊查詢創作者":
		go handlers.ErogsFuzzySearchCreator(s, i, nil)
	case "erogs模糊查詢音樂":
		go handlers.ErogsFuzzySearchMusic(s, i, nil)
	case "erogs模糊查詢遊戲":
		go handlers.ErogsFuzzySearchGame(s, i, nil)
	}
}

func onInteractionMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cid := strings.SplitN(i.MessageComponentData().CustomID, "_", 3)
	page, err := strconv.Atoi(cid[1])
	if err != nil {
		utils.InteractionRespond(s, i, "該功能目前異常，請稍後再嘗試")
		return
	}
	cidStruct := models.VndbInteractionCustomID{
		CommandName: cid[0],
		Page:        page,
		Key:         cid[2],
	}

	switch cidStruct.CommandName {
	case "SearchBrand":
		go handlers.VndbFuzzySearchBrand(s, i, &cidStruct)
	case "ErogsFuzzySearchCreator":
		go handlers.ErogsFuzzySearchCreator(s, i, &cidStruct)
	}

}
