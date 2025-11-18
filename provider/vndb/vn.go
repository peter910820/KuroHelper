package vndb

import (
	"encoding/json"
	"strings"

	kurohelpererrors "kurohelper/errors"
)

func GetVn(keyword string, isID bool, isRandom bool) (*BasicResponse[GetVnUseIDResponse], error) {
	req := VndbCreate()
	reqSort := "searchrank"
	req.Filters = []any{"search", "=", keyword}
	if isID {
		req.Filters = []any{"id", "=", keyword}
		reqSort = ""
		reqResults := 1
		req.Results = &reqResults
	}
	if isRandom {
		req.Filters = []any{"and", []any{"id", ">=", keyword}, []any{"votecount", ">=", "100"}}
		reqSort = ""
		reqResults := 1
		req.Results = &reqResults
	}
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

	req.Sort = &reqSort
	req.Fields = strings.Join(allFields, ", ")

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

func GetVnID(keyword string) (*[]GetVnIDUseListResponse, error) {
	req := VndbCreate()
	reqSort := "searchrank"
	req.Filters = []any{"search", "=", keyword}
	req.Fields = "id, title, alttitle, developers.name, developers.original, developers.aliases"
	req.Sort = &reqSort

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
