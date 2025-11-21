package vndb

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	kurohelpererrors "kurohelper/errors"
)

// GetVNByFuzzy 使用關鍵字模糊搜尋遊戲
func GetVNByFuzzy(keyword string) (*BasicResponse[GetVnUseIDResponse], error) {
	req := VndbCreate() // 建立基本request結構

	// 依照關鍵字的相關度排序
	reqSort := "searchrank"
	req.Sort = &reqSort

	// 限制回傳一筆結果
	reqResults := 1
	req.Results = &reqResults

	// 指定要取得的欄位
	titleFields := "title, alttitle"
	imageFields := "image.url"
	developersFields := "developers.name, developers.original, developers.aliases"
	nameFields := "titles.lang, titles.title, titles.official, titles.main"
	staffFields := "staff.name, staff.role, staff.aliases.name, staff.aliases.ismain"
	characterFields := "va.character.original, va.character.name, va.character.vns.role, va.character.vns.id"
	lengthFields := "length_minutes, length_votes"
	scoreFields := "average, rating, votecount"
	relationsFields := "relations.titles.title, relations.titles.main"

	allFields := []string{
		titleFields,
		imageFields,
		developersFields,
		nameFields,
		staffFields,
		characterFields,
		lengthFields,
		scoreFields,
		relationsFields,
	}
	req.Fields = strings.Join(allFields, ", ")

	// 設定搜尋條件
	req.Filters = []any{"search", "=", keyword}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	body, err := sendPostRequest("/vn", jsonData)
	if err != nil {
		return nil, err
	}

	var res BasicResponse[GetVnUseIDResponse]
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Results) == 0 {
		return nil, kurohelpererrors.ErrSearchNoContent
	}

	return &res, nil
}

// GetVNByID 使用 ID 精確搜尋遊戲
func GetVNByID(id string) (*BasicResponse[GetVnUseIDResponse], error) {
	req := VndbCreate() // 建立基本request結構

	// 不需要排序
	reqSort := ""
	req.Sort = &reqSort

	// 限制回傳一筆結果
	reqResults := 1
	req.Results = &reqResults

	// 指定要取得的欄位
	titleFields := "title, alttitle"
	imageFields := "image.url"
	developersFields := "developers.name, developers.original, developers.aliases"
	nameFields := "titles.lang, titles.title, titles.official, titles.main"
	staffFields := "staff.name, staff.role, staff.aliases.name, staff.aliases.ismain"
	characterFields := "va.character.original, va.character.name, va.character.vns.role, va.character.vns.id"
	lengthFields := "length_minutes, length_votes"
	scoreFields := "average, rating, votecount"
	relationsFields := "relations.titles.title, relations.titles.main"

	allFields := []string{
		titleFields,
		imageFields,
		developersFields,
		nameFields,
		staffFields,
		characterFields,
		lengthFields,
		scoreFields,
		relationsFields,
	}
	req.Fields = strings.Join(allFields, ", ")

	// 設定搜尋條件
	req.Filters = []any{"id", "=", id}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	body, err := sendPostRequest("/vn", jsonData)
	if err != nil {
		return nil, err
	}

	var res BasicResponse[GetVnUseIDResponse]
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Results) == 0 {
		return nil, kurohelpererrors.ErrSearchNoContent
	}

	return &res, nil
}

// GetRandomVN 隨機搜尋遊戲
func GetRandomVN() (*BasicResponse[GetVnUseIDResponse], error) {
	// 獲取遊戲id總數
	resStat, err := GetStats()
	if err != nil {
		return nil, err
	}

	req := VndbCreate() // 建立基本request結構

	// 不需要排序
	reqSort := ""
	req.Sort = &reqSort

	// 限制回傳一筆結果
	reqResults := 1
	req.Results = &reqResults

	// 指定要取得的欄位
	titleFields := "title, alttitle"
	imageFields := "image.url"
	developersFields := "developers.name, developers.original, developers.aliases"
	nameFields := "titles.lang, titles.title, titles.official, titles.main"
	staffFields := "staff.name, staff.role, staff.aliases.name, staff.aliases.ismain"
	characterFields := "va.character.original, va.character.name, va.character.vns.role, va.character.vns.id"
	lengthFields := "length_minutes, length_votes"
	scoreFields := "average, rating, votecount"
	relationsFields := "relations.titles.title, relations.titles.main"

	allFields := []string{
		titleFields,
		imageFields,
		developersFields,
		nameFields,
		staffFields,
		characterFields,
		lengthFields,
		scoreFields,
		relationsFields,
	}

	req.Fields = strings.Join(allFields, ", ")

	var res BasicResponse[GetVnUseIDResponse]
	for {
		// 產生隨機ID
		randomVNID := fmt.Sprintf("v%d", rand.Intn(resStat.VN))

		// 設定搜尋條件(隨機ID且投票數>=100)
		req.Filters = []any{"and", []any{"id", ">=", randomVNID}, []any{"votecount", ">=", "100"}}

		jsonData, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}

		body, err := sendPostRequest("/vn", jsonData)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &res)
		if err != nil {
			return nil, err
		}

		// 如果找到結果就回傳
		if len(res.Results) > 0 {
			break
		}
	}

	return &res, nil
}

// GetVnID 使用關鍵字搜尋遊戲ID列表(用於列表顯示)
func GetVnID(keyword string) (*[]GetVnIDUseListResponse, error) {
	// 建立基本request結構
	req := VndbCreate()

	// 依照關鍵字的相關度排序
	reqSort := "searchrank"
	req.Sort = &reqSort

	// 指定要取得的欄位
	req.Fields = "id, title, alttitle, developers.name, developers.original, developers.aliases"

	// 設定搜尋條件
	req.Filters = []any{"search", "=", keyword}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	body, err := sendPostRequest("/vn", jsonData)
	if err != nil {
		return nil, err
	}

	var res BasicResponse[GetVnIDUseListResponse]
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Results) == 0 {
		return nil, kurohelpererrors.ErrSearchNoContent
	}

	return &res.Results, nil
}
