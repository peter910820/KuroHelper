package bangumi

type CharacterRequest struct {
	Keyword string `json:"keyword"`
	Filter  Filter `json:"filter"`
}

// CharacterSearchResponse 是 Bangumi API 的原始回應
type CharacterSearchResponse struct {
	Data   []CharacterResponse `json:"data"`
	Total  int                 `json:"total"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
}
type CharacterResponse struct {
	BirthMon  any             `json:"birth_mon"`
	Gender    any             `json:"gender"`
	BirthDay  any             `json:"birth_day"`
	BirthYear any             `json:"birth_year"`
	BloodType any             `json:"blood_type"`
	Images    CharacterImages `json:"images"`
	Summary   string          `json:"summary"`
	Name      string          `json:"name"`
	Infobox   []CharacterInfo `json:"infobox"`
	Stat      CharacterStat   `json:"stat"`
	ID        int             `json:"id"`
	Locked    bool            `json:"locked"`
	Type      int             `json:"type"`
	NSFW      bool            `json:"nsfw"`
}

type CharacterRelatedPersonResponse struct {
	Images        CharacterImages `json:"images"`
	ActorName     string          `json:"name"`
	SubjectName   string          `json:"subject_name"`
	SubjectNameCN string          `json:"subject_name_cn"`
	SubjectType   int             `json:"subject_type"`
	SubjectID     int             `json:"subject_id"`
	Role          string          `json:"staff"`
	ID            int             `json:"id"`
	Type          int             `json:"type"`
}

type CharacterImages struct {
	Small  string `json:"small"`
	Grid   string `json:"grid"`
	Large  string `json:"large"`
	Medium string `json:"medium"`
}

type CharacterStat struct {
	Comments int `json:"comments"`
	Collects int `json:"collects"`
}
type CharacterInfo struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type CharacterAlias struct {
	Key   string `json:"k"`
	Value string `json:"v"`
}

type Filter struct {
	NSFW bool `json:"nsfw"`
}

type Character struct { // 彙總角色資訊
	ID        int
	Name      string
	NameCN    string
	Aliases   []string
	Age       string
	Image     string
	Summary   string
	Gender    string
	BirthDay  string
	Height    string
	Weight    string
	BWH       string
	BloodType string
	Other     []string
	Game      []string
	CV        []string
}
