package bangumi

import (
	"encoding/json"
	"fmt"
	kurohelpererrors "kurohelper/errors"
	"os"
	"sort"

	"github.com/sirupsen/logrus"
)

func GetCharacterByFuzzy(keyword string) (*Character, error) {
	reqCharacter := BangumiCharacterCreate()
	reqCharacter.Keyword = keyword
	jsonCharacter, err := json.Marshal(reqCharacter)
	if err != nil {
		return nil, err
	}

	// 獲取角色基本資訊
	limit := 1
	offset := 0
	url := fmt.Sprintf("%s/v0/search/characters?limit=%d&offset=%d", os.Getenv("BANGUMI_ENDPOINT"), limit, offset)
	r, err := sendPostRequest(url, jsonCharacter)
	if err != nil {
		return nil, err
	}
	logrus.Printf("API Response: %s", string(r))
	var res CharacterSearchResponse
	err = json.Unmarshal(r, &res)
	if err != nil {
		return nil, err
	}
	// 檢查是否有結果
	if len(res.Data) == 0 {
		return nil, kurohelpererrors.ErrSearchNoContent
	}
	// sort.Slice(res.Data, func(i, j int) bool {
	// 	return res.Data[i].Stat.Collects > res.Data[j].Stat.Collects
	// })

	resCharacter := res.Data[0]
	// 獲取角色相關資訊

	url = fmt.Sprintf("%s/v0/characters/%d/persons", os.Getenv("BANGUMI_ENDPOINT"), resCharacter.ID)
	r, err = sendGetRequest(url)
	if err != nil {
		return nil, err
	}
	var resCharacterRelatedPerson []CharacterRelatedPersonResponse
	err = json.Unmarshal(r, &resCharacterRelatedPerson)
	if err != nil {
		return nil, err
	}
	character := ParseCharacterResponse(&resCharacter, &resCharacterRelatedPerson) // 將infobox裡面的資訊提取出來
	return character, nil
}

func ParseCharacterResponse(resCharacter *CharacterResponse, resCharacterRelatedPerson *[]CharacterRelatedPersonResponse) *Character {
	character := NewCharacter()
	// 解析角色基本資訊
	character.ID = resCharacter.ID
	character.Name = resCharacter.Name
	character.Image = resCharacter.Images.Medium // 使用 Medium 大小的圖片
	character.Summary = resCharacter.Summary

	for _, info := range resCharacter.Infobox {
		switch info.Key {
		case "简体中文名":
			character.NameCN = getStringValue(info.Value)
		case "别名":
			character.Aliases = parseAliases(info.Value)
		case "性别":
			character.Gender = getStringValue(info.Value)
		case "生日":
			character.BirthDay = getStringValue(info.Value)
		case "血型":
			character.BloodType = getStringValue(info.Value)
		case "身高":
			character.Height = getStringValue(info.Value)
		case "体重":
			character.Weight = getStringValue(info.Value)
		case "BWH":
			character.BWH = getStringValue(info.Value)
		case "年龄":
			character.Age = getStringValue(info.Value)
		default:
			otherValue := getStringValue(info.Value)
			if otherValue != "wrong type" {
				character.Other = append(character.Other, fmt.Sprintf("%s: %s", info.Key, otherValue))
			}
		}
	}

	// 排序：主角 > 配角 > 客串 > 其他
	sort.Slice(*resCharacterRelatedPerson, func(i, j int) bool {
		return getRolePriority((*resCharacterRelatedPerson)[i].Role) < getRolePriority((*resCharacterRelatedPerson)[j].Role)
	})

	// 使用 map 去重聲優
	cvSet := make(map[string]bool)
	for _, item := range *resCharacterRelatedPerson {
		cvSet[item.ActorName] = true
		character.Game = append(character.Game, fmt.Sprintf("%s (%s)", item.SubjectName, item.Role))
	}
	for name := range cvSet {
		character.CV = append(character.CV, name)
	}

	if len(character.CV) == 0 {
		character.CV = []string{"未收錄"}
	}
	if len(character.Game) == 0 {
		character.Game = []string{"未收錄"}
	}
	if len(character.Other) == 0 {
		character.Other = []string{"無"}
	}
	if len(character.Aliases) == 0 {
		character.Aliases = []string{"無"}
	}

	return character
}

// getStringValue 安全地從 any 類型中提取字串值
func getStringValue(value any) string {
	if str, ok := value.(string); ok {
		return str
	}
	return "wrong type"
}

// parseAliases 解析別名陣列為 map
func parseAliases(value any) []string {
	aliases := []string{}

	// 先將 any 轉成 JSON bytes，再解析成結構體
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return aliases
	}

	var aliasArr []CharacterAlias
	err = json.Unmarshal(jsonBytes, &aliasArr)
	if err != nil {
		return aliases
	}

	for _, alias := range aliasArr {
		if alias.Key != "" && alias.Value != "" {
			aliases = append(aliases, fmt.Sprintf("%s:%s", alias.Key, alias.Value))
		}
	}

	return aliases
}

// getRolePriority 返回角色類型的優先級（數字越小越優先）
func getRolePriority(role string) int {
	switch role {
	case "主角":
		return 1
	case "配角":
		return 2
	case "客串":
		return 3
	default:
		return 999 // 其他角色類型排在最後
	}
}
