package erogs

// 只抓一筆(LIMIT 1)
type FuzzySearchCreatorResponse struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	TwitterUsername string `json:"twitter_username"`
	Blog            string `json:"blog"`
	Pixiv           *int   `json:"pixiv"`
	Games           []Game `json:"games"` // 參與過的遊戲
}

type Game struct {
	Gamename string     `json:"gamename"`
	SellDay  string     `json:"sellday"`
	Median   int        `json:"median"`
	CountAll int        `json:"count_all"`
	Shokushu []Shokushu `json:"shokushu"` // 有可能一個遊戲有多種身分
}

type Shokushu struct {
	Shubetu           int    `json:"shubetu"`
	ShubetuDetail     int    `json:"shubetu_detail"`
	ShubetuDetailName string `json:"shubetu_detail_name"` // *string
}

type FuzzySearchMusicResponse struct {
	ID             int            `json:"music_id"`
	MusicName      string         `json:"musicname"`
	PlayTime       string         `json:"playtime"`
	ReleaseDate    string         `json:"releasedate"`
	AvgTokuten     float64        `json:"avgtokuten"`
	TokutenCount   int            `json:"tokutencount"`
	Singers        []string       `json:"singers"`
	Lyrics         []string       `json:"lyrics"`
	Arrangments    []string       `json:"arrangments"`
	Compositions   []string       `json:"compositions"`
	GameCategories []GameCategory `json:"game_categories"`
	Album          []string       `json:"albums"`
}
type GameCategory struct {
	GameName string `json:"game_name"`
	Category string `json:"category"`
}
