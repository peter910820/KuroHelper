package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	kurohelpererrors "kurohelper/errors"
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

// 事件是InteractionApplicationCommand(使用斜線命令)的處理
func onInteractionApplicationCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "vndb統計資料":
		go handlers.VndbStats(s, i)
	case "查詢指定遊戲":
		go handlers.VndbSearchGameByID(s, i)
	case "查詢遊戲":
		go handlers.SearchGame(s, i, nil)
	case "查詢公司品牌":
		go handlers.SearchBrand(s, i, nil)
	case "查詢創作者":
		go handlers.SearchCreator(s, i, nil)
	case "查詢音樂":
		go handlers.SearchMusic(s, i, nil)
	case "查詢角色":
		go handlers.SearchCharacter(s, i, nil)
	case "加已玩":
		go handlers.AddHasPlayed(s, i, nil)
	case "清除快取":
		go handlers.CleanCache(s, i)
	case "隨機遊戲":
		go handlers.RandomGame(s, i)
	case "個人資料":
		go handlers.GetUserinfo(s, i)
	}
}

// 事件是InteractionMessageComponent(點擊按鈕)的處理
func onInteractionMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cidStringSlice := strings.Split(i.MessageComponentData().CustomID, "|")
	// 安全檢查，確保CID建立邏輯有誤的話不會出問題
	if len(cidStringSlice) > 1 {
		cid := utils.NewCID(cidStringSlice)
		switch cid.GetCommandName() {
		// case CustomIDTypeAddWish:
		case "查詢遊戲":
			go handlers.SearchGame(s, i, &cid)
		case "查詢公司品牌":
			go handlers.SearchBrand(s, i, &cid)
		case "查詢創作者":
			go handlers.SearchCreator(s, i, &cid)
		case "查詢音樂":
			go handlers.SearchMusic(s, i, &cid)
		case "查詢角色":
			go handlers.SearchCharacter(s, i, &cid)
		case "加已玩":
			go handlers.AddHasPlayed(s, i, &cid)
		}
	} else {
		utils.HandleError(kurohelpererrors.ErrCIDWrongFormat, s, i)
		return
	}
}
