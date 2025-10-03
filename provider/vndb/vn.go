package vndb

import (
	"encoding/json"
	"strings"

	kurohelpererrors "kurohelper/errors"
)

func GetVnUseID(brandid string) (*BasicResponse[GetVnUseIDResponse], error) {
	req := VndbCreate()

	req.Filters = []interface{}{
		"id", "=", brandid,
	}

	titleFields := "title, alttitle"
	imageFields := "image.url"
	developersFields := "developers.name, developers.original, developers.aliases"
	nameFields := "titles.lang, titles.title, titles.official, titles.main"
	staffFields := "staff.name, staff.role, staff.aliases.name"
	characterFields := "va.character.original, va.character.vns.role"
	lengthFields := "length_minutes, length_votes"
	scoreFields := "average, rating, votecount"
	relationsFields := "relations.titles.title"

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
