package vndb

// [VNDB]Request結構
//
// 對VNDB來講沒有必填項目，註解的必填項目是對於該專案來講的必填項目
// 所以預設值的部分可以完全不傳
//
// 這邊結構是根據需要的去對應，不是VNDB的完整結構
type BasicRequest struct {
	Filters           []interface{} `json:"filters"` // 必填
	Fields            string        `json:"fields"`  // 必填
	Sort              *string       `json:"sort,omitempty"`
	Reverse           *bool         `json:"reverse,omitempty"`
	Results           *int          `json:"results,omitempty"`
	Page              *int          `json:"page,omitempty"`
	Count             *bool         `json:"count,omitempty"`
	CompactFilters    *bool         `json:"compact_filters,omitempty"`
	NormalizedFilters *bool         `json:"normalized_filters,omitempty"`
}

// [VNDB]Response結構
type BasicResponse[T any] struct {
	Results           []T           `json:"results"`
	More              bool          `json:"more"`
	Count             int           `json:"count"`
	CompactFilters    string        `json:"compact_filters"`
	NormalizedFilters []interface{} `json:"normalized_filters"`
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
}

// 創作者結構
type StaffResponse struct {
	ID      string               `json:"id"`
	Name    string               `json:"name"`
	Role    string               `json:"role"`    // 角色類型
	Aliases []StaffAliasResponse `json:"aliases"` // 別名
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
	Character CharacterResponse `json:"character"`
}

type CharacterResponse struct {
	ID       string        `json:"id"`
	Original string        `json:"original"`
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
