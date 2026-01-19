package utils

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type SelectMenuItem struct {
	Title  string
	VndbID string
}

var (
	ErrMakeChangePageComponentIndexZero = errors.New("utils: make change page component page index parameters can not be zero")
)

func MakeSelectMenuComponent(cacheID string, gameData []SelectMenuItem) *discordgo.ActionsRow {
	menuOptions := []discordgo.SelectMenuOption{}

	for _, gd := range gameData {
		menuOptions = append(menuOptions, discordgo.SelectMenuOption{
			Label: gd.Title,
			Value: gd.VndbID,
		})
	}

	return &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{
				CustomID:    MakeSelectMenuCIDV2(cacheID),
				Placeholder: "選擇遊戲查看詳細",
				Options:     menuOptions,
			},
		},
	}
}

// 製作回到主頁的Component
func MakeBackToHomeComponent(cacheID string) *discordgo.ActionsRow {
	return &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "🏠回到主頁",
				Style:    discordgo.PrimaryButton,
				CustomID: MakeBackToHomeCIDV2(cacheID),
			},
		},
	}
}

// 製作翻頁Component
func MakeChangePageComponent(currentPage int, totalPage int, cacheID string) (*discordgo.ActionsRow, error) {
	if currentPage == 0 || totalPage == 0 {
		return nil, ErrMakeChangePageComponentIndexZero
	}

	// 中間的顯示頁數按鈕(不可點擊)
	tabButton := discordgo.Button{
		Label:    fmt.Sprintf("%d/%d", currentPage, totalPage),
		Style:    discordgo.SecondaryButton,
		Disabled: true,
		CustomID: MakePageCIDV2(currentPage, cacheID, true),
	}

	previousDisabled := false
	nextDisabled := false

	if currentPage == totalPage {
		nextDisabled = true
	}

	if currentPage == 1 {
		previousDisabled = true
	}

	// 上一頁按鈕
	previousButton := discordgo.Button{
		Label:    "◀️",
		Style:    discordgo.SecondaryButton,
		Disabled: previousDisabled,
		CustomID: MakePageCIDV2(currentPage-1, cacheID, false),
	}

	// 下一頁按鈕
	nextButton := discordgo.Button{
		Label:    "▶️",
		Style:    discordgo.SecondaryButton,
		Disabled: nextDisabled,
		CustomID: MakePageCIDV2(currentPage+1, cacheID, false),
	}

	return &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			previousButton,
			tabButton,
			nextButton,
		},
	}, nil
}
