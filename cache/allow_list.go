package cache

import (
	"github.com/sirupsen/logrus"

	kurohelperdb "github.com/peter910820/kurohelper-db/v2"
)

// 先不管資料庫其他參數(預留用)
var (
	GuildDiscordAllowList = make(map[string]struct{})
	DmDiscordAllowList    = make(map[string]struct{})
)

func InitAllowList() {
	guildDiscordAllowList, err := kurohelperdb.GetDiscordAllowListByKind("guild")
	if err != nil {
		logrus.Fatal(err)
	}

	dmDiscordAllowList, err := kurohelperdb.GetDiscordAllowListByKind("dm")
	if err != nil {
		logrus.Fatal(err)
	}

	// 存進快取
	for _, g := range guildDiscordAllowList {
		GuildDiscordAllowList[g.ID] = struct{}{}
	}
	for _, d := range dmDiscordAllowList {
		GuildDiscordAllowList[d.ID] = struct{}{}
	}
}
