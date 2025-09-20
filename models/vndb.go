package models

// [VNDB]Request結構
//
// 對VNDB來講沒有必填項目，註解的必填項目是對於該專案來講的必填項目
// 所以預設值的部分可以完全不傳
//
// 這邊結構是根據需要的去對應，不是VNDB的完整結構
type VndbRequest struct {
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
type VndbResponse[T any] struct {
	Results           []T           `json:"results"`
	More              bool          `json:"more"`
	Count             int           `json:"count"`
	CompactFilters    string        `json:"compact_filters"`
	NormalizedFilters []interface{} `json:"normalized_filters"`
}

// [VNDB]使用ID查詢指定遊戲Response
type VndbGetVnUseIDResponse struct {
	ID            string                  `json:"id"`
	Title         string                  `json:"title"`
	Alttitle      string                  `json:"alttitle"`
	Average       float64                 `json:"average"`
	Rating        float64                 `json:"rating"`
	Votecount     int                     `json:"votecount"`
	LengthMinutes int                     `json:"length_minutes"`
	LengthVotes   int                     `json:"length_votes"`
	Developers    []VndbDeveloperResponse `json:"developers"`
	Relations     []VndbRelationResponse  `json:"relations"`
	Staff         []VndbStaffResponse     `json:"staff"`
	Titles        []VndbTitleResponse     `json:"titles"`
	Va            []VndbVaResponse        `json:"va"`
	Image         VndbImageResponse       `json:"image"`
}

// [VNDB]查詢品牌Response
type VndbProducerSearchResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Original    *string                `json:"original"`
	Aliases     []string               `json:"aliases"`
	Lang        string                 `json:"lang"`
	Type        string                 `json:"type"`
	Description *string                `json:"description"`
	Extlinks    []VndbExtlinksResponse `json:"extlinks"`
}

/* basic type start */

// [VNDB]品牌(發行單位)Response
type VndbDeveloperResponse struct {
	Aliases  []string `json:"aliases"`
	Name     string   `json:"name"`
	Original string   `json:"original"`
}

// [VNDB]關聯Response
type VndbRelationResponse struct {
	ID     string                      `json:"id"`
	Titles []VndbRelationTitleResponse `json:"titles"`
}

// [VNDB]關聯的標題Response
type VndbRelationTitleResponse struct {
	Title string `json:"title"`
}

type VndbStaffResponse struct {
	ID      string                   `json:"id"`
	Name    string                   `json:"name"`
	Role    string                   `json:"role"`
	Aliases []VndbStaffAliasResponse `json:"aliases"`
}

type VndbStaffAliasResponse struct {
	Name string `json:"name"`
}

type VndbTitleResponse struct {
	Lang     string `json:"lang"`
	Main     bool   `json:"main"`
	Official bool   `json:"official"`
	Title    string `json:"title"`
}

type VndbVaResponse struct {
	Character VndbCharacterResponse `json:"character"`
}

type VndbCharacterResponse struct {
	ID       string            `json:"id"`
	Original string            `json:"original"`
	Vns      []VndbVnsResponse `json:"vns"`
}

type VndbVnsResponse struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type VndbImageResponse struct {
	Url string `json:"url"`
}

// [VNDB]外部連結Response
type VndbExtlinksResponse struct {
	Url   string `json:"url"`
	Label string `json:"label"`
	Name  string `json:"name"`
	ID    string `json:"id"`
}

/* basic type end */

// vndb request factory
func VndbCreate() *VndbRequest {
	results := 100
	return &VndbRequest{
		Results: &results,
	}
}
