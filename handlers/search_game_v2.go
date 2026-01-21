package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"kurohelper/cache"
	"kurohelper/store"
	"kurohelper/utils"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kuro-helper/kurohelper-core/v3/erogs"
	"github.com/kuro-helper/kurohelper-core/v3/seiya"
	"github.com/kuro-helper/kurohelper-core/v3/vndb"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	kurohelperdb "github.com/kuro-helper/kurohelper-db/v3"
)

const (
	searchGameListItemsPerPage = 10
	searchGameListCachePrefix  = "G@"
)

var (
	searchGameListColor = 0xF8F8DF
	searchGameColor     = 0x04108e
)

// 查詢遊戲列表Handler(新版API)
func SearchGameV2(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.CIDV2) {
	if cid == nil {
		erogsSearchGameListV2(s, i)
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
		switch cid.GetBehaviorID() {
		case utils.PageBehavior:
			erogsSearchGameListWithCIDV2(s, i, cid)
		case utils.SelectMenuBehavior:
			erogsSearchGameWithSelectMenuCIDV2(s, i, cid)
		case utils.BackToHomeBehavior:
			erogsSearchGameWithBackToHomeCIDV2(s, i, cid)
		}
	}
}

// 查詢遊戲列表
func erogsSearchGameListV2(s *discordgo.Session, i *discordgo.InteractionCreate) {
	keyword, err := utils.GetOptions(i, "keyword")
	if err != nil {
		utils.HandleErrorV2(err, s, i, utils.InteractionRespondV2)
		return
	}

	logrus.WithField("interaction", i).Infof("erogs查詢遊戲列表: %s", keyword)

	// 將 keyword 轉成 base64 作為快取鍵
	cacheKey := searchGameListCachePrefix + base64.RawURLEncoding.EncodeToString([]byte(keyword))

	// 檢查快取是否存在
	cacheValue, err := cache.ErogsGameListStore.Get(cacheKey)
	if err == nil {
		// 快取存在，直接使用，不需要延遲傳送
		components, err := buildSearchGameComponents(&cacheValue, 1, cacheKey)
		if err != nil {
			utils.HandleErrorV2(err, s, i, utils.InteractionRespondV2)
			return
		}
		utils.InteractionRespondV2(s, i, components)
		return
	}

	// 快取不存在，需要查詢資料
	// 先發送延遲回應
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	res, err := erogs.GetGameListByFuzzy(keyword)
	if err != nil {
		utils.HandleErrorV2(err, s, i, utils.WebhookEditRespond)
		return
	}

	// 將查詢結果存入快取
	cache.ErogsGameListStore.Set(cacheKey, *res)

	components, err := buildSearchGameComponents(res, 1, cacheKey)
	if err != nil {
		utils.HandleErrorV2(err, s, i, utils.WebhookEditRespond)
		return
	}

	utils.WebhookEditRespond(s, i, components)
}

// 查詢遊戲列表(有CID版本)
func erogsSearchGameListWithCIDV2(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.CIDV2) {
	if cid.GetBehaviorID() != utils.PageBehavior {
		utils.HandleErrorV2(errors.New("handlers: cid behavior id error"), s, i, utils.InteractionRespondEditComplex)
		return
	}

	pageCID, err := cid.ToPageCIDV2()
	if err != nil {
		utils.HandleErrorV2(err, s, i, utils.InteractionRespondEditComplex)
		return
	}

	cacheValue, err := cache.ErogsGameListStore.Get(pageCID.CacheId)
	if err != nil {
		utils.HandleErrorV2(err, s, i, utils.InteractionRespondEditComplex)
		return
	}

	components, err := buildSearchGameComponents(&cacheValue, pageCID.Value, pageCID.CacheId)
	if err != nil {
		utils.HandleErrorV2(err, s, i, utils.InteractionRespondEditComplex)
		return
	}

	utils.WebhookEditRespond(s, i, components)
}

