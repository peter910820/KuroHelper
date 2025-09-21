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
