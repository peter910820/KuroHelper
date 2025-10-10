package bot

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

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
		go handlers.FuzzySearchBrand(s, i, nil)
	case "查詢創作者":
		go handlers.FuzzySearchCreator(s, i, nil)
	case "查詢音樂":
		go handlers.FuzzySearchMusic(s, i, nil)
	case "查詢遊戲":
		go handlers.FuzzySearchGame(s, i, nil)
	case "隨機遊戲":
		go handlers.RandomGameHandler(s, i)
	case "加已玩":
		go handlers.AddHasPlayedHandler(s, i, nil)
	case "清除快取":
		go handlers.CleanCache(s, i)
	case "查詢角色":
		go handlers.FuzzySearchCharacter(s, i, nil)
	}
}

func onInteractionMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	newCID := strings.Split(i.MessageComponentData().CustomID, "|")
	if len(newCID) > 1 {
		newOnInteractionMessageComponent(s, i, newCID)
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
		go handlers.VndbFuzzySearchBrand(s, i, &cidStruct)
	case "查詢公司品牌(erogs)":
		go handlers.ErogsFuzzySearchBrand(s, i, &cidStruct)
	case "查詢創作者":
		go handlers.ErogsFuzzySearchCreator(s, i, &cidStruct)
	case "查詢遊戲列表":
		go handlers.FuzzySearchGame(s, i, &cidStruct)
	case "查詢音樂列表":
		go handlers.ErogsFuzzySearchMusicList(s, i, &cidStruct)
	case "查詢創作者列表":
		go handlers.ErogsFuzzySearchCreatorList(s, i, &cidStruct)
	case "查詢角色列表":
		go handlers.ErogsFuzzySearchCharacterList(s, i, &cidStruct)
	}
}

func newOnInteractionMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate, newCID []string) {
	value, err := strconv.Atoi(newCID[1])
	if err != nil {
		logrus.Fatal(err)
	}

	b, err := strconv.ParseBool(newCID[3])
	if err != nil {
		logrus.Fatal(err)
	}

	CIDType := utils.CustomIDType(value)
	switch CIDType {
	// case CustomIDTypeAddWish:
	case utils.CustomIDTypeAddHasPlayed:
		// 加已玩
		go handlers.AddHasPlayedHandler(s, i, &utils.NewCustomID[utils.AddHasPlayedArgs]{CommandName: newCID[0], Value: utils.AddHasPlayedArgs{CacheID: newCID[2], ConfirmMark: b}})
	default:
		logrus.Fatal(err)
	}
}
