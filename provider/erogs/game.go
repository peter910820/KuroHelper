package erogs

import (
	"encoding/json"
	"fmt"
)

func GetGameByFuzzy(search string, opt string) (*FuzzySearchGameResponse, error) {
	if opt == "" {
		search = zhtwToJp(search)
	}
	sql, err := buildFuzzySearchGameSQL(search)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res FuzzySearchGameResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return &res, nil
}
