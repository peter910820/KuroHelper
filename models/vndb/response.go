package vndb

// 使用ID查詢指定遊戲Response
type GetVnUseIDResponse struct {
	ID            string              `json:"id"`
	Title         string              `json:"title"`
	Alttitle      string              `json:"alttitle"`
	Average       float64             `json:"average"`
	Rating        float64             `json:"rating"`
	Votecount     int                 `json:"votecount"`
	LengthMinutes int                 `json:"length_minutes"`
	LengthVotes   int                 `json:"length_votes"`
	Developers    []DeveloperResponse `json:"developers"`
	Relations     []RelationResponse  `json:"relations"`
	Staff         []StaffResponse     `json:"staff"`
	Titles        []TitleResponse     `json:"titles"`
	Va            []VaResponse        `json:"va"`
	Image         ImageResponse       `json:"image"`
}

// producer Response
type ProducerSearchResponse struct {
	Producer BasicResponse[ProducerSearchProducerResponse]
	Vn       BasicResponse[ProducerSearchVnResponse]
}

// 查詢品牌API(Producer)
type ProducerSearchProducerResponse struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Original    string             `json:"original"` // *string
	Aliases     []string           `json:"aliases"`
	Lang        string             `json:"lang"`
	Type        string             `json:"type"`
	Description string             `json:"description"` // *string
	Extlinks    []ExtlinksResponse `json:"extlinks"`
}

// 查詢品牌API(vn)
type ProducerSearchVnResponse struct {
	Title         string  `json:"title"`
	Alttitle      string  `json:"alttitle"`
	Average       float64 `json:"average"`
	Rating        float64 `json:"rating"`
	Votecount     int     `json:"votecount"`
	LengthMinutes int     `json:"length_minutes"`
	LengthVotes   int     `json:"length_votes"`
}

// staff Response
//
// 統一字串不使用指標
type StaffSearchResponse struct {
	ID          string               `json:"id"`          // vndbid
	AID         int                  `json:"aid"`         // alias id
	IsMain      bool                 `json:"ismain"`      // 是否是主要名字
	Name        string               `json:"name"`        // 羅馬拼音名字
	Original    string               `json:"original"`    // 原文名, 可能為 null
	Lang        string               `json:"lang"`        // 主要語言
	Gender      string               `json:"gender"`      // 性別, 可能為 null
	Description string               `json:"description"` // 可能有格式化代碼
	ExtLinks    []ExtlinksResponse   `json:"extlinks"`    // 外部連結
	Aliases     []StaffAliasResponse `json:"aliases"`     // 別名清單
}