// 產生查詢遊戲列表的Components
func buildSearchGameComponents(res *[]erogs.FuzzySearchListResponse, currentPage int, cacheID string) ([]discordgo.MessageComponent, error) {
	totalItems := len(*res)
	totalPages := (totalItems + searchGameListItemsPerPage - 1) / searchGameListItemsPerPage

	divider := true
	containerComponents := []discordgo.MessageComponent{
		discordgo.TextDisplay{
			Content: fmt.Sprintf("# 遊戲搜尋\n遊戲筆數: **%d**", totalItems),
		},
		discordgo.Separator{Divider: &divider},
	}

	// 計算當前頁的範圍
	start := (currentPage - 1) * searchGameListItemsPerPage
	end := min(start+searchGameListItemsPerPage, totalItems)
	pagedResults := (*res)[start:end]

	gameMenuItems := []utils.SelectMenuItem{}

	// 產生遊戲列表組件
	listData := make([]string, 0, len(pagedResults))
	for _, r := range pagedResults {
		listData = append(listData, fmt.Sprintf("**e%-5s**　%s (%s)", strconv.Itoa(r.ID), r.Name, r.Category))

		gameMenuItems = append(gameMenuItems, utils.SelectMenuItem{
			Title: r.Name,
			ID:    strconv.Itoa(r.ID),
		})
	}

	// 將列表內容直接使用 TextDisplay 顯示（不使用 Section 避免 accessory 問題）
	if len(listData) > 0 {
		containerComponents = append(containerComponents, discordgo.TextDisplay{
			Content: strings.Join(listData, "\n"),
		})
	}

	// 產生選單組件
	selectMenuComponents := utils.MakeSelectMenuComponent(cacheID, gameMenuItems)

	// 產生翻頁組件
	pageComponents, err := utils.MakeChangePageComponent(currentPage, totalPages, cacheID)
	if err != nil {
		return nil, err
	}

	containerComponents = append(containerComponents,
		discordgo.Separator{Divider: &divider},
		selectMenuComponents,
		pageComponents,
	)

	// 組成完整組件回傳
	return []discordgo.MessageComponent{
		discordgo.Container{
			AccentColor: &searchGameListColor,
			Components:  containerComponents,
		},
	}, nil
}

