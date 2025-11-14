package vndb

import (
	"encoding/json"
	kurohelpererrors "kurohelper/errors"
	"regexp"
	"strings"
)

func GetCharacterByFuzzy(keyword string, isID bool) (*CharacterSearchResponse, error) {
	reqCharacter := VndbCreate()
	reqCharacterSort := "searchrank"
	reqCharacter.Filters = []any{"search", "=", keyword}
	if isID {
		reqCharacter.Filters = []any{"id", "=", keyword}
		reqCharacterSort = ""
	}

	reqCharacterResults := 1
	reqCharacter.Sort = &reqCharacterSort
	reqCharacter.Results = &reqCharacterResults
	basicFields := "id, name, original, aliases, description, image.url, blood_type, height, weight, bust, waist, hips, cup, age, birthday, sex, gender"
	vnsFields := "vns.title, vns.alttitle, vns.spoiler, vns.role, vns.titles.title, vns.titles.main"
	allFields := []string{
		basicFields,
		vnsFields,
	}
	reqCharacter.Fields = strings.Join(allFields, ", ")
	jsonCharacter, err := json.Marshal(reqCharacter)
	if err != nil {
		return nil, err
	}

	r, err := sendPostRequest("/character", jsonCharacter)
	if err != nil {
		return nil, err
	}

	var resCharacters BasicResponse[CharacterSearchResponse]
	err = json.Unmarshal(r, &resCharacters)
	if err != nil {
		return nil, err
	}
	if len(resCharacters.Results) == 0 {
		return nil, kurohelpererrors.ErrSearchNoContent
	}

	reqVn := VndbCreate()
	characterIDFilter := []interface{}{"id", "=", resCharacters.Results[0].ID}
	reqVn.Filters = []interface{}{"character", "=", characterIDFilter}
	reqVn.Fields = "va.staff.name, va.staff.original, va.character.id"
	jsonVn, err := json.Marshal(reqVn)
	if err != nil {
		return nil, err
	}
	r, err = sendPostRequest("/vn", jsonVn)
	if err != nil {
		return nil, err
	}
	var resVn BasicResponse[GetVnUseIDResponse]
	err = json.Unmarshal(r, &resVn)
	if err != nil {
		return nil, err
	}
	var vasMap = make(map[string]bool) // 去重
	var vas []string
	for _, vn := range resVn.Results {
		for _, va := range vn.Va {
			if va.Character.ID == resCharacters.Results[0].ID {
				if va.Staff.Original != "" {
					vasMap[va.Staff.Original] = true
				} else {
					vasMap[va.Staff.Name] = true
				}
			}
		}
	}
	if len(vasMap) == 0 {
		resCharacters.Results[0].Vas = []string{"未收錄"}
	} else {
		for va := range vasMap {
			vas = append(vas, va)
		}
		resCharacters.Results[0].Vas = vas
	}

	return &resCharacters.Results[0], nil
}

func ConvertBBCodeToMarkdown(text string) string {
	// 1. 處理配對的 URL 標籤
	reURL := regexp.MustCompile(`\[url=(.+?)\](.+?)\[/url\]`)
	text = reURL.ReplaceAllString(text, "[$2]($1)")

	// 2. 處理配對的 spoiler 標籤（支援多行）
	reSpoiler := regexp.MustCompile(`(?s)\[spoiler\](.+?)\[/spoiler\]`)
	text = reSpoiler.ReplaceAllString(text, "||$1||")

	// 3. 清理未配對的殘留標籤
	text = strings.ReplaceAll(text, "[spoiler]", "")
	text = strings.ReplaceAll(text, "[/spoiler]", "")

	// 4. 將角色ID轉換成連結[Sara](/c40662)
	reCharacterID := regexp.MustCompile(`\[(.+?)\]\(/c(\d+?)\)`)
	text = reCharacterID.ReplaceAllString(text, "[$1](https://vndb.org/c$2)")
	return strings.TrimSpace(text)
}

func GetCharacterListByFuzzy(keyword string) (*[]CharacterSearchResponse, error) {
	reqCharacter := VndbCreate()
	reqCharacter.Filters = []any{"search", "=", keyword}
	reqCharacterSort := "searchrank"
	reqCharacter.Sort = &reqCharacterSort
	basicFields := "id, name, original"
	vnsFields := "vns.title, vns.alttitle, vns.spoiler, vns.role, vns.titles.title, vns.titles.main"
	allFields := []string{
		basicFields,
		vnsFields,
	}
	reqCharacter.Fields = strings.Join(allFields, ", ")
	jsonCharacter, err := json.Marshal(reqCharacter)
	if err != nil {
		return nil, err
	}
	r, err := sendPostRequest("/character", jsonCharacter)
	if err != nil {
		return nil, err
	}
	var resCharacters BasicResponse[CharacterSearchResponse]
	err = json.Unmarshal(r, &resCharacters)
	if err != nil {
		return nil, err
	}

	return &resCharacters.Results, nil
}
