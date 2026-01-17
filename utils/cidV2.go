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
		cacheId    string
		behaviorID BehaviorID
		value      string
	}

	BehaviorID string

	// 翻頁CID
	PageCIDV2 struct {
		CacheId    string
		BehaviorID BehaviorID
		Value      int
	}
)

const (
	// PageBehavior Value會是int
	PageBehavior BehaviorID = "P"
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

/*
 * CID產生相關
 */

// 產生page的CID
func MakePageCIDV2(index int, cacheID string, disable bool) string {
	if disable {
		return fmt.Sprintf("%s:P:99", cacheID)
	}
	return fmt.Sprintf("%s:P:%d", cacheID, index)
}
