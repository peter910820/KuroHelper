package database

import "time"

type (
	ZhtwToJp struct {
		ZhTw      string    `gorm:"primaryKey;size:1"` // 繁體中文漢字
		Jp        string    `gorm:"size:1;not null"`   // 日文漢字
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	// 誠也對應表，專門針對極端狀況去對應
	SeiyaCorrespond struct {
		GameName  string    `gorm:"primaryKey"`
		SeiyaURL  string    `gorm:"not null"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}
)

func (ZhtwToJp) TableName() string {
	return "zhtw_to_jp"
}
