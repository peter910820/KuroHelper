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
	GameName  string `json:"game_name"`
	GameModel string `json:"game_model"`
	Category  string `json:"category"`
}

type FuzzySearchGameResponse struct {
	ID                               int              `json:"id"`
	BrandID                          int              `json:"brandid"`
	BrandName                        string           `json:"brandname"`
	Gamename                         string           `json:"gamename"`
	SellDay                          string           `json:"sellday"`
	Model                            string           `json:"model"`
	Median                           string           `json:"median"`
	TokutenCount                     string           `json:"count2"`
	TotalPlayTimeMedian              string           `json:"total_play_time_median"`
	TimeBeforeUnderstandingFunMedian string           `json:"time_before_understanding_fun_median"`
	Okazu                            string           `json:"okazu"`
	Erogame                          string           `json:"erogame"`
	Genre                            string           `json:"genre"`
	BannerUrl                        string           `json:"banner_url"`
	SteamId                          string           `json:"steam"`
	VndbId                           string           `json:"vndb"`
	Shoukai                          string           `json:"shoukai"`
	Junni                            int              `json:"junni"`
	CreatorShubetu                   []Creatorshubetu `json:"shubetu_detail"`
}

type Creatorshubetu struct {
	ShubetuType       int    `json:"shubetu_type"`
	CreatorName       string `json:"creater_name"`
	ShubetuDetailType int    `json:"shubetu_detail_type"`
	ShubetuDetailName string `json:"shubetu_detail_name"`
}

type FuzzySearchBrandResponse struct {
	ID            int         `json:"id"`
	BrandName     string      `json:"brandname"`
	BrandFurigana string      `json:"brandfurigana"`
	URL           string      `json:"url"`
	Kind          string      `json:"kind"`
	Lost          bool        `json:"lost"`
	DirectLink    bool        `json:"directlink"` // 網站可不可用
	Median        int         `json:"median"`     // 該品牌的遊戲評分中位數(一天更新一次)
	Twitter       string      `json:"twitter"`
	Count2        int         `json:"count2"`
	CountAll      int         `json:"count_all"`
	Average2      int         `json:"average2"`
	Stdev         int         `json:"stdev"` // 標準偏差值(更新週期官方沒寫明確)
	GameList      []BrandGame `json:"gamelist"`
}

type BrandGame struct {
	ID       int    `json:"id"`
	GameName string `json:"gamename"`
	Furigana string `json:"furigana"`
	SellDay  string `json:"sellday"`
	Model    string `json:"model"`
	Median   int    `json:"median"` // 分數中位數(一天更新一次)
	Stdev    int    `json:"stdev"`  // 分數標準偏差值(一天更新一次)
	Count2   int    `json:"count2"` // 分數計算的樣本數
	VNDB     string `json:"vndb"`   // *string
}

type FuzzySearchListResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Model    string `json:"model"`
}

type FuzzySearchCharacterResponse struct {
	ID            int    `json:"id"`
	CharacterName string `json:"name"`
	Sex           string `json:"sex"`
	BloodType     string `json:"bloodtype"`
	Birthday      string `json:"birthday"`
	GameName      string `json:"gamename"`
	URL           string `json:"url"`
	FormalExplain string `json:"formal_explanation"`
	Age           string `json:"age"`
	Bust          string `json:"bust"`
	Waist         string `json:"waist"`
	Hip           string `json:"hip"`
	Height        string `json:"height"`
	Weight        string `json:"weight"`
	Cup           string `json:"cup"`
	Role          int    `json:"role"`
	CreatorName   string `json:"creater_name"`
}
