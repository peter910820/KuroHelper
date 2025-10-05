package erogs

import (
	"encoding/json"
	"fmt"
)

func GetGameByFuzzy(search string) (*FuzzySearchGameResponse, error) {
	searchJP := zhtwToJp(search)
	sql, err := buildFuzzySearchGameSQL(search, searchJP)
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
