package vndb

// [VNDB]Request結構
//
// 對VNDB來講沒有必填項目，註解的必填項目是對於該專案來講的必填項目
// 所以預設值的部分可以完全不傳
//
// 這邊結構是根據需要的去對應，不是VNDB的完整結構
type BasicRequest struct {
	Filters           []any   `json:"filters"` // 必填
	Fields            string  `json:"fields"`  // 必填
	Sort              *string `json:"sort,omitempty"`
	Reverse           *bool   `json:"reverse,omitempty"`
	Results           *int    `json:"results,omitempty"`
	Page              *int    `json:"page,omitempty"`
	Count             *bool   `json:"count,omitempty"`
	CompactFilters    *bool   `json:"compact_filters,omitempty"`
	NormalizedFilters *bool   `json:"normalized_filters,omitempty"`
}

// [VNDB]Response結構
type BasicResponse[T any] struct {
	Results           []T    `json:"results"`
	More              bool   `json:"more"`
	Count             int    `json:"count"`
	CompactFilters    string `json:"compact_filters"`
	NormalizedFilters []any  `json:"normalized_filters"`
}

// [VNDB]品牌(發行單位)Response
type DeveloperResponse struct {
	Aliases  []string `json:"aliases"`
	Name     string   `json:"name"`
	Original string   `json:"original"`
}

// [VNDB]關聯Response
type RelationResponse struct {
	ID     string                  `json:"id"`
	Titles []RelationTitleResponse `json:"titles"`
}

// [VNDB]關聯的標題Response
type RelationTitleResponse struct {
	Title string `json:"title"`
	Main  bool   `json:"main"`
}

// 創作者結構
type StaffResponse struct {
	ID       string               `json:"id"`
	Name     string               `json:"name"`
	Original string               `json:"original"` // 原文名
	Role     string               `json:"role"`     // 角色類型
	Aliases  []StaffAliasResponse `json:"aliases"`  // 別名
}

// 創作者別名結構
type StaffAliasResponse struct {
	AID    int    `json:"aid"`
	Name   string `json:"name"`
	Latin  string `json:"latin"`
	IsMain bool   `json:"ismain"` // 是否是主要別名
}

type TitleResponse struct {
	Lang     string `json:"lang"`
	Main     bool   `json:"main"`
	Official bool   `json:"official"`
	Title    string `json:"title"`
}

type VaResponse struct {
	Staff     StaffResponse     `json:"staff"`
	Character CharacterResponse `json:"character"`
}

type CharacterResponse struct {
	ID       string        `json:"id"`
	Original string        `json:"original"`
	Name     string        `json:"name"`
	Vns      []VnsResponse `json:"vns"`
}

type VnsResponse struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type ImageResponse struct {
	Url string `json:"url"`
}

// [VNDB]外部連結Response
type ExtlinksResponse struct {
	Url   string `json:"url"`
	Label string `json:"label"`
	Name  string `json:"name"`
	ID    string `json:"id"`
}

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

type GetVnIDUseListResponse struct {
	ID         string              `json:"id"`
	Title      string              `json:"title"`
	Alttitle   string              `json:"alttitle"`
	Developers []DeveloperResponse `json:"developers"`
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

type Stats struct {
	Chars     int `json:"chars"`
	Producers int `json:"producers"`
	Releases  int `json:"releases"`
	Staff     int `json:"staff"`
	Tags      int `json:"tags"`
	Traits    int `json:"traits"`
	VN        int `json:"vn"`
}

type CharacterSearchResponse struct { // 角色搜尋Response結構
	ID          string                      `json:"id"`          // vndbid
	Name        string                      `json:"name"`        // 名稱
	Original    string                      `json:"original"`    // 原文名稱，可能為 null
	Aliases     []string                    `json:"aliases"`     // 別名列表
	Description string                      `json:"description"` // 描述，可能為 null，可能包含格式化代碼
	Image       CharacterImage              `json:"image"`       // 圖片，可能為 null
	BloodType   string                      `json:"blood_type"`  // 血型："a", "b", "ab" 或 "o"，可能為 null
	Height      int                         `json:"height"`      // 身高（公分），可能為 null
	Weight      int                         `json:"weight"`      // 體重（公斤），可能為 null
	Bust        int                         `json:"bust"`        // 胸圍（公分），可能為 null
	Waist       int                         `json:"waist"`       // 腰圍（公分），可能為 null
	Hips        int                         `json:"hips"`        // 臀圍（公分），可能為 null
	Cup         string                      `json:"cup"`         // 罩杯："AAA", "AA" 或任何單一字母，可能為 null
	Age         *int                        `json:"age"`         // 年齡（歲），可能為 null
	Birthday    [2]int                      `json:"birthday"`    // 生日 [月, 日]，可能為 null
	Sex         [2]string                   `json:"sex"`         // 性別 [表面性別, 真實性別]，可能為 null，值："m", "f", "b", "n"
	Gender      [2]string                   `json:"gender"`      // 自我性別認同 [非劇透, 劇透]，可能為 null，值："m", "f", "o", "a"
	VNs         []CharacterSearchVnResponse `json:"vns"`
	Vas         []string
	//	Traits      []CharacterSearchTraitResponse `json:"traits"`
}

// 角色圖片結構（與視覺小說圖片欄位相同，但不包含縮圖）
type CharacterImage struct {
	URL string `json:"url"` // 圖片 URL
}

type CharacterSearchVnResponse struct { // 獲得Role欄位
	Title    string          `json:"title"` // 羅馬拼音
	Alttitle string          `json:"alttitle"`
	Titles   []TitleResponse `json:"titles"`
	Spoiler  int             `json:"spoiler"` // 劇透等級
	Role     string          `json:"role"`    // main/primary/side/appears
}

type VnSearchCharacterResponse struct { // VN查角色Response結構
	ID       string `json:"id"`       // vndbid
	Role     string `json:"role"`     // main/primary/side/appears
	Name     string `json:"name"`     // 角色名稱
	Original string `json:"original"` // 角色原文名稱，可能為 null
}

type CharacterSearchTraitResponse struct { // 角色特徵結構
	ID          string   `json:"id"`          // vndbid
	Name        string   `json:"name"`        // 特徵名稱（應與 group_name 一起顯示）
	Aliases     []string `json:"aliases"`     // 別名列表
	Description string   `json:"description"` // 描述，可能包含格式化代碼
	Searchable  bool     `json:"searchable"`  // 是否可搜尋
	Applicable  bool     `json:"applicable"`  // 是否適用
	Sexual      bool     `json:"sexual"`      // 是否為性相關特徵
	GroupID     string   `json:"group_id"`    // vndbid，所屬群組 ID
	GroupName   string   `json:"group_name"`  // 群組名稱（頂層父特徵）
	Spoiler     int      `json:"spoiler"`     // 劇透等級
	Lie         bool     `json:"lie"`
}
