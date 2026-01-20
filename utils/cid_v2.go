package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type (
	// CID原型
	//
	// 這個原型型別只是方便後續轉換，直接拿本體使用會有型別不安全問題
	CIDV2 struct {
		behaviorID BehaviorID
		cacheId    string
		value      string
	}

	BehaviorID string

	// 翻頁CID
	PageCIDV2 struct {
		BehaviorID BehaviorID
		CacheId    string
		Value      int
	}

	// 選單CID
	SelectMenuCIDV2 struct {
		BehaviorID BehaviorID
		CacheId    string
		Value      string
	}

	// 回到主頁CID
	BackToHomeCIDV2 struct {
		BehaviorID BehaviorID
		CacheId    string
		// 回到主頁CID不需要Value
	}
)

const (
	// PageBehavior Value會是int
	PageBehavior BehaviorID = "P"
	// SelectMenuBehavior Value會是string(選擇後從Discord API獲得)
	SelectMenuBehavior BehaviorID = "S"
	// BackToHomeBehavior 不會有Value
	BackToHomeBehavior BehaviorID = "H"
)

var (
	ErrCIDV2ParseFailed      = errors.New("utils: cidv2 parse failed")
	ErrCIDV2ParseValueFailed = errors.New("utils: cidv2 parse value failed")
)

// 將字串嘗試轉型成CIDV2原型格式
//
// 檢查CIDV2的格式是否正確
func ParseCIDV2(target string) (*CIDV2, error) {
	parts := strings.Split(target, ":")
	if len(parts) != 3 {
		return nil, ErrCIDV2ParseFailed
	}

	return &CIDV2{
		cacheId:    parts[0],
		behaviorID: BehaviorID(parts[1]),
		value:      parts[2],
	}, nil
}

// 從CIDV2獲取BehaviorID
func (c CIDV2) GetBehaviorID() BehaviorID {
	return c.behaviorID
}

func (c CIDV2) ToPageCIDV2() (*PageCIDV2, error) {
	v, err := strconv.Atoi(c.value)
	if err != nil {
		return nil, ErrCIDV2ParseValueFailed
	}

	return &PageCIDV2{
		CacheId:    c.cacheId,
		BehaviorID: c.behaviorID,
		Value:      v,
	}, nil
}

func (c CIDV2) ToSelectMenuCIDV2() *SelectMenuCIDV2 {
	return &SelectMenuCIDV2{
		CacheId:    c.cacheId,
		BehaviorID: c.behaviorID,
		Value:      c.value,
	}
}

func (c CIDV2) ToBackToHomeCIDV2() *BackToHomeCIDV2 {
	return &BackToHomeCIDV2{
		CacheId:    c.cacheId,
		BehaviorID: c.behaviorID,
	}
}

// 修改Value值(SelectMenuBehavior時使用)
func (c *CIDV2) ChangeValue(value string) {
	c.value = value
}

/*
 * CID產生相關
 */

// 產生page的CID
//
// CID標示符是P
func MakePageCIDV2(index int, cacheID string, disable bool) string {
	if disable {
		return fmt.Sprintf("%s:P:99", cacheID)
	}
	return fmt.Sprintf("%s:P:%d", cacheID, index)
}

// 產生select menu的CID
//
// 產生select menu的CID時不需要先預留Value，Value會在選單選擇時才設定(Discord會自動設定)
//
// CID標示符是S
func MakeSelectMenuCIDV2(cacheID string) string {
	return fmt.Sprintf("%s:S:", cacheID)
}

// 產生回到主頁的CID
//
// CID標示符是H
func MakeBackToHomeCIDV2(cacheID string) string {
	return fmt.Sprintf("%s:H:", cacheID)
}
