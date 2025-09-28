package erogs

type SearchBrandResponse struct {
	ID            int    `json:"id"`
	BrandName     string `json:"brandname"`
	BrandFurigana string `json:"brandfurigana"`
	URL           string `json:"url"`
	Kind          string `json:"kind"`
	Lost          bool   `json:"lost"`
	DirectLink    bool   `json:"directlink"` // 網站可不可用
	Median        int    `json:"median"`     // 該品牌的遊戲評分中位數(一天更新一次)
	Twitter       string `json:"twitter"`
	Count2        int    `json:"count2"`
	CountAll      int    `json:"count_all"`
	Average2      int    `json:"average2"`
	Stdev         int    `json:"stdev"` // 標準偏差值
	// GameList      Game   `json:"game"`
}
