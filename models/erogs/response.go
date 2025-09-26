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
	CountAll int        `json:"count2"`
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
	AvgTokuten     float64        `json:"avg_tokuten"`
	TokutenCount   int            `json:"tokuten_count"`
	Singers        string         `json:"singer_name"`
	Lyrics         string         `json:"lyric_name"`
	Arrangments    string         `json:"arrangement_name"`
	Compositions   string         `json:"composition_name"`
	GameCategories []GameCategory `json:"game_categories"`
	Album          string         `json:"album_name"`
}
type GameCategory struct {
	GameName string `json:"game_name"`
	Category string `json:"category"`
}
