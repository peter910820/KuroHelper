package handlers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"kurohelper/cache"
	kurohelpererrors "kurohelper/errors"
	"kurohelper/provider/erogs"
	"kurohelper/utils"
)

// 查詢音樂Handler
func SearchMusic(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.NewCID) {
	// 長時間查詢
	if cid == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
	}

	if i.Type == discordgo.InteractionApplicationCommand {
		opt, err := utils.GetOptions(i, "列表搜尋")
		if err != nil && errors.Is(err, kurohelpererrors.ErrOptionTranslateFail) {
			utils.HandleError(err, s, i)
			return
		}
		if opt == "" {
			erogsSearchMusic(s, i)
		} else {
			erogsSearchMusicList(s, i, cid)
		}
	} else {
		if !cid.GetCommandNameIsList() {
			erogsSearchMusic(s, i)
		} else {
			erogsSearchMusicList(s, i, cid)
		}

	}
}

// erogs查詢音樂處理
func erogsSearchMusic(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	var res *erogs.FuzzySearchMusicResponse
	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}
	idSearch, _ := regexp.MatchString(`^e\d+$`, keyword)
	res, err = erogs.GetMusicByFuzzy(keyword, idSearch)
	if err != nil {
		utils.HandleError(err, s, i)
		return
	}

	logrus.Printf("erogs查詢音樂: %s", keyword)

	musicData := make([]string, 0, len(res.GameCategories))
	for _, m := range res.GameCategories {
		musicData = append(musicData, m.GameName+" ("+m.GameModel+")"+" ("+m.Category+")")
	}

	singerList := strings.Split(res.Singers, ",")
	arrangementList := strings.Split(res.Arrangments, ",")
	lyricList := strings.Split(res.Lyrics, ",")
	compositionList := strings.Split(res.Compositions, ",")
	albumList := strings.Split(res.Album, ",")
	if res.PlayTime == "00:00:00" {
		res.PlayTime = "未收錄"
	}
	if res.ReleaseDate == "0001-01-01" {
		res.ReleaseDate = "未收錄"
	}

	embed := &discordgo.MessageEmbed{
		Title: res.MusicName,
		Color: 0xF8F8DF,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "音樂時長",
				Value:  res.PlayTime,
				Inline: true,
			},
			{
				Name:   "發行日期",
				Value:  res.ReleaseDate,
				Inline: true,
			},
			{
				Name:   "平均分數/樣本數",
				Value:  fmt.Sprintf("%.2f / %d", res.AvgTokuten, res.TokutenCount),
				Inline: true,
			},
			{
				Name:   "歌手",
				Value:  strings.Join(singerList, "\n"),
				Inline: false,
			},
			{
				Name:   "作詞",
				Value:  strings.Join(lyricList, "\n"),
				Inline: true,
			},
			{
				Name:   "作曲",
				Value:  strings.Join(compositionList, "\n"),
				Inline: true,
			},
			{
				Name:   "編曲",
				Value:  strings.Join(arrangementList, "\n"),
				Inline: true,
			},
			{
				Name:   "遊戲收錄",
				Value:  strings.Join(musicData, "\n"),
				Inline: false,
			},
			{
				Name:   "專輯",
				Value:  strings.Join(albumList, "\n"),
				Inline: false,
			},
		},
	}
	utils.InteractionEmbedRespond(s, i, embed, nil, true)
}

// erogs查詢音樂列表搜尋處理
func erogsSearchMusicList(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.NewCID) {
	var res *[]erogs.FuzzySearchListResponse
	var messageComponent []discordgo.MessageComponent
	var hasMore bool
	var count int
	var pageIndex int
	if cid == nil {
		keyword, err := utils.GetOptions(i, "keyword")
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		res, err = erogs.GetMusicListByFuzzy(keyword)
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}

		idStr := uuid.New().String()
		cache.Set(idStr, *res)

		// 計算筆數
		count = len(*res)

		hasMore = pagination(res, 0, false)

		if hasMore {
			cidCommandName := utils.MakeCIDCommandName(i.ApplicationCommandData().Name, true, "")
			messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("▶️", idStr, 1, cidCommandName)}
		}
	} else {
		// 處理CID
		pageCID := utils.PageCID{
			NewCID: *cid,
		}
		cacheValue, err := cache.Get(pageCID.GetCacheID())
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		resValue := cacheValue.([]erogs.FuzzySearchListResponse)
		res = &resValue

		// 計算筆數
		count = len(*res)

		// 資料分頁
		pageIndex, err = pageCID.GetPageIndex()
		if err != nil {
			utils.HandleError(err, s, i)
			return
		}
		hasMore = pagination(res, pageIndex, true)
		cidCommandName := utils.MakeCIDCommandName(cid.GetCommandName(), true, "")
		if hasMore {
			if pageIndex == 0 {
				messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("▶️", pageCID.GetCacheID(), 1, cidCommandName)}
			} else {
				messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("◀️", pageCID.GetCacheID(), pageIndex-1, cidCommandName)}
				messageComponent = append(messageComponent, utils.MakeCIDPageComponent("▶️", pageCID.GetCacheID(), pageIndex+1, cidCommandName))
			}
		} else {
			messageComponent = []discordgo.MessageComponent{utils.MakeCIDPageComponent("◀️", pageCID.GetCacheID(), pageIndex-1, cidCommandName)}
		}
	}
	actionsRow := utils.MakeActionsRow(messageComponent)

	listData := make([]string, 0, len(*res))
	for _, r := range *res {
		categoryData := strings.Split(r.Category, ",")
		listData = append(listData, fmt.Sprintf("e%-5s　%s (%s)", strconv.Itoa(r.ID), r.Name, strings.Join(categoryData, "/")))
	}
	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("音樂列表搜尋 (%d筆)", count),
		Color: 0xF8F8DF,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "ID/名稱",
				Value:  strings.Join(listData, "\n"),
				Inline: false,
			},
		},
	}

	if cid == nil {
		utils.InteractionEmbedRespond(s, i, embed, actionsRow, true)
	} else {
		utils.EditEmbedRespond(s, i, embed, actionsRow)
	}

}
