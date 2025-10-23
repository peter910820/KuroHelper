package cache

import (
	kurohelperdb "github.com/peter910820/kurohelper-db"
	"github.com/peter910820/kurohelper-db/repository"
	"github.com/sirupsen/logrus"
)

// 先不管資料庫其他參數(預留用)
var (
	GuildDiscordAllowList = make(map[string]struct{})
	DmDiscordAllowList    = make(map[string]struct{})
)

func InitAllowList() {
	guildDiscordAllowList, err := repository.GetDiscordAllowListByKind(kurohelperdb.Dbs, "guild")
	if err != nil {
		logrus.Fatal(err)
	}

	dmDiscordAllowList, err := repository.GetDiscordAllowListByKind(kurohelperdb.Dbs, "dm")
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
