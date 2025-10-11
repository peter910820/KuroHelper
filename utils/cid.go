package utils

import (
	"fmt"
	kurohelpererrors "kurohelper/errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type (
	NewCID []string

	CustomIDType int

	CustomIDCommandName string
)

type (
	PageCID struct {
		NewCID
	}

	SortCID struct {
		NewCID
	}

	AddWishCID struct {
		NewCID
	}

	AddHasPlayedCID struct {
		NewCID
	}
)

const (
	CustomIDTypePage CustomIDType = iota + 1
	CustomIDTypeSort
	CustomIDTypeAddWish
	CustomIDTypeAddHasPlayed
)

// 獲取CID中的CommandName欄位
func (cid NewCID) GetCommandName() string {
	return strings.Split([]string(cid)[0], "/")[0]
}

// 獲取CID中的CommandName欄位，並回傳是否是列表行為
//
// 這邊是安全行為，如果是沒有列表行為的狀況這邊會單純回傳False
func (cid NewCID) GetCommandNameIsList() bool {
	commandName := strings.Split([]string(cid)[0], "/")
	if len(commandName) <= 1 {
		return false
	} else {
		if commandName[1] == "list" {
			return true
		} else {
			return false
		}
	}
}

// 獲取CID中的CommandName欄位，並回傳供應商
//
// 這邊是安全行為，如果是沒有列表行為的狀況這邊會單純回傳空字串
func (cid NewCID) GetCommandNameProvider() string {
	commandName := strings.Split([]string(cid)[0], "/")
	if len(commandName) <= 2 {
		return ""
	} else {
		return commandName[2]
	}
}

// 獲取CID類型
func (cid NewCID) GetCIDType() (CustomIDType, error) {
	value, err := strconv.Atoi([]string(cid)[1])
	if err != nil {
		return 0, kurohelpererrors.ErrCIDGetParameterFailed
	}
	return CustomIDType(value), nil
}

// 獲取CID中的CacheID欄位，之後用於查找Cache
func (cid NewCID) GetCacheID() string {
	return cid[2]
}

func (cid PageCID) GetPageIndex() (int, error) {
	value, err := strconv.Atoi([]string(cid.NewCID)[3])
	if err != nil {
		return 0, kurohelpererrors.ErrCIDGetParameterFailed
	}
	return value, nil
}

func (cid AddWishCID) GetConfirmMark() (bool, error) {
	value, err := strconv.ParseBool([]string(cid.NewCID)[3])
	if err != nil {
		return false, kurohelpererrors.ErrCIDGetParameterFailed
	}
	return value, nil
}

func (cid AddHasPlayedCID) GetConfirmMark() (bool, error) {
	value, err := strconv.ParseBool([]string(cid.NewCID)[3])
	if err != nil {
		return false, kurohelpererrors.ErrCIDGetParameterFailed
	}
	return value, nil
}

// 建立CID索引為0的字串(CommandName)
func MakeCIDCommandName(commandName string, isList bool, provider string) CustomIDCommandName {
	if isList {
		return CustomIDCommandName(commandName + "/list/" + provider)
	} else {
		return CustomIDCommandName(commandName + "//" + provider)
	}

}

// 建立事件為Page的CID
func MakeCIDPageComponent(label string, id string, value int, commandName CustomIDCommandName) *discordgo.Button {
	return &discordgo.Button{
		Label:    label,
		Style:    discordgo.PrimaryButton,
		CustomID: fmt.Sprintf("%s|%d|%s|%d", string(commandName), CustomIDTypePage, id, value),
	}
}

// 建立事件為AddHasPlayed的CID
func MakeCIDAddHasPlayedComponent(label string, id string, confirmMark bool, commandName CustomIDCommandName) *discordgo.Button {
	return &discordgo.Button{
		Label:    label,
		Style:    discordgo.PrimaryButton,
		CustomID: fmt.Sprintf("%s|%d|%s|%t", string(commandName), CustomIDTypeAddHasPlayed, id, confirmMark),
	}
}