// 查詢單一遊戲資料(有CID版本，從選單選擇)
func erogsSearchGameWithSelectMenuCIDV2(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.CIDV2) {
	if cid.GetBehaviorID() != utils.SelectMenuBehavior {
		utils.HandleErrorV2(errors.New("handlers: cid behavior id error"), s, i, utils.InteractionRespondEditComplex)
		return
	}

	selectMenuCID := cid.ToSelectMenuCIDV2()

	utils.WebhookEditRespond(s, i, []discordgo.MessageComponent{
		discordgo.Container{
			Components: []discordgo.MessageComponent{
				discordgo.TextDisplay{
					Content: "# ⌛ 正在跳轉，請稍候...",
				},
			},
		},
	})

	gameID := "e" + selectMenuCID.Value

	res, err := erogs.GetGameByFuzzy(gameID, true)
	if err != nil {
		utils.HandleErrorV2(err, s, i, utils.InteractionRespondEditComplex)
		return
	}

	logrus.WithField("guildID", i.GuildID).Infof("erogs查詢遊戲: %s", gameID)

	// 處理使用者資訊
	userID := utils.GetUserID(i)
	var userData string
	_, err = kurohelperdb.GetUserHasPlayed(userID, res.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleErrorV2(err, s, i, utils.InteractionRespondEditComplex)
			return
		}
	} else {
		userData += "✅"
	}
	_, err = kurohelperdb.GetUserInWish(userID, res.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleErrorV2(err, s, i, utils.InteractionRespondEditComplex)
			return
		}
	} else {
		userData += "❤️"
	}

	// 獲取 VNDB 資料
	vndbRating := 0.0
	vndbVotecount := 0
	var resVndb *vndb.BasicResponse[vndb.GetVnUseIDResponse]
	if strings.TrimSpace(res.VndbId) != "" {
		resVndb, err = vndb.GetVNByID(res.VndbId)
		if err != nil {
			utils.HandleErrorV2(err, s, i, utils.InteractionRespondEditComplex)
			return
		}
		vndbRating = resVndb.Results[0].Rating
		vndbVotecount = resVndb.Results[0].Votecount
	}

	// 處理 shubetu 資料
	shubetuData := make(map[int]map[int][]string) // map[shubetu_type]map[shubetu_detail]][]creator name + shubetu_detail_name

	for typeIdx := 1; typeIdx <= 6; typeIdx++ {
		shubetuData[typeIdx] = make(map[int][]string)
		for detailIdx := 1; detailIdx <= 3; detailIdx++ {
			shubetuData[typeIdx][detailIdx] = make([]string, 0)
		}
	}
	for _, shubetu := range res.CreatorShubetu {
		shubetuType := shubetu.ShubetuType
		detailType := shubetu.ShubetuDetailType
		creatorData := ""
		if shubetu.ShubetuDetailName == "" {
			creatorData = shubetu.CreatorName
		} else {
			creatorData = shubetu.CreatorName + " (" + shubetu.ShubetuDetailName + ")"
		}
		if shubetu.ShubetuType != 5 {
			shubetuData[shubetuType][1] = append(shubetuData[shubetuType][1], creatorData)
		} else {
			if detailType == 2 || detailType == 3 {
				shubetuData[shubetuType][2] = append(shubetuData[shubetuType][2], creatorData)
			} else {
				shubetuData[shubetuType][1] = append(shubetuData[shubetuType][1], creatorData)
			}
		}
	}

	// 處理其他資訊
	switch res.Okazu {
	case "true":
		res.Okazu = "拔作"
	case "false":
		res.Okazu = "非拔作"
	default:
		res.Okazu = ""
	}

	switch res.Erogame {
	case "true":
		res.Erogame = "18禁"
	case "false":
		res.Erogame = "全年齡"
	default:
		res.Erogame = ""
	}

	otherInfo := ""
	if res.Erogame == "" && res.Okazu == "" {
		otherInfo = "無"
	} else if res.Erogame == "" || res.Okazu == "" {
		otherInfo = res.Erogame + res.Okazu
	} else {
		otherInfo = res.Okazu + " / " + res.Erogame
	}

	// 處理排名和顏色
	junni := 0x04108e
	rank := ""
	if res.Junni == 0 || res.Junni > 500 {
		junni = 0x04108e // Default
	} else if res.Junni <= 50 {
		junni = 0xFFD700 // Gold
		rank = "批評空間 TOP 50"
	} else if res.Junni <= 100 {
		junni = 0xC0C0C0 // Silver
		rank = "批評空間 TOP 100"
	} else {
		junni = 0xCD7F32 // Bronze
		rank = "批評空間 TOP 500"
	}

	// 用批評空間回來的遊戲名對誠也做模糊搜尋
	seiyaURL := seiya.GetGuideURL(res.Gamename)
	if seiyaURL != "" {
		rank += "  " + fmt.Sprintf("[誠也攻略](%s)", seiyaURL)
	}
	erogsURL := "https://erogamescape.dyndns.org/~ap2/ero/toukei_kaiseki/game.php?game=" + fmt.Sprint(res.ID)
	rank += "  " + fmt.Sprintf("[批評空間](%s)", erogsURL)
	if res.VndbId != "" {
		vndbURL := "https://vndb.org/" + res.VndbId
		rank += "  " + fmt.Sprintf("[VNDB](%s)", vndbURL)
	}

	vndbData := "無"
	if vndbVotecount != 0 {
		vndbData = fmt.Sprintf("%.1f/%d", vndbRating, vndbVotecount)
	}

	// 過濾圖片
	imageURL := res.BannerUrl
	if i.GuildID != "" {
		// guild
		if _, ok := store.GuildDiscordAllowList[i.GuildID]; !ok {
			imageURL = ""
		}
	} else {
		// DM
		if _, ok := store.GuildDiscordAllowList[userID]; !ok {
			imageURL = ""
		}
	}

	// 構建 Components V2 格式
	divider := true
	contentParts := []string{}

	// 品牌名稱
	if strings.TrimSpace(res.BrandName) != "" {
		contentParts = append(contentParts, fmt.Sprintf("**品牌名稱**\n%s", res.BrandName))
	}

	// 排名和連結
	if strings.TrimSpace(rank) != "" {
		contentParts = append(contentParts, rank)
	}

	// 劇本
	if len(shubetuData[2][1]) > 0 {
		contentParts = append(contentParts, fmt.Sprintf("**劇本**\n%s", strings.Join(shubetuData[2][1], " / ")))
	}

	// 原畫
	if len(shubetuData[1][1]) > 0 {
		contentParts = append(contentParts, fmt.Sprintf("**原畫**\n%s", strings.Join(shubetuData[1][1], " / ")))
	}

	// 主角群CV
	if len(shubetuData[5][1]) > 0 {
		contentParts = append(contentParts, fmt.Sprintf("**主角群CV**\n%s", strings.Join(shubetuData[5][1], " / ")))
	}

	// 配角群CV
	if len(shubetuData[5][2]) > 0 {
		contentParts = append(contentParts, fmt.Sprintf("**配角群CV**\n%s", strings.Join(shubetuData[5][2], " / ")))
	}

	// 歌手
	if len(shubetuData[6][1]) > 0 {
		contentParts = append(contentParts, fmt.Sprintf("**歌手**\n%s", strings.Join(shubetuData[6][1], " / ")))
	}

	// 音樂
	if len(shubetuData[3][1]) > 0 {
		contentParts = append(contentParts, fmt.Sprintf("**音樂**\n%s", strings.Join(shubetuData[3][1], " / ")))
	}

	// 分數資訊
	evaluationText := fmt.Sprintf("**批評空間分數/樣本數**\n%s / %s", res.Median, res.TokutenCount)
	vndbText := fmt.Sprintf("**vndb分數/樣本數**\n%s", vndbData)
	contentParts = append(contentParts, evaluationText, vndbText)

	// 遊玩時數
	if strings.TrimSpace(res.TotalPlayTimeMedian) != "" {
		contentParts = append(contentParts, fmt.Sprintf("**遊玩時數**\n%s", res.TotalPlayTimeMedian))
	}

	// 開始理解遊戲樂趣時數
	if strings.TrimSpace(res.TimeBeforeUnderstandingFunMedian) != "" {
		contentParts = append(contentParts, fmt.Sprintf("**開始理解遊戲樂趣時數**\n%s", res.TimeBeforeUnderstandingFunMedian))
	}

	// 發行機種
	if strings.TrimSpace(res.Model) != "" {
		contentParts = append(contentParts, fmt.Sprintf("**發行機種**\n%s", res.Model))
	}

	// 類型
	if strings.TrimSpace(res.Genre) != "" {
		contentParts = append(contentParts, fmt.Sprintf("**類型**\n%s", res.Genre))
	}

	// 其他資訊
	contentParts = append(contentParts, fmt.Sprintf("**其他資訊**\n%s", otherInfo))

	// 合併所有內容
	fullContent := strings.Join(contentParts, "\n\n")

	containerComponents := []discordgo.MessageComponent{
		discordgo.TextDisplay{
			Content: fmt.Sprintf("# %s**%s(%s)**", userData, res.Gamename, res.SellDay),
		},
		discordgo.Separator{Divider: &divider},
		discordgo.TextDisplay{
			Content: fullContent,
		},
		discordgo.Separator{Divider: &divider},
	}

	// 如果有圖片，放在最下面（使用 MediaGallery）
	if strings.TrimSpace(imageURL) != "" {
		containerComponents = append(containerComponents, discordgo.MediaGallery{
			Items: []discordgo.MediaGalleryItem{
				{
					Media: discordgo.UnfurledMediaItem{
						URL: imageURL,
					},
				},
			},
		})
	}

	containerComponents = append(containerComponents, utils.MakeBackToHomeComponent(selectMenuCID.CacheId))

	components := []discordgo.MessageComponent{
		discordgo.Container{
			AccentColor: &junni,
			Components:  containerComponents,
		},
	}

	utils.InteractionRespondEditComplex(s, i, components)
}

// 返回遊戲列表主頁(有CID版本)
func erogsSearchGameWithBackToHomeCIDV2(s *discordgo.Session, i *discordgo.InteractionCreate, cid *utils.CIDV2) {
	if cid.GetBehaviorID() != utils.BackToHomeBehavior {
		utils.HandleErrorV2(errors.New("handlers: cid behavior id error"), s, i, utils.InteractionRespondEditComplex)
		return
	}

	backToHomeCID := cid.ToBackToHomeCIDV2()

	cacheValue, err := cache.ErogsGameListStore.Get(backToHomeCID.CacheId)
	if err != nil {
		utils.HandleErrorV2(err, s, i, utils.InteractionRespondEditComplex)
		return
	}

	components, err := buildSearchGameComponents(&cacheValue, 1, backToHomeCID.CacheId)
	if err != nil {
		utils.HandleErrorV2(err, s, i, utils.InteractionRespondEditComplex)
		return
	}
	utils.InteractionRespondEditComplex(s, i, components)
}
