package handlers

import (
	"errors"
	"fmt"
	"kurohelper/cache"
	"kurohelper/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/kuro-helper/kurohelper-core/v3/vndb"
	"github.com/sirupsen/logrus"
)

const (
	searchBrandItemsPerPage = 7
	searchBrandCachePrefix  = "B@"
)

var (
	searchBrandColor = 0x00AA90
)

// 查詢公司品牌Handler(新版API)
func SearchBrandV2(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.CIDV2) {
	if cid == nil {
		vndbSearchBrandV2(s, i)
	} else {
		vndbSearchBrandWithCIDV2(s, i, cid)
	}
}

// 查詢公司品牌
func vndbSearchBrandV2(s *discordgo.Session, i *discordgo.InteractionCreate) {
	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		utils.HandleErrorWithInteractionApplicationCommand(err, s, i)
		return
	}

	logrus.WithField("interaction", i).Infof("vndb查詢公司品牌: %s", keyword)

	res, err := vndb.GetProducerByFuzzy(keyword, "")
	if err != nil {
		utils.HandleErrorWithInteractionApplicationCommand(err, s, i)
		return
	}

	// 產生快取ID
	cacheID := searchBrandCachePrefix + uuid.New().String()
	cache.SearchBrandCache.Set(cacheID, res)

	components, err := buildSearchBrandComponents(res, 1, cacheID)
	if err != nil {
		utils.HandleErrorWithInteractionApplicationCommand(err, s, i)
		return
	}
	utils.InteractionRespondV2(s, i, components)
}

// 產生查詢公司品牌(有CID版本)
//
// 目前只有翻頁事件
func vndbSearchBrandWithCIDV2(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.CIDV2) {
	if cid.GetBehaviorID() != utils.PageBehavior {
		utils.HandleErrorV2(errors.New("handlers: cid behavior id error"), s, i)
		return
	}

	pageCID, err := cid.ToPageCIDV2()
	if err != nil {
		utils.HandleErrorV2(err, s, i)
		return
	}

	cacheValue, err := cache.SearchBrandCache.Get(pageCID.CacheId)
	if err != nil {
		utils.HandleErrorV2(err, s, i)
		return
	}

	components, err := buildSearchBrandComponents(cacheValue, pageCID.Value, pageCID.CacheId)
	if err != nil {
		utils.HandleErrorV2(err, s, i)
		return
	}
	utils.InteractionRespondEditComplex(s, i, components)
}

// 產生查詢公司品牌的Components
func buildSearchBrandComponents(res *vndb.ProducerSearchResponse, currentPage int, cacheID string) ([]discordgo.MessageComponent, error) {
	producerName := res.Producer.Results[0].Name
	totalItems := len(res.Vn.Results)
	totalPages := (totalItems + searchBrandItemsPerPage - 1) / searchBrandItemsPerPage

	divider := true
	containerComponents := []discordgo.MessageComponent{
		discordgo.TextDisplay{
			Content: fmt.Sprintf("# %s\n遊戲筆數: **%d**\n⭐: vndb分數 📊:投票人數 🕒: 遊玩時間(小時)", producerName, totalItems),
		},
		discordgo.Separator{Divider: &divider},
	}

	// 計算當前頁的範圍
	start := (currentPage - 1) * searchBrandItemsPerPage
	end := min(start+searchBrandItemsPerPage, totalItems)
	pagedResults := res.Vn.Results[start:end]

	// 產生遊戲組件
	for idx, item := range pagedResults {
		itemNum := start + idx + 1
		title := item.Title
		if strings.TrimSpace(item.Alttitle) != "" {
			title = item.Alttitle
		}
		hours := item.LengthMinutes / 60
		itemContent := fmt.Sprintf("**%d. %s**\n⭐**%.1f**/📊**%d**/🕒**%02d**", itemNum, title, item.Rating, item.Votecount, hours)

		if strings.TrimSpace(item.Image.Thumbnail) == "" {
			containerComponents = append(containerComponents, discordgo.Section{
				Components: []discordgo.MessageComponent{
					discordgo.TextDisplay{
						Content: itemContent,
					},
				},
				Accessory: &discordgo.Thumbnail{
					Media: discordgo.UnfurledMediaItem{
						URL: "https://image.kurohelper.com/docs/neneGIF.gif",
					},
				},
			})
		} else {
			containerComponents = append(containerComponents, discordgo.Section{
				Components: []discordgo.MessageComponent{
					discordgo.TextDisplay{
						Content: itemContent,
					},
				},
				Accessory: &discordgo.Thumbnail{
					Media: discordgo.UnfurledMediaItem{
						URL: item.Image.Thumbnail,
					},
				},
			})
		}
	}

	// 產生翻頁組件
	pageComponents, err := utils.MakeChangePageComponent(currentPage, totalPages, cacheID)
	if err != nil {
		return nil, err
	} else {
		containerComponents = append(containerComponents,
			discordgo.Separator{Divider: &divider},
			pageComponents,
		)
	}

	// 組成完整組件回傳
	return []discordgo.MessageComponent{
		discordgo.Container{
			AccentColor: &searchBrandColor,
			Components:  containerComponents,
		},
	}, nil
}
