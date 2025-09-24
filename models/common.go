package models

import "time"

// 快取結構
type Cache struct {
	Value    interface{}
	ExpireAt time.Time
}

type VndbInteractionCustomID struct {
	CommandName string
	Page        int
	Key         string
}

type ErogsInteractionCustomID struct {
	CommandName string
	Type        int // 1: 翻頁 2: 排序
	Key         string
	Value       string
}
