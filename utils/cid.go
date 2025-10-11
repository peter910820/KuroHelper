package utils

import (
	kurohelpererrors "kurohelper/errors"
	"strconv"
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
	return []string(cid)[0]
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

func (cid AddWishCID) GetPageIndex() (bool, error) {

	value, err := strconv.ParseBool([]string(cid.NewCID)[3])
	if err != nil {
		return false, kurohelpererrors.ErrCIDGetParameterFailed
	}
	return value, nil
}

func (cid AddHasPlayedCID) GetPageIndex() (bool, error) {

	value, err := strconv.ParseBool([]string(cid.NewCID)[3])
	if err != nil {
		return false, kurohelpererrors.ErrCIDGetParameterFailed
	}
	return value, nil
}
