package vndb

import (
	"encoding/json"
	"strings"

	internalerrors "kurohelper/errors"
	vndbmodels "kurohelper/models/vndb"
)

func GetStaffByFuzzy(keyword string, roleType string) (*vndbmodels.BasicResponse[vndbmodels.StaffSearchResponse], error) {
	req := vndbmodels.VndbCreate()

	filters := []interface{}{}
	if roleType != "" {
		filters = append(filters, "and")
		// 傳進來的直接就是API篩選項規格的字串
		filters = append(filters, []string{"type", "=", roleType})
		filters = append(filters, []string{"search", "=", keyword})
	} else {
		filters = []interface{}{"search", "=", keyword}
	}

	req.Filters = filters

	basicFields := "id, aid, ismain, name, original, lang, gender, description"
	extlinksFields := "extlinks{url, label, name, id}"
	aliasesFields := "aliases{aid, name, latin, ismain}"

	allFields := []string{
		basicFields,
		extlinksFields,
		aliasesFields,
	}

	req.Fields = strings.Join(allFields, ", ")

	jsonStaff, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := sendPostRequest("/staff", jsonStaff)
	if err != nil {
		return nil, err
	}

	var res vndbmodels.BasicResponse[vndbmodels.StaffSearchResponse]
	err = json.Unmarshal(r, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Results) == 0 {
		return nil, internalerrors.ErrVndbNoResult
	}

	return &res, nil
}
