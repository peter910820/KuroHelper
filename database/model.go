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

type (
	User struct {
		ID        string `gorm:"primaryKey"`
		Name      string
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	GameErogs struct {
		ID           int `gorm:"primaryKey;autoIncrement:false"`
		BrandErogsID int
		Title        string    `gorm:"unique"` // 遊戲名稱(批評空間)
		CreatedAt    time.Time `gorm:"autoCreateTime"`
		UpdatedAt    time.Time `gorm:"autoUpdateTime"`

		BrandErogs *BrandErogs `gorm:"foreignKey:BrandErogsID;references:ID"` // 單向 preload
	}

	BrandErogs struct {
		ID        int       `gorm:"primaryKey;autoIncrement:false"`
		Name      string    `gorm:"unique"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	UserGameErogs struct {
		UserID      string `gorm:"primaryKey"`
		GameErogsID int    `gorm:"primaryKey;autoIncrement:false"`
		HasPlayed   bool
		InWish      bool
		CreatedAt   time.Time `gorm:"autoCreateTime"`
		UpdatedAt   time.Time `gorm:"autoUpdateTime"`

		GameErogs *GameErogs `gorm:"foreignKey:GameErogsID;references:ID"` // 單向 preload
	}
)

func (ZhtwToJp) TableName() string {
	return "zhtw_to_jp"
}
