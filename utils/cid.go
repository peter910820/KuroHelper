package utils

import (
	kurohelpererrors "kurohelper/errors"
	"strconv"
	"strings"
)

type (
	NewCID []string

	CustomIDType int
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

// 獲取CID中的CommandName欄位，並確認是否是列表行為
//
// 這邊是安全行為，如果是沒有列表行為的狀況這邊會單純回傳False
func (cid NewCID) GetCommandNameIsList() bool {
	commandName := strings.Split([]string(cid)[0], "/")
	if len(commandName) == 1 {
		return false
	} else {
		if commandName[1] == "list" {
			return true
		} else {
			return false
		}
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
