package erogs

import (
	"encoding/json"
)

func GetCreatorByFuzzy(search string, opt string) (*FuzzySearchCreatorResponse, error) {
	if opt == "" {
		search = zhtwToJp(search)
	}
	sql, err := buildFuzzySearchCreatorSQL(search)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res FuzzySearchCreatorResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
